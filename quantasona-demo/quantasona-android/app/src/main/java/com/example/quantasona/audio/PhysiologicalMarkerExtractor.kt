package com.example.quantasona.audio

import kotlin.math.abs
import kotlin.math.sqrt

/**
 * Replaces MfccExtractor. 
 * Specifically extracts physiological stress markers (Jitter, Shimmer, Vocal Energy, Pitch Variability)
 * from lossless audio frames, as authorized by Gemma-4-e4b for RLDD analysis.
 */
class PhysiologicalMarkerExtractor(
    private val sampleRate: Int = 16_000
) {

    /**
     * Returns an array containing [Energy, ZeroCrossingRate (Pitch Proxy), Jitter Proxy, Shimmer Proxy]
     */
    fun extract(frame: FloatArray): FloatArray {
        if (frame.isEmpty()) return floatArrayOf(0f, 0f, 0f, 0f)

        val energy = calculateEnergy(frame)
        val zcr = calculateZeroCrossingRate(frame)
        val jitter = calculateJitterProxy(frame)
        val shimmer = calculateShimmerProxy(frame)

        return floatArrayOf(energy, zcr, jitter, shimmer)
    }

    // Measure of Vocal Energy (RMS)
    private fun calculateEnergy(frame: FloatArray): Float {
        var sum = 0f
        for (sample in frame) {
            sum += sample * sample
        }
        return sqrt(sum / frame.size)
    }

    // Proxy for Pitch Variability (F0 changes)
    private fun calculateZeroCrossingRate(frame: FloatArray): Float {
        var crossings = 0
        for (i in 1 until frame.size) {
            if ((frame[i] >= 0 && frame[i - 1] < 0) || (frame[i] < 0 && frame[i - 1] >= 0)) {
                crossings++
            }
        }
        return crossings.toFloat() / frame.size
    }

    // Proxy for Jitter (Frequency instability)
    // Measures variation in distance between zero-crossings
    private fun calculateJitterProxy(frame: FloatArray): Float {
        val zeroCrossingIndices = mutableListOf<Int>()
        for (i in 1 until frame.size) {
            if ((frame[i] >= 0 && frame[i - 1] < 0) || (frame[i] < 0 && frame[i - 1] >= 0)) {
                zeroCrossingIndices.add(i)
            }
        }
        
        if (zeroCrossingIndices.size < 3) return 0f

        var jitterSum = 0f
        var previousPeriod = zeroCrossingIndices[1] - zeroCrossingIndices[0]
        
        for (i in 2 until zeroCrossingIndices.size) {
            val currentPeriod = zeroCrossingIndices[i] - zeroCrossingIndices[i - 1]
            jitterSum += abs(currentPeriod - previousPeriod)
            previousPeriod = currentPeriod
        }
        
        return jitterSum / (zeroCrossingIndices.size - 1)
    }

    // Proxy for Shimmer (Amplitude instability)
    // Measures peak-to-peak amplitude variation
    private fun calculateShimmerProxy(frame: FloatArray): Float {
        val windowSize = 160 // 10ms at 16kHz
        if (frame.size < windowSize * 2) return 0f

        val peaks = mutableListOf<Float>()
        for (i in frame.indices step windowSize) {
            val endIdx = minOf(i + windowSize, frame.size)
            var maxVal = 0f
            for (j in i until endIdx) {
                if (abs(frame[j]) > maxVal) {
                    maxVal = abs(frame[j])
                }
            }
            peaks.add(maxVal)
        }

        if (peaks.size < 2) return 0f

        var shimmerSum = 0f
        for (i in 1 until peaks.size) {
            shimmerSum += abs(peaks[i] - peaks[i-1])
        }
        
        return shimmerSum / (peaks.size - 1)
    }
}
