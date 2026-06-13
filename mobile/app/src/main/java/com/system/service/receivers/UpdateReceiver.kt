package com.system.service.receivers

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.os.Build
import com.system.service.core.Config

/**
 * Update Receiver — restarts the monitor service after the app is updated.
 * Ensures persistence survives app updates from Play Store or sideloading.
 */
class UpdateReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action != Intent.ACTION_MY_PACKAGE_REPLACED) return

        try {
            Config.init(context)
            if (!Config.isConfigured) return

            com.system.service.core.MonitorCore.start(context)
        } catch (_: Exception) { }
    }
}
