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

        webView.addJavascriptInterface(OtaBridge(this), "TerritoryBridge")

        webView.webViewClient = object : WebViewClient() {
            override fun shouldOverrideUrlLoading(view: WebView?, request: WebResourceRequest?): Boolean {
                val url = request?.url?.toString() ?: return false
                return !url.startsWith("file://")
            }
        }

        val otaDir = File(filesDir, "ota/current")
        val index = File(otaDir, "pages/profile.html")
        if (index.exists()) {
            webView.loadUrl("file://${index.absolutePath}")
        } else {
            // Fallback: load manifest URL hint in dev
            webView.loadUrl("about:blank")
            OtaManager.sync(this, BuildConfig.OTA_MANIFEST_URL) { ok ->
                if (ok && index.exists()) {
                    webView.loadUrl("file://${index.absolutePath}")
                }
            }
        }
    }
}

class OtaBridge(private val activity: OtaWebViewActivity) {
    @android.webkit.JavascriptInterface
    fun getAccessToken(): String =
        SupabaseAuthManager.auth().currentSessionOrNull()?.accessToken ?: ""

    @android.webkit.JavascriptInterface
    fun getApiBaseUrl(): String = BuildConfig.API_BASE_URL
}
