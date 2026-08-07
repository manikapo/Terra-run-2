package com.territoryrun.app.ui.auth

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.territoryrun.app.auth.SupabaseAuthManager

@Composable
fun LoginScreen(onLoggedIn: () -> Unit) {
    LaunchedEffect(Unit) {
        SupabaseAuthManager.updateApiToken()
        if (SupabaseAuthManager.isLoggedIn()) onLoggedIn()
    }

    Column(
        modifier = Modifier.fillMaxSize().padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text("Territory Run", style = MaterialTheme.typography.headlineLarge)
        Text("Capture turf on every run", modifier = Modifier.padding(top = 8.dp, bottom = 32.dp))
        Button(onClick = {
            // TODO Phase 1: Integrate Google Credential Manager
            // For dev: paste Supabase session via dashboard or use Google Sign-In SDK
            onLoggedIn()
        }) {
            Text("Sign in with Google")
        }
        Text(
            "Configure Google OAuth in Supabase + add Sign-In SDK",
            modifier = Modifier.padding(top = 16.dp),
            style = MaterialTheme.typography.bodySmall
        )
    }
}
