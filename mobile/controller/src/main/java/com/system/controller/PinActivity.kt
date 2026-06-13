package com.system.controller

import android.content.Intent
import android.os.Bundle
import android.webkit.JavascriptInterface
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity

/**
 * PIN Lock Screen — entry gate for the Controller APK.
 *
 * Authenticates using the Monitor app's encryption password.
 * The user enters the SAME password they set in the Monitor's setup.
 * This is verified via ContentProvider (verify_password call method).
 *
 * If the Monitor has no password set, automatically bypasses to dashboard.
 * If the Monitor is not installed, shows an error message.
 */
class PinActivity : AppCompatActivity() {

    private lateinit var webView: WebView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Check if monitor is installed and has a password
        val monitorInstalled = MonitorClient.isMonitorInstalled(this)
        if (monitorInstalled) {
            val status = MonitorClient.getStatus(this)
            if (!status.hasPassword) {
                // No password set — skip PIN screen entirely
                launchDashboard()
                return
            }
        }

        webView = WebView(this).apply {
            settings.javaScriptEnabled = true
            settings.domStorageEnabled = true
            webViewClient = WebViewClient()
            addJavascriptInterface(PinBridge(), "Android")
        }
        setContentView(webView)
        loadPinScreen()
    }

    inner class PinBridge {
        @JavascriptInterface
        fun isMonitorInstalled(): Boolean = MonitorClient.isMonitorInstalled(this@PinActivity)

        @JavascriptInterface
        fun submitPassword(password: String): Boolean {
            if (!isMonitorInstalled()) return false
            val verified = MonitorClient.verifyPassword(this@PinActivity, password)
            if (verified) {
                runOnUiThread { launchDashboard() }
            }
            return verified
        }

        @JavascriptInterface
        fun vibrate(ms: Long) {
            try {
                val vibrator = getSystemService(android.os.Vibrator::class.java)
                if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
                    vibrator?.vibrate(android.os.VibrationEffect.createOneShot(ms, android.os.VibrationEffect.DEFAULT_AMPLITUDE))
                } else {
                    @Suppress("DEPRECATION")
                    vibrator?.vibrate(ms)
                }
            } catch (_: Exception) {}
        }
    }

    private fun launchDashboard() {
        startActivity(Intent(this, StatusActivity::class.java))
        finish()
    }

    private fun loadPinScreen() {
        val monitorInstalled = MonitorClient.isMonitorInstalled(this)
        val html = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
            <title>Authentication</title>
            <style>
                :root{--primary:#3b82f6;--bg:#0f172a;--card:#1e293b;--text:#f8fafc;--dim:#94a3b8;--border:rgba(255,255,255,0.1);--danger:#ef4444;--success:#10b981;--glass:rgba(30,41,59,0.85)}
                *{margin:0;padding:0;box-sizing:border-box}
                body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:var(--bg);color:var(--text);display:flex;flex-direction:column;align-items:center;justify-content:center;min-height:100vh;padding:24px;overflow:hidden;position:relative}

                /* Main content */
                .content{position:relative;z-index:2;display:flex;flex-direction:column;align-items:center;width:100%;max-width:340px}

                /* Glass card */
                .glass-card{background:var(--glass);backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);border:1px solid var(--border);border-radius:24px;padding:36px 28px;width:100%;text-align:center;box-shadow:0 8px 32px rgba(0,0,0,0.2);animation:cardEntry 0.6s cubic-bezier(0.16,1,0.3,1)}
                @keyframes cardEntry{from{opacity:0;transform:translateY(20px) scale(0.98)}to{opacity:1;transform:translateY(0) scale(1)}}

                /* Lock icon */
                .lock-container{width:72px;height:72px;border-radius:20px;background:rgba(59,130,246,0.1);border:1px solid rgba(59,130,246,0.2);display:flex;align-items:center;justify-content:center;margin:0 auto 24px;position:relative;transition:all 0.4s cubic-bezier(0.16,1,0.3,1);color:var(--primary)}
                .lock-container .lock-icon{font-size:32px;transition:all 0.3s}
                
                /* Success state */
                .lock-container.success{background:rgba(16,185,129,0.15);border-color:rgba(16,185,129,0.3);transform:scale(1.05);color:var(--success)}
                .lock-container.error{animation:lockShake 0.4s;border-color:rgba(239,68,68,0.4);background:rgba(239,68,68,0.1);color:var(--danger)}
                @keyframes lockShake{0%,100%{transform:translateX(0)}25%{transform:translateX(-6px)}75%{transform:translateX(6px)}}

                h1{font-size:22px;font-weight:700;margin-bottom:6px;letter-spacing:-0.02em}
                .sub{color:var(--dim);font-size:13px;margin-bottom:28px;text-align:center;line-height:1.5}

                /* Input */
                .input-wrap{width:100%;margin-bottom:20px;position:relative}
                .input-wrap input{width:100%;padding:16px;background:rgba(15,23,42,0.4);border:1px solid var(--border);border-radius:12px;color:var(--text);font-size:16px;font-weight:500;outline:none;text-align:center;letter-spacing:0.1em;transition:all 0.2s}
                .input-wrap input:focus{border-color:var(--primary);box-shadow:0 0 0 3px rgba(59,130,246,0.15)}
                .input-wrap input::placeholder{color:var(--dim);letter-spacing:normal}
                .input-wrap input.error{border-color:var(--danger);box-shadow:0 0 0 3px rgba(239,68,68,0.15)}
                .input-wrap input.success{border-color:var(--success);box-shadow:0 0 0 3px rgba(16,185,129,0.15)}

                /* Button */
                .submit-btn{width:100%;padding:16px;background:var(--primary);border:none;border-radius:12px;color:#ffffff;font-size:15px;font-weight:600;cursor:pointer;transition:all 0.2s;letter-spacing:0.01em}
                .submit-btn:active{transform:scale(0.98);opacity:0.9}
                .submit-btn:disabled{opacity:0.5;cursor:not-allowed;transform:none}
                .submit-btn.success{background:var(--success)}

                /* Message */
                .msg{margin-top:16px;font-size:12px;color:var(--dim);text-align:center;min-height:20px;line-height:1.4}
                .msg.error{color:var(--danger)}
                .msg.ok{color:var(--success)}

                /* Lockout timer */
                .lockout-bar{width:100%;height:4px;background:rgba(239,68,68,0.1);border-radius:2px;margin-top:12px;overflow:hidden;opacity:0;transition:opacity 0.3s}
                .lockout-bar.active{opacity:1}
                .lockout-bar .fill{height:100%;background:var(--danger);border-radius:2px;transition:width 0.1s linear;width:0}

                /* No monitor error */
                .no-monitor{background:var(--glass);backdrop-filter:blur(20px);border:1px solid var(--border);border-radius:20px;padding:32px;text-align:center;max-width:320px}
                .no-monitor .icon{font-size:40px;margin-bottom:16px}
                .no-monitor h2{font-size:18px;font-weight:600;color:var(--danger);margin-bottom:12px}
                .no-monitor p{font-size:13px;color:var(--dim);line-height:1.6}

                /* Dots indicator */
                .dots{display:flex;gap:8px;justify-content:center;margin-bottom:24px}
                .dot{width:8px;height:8px;border-radius:50%;background:rgba(148,163,184,0.3);transition:all 0.2s}
                .dot.filled{background:var(--primary)}
                .dot.error{background:var(--danger)}
            </style>
        </head>
        <body>
            <div class="bg-gradient"></div>
            ${if (!monitorInstalled) """
            <div class="content">
                <div class="no-monitor">
                    <div class="icon">
                        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--danger);"><path d="M8.5 2v4"/><path d="M15.5 2v4"/><path d="M16 8v2a4 4 0 0 1-4 4H4a4 4 0 0 1-4-4V8h16z"/><path d="M12 14v8"/></svg>
                    </div>
                    <h2>Monitor Not Found</h2>
                    <p>Install the System Service (Monitor) APK first, then set up an encryption password.<br><br>The Controller uses the same password to authenticate.</p>
                </div>
            </div>
            """ else """
            <div class="content">
                <div class="glass-card">
                    <div class="lock-container" id="lockBox">
                        <span class="lock-icon" id="lockIcon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                        </span>
                    </div>
                    <h1>Authentication</h1>
                    <div class="sub">Enter the encryption password<br>configured in the Monitor app</div>
                    <div class="dots" id="dots">
                        <div class="dot"></div><div class="dot"></div><div class="dot"></div><div class="dot"></div><div class="dot"></div>
                    </div>
                    <div class="input-wrap">
                        <input type="password" id="pwd" placeholder="Enter password" autocomplete="off" />
                    </div>
                    <button class="submit-btn" id="submitBtn" onclick="submit()">
                        Authenticate
                    </button>
                    <div class="msg" id="msg"></div>
                    <div class="lockout-bar" id="lockoutBar"><div class="fill" id="lockoutFill"></div></div>
                </div>
            </div>
            <script>
                var attempts = 0;
                var maxAttempts = 5;
                var locked = false;
                var lockoutInterval = null;

                var pwdInput = document.getElementById('pwd');
                var dots = document.querySelectorAll('.dot');

                pwdInput.addEventListener('keydown', function(e) {
                    if (e.key === 'Enter') submit();
                });

                // Update dots as user types
                pwdInput.addEventListener('input', function() {
                    var len = Math.min(pwdInput.value.length, dots.length);
                    dots.forEach(function(d, i) {
                        d.className = 'dot' + (i < len ? ' filled' : '');
                    });
                });

                function submit() {
                    if (locked) return;
                    var pwd = pwdInput.value;
                    if (!pwd) {
                        showMsg('Enter a password', 'error');
                        return;
                    }
                    var btn = document.getElementById('submitBtn');
                    btn.disabled = true;
                    btn.textContent = 'Verifying...';

                    setTimeout(function() {
                        var ok = Android.submitPassword(pwd);
                        if (ok) {
                            // Success animation
                            Android.vibrate(50);
                            document.getElementById('lockBox').className = 'lock-container success';
                            document.getElementById('lockIcon').innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 9.9-1"/></svg>';
                            pwdInput.className = 'success';
                            btn.className = 'submit-btn success';
                            btn.textContent = '✓ Authenticated';
                            showMsg('Access granted — loading dashboard...', 'ok');
                            dots.forEach(function(d) { d.className = 'dot filled'; });
                        } else {
                            attempts++;
                            Android.vibrate(100);

                            // Error animation
                            document.getElementById('lockBox').className = 'lock-container error';
                            pwdInput.className = 'error';
                            setTimeout(function() {
                                document.getElementById('lockBox').className = 'lock-container';
                                pwdInput.className = '';
                            }, 600);

                            pwdInput.value = '';
                            dots.forEach(function(d) {
                                d.className = 'dot error';
                                setTimeout(function() { d.className = 'dot'; }, 600);
                            });

                            if (attempts >= maxAttempts) {
                                startLockout();
                            } else {
                                showMsg('Wrong password — ' + (maxAttempts - attempts) + ' attempts remaining', 'error');
                                btn.disabled = false;
                                btn.textContent = 'Authenticate';
                            }
                        }
                    }, 300);
                }

                function startLockout() {
                    locked = true;
                    var btn = document.getElementById('submitBtn');
                    var bar = document.getElementById('lockoutBar');
                    var fill = document.getElementById('lockoutFill');
                    btn.disabled = true;
                    bar.className = 'lockout-bar active';
                    var totalMs = 30000;
                    var elapsed = 0;

                    showMsg('Too many attempts — locked for 30s', 'error');
                    btn.textContent = 'Locked (30s)';

                    lockoutInterval = setInterval(function() {
                        elapsed += 100;
                        var pct = (elapsed / totalMs) * 100;
                        fill.style.width = pct + '%';
                        var remaining = Math.ceil((totalMs - elapsed) / 1000);
                        btn.textContent = 'Locked (' + remaining + 's)';

                        if (elapsed >= totalMs) {
                            clearInterval(lockoutInterval);
                            locked = false;
                            attempts = 0;
                            btn.disabled = false;
                            btn.textContent = 'Authenticate';
                            bar.className = 'lockout-bar';
                            fill.style.width = '0';
                            showMsg('', '');
                        }
                    }, 100);
                }

                function showMsg(text, type) {
                    var el = document.getElementById('msg');
                    el.textContent = text;
                    el.className = 'msg' + (type ? ' ' + type : '');
                }
            </script>
            """}
        </body>
        </html>
        """.trimIndent()
        webView.loadDataWithBaseURL(null, html, "text/html", "UTF-8", null)
    }
}
