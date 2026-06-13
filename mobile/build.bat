@echo off
setlocal
title ShadowLog Mobile Build
echo ====================================================
echo  ShadowLog Mobile — Android Build System
echo ====================================================
echo.

:: Check for Android SDK
if not defined ANDROID_HOME (
    if not defined ANDROID_SDK_ROOT (
        echo ERROR: ANDROID_HOME or ANDROID_SDK_ROOT must be set.
        echo Set it to your Android SDK installation path.
        echo Example: set ANDROID_HOME=C:\Users\%USERNAME%\AppData\Local\Android\Sdk
        pause
        exit /b 1
    )
)

:: Navigate to mobile directory
cd /d "%~dp0"

echo [1/3] Cleaning previous build...
if exist app\build rmdir /s /q app\build

echo [2/3] Building Release APK...
call gradlew.bat assembleRelease
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ERROR: Release build failed.
    echo Make sure you have:
    echo   1. Android SDK installed (API 35)
    echo   2. JDK 17+ installed
    echo   3. ANDROID_HOME or ANDROID_SDK_ROOT set
    pause
    exit /b 1
)

echo [3/3] Locating output APK...
set APK_PATH=app\build\outputs\apk\release\app-release.apk
if exist %APK_PATH% (
    echo.
    echo ====================================================
    echo  BUILD SUCCESSFUL
    echo ====================================================
    echo  APK: %APK_PATH%
    echo.
    echo  To install on a connected device:
    echo    adb install -r %APK_PATH%
    echo ====================================================
) else (
    echo WARNING: APK not found at expected path.
    echo Check build output above for errors.
)

echo.
pause
