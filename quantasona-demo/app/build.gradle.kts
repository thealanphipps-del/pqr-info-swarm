plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "com.swend.mesh"
    compileSdk = 34

    defaultConfig {
        applicationId = "com.swend.mesh"
        minSdk = 26
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"

        // Required for foreground services
        manifestPlaceholders["foregroundServiceType"] = "dataSync"
    }

    buildTypes {
        release {
            isMinifyEnabled = false
        }
    }

    buildFeatures {
        viewBinding = true
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = "17"
    }
}

dependencies {

    // Material 3 UI
    implementation("com.google.android.material:material:1.12.0")

    // Kotlin coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.8.1")

    // OkHttp WebSocket
    implementation("com.squareup.okhttp3:okhttp:4.12.0")

    // ML Kit: On-device LLM + Embeddings
    implementation("com.google.mlkit:language-id:17.0.5")
    implementation("com.google.mlkit:smart-reply:17.0.3")
    implementation("com.google.mlkit:translate:17.0.3")

    // (These ML Kit modules internally use the system LLM + embeddings on S24/S25)
}
