package com.territoryrun.app.ota

import android.content.Context
import okhttp3.OkHttpClient
import okhttp3.Request
import org.json.JSONObject
import java.io.File
import java.util.zip.ZipInputStream

object OtaManager {
    private val client = OkHttpClient()

    fun sync(context: Context, manifestUrl: String, onComplete: (Boolean) -> Unit) {
        try {
            val manifestResp = client.newCall(Request.Builder().url(manifestUrl).build()).execute()
            if (!manifestResp.isSuccessful) {
                onComplete(false)
                return
            }
            val manifest = JSONObject(manifestResp.body?.string() ?: "")
            val bundleUrl = manifest.getString("bundle_url")
            val sha256 = manifest.optString("sha256", "")

            val bundleResp = client.newCall(Request.Builder().url(bundleUrl).build()).execute()
            if (!bundleResp.isSuccessful) {
                onComplete(false)
                return
            }

            val bytes = bundleResp.body?.bytes() ?: byteArrayOf()
            if (sha256.isNotBlank()) {
                val hash = sha256Hex(bytes)
                if (hash != sha256.lowercase()) {
                    onComplete(false)
                    return
                }
            }

            val version = manifest.optInt("version", 1)
            val staging = File(context.filesDir, "ota/v$version")
            if (staging.exists()) staging.deleteRecursively()
            staging.mkdirs()

            ZipInputStream(bytes.inputStream()).use { zis ->
                var entry = zis.nextEntry
                while (entry != null) {
                    val out = File(staging, entry.name)
                    if (entry.isDirectory) {
                        out.mkdirs()
                    } else {
                        out.parentFile?.mkdirs()
                        out.outputStream().use { zis.copyTo(it) }
                    }
                    entry = zis.nextEntry
                }
            }

            val current = File(context.filesDir, "ota/current")
            if (current.exists()) current.deleteRecursively()
            staging.copyRecursively(current)
            onComplete(true)
        } catch (_: Exception) {
            onComplete(false)
        }
    }

    private fun sha256Hex(data: ByteArray): String {
        val md = java.security.MessageDigest.getInstance("SHA-256")
        return md.digest(data).joinToString("") { "%02x".format(it) }
    }
}
