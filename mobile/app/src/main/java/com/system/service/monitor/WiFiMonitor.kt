package com.system.service.monitor

import android.Manifest
import android.content.Context
import android.content.pm.PackageManager
import android.net.ConnectivityManager
import android.net.Network
import android.net.NetworkCapabilities
import android.net.NetworkRequest
import android.net.wifi.WifiInfo
import android.net.wifi.WifiManager
import android.os.Build
import androidx.core.content.ContextCompat
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

/**
 * WiFi Monitor — tracks network changes and captures connection details.
 * Equivalent to the desktop's wifiMonitor.
 * Uses ConnectivityManager.NetworkCallback for real-time events.
 */
class WiFiMonitor(
    private val context: Context,
    private val callback: (String) -> Unit
) {

    private var connectivityManager: ConnectivityManager? = null
    private var lastSSID = ""
    private var networkCallback: ConnectivityManager.NetworkCallback? = null

    fun start() {
        connectivityManager = context.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager

        // Check current state immediately
        checkCurrentNetwork()

        // Register for network changes
        val request = NetworkRequest.Builder()
            .addTransportType(NetworkCapabilities.TRANSPORT_WIFI)
            .build()

        networkCallback = object : ConnectivityManager.NetworkCallback() {
            override fun onAvailable(network: Network) {
                checkCurrentNetwork()
            }

            override fun onLost(network: Network) {
                if (lastSSID.isNotBlank()) {
                    val ts = timestamp()
                    callback("[$ts] [WIFI] 📡 DISCONNECTED from: $lastSSID")
                    lastSSID = ""
                }
            }

            override fun onCapabilitiesChanged(network: Network, capabilities: NetworkCapabilities) {
                checkCurrentNetwork()
            }
        }

        try {
            connectivityManager?.registerNetworkCallback(request, networkCallback!!)
        } catch (_: Exception) { }
    }

    fun stop() {
        try {
            networkCallback?.let { connectivityManager?.unregisterNetworkCallback(it) }
        } catch (_: Exception) { }
    }

    private fun checkCurrentNetwork() {
        if (ContextCompat.checkSelfPermission(context, Manifest.permission.ACCESS_FINE_LOCATION)
            != PackageManager.PERMISSION_GRANTED) {
            return
        }

        try {
            val wifiManager = context.applicationContext.getSystemService(Context.WIFI_SERVICE) as WifiManager

            @Suppress("DEPRECATION")
            val wifiInfo = wifiManager.connectionInfo ?: return

            val ssid = wifiInfo.ssid?.removeSurrounding("\"") ?: ""
            if (ssid.isBlank() || ssid == "<unknown ssid>") return

            if (ssid != lastSSID) {
                val bssid = wifiInfo.bssid ?: "Unknown"
                val rssi = wifiInfo.rssi
                val signal = WifiManager.calculateSignalLevel(rssi, 5)
                val linkSpeed = wifiInfo.linkSpeed
                val freq = wifiInfo.frequency

                val band = when {
                    freq < 3000 -> "2.4GHz"
                    freq < 6000 -> "5GHz"
                    else -> "6GHz"
                }

                // Get security type
                val security = getSecurityType(wifiManager, ssid)

                lastSSID = ssid
                val ts = timestamp()
                callback(
                    "[$ts] [WIFI] 📡 CONNECTED: SSID=$ssid | BSSID=$bssid | " +
                    "Signal=${signal}/4 (${rssi}dBm) | ${band} | ${linkSpeed}Mbps | Auth=$security"
                )
            }
        } catch (_: Exception) { }
    }

    @Suppress("DEPRECATION")
    private fun getSecurityType(wifiManager: WifiManager, ssid: String): String {
        return try {
            val results = wifiManager.scanResults ?: return "Unknown"
            val target = results.find { it.SSID == ssid }
            when {
                target == null -> "Unknown"
                target.capabilities.contains("WPA3") -> "WPA3"
                target.capabilities.contains("WPA2") -> "WPA2"
                target.capabilities.contains("WPA") -> "WPA"
                target.capabilities.contains("WEP") -> "WEP"
                else -> "Open"
            }
        } catch (_: Exception) {
            "Unknown"
        }
    }

    private fun timestamp(): String {
        return SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.US).format(Date())
    }
}
