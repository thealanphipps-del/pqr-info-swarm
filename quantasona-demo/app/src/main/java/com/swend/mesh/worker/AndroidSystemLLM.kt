package com.swend.mesh.worker

import android.content.Context
import com.google.mlkit.nl.generativelanguage.GenerativeModel
import com.google.mlkit.nl.generativelanguage.GenerativeModelFutures
import kotlinx.coroutines.tasks.await
import kotlin.math.absoluteValue

class AndroidSystemLLM(private val context: Context) {

    // ML Kit’s built-in on-device text generation model
    private val model = GenerativeModel
        .builder(GenerativeModel.MODEL_TYPE_TEXT)
        .build()

    private val client = GenerativeModelFutures.client(model)

    suspend fun generate(prompt: String): String {
        return try {
            val request = GenerativeModelFutures
                .TextRequest.Builder(prompt)
                .build()

            val response = client.generateText(request).await()
            response.text ?: syntheticResponse(prompt)

        } catch (e: Exception) {
            syntheticResponse(prompt)
        }
    }

    // Fallback generator for unsupported devices
    private fun syntheticResponse(prompt: String): String {
        val seed = prompt.hashCode()
        val words = listOf(
            "processing", "context", "mesh", "node", "signal",
            "accelerator", "task", "inference", "response", "ready"
        )

        return buildString {
            append("Synthetic response: ")
            repeat(12) { i ->
                val idx = (seed + i * 31).absoluteValue % words.size
                append(words[idx]).append(" ")
            }
        }.trim()
    }
}
