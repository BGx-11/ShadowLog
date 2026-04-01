@echo off
setlocal
title ShadowLog Builder

echo ----------------------------------------------------
echo ShadowLog Optimized Build System
echo ----------------------------------------------------

:: Ensure we are in the project root
if not exist main.go (
    echo Error: main.go not found. Please run from the project root.
    pause
    exit /b
)

:: Define build flags
set LDFLAGS=-H=windowsgui -s -w

echo [1/3] Building Core Monitor (ShadowLog.exe)...
go build -ldflags "%LDFLAGS%" -o ShadowLog.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build ShadowLog.exe
    exit /b
)

echo [2/3] Building Forensic Decryptor (Decryptor.exe)...
go build -ldflags "%LDFLAGS%" -o Decryptor.exe decryptor/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Decryptor.exe
    exit /b
)

echo [3/3] Building System Uninstaller (Uninstaller.exe)...
go build -ldflags "%LDFLAGS%" -o Uninstaller.exe uninstaller/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Uninstaller.exe
    exit /b
)

echo [4/4] Packaging Release (ShadowLog_Release.zip)...
if exist ShadowLog_Release.zip del ShadowLog_Release.zip
powershell -Command "Compress-Archive -Path ShadowLog.exe, Decryptor.exe, Uninstaller.exe, README.md -DestinationPath ShadowLog_Release.zip"
if %ERRORLEVEL% NEQ 0 (
    echo Failed to create ShadowLog_Release.zip
    exit /b
)

echo.
echo ----------------------------------------------------
echo Build Complete: All binaries generated successfully.
echo - ShadowLog.exe
echo - Decryptor.exe
echo - Uninstaller.exe
echo - ShadowLog_Release.zip
echo ----------------------------------------------------
echo.
pause
