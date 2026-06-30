package com.example.quantasona.di

import android.content.Context
import com.example.quantasona.audio.AudioProcessor
import com.example.quantasona.mesh.GrpcMeshClient
// import dagger.Module
// import dagger.Provides
// import dagger.hilt.InstallIn
// import dagger.hilt.android.qualifiers.ApplicationContext
// import dagger.hilt.components.SingletonComponent
// import javax.inject.Singleton

/**
 * AppModule provides singleton instances of core, complex services 
 * that require application context or specific initialization logic (e.g., gRPC setup).
 */
// @Module
// @InstallIn(SingletonComponent::class)
object AppModule {

    // --- Audio Pipeline Dependency ---
    // @Provides
    // @Singleton
    fun provideAudioProcessor(/*@ApplicationContext*/ context: Context): AudioProcessor {
        // In a real scenario, this would initialize the microphone manager, 
        // TFLite interpreter, and audio buffer handlers.
        return AudioProcessor(context)
    }

    // --- Mesh Client Dependency ---
    // @Provides
    // @Singleton
    fun provideGrpcMeshClient(): GrpcMeshClient {
        // This simulates the complex setup of a gRPC channel connection 
        // to the backend mesh endpoint.
        println("INFO: Initializing Sovereign Mesh gRPC client...")
        return GrpcMeshClient()
    }
}
