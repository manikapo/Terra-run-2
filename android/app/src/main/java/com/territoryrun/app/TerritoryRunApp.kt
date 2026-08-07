package com.territoryrun.app

import android.app.Application
import com.mapbox.maps.MapboxOptions

class TerritoryRunApp : Application() {
    override fun onCreate() {
        super.onCreate()
        MapboxOptions.accessToken = BuildConfig.MAPBOX_ACCESS_TOKEN
    }
}
