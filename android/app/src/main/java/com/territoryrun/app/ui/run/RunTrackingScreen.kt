package com.territoryrun.app.ui.run

import android.Manifest
import android.content.Intent
import android.os.Build
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.DisposableEffect
import androidx.compose.runtime.mutableDoubleStateOf
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import com.territoryrun.app.data.local.ActivityLocalEntity
import com.territoryrun.app.data.local.DatabaseProvider
import com.territoryrun.app.data.local.GpsPointEntity
import com.territoryrun.app.data.remote.ApiClient
import com.territoryrun.app.data.remote.CompleteActivityBody
import com.territoryrun.app.data.remote.CreateActivityBody
import com.territoryrun.app.data.remote.GpsPointDto
import com.territoryrun.app.data.remote.UploadPointsBody
import com.territoryrun.app.gps.GpsMath
import com.territoryrun.app.gps.GpsTracker
import com.territoryrun.app.gps.TrackedPoint
import com.territoryrun.app.service.RunTrackingService
import com.territoryrun.app.territory.TerritoryEngine
import kotlinx.coroutines.launch
import java.time.Instant
import java.util.UUID

@Composable
fun RunTrackingScreen(onFinished: () -> Unit) {
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    val db = remember { DatabaseProvider.get(context) }
    val tracker = remember { GpsTracker(context) }
    val territory = remember { TerritoryEngine() }

    val recording = remember { mutableStateOf(false) }
    val paused = remember { mutableStateOf(false) }
    val activityId = remember { mutableStateOf<String?>(null) }
    val distanceM = remember { mutableDoubleStateOf(0.0) }
    val cellsCaptured = remember { mutableIntStateOf(0) }
    val lastPoint = remember { mutableStateOf<TrackedPoint?>(null) }
    val status = remember { mutableStateOf("Ready") }

    val permissions = buildList {
        add(Manifest.permission.ACCESS_FINE_LOCATION)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            add(Manifest.permission.ACCESS_BACKGROUND_LOCATION)
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            add(Manifest.permission.POST_NOTIFICATIONS)
        }
    }

    val permissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { /* continue */ }

    DisposableEffect(recording.value, paused.value) {
        if (!recording.value || paused.value) {
            onDispose { }
            return@DisposableEffect onDispose { }
        }

        val job = scope.launch {
            tracker.track(true).collect { point ->
                lastPoint.value = point
                val prev = lastPoint.value
                // distance from previous stored point in flow - use db last
                val points = db.gpsPointDao().forActivity(activityId.value ?: return@collect)
                if (points.isNotEmpty()) {
                    val p = points.last()
                    distanceM.doubleValue += GpsMath.haversineM(p.lat, p.lon, point.lat, point.lon)
                }
                territory.processPoint(point.lat, point.lon)?.let { cellsCaptured.intValue = territory.capturedCount() }

                db.gpsPointDao().insert(
                    GpsPointEntity(
                        activityId = activityId.value ?: return@collect,
                        lat = point.lat,
                        lon = point.lon,
                        timestamp = point.timestamp,
                        accuracy = point.accuracy,
                        speed = point.speed
                    )
                )
            }
        }
        onDispose { job.cancel() }
    }

    Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
        Text("Run Tracking", style = MaterialTheme.typography.headlineMedium)
        Text("Distance: ${(distanceM.doubleValue / 1000).format(2)} km")
        Text("Cells captured: ${cellsCaptured.intValue}")
        Text("Status: ${status.value}")

        if (!recording.value) {
            Button(onClick = {
                permissionLauncher.launch(permissions.toTypedArray())
                val id = UUID.randomUUID().toString()
                activityId.value = id
                tracker.reset()
                territory.reset()
                distanceM.doubleValue = 0.0
                cellsCaptured.intValue = 0
                recording.value = true
                paused.value = false
                status.value = "Recording"

                context.startForegroundService(
                    Intent(context, RunTrackingService::class.java).apply {
                        action = RunTrackingService.ACTION_START
                    }
                )

                scope.launch {
                    try {
                        val started = Instant.now().toString()
                        db.activityLocalDao().insert(ActivityLocalEntity(id = id, startedAt = Instant.now().epochSecond))
                        ApiClient.api.createActivity(
                            CreateActivityBody(started_at = started, idempotency_key = id)
                        )
                    } catch (e: Exception) {
                        status.value = "API error: ${e.message}"
                    }
                }
            }) { Text("Start") }
        } else {
            Button(onClick = { paused.value = !paused.value }) {
                Text(if (paused.value) "Resume" else "Pause")
            }
            Button(onClick = {
                recording.value = false
                status.value = "Syncing…"
                context.startService(
                    Intent(context, RunTrackingService::class.java).apply {
                        action = RunTrackingService.ACTION_STOP
                    }
                )
                scope.launch {
                    try {
                        val id = activityId.value ?: return@launch
                        val points = db.gpsPointDao().forActivity(id)
                        val dtos = points.map {
                            GpsPointDto(it.lat, it.lon, it.timestamp, it.accuracy.toDouble(), it.speed.toDouble())
                        }
                        if (dtos.isNotEmpty()) {
                            ApiClient.api.uploadPoints(id, UploadPointsBody(dtos))
                        }
                        val ended = Instant.now().toString()
                        val duration = if (points.size >= 2) {
                            (points.last().timestamp - points.first().timestamp).toInt()
                        } else 0
                        val resp = ApiClient.api.completeActivity(
                            id,
                            CompleteActivityBody(
                                ended_at = ended,
                                distance_m = distanceM.doubleValue,
                                duration_s = duration,
                                h3_cells_hint = territory.capturedCells()
                            )
                        )
                        db.activityLocalDao().complete(
                            id = id,
                            endedAt = Instant.now().epochSecond,
                            status = resp.status,
                            distanceM = distanceM.doubleValue,
                            durationSec = duration,
                            captureScore = resp.capture_result?.capture_score ?: 0,
                            h3CellsJson = territory.capturedCells().toString()
                        )
                        status.value = "Done — ${resp.capture_result?.new_cells ?: 0} new cells"
                        onFinished()
                    } catch (e: Exception) {
                        status.value = "Sync failed: ${e.message}"
                    }
                }
            }, modifier = Modifier.padding(top = 8.dp)) {
                Text("Finish & Sync")
            }
        }
    }
}

private fun Double.format(digits: Int) = "%.${digits}f".format(this)
