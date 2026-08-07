package com.territoryrun.app.ui.home

import android.content.Intent
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import com.territoryrun.app.data.remote.ApiClient
import com.territoryrun.app.data.remote.UserProfileDto
import com.territoryrun.app.ota.OtaWebViewActivity

@Composable
fun HomeScreen(onStartRun: () -> Unit, onOpenProfile: () -> Unit) {
    val profile = remember { mutableStateOf<UserProfileDto?>(null) }
    val error = remember { mutableStateOf<String?>(null) }
    val context = LocalContext.current

    LaunchedEffect(Unit) {
        try {
            profile.value = ApiClient.api.me()
        } catch (e: Exception) {
            error.value = e.message
        }
    }

    Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
        Text("Territory Run", style = MaterialTheme.typography.headlineMedium)
        profile.value?.let { p ->
            Text("Cells owned: ${p.cells_owned}")
            Text("Capture score: ${p.total_capture_score}")
            Text("Activities: ${p.activity_count}")
        }
        error.value?.let { Text("API: $it", color = MaterialTheme.colorScheme.error) }

        Button(onClick = onStartRun, modifier = Modifier.padding(top = 16.dp)) {
            Text("Start Run")
        }
        Button(
            onClick = {
                context.startActivity(Intent(context, OtaWebViewActivity::class.java))
                onOpenProfile()
            },
            modifier = Modifier.padding(top = 8.dp)
        ) {
            Text("Profile (OTA)")
        }
    }
}
