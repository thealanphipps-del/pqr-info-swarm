package com.example.quantasona.pipeline

import com.example.quantasona.audio.AudioProcessor
import com.example.quantasona.audio.MfccExtractor
import com.example.quantasona.mesh.AndroidMeshClient
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.collectLatest
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import quantasona.epistemic.EpistemicDelta
import com.example.quantasona.inference.BiomarkerModel
import kotlin.math.sqrt

class BiomarkerPipeline(
    private val audioProcessor: AudioProcessor,
    private val meshClient: AndroidMeshClient,
    private val biomarkerModel: BiomarkerModel
) {
    private val scope = CoroutineScope(Dispatchers.Main.immediate)
    private val markerExtractor = PhysiologicalMarkerExtractor(16000)

    private val mfccWindow = ArrayList<FloatArray>()

    fun start() {
        audioProcessor.start(scope)
        scope.launch {
            audioProcessor.frames.collectLatest { frame ->
                try {
                    val markers = markerExtractor.extract(frame)
                    mfccWindow.add(markers)
                    if (mfccWindow.size >= 32) {
                        val embedding = biomarkerModel.infer(mfccWindow)
                        
                        val delta = buildDelta(embedding)
                        
                        withContext(Dispatchers.IO) {
                            meshClient.submit(delta)
                        }
                        mfccWindow.clear()
                    }
                } catch (e: kotlinx.coroutines.CancellationException) {
                    throw e
                } catch (e: Exception) {
                    // Graceful recovery from faults
                    android.util.Log.e("BiomarkerPipeline", "Pipeline fault intercepted: ${e.message}. Recovering...", e)
                    mfccWindow.clear()
                }
            }
        }
    }

    fun forceFlush() {
        if (mfccWindow.isNotEmpty()) {
            scope.launch {
                try {
                    val embedding = biomarkerModel.infer(mfccWindow)
                    val delta = buildDelta(embedding)
                    withContext(Dispatchers.IO) {
                        meshClient.submit(delta)
                    }
                    mfccWindow.clear()
                } catch (e: Exception) {
                    android.util.Log.e("BiomarkerPipeline", "Flush fault: ${e.message}", e)
                    mfccWindow.clear()
                }
            }
        }
    }

    private fun buildFaultyDelta(): EpistemicDelta {
        return EpistemicDelta.newBuilder()
            .setSigmaId("") // Corrupted empty ID
            .setSemanticWeight(-999.0) // Corrupted weight
            .setConfidence(-1.0) // Invalid confidence
            .setProvenance("corrupted-node")
            .setTimestamp("")
            .setDeltaTypeValue(999) // Invalid type
            .setRelationTypeValue(999) // Invalid relation
            .setDeltaId("invalid-seq-id-!@#$") // Faulty sequence ID
            .build()
    }

    private var lastEmbedding: FloatArray? = null

    private fun buildDelta(embedding: FloatArray): EpistemicDelta {
        val sigmaId = hashEmbedding(embedding)
        val semanticWeight = computeDeviation(embedding, lastEmbedding)
        lastEmbedding = embedding

        return EpistemicDelta.newBuilder()
            .setSigmaId("android://biomarker/$sigmaId")
            .setSemanticWeight(semanticWeight)
            .setConfidence(0.9)
            .setProvenance("android-node-001")
            .setTimestamp(System.currentTimeMillis().toString())
            .setDeltaTypeValue(1)      // OBSERVATION
            .setRelationTypeValue(3)   // INTRODUCES
            .setDeltaId("delta-" + System.currentTimeMillis())
            .build()
    }

    private fun hashEmbedding(vec: FloatArray): String {
        var h = 1125899906842597L
        for (v in vec) {
            val bits = java.lang.Float.floatToIntBits(v)
            h = 31 * h + bits
        }
        return h.toString(16)
    }

    private fun computeDeviation(cur: FloatArray, prev: FloatArray?): Double {
        if (prev == null || prev.size != cur.size) return 0.0
        var sum = 0.0
        for (i in cur.indices) {
            val d = (cur[i] - prev[i]).toDouble()
            sum += d * d
        }
        return sqrt(sum)
    }
}
