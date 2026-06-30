package com.example.patentdemo.ui

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.lifecycle.viewmodel.compose.viewModel

@Composable
fun PatentDemoNavHost() {
    val navController = rememberNavController()
    val viewModel: AnalysisViewModel = viewModel()

    NavHost(
        navController = navController,
        startDestination = "record"
    ) {
        composable("record") {
            RecordScreen(
                onRecorded = { pcm ->
                    viewModel.pcmData = pcm
                    navController.navigate("fft")
                }
            )
        }

        composable("fft") {
            FftScreen(
                pcmData = viewModel.pcmData,
                onFftComputed = { fft ->
                    viewModel.fftData = fft
                    navController.navigate("signature")
                }
            )
        }

        composable("signature") {
            SignatureScreen(
                fftData = viewModel.fftData,
                onZeroSlopeDetected = { bands ->
                    viewModel.zeroSlopeBands = bands
                    navController.navigate("profile")
                }
            )
        }

        composable("profile") {
            ProfileScreen(
                zeroSlopeBands = viewModel.zeroSlopeBands,
                onProfileMatched = { match ->
                    viewModel.profileMatch = match
                    navController.navigate("report")
                }
            )
        }

        composable("report") {
            ReportScreen(
                match = viewModel.profileMatch
            )
        }
    }
}
