package com.territoryrun.app.ui

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.territoryrun.app.ui.auth.LoginScreen
import com.territoryrun.app.ui.home.HomeScreen
import com.territoryrun.app.ui.run.RunTrackingScreen

@Composable
fun TerritoryRunNavHost() {
    val nav = rememberNavController()
    NavHost(navController = nav, startDestination = "login") {
        composable("login") {
            LoginScreen(onLoggedIn = { nav.navigate("home") { popUpTo("login") { inclusive = true } } })
        }
        composable("home") {
            HomeScreen(
                onStartRun = { nav.navigate("run") },
                onOpenProfile = { /* OTA WebView from HomeScreen */ }
            )
        }
        composable("run") {
            RunTrackingScreen(onFinished = { nav.popBackStack() })
        }
    }
}
