package com.territoryrun.app.auth

import android.content.Context
import com.territoryrun.app.BuildConfig
import io.github.jan.supabase.createSupabaseClient
import io.github.jan.supabase.gotrue.Auth
import io.github.jan.supabase.gotrue.auth
import io.github.jan.supabase.gotrue.providers.Google
import io.github.jan.supabase.gotrue.providers.builtin.IDToken
import io.github.jan.supabase.postgrest.Postgrest
import com.territoryrun.app.data.remote.ApiClient

object SupabaseAuthManager {
    private var client = createSupabaseClient(
        supabaseUrl = BuildConfig.SUPABASE_URL,
        supabaseKey = BuildConfig.SUPABASE_ANON_KEY
    ) {
        install(Auth)
        install(Postgrest)
    }

    fun auth() = client.auth

    suspend fun accessToken(): String? =
        auth().currentSessionOrNull()?.accessToken

    fun updateApiToken() {
        ApiClient.setTokenProvider { auth().currentSessionOrNull()?.accessToken ?: "" }
    }

    /**
     * Sign in with Google ID token from Credential Manager / Google Sign-In.
     * Wire GoogleSignIn in LoginScreen and pass idToken here.
     */
    suspend fun signInWithGoogleIdToken(idToken: String, nonce: String? = null) {
        auth().signInWith(IDToken) {
            this.idToken = idToken
            provider = Google
            if (nonce != null) this.nonce = nonce
        }
        updateApiToken()
    }

    suspend fun signOut() {
        auth().signOut()
        ApiClient.setTokenProvider { "" }
    }

    fun isLoggedIn(): Boolean = auth().currentSessionOrNull() != null
}
