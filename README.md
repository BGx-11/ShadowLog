# Shadow Log: Discrete Activity Analytics

![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Build](https://img.shields.io/badge/Build-Optimized-success?style=for-the-badge)

Shadow Log is a high-performance, native systems monitoring framework engineered for advanced cybersecurity research and authorized security demonstrations. Built entirely in Go with zero runtime dependencies, it provides an enterprise-grade approach to discrete activity capture and secure exfiltration.

---

> [!CAUTION]
> ### Legal and Ethical Warning
> **FOR AUTHORIZED EDUCATIONAL AND RESEARCH USE ONLY.**
> Unauthorized deployment of monitoring tools on systems you do not own or have explicit, documented permission to monitor is illegal and a violation of privacy laws (such as the Computer Fraud and Abuse Act). This software is provided "as-is" for security professionals and students to study low-level system hooks and secure exfiltration techniques. The developer assumes no liability for misuse, data loss, or legal consequences resulting from the use of this tool.

---

## Technical Features

### 🛡️ Advanced Stealth and OPSEC
-   **Unified Identity**: The system identifies as **`Shadow Log`** across all telemetry, registry entries, and scheduled tasks for transparent system management.
-   **Note**: During installation or uninstallation, the terminal may flash briefly. This is normal behavior for system service registration.

### 🌐 Universal Precision Capture
-   **International Keyboard Support**: Leverages the Windows `ToUnicode` API to ensure 100% accuracy across all regions, correctly capturing symbols, accents, and shifted characters.
-   **Context-Aware Metadata**: Every captured event is automatically correlated with the foreground `ProcessID`, Application Name, and active Window Title.
-   **Zero-Lag Performance**: Implemented Active Window Handle Caching; the system only queries process metadata when the foreground window actually changes, reducing API calls and CPU overhead to negligible levels.

### 🔐 Hardware-Bound Security
-   **Machine-Linked Encryption**: Logs are encrypted using **AES-256-GCM**. The encryption key is dynamically derived from the host's unique `MachineGuid`, meaning exfiltrated data is forensicly useless if moved to another system without the original GUID.
-   **Asynchronous Exfiltration**: A robust, buffered worker pool handles log transmission via Discord Webhooks, featuring masked browser User-Agents and intelligent network retry logic.

## 🚀 Deployment Guide

### Option 1: Quick Start (Security Analysts)
If you are an analyst and do not wish to compile from source, download the latest **`ShadowLog_Release.zip`** from the project root.
1.  **Extract**: Unzip the archive on the host machine.
2.  **Initialize**: Run `ShadowLog.exe` as Administrator to start the setup wizard.
    - *Note: A terminal window may flash a few times while hooks are being registered.*
3.  **Configure**: Enter your Discord Webhook URL and click **Test Connection**. This is required to receive your unique **Machine GUID**, which you will need later to decrypt logs.
4.  **Deploy**: Click **Initialize Monitor**. The service will initialize and run silently as `Shadow Log`.

### Option 2: Build from Source (Developers)

### Prerequisites
-   **Go 1.21+**: Required for compilation.
-   **Windows 10/11**: Target environment for native hooks.
-   **Administrator Rights**: Required for the initialization phase to register system hooks.

### Compilation
Generate optimized, stripped binaries using the following commands:

```powershell
# Build Core Monitor & Setup
go build -ldflags "-H=windowsgui -s -w" -o ShadowLog.exe main.go

# Build Forensic Decryptor
go build -ldflags "-H=windowsgui -s -w" -o Decryptor.exe decryptor/main.go

# Build System Uninstaller
go build -ldflags "-H=windowsgui -s -w" -o Uninstaller.exe uninstaller/main.go
```

---

## Forensic Reconstruction

1.  **Recovery**: Run `Decryptor.exe` on the host machine.
2.  **Decryption**: Input the unique **Machine GUID** sent during the test pulse to reconstruct and read local encrypted logs.
3.  **Decommission**: Run `Uninstaller.exe` to purge all processes, registry keys, scheduled tasks, and forensic system artifacts.

---

## System Specifications
- **Service Name**: `Shadow Log`
- **Mutex ID**: `Shadow_Log_Sync`
- **Registry Entry**: `HKCU\Software\Microsoft\Windows\CurrentVersion\Run\Shadow Log`
- **Scheduled Task**: `Shadow Log Reporting`
- **Encryption Standard**: AES-256-GCM (Device-Bound)

---

**Devansh Agarwal (BGx)**  
*Cyber Sec Student & Student Developer*  
[Portfolio: iambgx.in](https://iambgx.in)
