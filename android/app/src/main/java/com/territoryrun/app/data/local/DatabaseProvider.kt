package com.territoryrun.app.data.local

import android.content.Context
import androidx.room.Room

object DatabaseProvider {
    @Volatile private var instance: AppDatabase? = null

    fun get(context: Context): AppDatabase =
        instance ?: synchronized(this) {
            instance ?: Room.databaseBuilder(
                context.applicationContext,
                AppDatabase::class.java,
                "territory_run.db"
            ).build().also { instance = it }
        }
}