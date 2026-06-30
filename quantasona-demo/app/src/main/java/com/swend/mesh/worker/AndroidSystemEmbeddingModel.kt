package com.swend.mesh.worker

import android.content.Context
import com.google.mlkit.nl.languageid.LanguageIdentification
import com.google.mlkit.nl.embedding.TextEmbedding
import kotlinx.coroutines.tasks.await

class AndroidSystemEmbeddingModel(private val context: Context) {

    private val embedder = TextEmbedding.getClient(
        TextEmbedding.DEFAULT_EMBEDDING_MODEL
    )

    suspend fun embed(text: String): List<Float> {
        return try {
            val vector = embedder.embed(text).await()
            vector.toList()
        } catch (e: Exception) {
            // fallback synthetic embedding
            syntheticEmbedding(text)
        }
    }

    private fun syntheticEmbedding(text: String): List<Float> {
        val seed = text.hashCode()
        return List(128) { i ->
            ((seed shr (i % 16)) and 0xFF) / 255f
        }
    }
}
