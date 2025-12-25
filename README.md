# 🕵️ ShadowLog: Advanced Activity Monitor

![Python](https://img.shields.io/badge/Python-3.x-blue?style=flat&logo=python)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Educational-orange)

**ShadowLog** is a robust, Python-based keystroke logging tool designed for educational purposes and security testing (Red Team/Blue Team scenarios). It features active window tracking, local file logging with rotation, and optional remote data exfiltration via Discord Webhooks.

---

## ⚠️ ETHICAL DISCLAIMER

**THIS TOOL IS FOR EDUCATIONAL AND AUTHORIZED USE ONLY.**

* **Do not** use this tool on any system you do not own or have explicit written permission to monitor.
* **Do not** use this tool for malicious purposes, stalking, or unauthorized surveillance.
* The developers are not responsible for any damage or legal consequences caused by the misuse of this software.

---

## 🚀 Features

* **⌨️ Keystroke Capture:** Logs all standard keys and special keys (Space, Enter, Backspace, etc.).
* **🖥️ Active Window Tracking:** Logs the title of the active window (e.g., `[Chrome - Gmail]`) to provide context.
* **📡 Remote Exfiltration:** Optionally sends log batches to Discord Webhooks.
* **📂 Local Storage:** Saves logs to a local directory with automatic file rotation.
* **⏯️ Live Control:** Toggle logging on/off instantly using the `F9` key.

---

## 📦 Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/YOUR_USERNAME/ShadowLog.git](https://github.com/YOUR_USERNAME/ShadowLog.git)
    cd ShadowLog
    ```

2.  **Install dependencies:**
    ```bash
    pip install -r requirements.txt
    ```

---

## ⚙️ Easy Configuration

To set up the tool, open `shadowlog.py` in any text editor and scroll to the bottom. Look for the **EASY CONFIGURATION** section.

### 1. Set Log Folder
Find the line `LOG_PATH` and paste the folder path where you want to save your logs.
```python
# Example:
LOG_PATH = Path(r"E:\KeyLogger\Logs")
2. Enable Discord Logging (Optional)
Find the line DISCORD_WEBHOOK_URL.

To Enable: Paste your Webhook URL inside the quotes.

To Disable: Set it to None.

Python

# Example:
DISCORD_WEBHOOK_URL = "[https://discord.com/api/webhooks/](https://discord.com/api/webhooks/)..."
🛠️ Usage
Run the script:

Bash

python shadowlog.py
(Note: On macOS/Linux, you may need sudo permissions to capture keystrokes).

Controls:

Start/Stop: The tool starts automatically. Press F9 to pause or resume logging.

Exit: Press CTRL+C in the terminal to stop the tool and save the final logs.

📂 Requirements
Python 3.8+

pynput

requests

psutil

pywin32 (Windows only)

📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
