package monitor

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// wifiMonitor tracks Wi-Fi network connections for location correlation.
type wifiMonitor struct {
	lastSSID string
	callback func(string)
}

// newWifiMonitor creates a new Wi-Fi network monitoring instance.
func newWifiMonitor(callback func(string)) *wifiMonitor {
	return &wifiMonitor{
		callback: callback,
	}
}

// start begins polling the Wi-Fi interface for network changes.
func (wm *wifiMonitor) start() {
	// Check immediately on startup
	wm.checkNetwork()
	for {
		time.Sleep(60 * time.Second)
		wm.checkNetwork()
	}
}

// checkNetwork runs netsh to get the current Wi-Fi connection details.
func (wm *wifiMonitor) checkNetwork() {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return
	}

	info := parseNetshOutput(string(out))
	if info.ssid == "" {
		// Check if we were previously connected (disconnection event).
		if wm.lastSSID != "" {
			ts := time.Now().Format("2006-01-02 15:04:05")
			logLine := fmt.Sprintf("[%s] [WIFI] 📡 DISCONNECTED from: %s", ts, wm.lastSSID)
			wm.lastSSID = ""
			if wm.callback != nil {
				wm.callback(logLine)
			}
		}
		return
	}

	// Check if SSID changed (new connection or roaming).
	if info.ssid != wm.lastSSID {
		ts := time.Now().Format("2006-01-02 15:04:05")
		logLine := fmt.Sprintf("[%s] [WIFI] 📡 CONNECTED: SSID=%s | BSSID=%s | Signal=%s | Auth=%s | Channel=%s",
			ts, info.ssid, info.bssid, info.signal, info.auth, info.channel)
		wm.lastSSID = info.ssid
		if wm.callback != nil {
			wm.callback(logLine)
		}
	}
}

// wifiInfo holds parsed Wi-Fi connection details.
type wifiInfo struct {
	ssid    string
	bssid   string
	signal  string
	auth    string
	channel string
}

// parseNetshOutput extracts Wi-Fi details from netsh wlan show interfaces output.
func parseNetshOutput(output string) wifiInfo {
	var info wifiInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch {
		case key == "ssid" || key == "    ssid":
			// Only set SSID if it's the actual SSID field, not BSSID
			if !strings.Contains(key, "bssid") {
				info.ssid = val
			}
		case key == "bssid" || strings.HasSuffix(key, "bssid"):
			info.bssid = val
		case key == "signal" || strings.HasSuffix(key, "signal"):
			info.signal = val
		case key == "authentication" || strings.HasSuffix(key, "authentication"):
			info.auth = val
		case key == "channel" || strings.HasSuffix(key, "channel"):
			info.channel = val
		}
	}

	return info
}
