package com.swend.mesh.worker

import android.content.Context
import com.swend.mesh.model.AndroidSystemEmbeddingModel
import com.swend.mesh.protocol.TaskProtocol
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class EmbeddingTaskHandler(context: Context) {

    private val model = AndroidSystemEmbeddingModel(context)

    suspend fun handle(task: TaskProtocol.Task): String = withContext(Dispatchers.Default) {
        val input = task.input?.toString() ?: ""

        return@withContext try {
            val vector = model.embed(input)

            TaskProtocol.buildResult(
                id = task.id,
                result = vector.toList(),   // JSON‑friendly
                status = "ok"
            )
        } catch (e: Exception) {
            TaskProtocol.buildResult(
                id = task.id,
                result = "Embedding error: ${e.message}",
                status = "error"
            )
        }
    }
}
