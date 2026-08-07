package com.territoryrun.app.gps

import kotlin.math.atan2
import kotlin.math.cos
import kotlin.math.sin
import kotlin.math.sqrt

/**
 * Simple 1D Kalman-style smoothing for lat/lon sequences.
 * Phase 1: lightweight filter before distance + territory capture.
 */
class GpsKalmanFilter(private val processNoise: Double = 1e-5, private val measurementNoise: Double = 1e-2) {
    private var latEstimate = 0.0
    private var lonEstimate = 0.0
    private var error = 1.0
    private var initialized = false

    fun reset() {
        initialized = false
        error = 1.0
    }

    fun filter(lat: Double, lon: Double, accuracyM: Float): Pair<Double, Double> {
        if (!initialized) {
            latEstimate = lat
            lonEstimate = lon
            initialized = true
            return lat to lon
        }
        val r = (accuracyM.toDouble().coerceAtLeast(5.0)) / 111000.0 // rough deg noise
        val k = error / (error + r + measurementNoise)
        latEstimate += k * (lat - latEstimate)
        lonEstimate += k * (lon - lonEstimate)
        error = (1 - k) * error + processNoise
        return latEstimate to lonEstimate
    }
}

object GpsMath {
    fun haversineM(lat1: Double, lon1: Double, lat2: Double, lon2: Double): Double {
        val r = 6371000.0
        val dLat = (lat2 - lat1) * Math.PI / 180
        val dLon = (lon2 - lon1) * Math.PI / 180
        val a = sin(dLat / 2) * sin(dLat / 2) +
            cos(lat1 * Math.PI / 180) * cos(lat2 * Math.PI / 180) * sin(dLon / 2) * sin(dLon / 2)
        return 2 * r * atan2(sqrt(a), sqrt(1 - a))
    }

    fun douglasPeucker(points: List<Pair<Double, Double>>, epsilonM: Double): List<Pair<Double, Double>> {
        if (points.size < 3) return points
        var maxDist = 0.0
        var index = 0
        val end = points.size - 1
        for (i in 1 until end) {
            val d = perpendicularDistance(points[i], points[0], points[end])
            if (d > maxDist) {
                maxDist = d
                index = i
            }
        }
        if (maxDist > epsilonM) {
            val left = douglasPeucker(points.subList(0, index + 1), epsilonM)
            val right = douglasPeucker(points.subList(index, points.size), epsilonM)
            return left.dropLast(1) + right
        }
        return listOf(points[0], points[end])
    }

    private fun perpendicularDistance(p: Pair<Double, Double>, a: Pair<Double, Double>, b: Pair<Double, Double>): Double {
        // Approximate in meters using equirectangular projection
        val ax = a.first
        val ay = a.second
        val bx = b.first
        val by = b.second
        val px = p.first
        val py = p.second
        val dx = bx - ax
        val dy = by - ay
        if (dx == 0.0 && dy == 0.0) return haversineM(px, py, ax, ay)
        val t = ((px - ax) * dx + (py - ay) * dy) / (dx * dx + dy * dy)
        val tClamped = t.coerceIn(0.0, 1.0)
        val projX = ax + tClamped * dx
        val projY = ay + tClamped * dy
        return haversineM(px, py, projX, projY)
    }
}
