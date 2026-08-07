package com.territoryrun.app.territory

import com.uber.h3core.H3Core

class TerritoryEngine {
    private val h3 = H3Core.newInstance()
    private val resolution = 10
    private val captured = mutableSetOf<Long>()

    fun reset() {
        captured.clear()
    }

    fun processPoint(lat: Double, lon: Double): String? {
        val cell = h3.geoToH3(lat, lon, resolution)
        if (captured.add(cell)) {
            return h3.h3ToString(cell)
        }
        return null
    }

    fun capturedCells(): List<String> = captured.map { h3.h3ToString(it) }

    fun capturedCount(): Int = captured.size
}
