package com.example.quantasona

import android.app.Application
// import dagger.hilt.android.HiltAndroidApp

/**
 * The Quantasona Application Class.
 * This serves as the root context for all dependency graph construction via Hilt.
 */
// @HiltAndroidApp
class QuantasonaApp : Application() {
    override fun onCreate() {
        super.onCreate()
        // Initialization logic for global services (e.g., logging, analytics) 
        // can be placed here if they require the application context early.
        println("Quantasona Node Edge Initialized: Sovereign Mesh Connection Established.")
    }
}
