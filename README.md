# 🕵️ ShadowLog: Advanced Activity Monitor

![Python](https://img.shields.io/badge/Python-3.x-blue?style=flat&logo=python)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Educational-orange)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-lightgrey)

**ShadowLog** is a Python-based activity monitoring tool built for **cybersecurity education and authorized security testing**.  
It demonstrates how keystroke monitoring, contextual logging, and controlled data transmission operate in **Red Team and Blue Team scenarios**.

Unlike consumer monitoring software, ShadowLog is designed to **teach detection, analysis, and defense techniques**, not surveillance.

---

## 📖 Table of Contents
- [Features](#-features)
- [How It Works](#-how-it-works)
- [Installation](#-installation)
- [Usage](#-usage)
- [Project Structure](#-project-structure)
- [Configuration](#-configuration)
- [Disclaimer](#-disclaimer)
- [License](#-license)

---

## 🚀 Features

- **Keystroke Capture**  
  Logs standard and special keys (Enter, Space, Backspace, etc.).

- **Active Window Tracking**  
  Records the currently focused application to provide context for each log entry.

- **Local Log Storage**  
  Saves logs to disk with automatic file rotation for stability.

- **Optional Remote Logging**  
  Supports sending log batches to a Discord Webhook (disabled by default).

- **Live Control Toggle**  
  Instantly pause or resume logging using a configurable hotkey (`F9`).

---

## 🧠 How It Works

ShadowLog follows a simple monitoring pipeline:

1. **Hook** – Registers a keyboard listener using system-level hooks.
2. **Context Capture** – Retrieves the active window title during key events.
3. **Log Management** – Writes structured logs locally with rotation support.
4. **Transmission (Optional)** – Sends log batches to a remote endpoint if enabled.

This architecture mirrors real-world monitoring techniques used by both attackers and defenders, making it valuable for **threat modeling and detection research**.

---

## 📦 Installation

### Prerequisites
- Python 3.8+
- `pip`

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/ShadowLog.git
   cd ShadowLog
   ```

2. **Install dependencies**

   ```bash
   pip install -r requirements.txt
   ```

---

## 🛠️ Usage

> Administrator / root privileges may be required for full keystroke capture.

**Windows**

```powershell
python shadowlog.py
```

**Linux / macOS**

```bash
sudo python3 shadowlog.py
```

---

## 🎛️ Controls

| Action                 | Key        |
| ---------------------- | ---------- |
| Pause / Resume Logging | `F9`       |
| Exit Safely            | `CTRL + C` |

Logs are flushed and saved automatically on exit.

---

## 📂 Project Structure

```text
ShadowLog/
│
├── shadowlog.py        # Main monitoring engine
├── requirements.txt   # Python dependencies
└── README.md          # Documentation
```

---

## ⚙️ Configuration

Edit `shadowlog.py` and locate the **EASY CONFIGURATION** section.

### Set Log Directory

```python
LOG_PATH = Path(r"E:\ShadowLog\Logs")
```

### Enable / Disable Remote Logging

```python
DISCORD_WEBHOOK_URL = "https://discord.com/api/webhooks/..."
# or
DISCORD_WEBHOOK_URL = None
```

Remote logging is **optional and disabled by default**.

---

## ⚠️ Disclaimer

**FOR EDUCATIONAL AND AUTHORIZED USE ONLY**

* Do not use on systems you do not own or have permission to monitor
* Do not use for spying, stalking, or malicious surveillance
* Always test in controlled environments (VMs recommended)

The author is **not responsible** for misuse or legal consequences.

---

## 📄 License

This project is licensed under the **MIT License**.

---

<p align="center">
<strong>Developed by BGx (Devansh Agarwal)</strong><br>
<em>Cybersecurity Enthusiast & Developer</em>
</p>
