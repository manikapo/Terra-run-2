package com.territoryrun.app.ota

import android.annotation.SuppressLint
import android.os.Bundle
import android.webkit.WebResourceRequest
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.activity.ComponentActivity
import com.territoryrun.app.BuildConfig
import com.territoryrun.app.auth.SupabaseAuthManager
import java.io.File

class OtaWebViewActivity : ComponentActivity() {
    @SuppressLint("SetJavaScriptEnabled")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val webView = WebView(this)
        setContentView(webView)

        webView.settings.javaScriptEnabled = true
        webView.settings.domStorageEnabled = true
        webView.settings.allowFileAccess = false

        webView.addJavascriptInterface(OtaBridge(), "TerritoryBridge")

        webView.webViewClient = object : WebViewClient() {
            override fun shouldOverrideUrlLoading(view: WebView?, request: WebResourceRequest?): Boolean {
                val url = request?.url?.toString() ?: return false
                return !url.startsWith("file://") && !url.startsWith("https://")
            }
        }

        val localProfile = File(filesDir, "ota/current/pages/profile.html")
        when {
            localProfile.exists() -> webView.loadUrl("file://${localProfile.absolutePath}")
            else -> {
                // Dev: load live URL from run.8me.in (no zip needed for first test)
                webView.loadUrl(BuildConfig.OTA_PROFILE_URL)
                // Background: download zip for offline / next launch
                OtaManager.sync(this, BuildConfig.OTA_MANIFEST_URL) { ok ->
                    if (ok && localProfile.exists()) {
                        // Optional: could switch to cached file on next open
                    }
                }
            }
        }
    }
}

class OtaBridge {
    @android.webkit.JavascriptInterface
    fun getAccessToken(): String =
        SupabaseAuthManager.auth().currentSessionOrNull()?.accessToken ?: ""

    @android.webkit.JavascriptInterface
    fun getApiBaseUrl(): String = BuildConfig.API_BASE_URL
}
