package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"shadowlog/config"

	"golang.org/x/sys/windows/registry"
)

var (
	lastHeartbeat = time.Now()
	heartbeatMu   sync.Mutex
)

// Status represents the current state of ShadowLog installation.
type Status struct {
	IsRunning   bool   // Whether the ShadowLog process is currently running.
	IsInstalled bool   // Whether ShadowLog is installed in startup.
	LogsCount   int    // Number of log lines in the local backup file.
}

func main() {
	// Uninstaller is a web-based tool to safely remove ShadowLog.
	// It checks current status and provides a removal option.
	// Currently Windows-only due to platform-specific commands.

	if runtime.GOOS != "windows" {
		fmt.Println("Uninstaller only supported on Windows.")
		return
	}

	// Synchronization for server lifecycle.
	var wg sync.WaitGroup
	wg.Add(1)

	// HTTP server on port 58293 to avoid conflict with setup (58291).
	mux := http.NewServeMux()
	server := &http.Server{Addr: ":58293", Handler: mux}

	// Watchdog goroutine to shutdown if heartbeat stops
	go func() {
		for {
			time.Sleep(5 * time.Second)
			heartbeatMu.Lock()
			elapsed := time.Since(lastHeartbeat)
			heartbeatMu.Unlock()
			if elapsed > 12*time.Second {
				fmt.Println("No heartbeat received for 12s. Shutting down Uninstaller...")
				server.Shutdown(context.Background())
				wg.Done()
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

	// Handle status display and uninstall action.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := getStatus()
		switch r.Method {
		case "GET":
			// Display status and uninstall form.
			tmpl := template.Must(template.New("uninstaller").Parse(htmlContent))
			tmpl.Execute(w, status)
		case "POST":
			// Perform complete removal.
			uninstall()
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, successHTML)
			// Shutdown server after showing success message.
			go func() {
				time.Sleep(2 * time.Second)
				server.Shutdown(context.Background())
				wg.Done()
			}()
		}
	})

	go func() {
		url := "http://localhost:58293"
		fmt.Printf("Opening uninstaller at: %s\n", url)
		openBrowser(url)
	}()

	go server.ListenAndServe()
	wg.Wait()
}

// getStatus checks the current installation and runtime status of ShadowLog.
// Returns a Status struct with process state, registry entry, and log count.
func getStatus() Status {
	var s Status

	// 1. Check if registry entry exists for startup.
	var exeName string
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		// Check multiple possible registry values to ensure deep cleaning of all versions
		targetValues := []string{"Windows Update Service", "Shadow Log", "OneDrive Helper", "onedrive hgelper", "OneDriveUpdate", "ShadowLog"}
		for _, v := range targetValues {
			if val, _, err := k.GetStringValue(v); err == nil {
				s.IsInstalled = true
				exeName = filepath.Base(val)
				fmt.Printf("Registry found: %s -> %s\n", v, val)
				break // Stop at the first one found
			}
		}
	}

	// 2. Fallback check for known process names.
	if exeName == "" {
		exeName = "WinUpdateSvc.exe"
	}

	// 3. Check if the process is running using tasklist.
	// We check for both the registry name and the current setup tool name.
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+exeName)
	out, _ := cmd.Output()
	s.IsRunning = strings.Contains(string(out), exeName)

	if !s.IsRunning {
		// Secondary check for Setup/Main tool names.
		// Monitor names across all branding iterations
		checkNames := []string{"WinUpdateSvc.exe", "ShadowLog.exe", "Shadow Log.exe", "ShadowLog_Setup.exe", "OneDriveHelper.exe", "onedrive hgelper.exe", "OneDriveService.exe"}
		for _, name := range checkNames {
			cmdProcess := exec.Command("tasklist", "/FI", "IMAGENAME eq "+name)
			processOut, _ := cmdProcess.Output()
			if strings.Contains(string(processOut), name) {
				s.IsRunning = true
				break
			}
		}
	}

	// 4. Count lines in local log file.
	logPath := config.GetLogPath()
	if data, err := os.ReadFile(logPath); err == nil {
		lines := strings.Count(string(data), "\n")
		// The first line is always the configuration, so we only count if lines > 1.
		if lines > 1 {
			s.LogsCount = lines - 1
		} else {
			s.LogsCount = 0
		}
	}

	return s
}

// uninstall performs a complete removal of ShadowLog from the system.
// Kills running processes, removes registry entries, and deletes all files.
func uninstall() {
	names := []string{
		"WinUpdateSvc.exe",
		"ShadowLog.exe", "Shadow Log.exe", "ShadowLog_Setup.exe", "shadowlog.exe",
		"OneDriveHelper.exe", "onedrive hgelper.exe", "OneDriveService.exe",
	}

	// Try to get name from registry for all possible versions.
	kReg, kErr := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if kErr == nil {
		targetValues := []string{"Windows Update Service", "Shadow Log", "OneDrive Helper", "onedrive hgelper", "OneDriveUpdate", "ShadowLog"}
		for _, v := range targetValues {
			if val, _, err := kReg.GetStringValue(v); err == nil {
				names = append(names, filepath.Base(val))
			}
		}
		kReg.Close()
	}

	for _, name := range names {
		// Use /F to force, /IM for image name, /T to kill child processes.
		cmd := exec.Command("taskkill", "/F", "/IM", name, "/T")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
	}

	// 2. Remove registry entry for startup.
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err == nil {
		defer k.Close()
		// Clean all possible registry values
		targetValues := []string{"Windows Update Service", "Shadow Log", "OneDrive Helper", "onedrive hgelper", "OneDriveUpdate", "ShadowLog"}
		for _, v := range targetValues {
			k.DeleteValue(v)
		}
	}

	// 2.5 Remove secondary persistence (Scheduled Task)
	// Clean all possible scheduled tasks
	tasks := []string{"Windows Update Telemetry", "Shadow Log Reporting", "OneDrive Helper Reporting", "onedrive hgelper update"}
	for _, t := range tasks {
		exec.Command("schtasks", "/delete", "/tn", t, "/f").Run()
	}

	// 3. Remove single binary storage file and sync lock.
	os.Remove(config.GetStoragePath())
	os.Remove(config.GetSyncPath())

	// 4. Remove registry-based configuration.
	config.DeleteConfigFromRegistry()
}

// openBrowser opens the given URL in the system's default web browser.
// Windows-specific implementation using rundll32.
func openBrowser(url string) {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Shadow Log: Complete Removal</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #ff3d71;
            --primary-glow: rgba(255, 61, 113, 0.3);
            --bg: #020617;
            --card-bg: rgba(15, 23, 42, 0.6);
            --card-border: rgba(255, 255, 255, 0.08);
            --text-main: #f8fafc;
            --text-dim: #94a3b8;
            --danger: #ff3d71;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background-color: var(--bg);
            background-image: 
                radial-gradient(circle at 50% 0%, rgba(255, 61, 113, 0.1) 0%, transparent 50%),
                radial-gradient(circle at 50% 100%, rgba(255, 61, 113, 0.05) 0%, transparent 50%);
            color: var(--text-main);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 24px;
        }

        .container { 
            max-width: 460px; 
            width: 100%;
            background: var(--card-bg);
            backdrop-filter: blur(40px);
            border: 1px solid var(--card-border);
            border-radius: 28px;
            padding: 56px;
            text-align: center;
            box-shadow: 0 40px 80px -20px rgba(0, 0, 0, 0.8);
            animation: fadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1);
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }
        
        .icon {
            width: 72px;
            height: 72px;
            background: rgba(255, 61, 113, 0.1);
            border-radius: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            margin: 0 auto 24px;
            color: var(--danger);
            box-shadow: 0 0 30px rgba(255, 61, 113, 0.2);
            border: 1px solid rgba(255, 61, 113, 0.2);
        }

        h1 { font-size: 1.75rem; font-weight: 800; margin-bottom: 8px; letter-spacing: -0.04em; color: #fff; }
        .subtitle { color: var(--text-dim); font-size: 0.9375rem; margin-bottom: 32px; font-weight: 500; }

        .warning-box {
            background: rgba(255, 61, 113, 0.05);
            border: 1px dashed rgba(255, 61, 113, 0.3);
            border-radius: 16px;
            padding: 20px;
            margin-bottom: 32px;
            text-align: left;
            font-size: 0.8125rem;
            line-height: 1.6;
            color: #ff9999;
        }
        .warning-box strong { color: var(--danger); font-weight: 800; display: block; margin-bottom: 4px; text-transform: uppercase; letter-spacing: 0.05em; }

        .terminal-note {
            font-size: 0.75rem;
            color: var(--text-dim);
            margin-top: 16px;
            font-style: italic;
        }

        .status-card { 
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 20px;
            padding: 28px;
            margin-bottom: 32px;
            text-align: left;
        }
        .status-line { 
            display: flex; 
            justify-content: space-between; 
            margin-bottom: 16px; 
            font-size: 0.875rem;
            align-items: center;
        }
        .status-line:last-child { margin-bottom: 0; }
        .status-val { font-weight: 800; font-size: 0.6875rem; letter-spacing: 0.08em; padding: 4px 10px; border-radius: 8px; }
        
        @keyframes pulse {
            0% { transform: scale(1); opacity: 1; }
            50% { transform: scale(1.1); opacity: 0.7; }
            100% { transform: scale(1); opacity: 1; }
        }
        .dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; margin-right: 6px; vertical-align: middle; }
        .running-dot { background: var(--primary); box-shadow: 0 0 10px var(--primary); animation: pulse 2s infinite; }
        .stopped-dot { background: var(--text-dim); }

        .running-label { background: rgba(0, 191, 165, 0.1); color: var(--primary); border: 1px solid rgba(0,191,165,0.2); }
        .stopped-label { background: rgba(255, 255, 255, 0.05); color: var(--text-dim); border: 1px solid rgba(255,255,255,0.1); }

        .btn {
            width: 100%;
            background: rgba(255, 61, 113, 0.1);
            color: var(--danger);
            padding: 18px;
            border-radius: 14px;
            border: 1px solid rgba(255, 61, 113, 0.3);
            font-weight: 800;
            font-size: 1rem;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 12px;
        }
        .btn:hover { 
            transform: translateY(-2px);
            background: var(--danger);
            color: #fff;
            box-shadow: 0 10px 30px -10px rgba(255, 61, 113, 0.6);
        }
        .btn:active { transform: translateY(0); }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; transform: none; }
    </style>
</head>
<body>
    <div class="container" id="mainContainer">
        <div class="icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
        </div>
        <h1>Shadow Log</h1>
        <p class="subtitle">Complete system purge</p>
        
        <div class="warning-box">
            <strong>⚠️ Irreversible Action</strong>
            This will permanently remove the cloud integration service, all cached telemetry data, and system registry entries.
        </div>
        
        <div class="status-card">
            <div class="status-line">
                <span style="color: var(--text-dim)">Service Status:</span>
                <span class="status-val {{ if .IsRunning }}running-label{{ else }}stopped-label{{ end }}">
                    {{ if .IsRunning }}ACTIVE AGENT{{ else }}INACTIVE{{ end }}
                </span>
            </div>
            <div class="status-line">
                <span style="color: var(--text-dim)">Persistence:</span>
                <span class="status-val" style="background: rgba(255,255,255,0.05); color: #fff;">{{ if .IsInstalled }}REGISTERED{{ else }}CLEAN{{ end }}</span>
            </div>
            <div class="status-line">
                <span style="color: var(--text-dim)">Local Cache:</span>
                <span class="status-val" style="background: rgba(255,255,255,0.05); color: #fff;">{{ .LogsCount }} Artifacts</span>
            </div>
        </div>

        <form method="POST" onsubmit="handleUninstall(event)">
            <button type="submit" id="uninstallBtn" class="btn">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                Begin System Cleanup
            </button>
            <p class="terminal-note">Note: The terminal might flash briefly during cleanup.</p>
        </form>
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

        function handleUninstall(e) {
            const btn = document.getElementById('uninstallBtn');
            btn.disabled = true;
            btn.innerHTML = '<svg style="animation: spin 1s linear infinite" width="20" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M21 12a9 9 0 1 1-6.219-8.56"></path></svg> Purging Service...';
        }
    </script>
    <style>
        @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
    </style>
</body>
</html>
`

const successHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Removal Successful</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #020617;
            --card-bg: rgba(15, 23, 42, 0.6);
            --card-border: rgba(255, 255, 255, 0.08);
            --text-main: #f8fafc;
            --text-dim: #94a3b8;
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
                radial-gradient(circle at 50% 50%, rgba(255, 61, 113, 0.05) 0%, transparent 50%);
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
            background: rgba(0, 191, 165, 0.1);
            color: #00bfa5;
            border-radius: 50%;
            display: flex;
            justify-content: center;
            align-items: center;
            margin: 0 auto 32px;
            box-shadow: 0 0 20px rgba(0, 191, 165, 0.2);
        }
        h1 { font-size: 1.5rem; font-weight: 800; margin-bottom: 12px; letter-spacing: -0.04em; }
        p { color: var(--text-dim); font-size: 0.9375rem; line-height: 1.6; font-weight: 500; }
        @keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
    </style>
</head>
<body>
    <div class="card">
        <div class="check-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
        </div>
        <h1>Cleanup Complete</h1>
        <p>Shadow Log has been successfully purged. All local data and artifacts have been removed.</p>
        <p style="margin-top: 24px; font-size: 0.75rem; opacity: 0.5;">The cleanup agent is shutting down. You may now close this tab.</p>
    </div>
</body>
</html>
`
