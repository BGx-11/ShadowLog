package com.system.service

import android.app.Application
import android.content.Context
import android.app.NotificationChannel
import android.app.NotificationManager
import android.os.Build
import com.system.service.core.Config

/**
 * Application class — initializes notification channels and config on cold start.
 * Named generically to avoid suspicion in system settings.
 */
class ShadowApp : Application() {

    companion object {
        const val NOTIF_CHANNEL_ID = "sys_svc_channel"
        const val NOTIF_ID = 9901

        @Volatile
        lateinit var appContext: Context
            private set
    }

    override fun onCreate() {
        super.onCreate()
        appContext = applicationContext
        createNotificationChannel()
        Config.init(this)
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                NOTIF_CHANNEL_ID,
                getString(R.string.notif_channel_name),
                NotificationManager.IMPORTANCE_MIN  // Lowest visibility
            ).apply {
                description = getString(R.string.notif_channel_desc)
                setShowBadge(false)
                lockscreenVisibility = android.app.Notification.VISIBILITY_SECRET
            }
            val nm = getSystemService(NotificationManager::class.java)
            nm.createNotificationChannel(channel)
        }
    }
}
