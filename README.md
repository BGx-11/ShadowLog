# 🕵️ ShadowLog — Advanced Activity Monitor

![Python](https://img.shields.io/badge/Python-3.x-blue?style=flat&logo=python)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Educational-orange)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-lightgrey)

**ShadowLog** is an advanced, Python-based activity monitoring tool built strictly for **educational use, cybersecurity learning, and authorized security testing**.  
It is designed to help students and professionals understand how keystroke monitoring, context-aware logging, and controlled data exfiltration work in **Red Team and Blue Team environments**.

> ⚠️ This project focuses on *awareness and defense*, not misuse.

---

## 📖 Table of Contents
- [Overview](#-overview)
- [Ethical Disclaimer](#-ethical-disclaimer)
- [Key Features](#-key-features)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Controls](#-controls)
- [Requirements](#-requirements)
- [License](#-license)

---

## 🧠 Overview

Traditional security tools often detect threats **after damage has occurred**.  
ShadowLog demonstrates how activity monitoring tools function internally, allowing learners to:

- Understand how keystrokes can be intercepted
- Learn how attackers correlate keystrokes with active windows
- Analyze how log rotation and data transmission work
- Build stronger detection and prevention strategies

This makes ShadowLog a **learning tool**, not a weapon.

---

## ⚠️ Ethical Disclaimer

**THIS SOFTWARE IS FOR EDUCATIONAL AND AUTHORIZED USE ONLY.**

- ❌ Do **NOT** use this tool on systems you do not own
- ❌ Do **NOT** use without explicit written permission
- ❌ Do **NOT** use for spying, stalking, or malicious surveillance

By using this software, **you take full responsibility** for complying with all applicable laws.  
The authors are **not liable** for misuse or legal consequences.

---

## 🚀 Key Features

- **⌨️ Keystroke Logging**  
  Captures standard and special keys (Enter, Backspace, Space, etc.)

- **🖥️ Active Window Awareness**  
  Logs the currently focused application window to provide context  
  *(Example: `[Chrome — Gmail]`)*

- **📂 Local File Logging**  
  Stores logs locally with **automatic file rotation** for stability

- **📡 Optional Remote Logging**  
  Supports sending log batches to a Discord Webhook (disabled by default)

- **⏯️ Live Control Toggle**  
  Instantly pause or resume logging using a hotkey (`F9`)

---

## 📦 Installation

### 1️⃣ Clone the Repository
```bash
git clone https://github.com/YOUR_USERNAME/ShadowLog.git
cd ShadowLog
````

### 2️⃣ Install Dependencies

```bash
pip install -r requirements.txt
```

---

## ⚙️ Configuration (Easy Setup)

Open `shadowlog.py` and scroll to the **EASY CONFIGURATION** section.

### 📁 Set Log Directory

Choose where logs will be stored locally.

```python
LOG_PATH = Path(r"E:\ShadowLog\Logs")
```

---

### 📡 Enable Discord Logging (Optional)

To enable remote logging, paste your webhook URL.
To disable it, set the value to `None`.

```python
DISCORD_WEBHOOK_URL = "https://discord.com/api/webhooks/..."
# or
DISCORD_WEBHOOK_URL = None
```

> 🔐 Tip: Keep webhook logging disabled during local testing.

---

## 🛠️ Usage

Run the tool using:

```bash
python shadowlog.py
```

> On Linux/macOS, elevated permissions may be required to capture keystrokes.

---

## 🎮 Controls

| Action                | Key        |
| --------------------- | ---------- |
| Start / Pause Logging | `F9`       |
| Exit Safely           | `CTRL + C` |

Logs are automatically saved before exit.

---

## 📂 Requirements

* Python **3.8+**
* `pynput`
* `requests`
* `psutil`
* `pywin32` *(Windows only)*

---

## 📄 License

This project is licensed under the **MIT License**.
You are free to modify and learn from it — **not misuse it**.

---

## 🧩 Final Note

ShadowLog exists to **teach how monitoring tools work so they can be detected, mitigated, and defended against**.

If you’re learning:

* Cybersecurity 🛡️
* Ethical Hacking 🧑‍💻
* Malware Analysis 🧬
* Blue Team Defense 🔵

— this project is for you.

```
