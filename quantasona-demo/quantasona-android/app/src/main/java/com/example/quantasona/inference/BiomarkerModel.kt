package com.example.quantasona.inference

import android.content.Context
import org.tensorflow.lite.Interpreter
import org.tensorflow.lite.support.common.FileUtil

class BiomarkerModel(context: Context) {

    private val interpreter: Interpreter

    private val windowSize = 32
    private val mfccDim = 13
    private val embeddingDim = 128

    init {
        val model = FileUtil.loadMappedFile(context, "biomarker_model.tflite")
        interpreter = Interpreter(model)
    }

    fun infer(window: List<FloatArray>): FloatArray {
        val input = Array(1) { Array(windowSize) { FloatArray(mfccDim) } }
        val frames = window.takeLast(windowSize).padStart(windowSize, mfccDim)

        for (i in frames.indices) {
            val frame = frames[i]
            for (j in frame.indices) {
                input[0][i][j] = frame[j]
            }
        }

        val output = Array(1) { FloatArray(embeddingDim) }
        interpreter.run(input, output)
        return output[0]
    }

    private fun List<FloatArray>.padStart(size: Int, dim: Int): List<FloatArray> {
        if (this.size >= size) return this.takeLast(size)
        val padCount = size - this.size
        val pad = List(padCount) { FloatArray(dim) { 0f } }
        return pad + this
    }
}
