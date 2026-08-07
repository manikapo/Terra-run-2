package com.territoryrun.app.gps

import android.annotation.SuppressLint
import android.content.Context
import android.os.Looper
import com.google.android.gms.location.LocationCallback
import com.google.android.gms.location.LocationRequest
import com.google.android.gms.location.LocationResult
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow

data class TrackedPoint(
    val lat: Double,
    val lon: Double,
    val timestamp: Long,
    val accuracy: Float,
    val speed: Float,
    val isMock: Boolean
)

class GpsTracker(private val context: Context) {
    private val fused = LocationServices.getFusedLocationProviderClient(context)
    private val kalman = GpsKalmanFilter()
    private var lastPoint: TrackedPoint? = null

    companion object {
        const val MAX_ACCURACY_M = 40f
        const val MAX_SPEED_MPS = 7f
        const val MIN_DISTANCE_M = 3f
    }

  fun reset() {
    kalman.reset()
    lastPoint = null
  }

  @SuppressLint("MissingPermission")
  fun track(recording: Boolean): Flow<TrackedPoint> = callbackFlow {
    if (!recording) {
      awaitClose { }
      return@callbackFlow
    }

    val request = LocationRequest.Builder(Priority.PRIORITY_HIGH_ACCURACY, 2000L)
      .setMinUpdateIntervalMillis(1000L)
      .setMinUpdateDistanceMeters(3f)
      .setMaxUpdateDelayMillis(5000L)
      .build()

    val callback = object : LocationCallback() {
      override fun onLocationResult(result: LocationResult) {
        val loc = result.lastLocation ?: return
        if (loc.accuracy > MAX_ACCURACY_M) return

        val isMock = loc.isFromMockProvider
        val (lat, lon) = kalman.filter(loc.latitude, loc.longitude, loc.accuracy)

        val point = TrackedPoint(
          lat = lat,
          lon = lon,
          timestamp = loc.time / 1000,
          accuracy = loc.accuracy,
          speed = loc.speed.coerceAtLeast(0f),
          isMock = isMock
        )

        val prev = lastPoint
        if (prev != null) {
          val dt = (point.timestamp - prev.timestamp).toFloat().coerceAtLeast(0.001f)
          val dist = GpsMath.haversineM(prev.lat, prev.lon, point.lat, point.lon)
          val speed = dist / dt
          if (speed > MAX_SPEED_MPS) return
          if (dist < MIN_DISTANCE_M) return
        }

        lastPoint = point
        trySend(point)
      }
    }

    fused.requestLocationUpdates(request, callback, Looper.getMainLooper())
    awaitClose { fused.removeLocationUpdates(callback) }
  }
}
