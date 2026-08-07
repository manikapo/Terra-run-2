package com.territoryrun.app.data.local

import androidx.room.Dao
import androidx.room.Database
import androidx.room.Entity
import androidx.room.Insert
import androidx.room.PrimaryKey
import androidx.room.Query
import androidx.room.RoomDatabase

@Entity(tableName = "gps_points")
data class GpsPointEntity(
    @PrimaryKey(autoGenerate = true) val id: Long = 0,
    val activityId: String,
    val lat: Double,
    val lon: Double,
    val timestamp: Long,
    val accuracy: Float,
    val speed: Float,
    val synced: Boolean = false
)

@Entity(tableName = "activities_local")
data class ActivityLocalEntity(
    @PrimaryKey val id: String,
    val startedAt: Long,
    val endedAt: Long? = null,
    val status: String = "recording",
    val distanceM: Double = 0.0,
    val durationSec: Int = 0,
    val captureScore: Int = 0,
    val h3CellsJson: String = "[]"
)

@Dao
interface GpsPointDao {
    @Insert
    suspend fun insert(point: GpsPointEntity)

    @Insert
    suspend fun insertAll(points: List<GpsPointEntity>)

    @Query("SELECT * FROM gps_points WHERE activityId = :activityId ORDER BY timestamp")
    suspend fun forActivity(activityId: String): List<GpsPointEntity>

    @Query("DELETE FROM gps_points WHERE activityId = :activityId")
    suspend fun deleteForActivity(activityId: String)
}

@Dao
interface ActivityLocalDao {
    @Insert
    suspend fun insert(activity: ActivityLocalEntity)

    @Query("UPDATE activities_local SET endedAt = :endedAt, status = :status, distanceM = :distanceM, durationSec = :durationSec, captureScore = :captureScore, h3CellsJson = :h3CellsJson WHERE id = :id")
    suspend fun complete(
        id: String,
        endedAt: Long,
        status: String,
        distanceM: Double,
        durationSec: Int,
        captureScore: Int,
        h3CellsJson: String
    )

    @Query("SELECT * FROM activities_local ORDER BY startedAt DESC LIMIT 20")
    suspend fun recent(): List<ActivityLocalEntity>
}

@Database(
    entities = [GpsPointEntity::class, ActivityLocalEntity::class],
    version = 1,
    exportSchema = false
)
abstract class AppDatabase : RoomDatabase() {
    abstract fun gpsPointDao(): GpsPointDao
    abstract fun activityLocalDao(): ActivityLocalDao
}
