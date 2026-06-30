package com.example.quantasonaapp.domain

import java.security.MessageDigest
import kotlin.random.Random

object TesseractGenerator {

    /**
     * Generates an 81-character cascading complexity tesseract hash based on the provided audio biometric data.
     * The resulting hash is strictly 81 alphanumeric characters, representing a unique 5D coordinate.
     */
    fun generateTesseractHash(audioData: ByteArray): String {
        // Step 1: Base cryptographic digest (SHA-384 provides 48 bytes -> ~64 base64 chars, 96 hex chars)
        val md = MessageDigest.getInstance("SHA-384")
        val digest = md.digest(audioData)

        // Step 2: Convert to alphanumeric cascading sequence
        val hexString = digest.joinToString("") { "%02x".format(it) }

        // Step 3: Enhance complexity with deterministic dimensional folding (simulating 5D projection)
        // We take the hex and mix it with a seeded random to ensure we get exactly 81 alphanumeric characters
        val seed = digest.fold(0L) { acc, byte -> acc + byte.toLong() }
        val random = java.util.Random(seed)

        val charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
        val stringBuilder = java.lang.StringBuilder()

        // We use the first 40 chars of the actual hex digest to represent the structural core
        for (i in 0 until 40) {
            val charFromHex = hexString[i % hexString.length]
            stringBuilder.append(charFromHex)
        }

        // The remaining 41 characters are the cascading complexity layers derived from the audio seed
        for (i in 40 until 81) {
            stringBuilder.append(charset[random.nextInt(charset.length)])
        }

        return stringBuilder.toString()
    }

    /**
     * Overload for testing or simulating voice analysis without raw bytes
     */
    fun generateTesseractHash(anomalyDescription: String): String {
        return generateTesseractHash(anomalyDescription.toByteArray(Charsets.UTF_8))
    }
}
