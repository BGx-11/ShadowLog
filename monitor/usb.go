package monitor

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// usbMonitor detects USB drive insertion and removal events.
type usbMonitor struct {
	knownDrives map[string]bool
	callback    func(string)
}

var (
	procGetLogicalDrives = kernel32.NewProc("GetLogicalDrives")
	procGetDriveTypeW    = kernel32.NewProc("GetDriveTypeW")
	procGetVolumeInfoW   = kernel32.NewProc("GetVolumeInformationW")
)

const (
	driveRemovable = 2 // DRIVE_REMOVABLE
	driveFixed     = 3 // DRIVE_FIXED (some USB drives report as fixed)
)

// newUSBMonitor creates a new USB drive monitoring instance.
func newUSBMonitor(callback func(string)) *usbMonitor {
	m := &usbMonitor{
		knownDrives: make(map[string]bool),
		callback:    callback,
	}
	// Snapshot current drives so we don't log existing ones.
	m.scanDrives(true)
	return m
}

// start begins polling for USB drive changes.
func (um *usbMonitor) start() {
	for {
		time.Sleep(5 * time.Second)
		um.scanDrives(false)
	}
}

// scanDrives checks all logical drives and detects insertions/removals.
func (um *usbMonitor) scanDrives(initialScan bool) {
	bitmask, _, _ := procGetLogicalDrives.Call()
	if bitmask == 0 {
		return
	}

	currentDrives := make(map[string]bool)

	for i := 0; i < 26; i++ {
		if bitmask&(1<<uint(i)) == 0 {
			continue
		}

		driveLetter := string(rune('A'+i)) + ":\\"
		driveLetterPtr, _ := syscall.UTF16PtrFromString(driveLetter)

		// Check drive type.
		driveType, _, _ := procGetDriveTypeW.Call(uintptr(unsafe.Pointer(driveLetterPtr)))
		if driveType != driveRemovable {
			continue // Only track removable drives
		}

		currentDrives[driveLetter] = true

		// Check if this is a new drive.
		if !initialScan && !um.knownDrives[driveLetter] {
			// New USB drive detected!
			label, serial := um.getVolumeInfo(driveLetter)
			ts := time.Now().Format("2006-01-02 15:04:05")
			logLine := fmt.Sprintf("[%s] [USB] 🔌 INSERTED: %s | Label: %s | Serial: %s",
				ts, driveLetter, label, serial)
			if um.callback != nil {
				um.callback(logLine)
			}
		}
	}

	// Check for removals.
	if !initialScan {
		for drive := range um.knownDrives {
			if !currentDrives[drive] {
				ts := time.Now().Format("2006-01-02 15:04:05")
				logLine := fmt.Sprintf("[%s] [USB] ⏏️ REMOVED: %s", ts, drive)
				if um.callback != nil {
					um.callback(logLine)
				}
			}
		}
	}

	// Update known drives.
	um.knownDrives = currentDrives
}

// getVolumeInfo retrieves the volume label and serial number for a drive.
func (um *usbMonitor) getVolumeInfo(driveLetter string) (string, string) {
	drivePtr, _ := syscall.UTF16PtrFromString(driveLetter)

	var volumeName [256]uint16
	var serialNumber uint32
	var maxComponentLen uint32
	var fsFlags uint32
	var fsName [256]uint16

	ret, _, _ := procGetVolumeInfoW.Call(
		uintptr(unsafe.Pointer(drivePtr)),
		uintptr(unsafe.Pointer(&volumeName[0])),
		uintptr(len(volumeName)),
		uintptr(unsafe.Pointer(&serialNumber)),
		uintptr(unsafe.Pointer(&maxComponentLen)),
		uintptr(unsafe.Pointer(&fsFlags)),
		uintptr(unsafe.Pointer(&fsName[0])),
		uintptr(len(fsName)),
	)

	label := "(No Label)"
	serial := "Unknown"

	if ret != 0 {
		name := syscall.UTF16ToString(volumeName[:])
		if strings.TrimSpace(name) != "" {
			label = name
		}
		serial = fmt.Sprintf("%08X", serialNumber)
	}

	return label, serial
}
