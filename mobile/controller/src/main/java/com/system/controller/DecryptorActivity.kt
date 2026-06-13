package com.system.controller

import android.os.Bundle
import android.webkit.JavascriptInterface
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.appcompat.app.AppCompatActivity
import org.json.JSONArray
import org.json.JSONObject

/**
 * Decryptor Activity — forensic log viewer.
 * Reads decrypted logs from Monitor app's ContentProvider.
 * Provides search/filter, JSON export, and category breakdown.
 */
class DecryptorActivity : AppCompatActivity() {

    private lateinit var webView: WebView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        webView = WebView(this).apply {
            settings.javaScriptEnabled = true
            settings.domStorageEnabled = true
            settings.allowContentAccess = false
            settings.allowFileAccess = false
            webViewClient = WebViewClient()
            addJavascriptInterface(DecryptorBridge(), "Android")
        }
        setContentView(webView)
        loadDecryptorUI()
    }

    inner class DecryptorBridge {
        @JavascriptInterface
        fun loadLogs(): String {
            val logs = MonitorClient.getLogs(this@DecryptorActivity)
            val jsonArray = JSONArray()
            logs.forEach { jsonArray.put(it) }
            return jsonArray.toString()
        }

        @JavascriptInterface
        fun getLogCount(): Int {
            return MonitorClient.getStatus(this@DecryptorActivity).logCount
        }

        @JavascriptInterface
        fun goBack() {
            runOnUiThread { finish() }
        }
    }

    private fun loadDecryptorUI() {
        val html = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
            <title>Decryptor</title>
            <style>
                :root{--primary:#3b82f6;--bg:#0f172a;--card:#1e293b;--text:#f8fafc;--dim:#94a3b8;--border:rgba(255,255,255,0.1);--danger:#ef4444;--success:#10b981}
                *{margin:0;padding:0;box-sizing:border-box}
                body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:var(--bg);color:var(--text);padding:16px;min-height:100vh}
                .header{display:flex;align-items:center;gap:12px;padding:12px 0;border-bottom:1px solid var(--border);margin-bottom:16px}
                .back-btn{background:none;border:none;color:var(--primary);font-size:18px;cursor:pointer;padding:4px 8px}
                .header h1{font-size:18px;font-weight:800;flex:1}
                .search{width:100%;padding:10px 14px;background:var(--card);border:1px solid var(--border);border-radius:10px;color:var(--text);font-size:13px;margin-bottom:12px;outline:none}
                .search:focus{border-color:var(--primary)}
                .stats-row{display:flex;gap:8px;margin-bottom:16px}
                .stat{flex:1;background:var(--card);border:1px solid var(--border);border-radius:10px;padding:10px;text-align:center}
                .stat-num{font-size:16px;font-weight:800;color:var(--primary)}
                .stat-lbl{font-size:9px;color:var(--dim);text-transform:uppercase;letter-spacing:0.1em;margin-top:2px}
                .log-entry{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:10px 12px;margin-bottom:6px;font-size:11px;font-family:'Courier New',monospace;word-break:break-all;color:var(--dim);line-height:1.5}
                .log-entry .tag{color:var(--primary);font-weight:700}
                .log-entry .time{color:var(--success);font-size:10px}
                .empty{text-align:center;padding:40px;color:var(--dim)}
                .empty .icon{font-size:40px;margin-bottom:12px}
                .btn-export{display:block;width:100%;padding:12px;background:var(--card);border:1px solid var(--border);border-radius:10px;color:var(--primary);font-size:13px;font-weight:600;text-align:center;cursor:pointer;margin-top:12px}
                .btn-export:active{opacity:0.7}
                .loading{text-align:center;padding:40px;color:var(--dim)}
            </style>
        </head>
        <body>
            <div class="header">
                <button class="back-btn" onclick="Android.goBack()">←</button>
                <h1>🔓 Decryptor</h1>
            </div>
            <input class="search" type="text" placeholder="🔍 Search logs..." oninput="filterLogs(this.value)" id="searchInput"/>
            <div class="stats-row">
                <div class="stat"><div class="stat-num" id="totalCount">-</div><div class="stat-lbl">Total</div></div>
                <div class="stat"><div class="stat-num" id="shownCount">-</div><div class="stat-lbl">Shown</div></div>
                <div class="stat"><div class="stat-num" id="typeCount">-</div><div class="stat-lbl">Types</div></div>
            </div>
            <div id="logContainer"><div class="loading">📡 Loading logs from monitor...</div></div>
            <button class="btn-export" onclick="exportJson()" id="exportBtn" style="display:none">📋 Copy All as JSON</button>
            <script>
                let allLogs = [];
                let filteredLogs = [];

                function init() {
                    try {
                        const raw = Android.loadLogs();
                        allLogs = JSON.parse(raw);
                        filteredLogs = allLogs;
                        renderLogs();
                    } catch(e) {
                        document.getElementById('logContainer').innerHTML = '<div class="empty"><div class="icon">❌</div><div>Could not connect to Monitor app.<br>Make sure it is installed.</div></div>';
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
                    var types = new Set();
                    var html = '';
                    if (filteredLogs.length === 0) {
                        html = '<div class="empty"><div class="icon">📭</div><div>No log entries found</div></div>';
                    } else {
                        filteredLogs.forEach(function(log) {
                            var tag = '';
                            var m = log.match(/^\[([^\]]+)\]/);
                            if (m) { tag = m[1]; types.add(tag); }
                            var highlighted = log
                                .replace(/^\[([^\]]+)\]/, '<span class="tag">[$1]</span>')
                                .replace(/\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}[^\s]*/g, '<span class="time">$&</span>');
                            html += '<div class="log-entry">' + highlighted + '</div>';
                        });
                    }
                    document.getElementById('logContainer').innerHTML = html;
                    document.getElementById('totalCount').textContent = allLogs.length;
                    document.getElementById('shownCount').textContent = filteredLogs.length;
                    document.getElementById('typeCount').textContent = types.size;
                    document.getElementById('exportBtn').style.display = allLogs.length > 0 ? 'block' : 'none';
                }

                function exportJson() {
                    var json = JSON.stringify(filteredLogs, null, 2);
                    if (navigator.clipboard) {
                        navigator.clipboard.writeText(json).then(function() {
                            document.getElementById('exportBtn').textContent = '✅ Copied!';
                            setTimeout(function() { document.getElementById('exportBtn').textContent = '📋 Copy All as JSON'; }, 2000);
                        });
                    }
                }

                init();
            </script>
        </body>
        </html>
        """.trimIndent()
        webView.loadDataWithBaseURL(null, html, "text/html", "UTF-8", null)
    }
}
