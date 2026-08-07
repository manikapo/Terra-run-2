package com.territoryrun.app.data.remote

import com.territoryrun.app.BuildConfig
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path
import java.util.concurrent.TimeUnit

data class CreateActivityBody(
    val started_at: String,
    val idempotency_key: String,
    val device_info: Map<String, String> = emptyMap()
)

data class CreateActivityResponse(val activity_id: String)

data class GpsPointDto(
    val lat: Double,
    val lon: Double,
    val ts: Long,
    val acc: Double,
    val speed: Double = 0.0
)

data class UploadPointsBody(val points: List<GpsPointDto>)

data class CompleteActivityBody(
    val ended_at: String,
    val distance_m: Double,
    val duration_s: Int,
    val h3_cells_hint: List<String> = emptyList()
)

data class CaptureResultDto(
    val new_cells: Int,
    val total_cells: Int,
    val capture_score: Int,
    val h3_cells: List<String>
)

data class CompleteActivityResponse(
    val activity_id: String,
    val status: String,
    val capture_result: CaptureResultDto?,
    val trust_score: Int,
    val message: String?
)

data class UserProfileDto(
    val user_id: String,
    val username: String,
    val display_name: String,
    val cells_owned: Int,
    val territories_captured: Int,
    val total_capture_score: Int,
    val activity_count: Int
)

interface TerritoryApi {
    @POST("api/v1/activities")
    suspend fun createActivity(@Body body: CreateActivityBody): CreateActivityResponse

    @POST("api/v1/activities/{id}/points")
    suspend fun uploadPoints(@Path("id") id: String, @Body body: UploadPointsBody): Map<String, Any>

    @POST("api/v1/activities/{id}/complete")
    suspend fun completeActivity(@Path("id") id: String, @Body body: CompleteActivityBody): CompleteActivityResponse

    @GET("api/v1/users/me")
    suspend fun me(): UserProfileDto
}

object ApiClient {
  private var tokenProvider: () -> String = { "" }

  fun setTokenProvider(provider: () -> String) {
    tokenProvider = provider
  }

  private val authInterceptor = Interceptor { chain ->
    val token = tokenProvider()
    val request = if (token.isNotBlank()) {
      chain.request().newBuilder().addHeader("Authorization", "Bearer $token").build()
    } else chain.request()
    chain.proceed(request)
  }

  private val client = OkHttpClient.Builder()
    .connectTimeout(30, TimeUnit.SECONDS)
    .readTimeout(60, TimeUnit.SECONDS)
    .addInterceptor(authInterceptor)
    .addInterceptor(HttpLoggingInterceptor().apply { level = HttpLoggingInterceptor.Level.BODY })
    .build()

  val api: TerritoryApi = Retrofit.Builder()
    .baseUrl(ensureTrailingSlash(BuildConfig.API_BASE_URL))
    .client(client)
    .addConverterFactory(GsonConverterFactory.create())
    .build()
    .create(TerritoryApi::class.java)

  private fun ensureTrailingSlash(url: String): String =
    if (url.endsWith("/")) url else "$url/"
}
