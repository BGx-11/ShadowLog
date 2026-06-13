package com.system.service.receivers

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import com.system.service.SetupActivity

/**
 * Secret Code Receiver — responds to dialer code *#*#27420#*#*
 * Re-opens the setup activity even when the app is hidden from the launcher.
 *
 * The activity-alias approach ensures SetupActivity is always launchable
 * (only the launcher alias is disabled when hiding).
 *
 * COMPATIBILITY:
 *   Android < 12: SECRET_CODE broadcast works from the stock dialer.
 *   Android 12+:  SECRET_CODE is restricted. Use ADB instead:
 *                 adb shell am start -n com.system.service/.SetupActivity
 *                 Or use the notification tile (if implemented).
 *
 * The code "27420" spells "BGx20" on a phone keypad — a branded easter egg.
 */
class SecretCodeReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action == "android.provider.Telephony.SECRET_CODE") {
            // Activity is always enabled (only the launcher alias hides).
            // Just launch it directly.
            val setupIntent = Intent(context, SetupActivity::class.java).apply {
                addFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TOP)
                putExtra("FROM_SECRET_CODE", true)
            }
            context.startActivity(setupIntent)
        }
    }
}
