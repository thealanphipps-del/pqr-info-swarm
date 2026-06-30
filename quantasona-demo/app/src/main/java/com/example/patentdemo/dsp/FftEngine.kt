package com.example.patentdemo.dsp

import kotlin.math.PI
import kotlin.math.cos
import kotlin.math.sin
import kotlin.math.sqrt

class FftEngine {
    
    fun computeFFT(windowed: FloatArray): FloatArray {
        if (windowed.isEmpty()) return FloatArray(0)
        
        // Pad to next power of 2 for Radix-2 FFT
        var n = 1
        while (n < windowed.size) {
            n = n shl 1
        }
        
        val paddedData = FloatArray(n)
        System.arraycopy(windowed, 0, paddedData, 0, windowed.size)
        
        // Perform FFT
        val real = paddedData
        val imag = FloatArray(n)
        fft(real, imag)
        
        // Compute magnitude spectrum
        val magnitude = FloatArray(n / 2)
        for (i in magnitude.indices) {
            magnitude[i] = sqrt(real[i] * real[i] + imag[i] * imag[i])
        }
        
        return magnitude
    }

    fun transformToFrequencyDomain(pcmData: ByteArray): FloatArray {
        if (pcmData.isEmpty()) return FloatArray(0)
        
        // Convert 16-bit PCM to FloatArray
        val floatData = FloatArray(pcmData.size / 2)
        for (i in floatData.indices) {
            val byte1 = pcmData[i * 2].toInt() and 0xFF
            val byte2 = pcmData[i * 2 + 1].toInt()
            val sample = (byte2 shl 8) or byte1
            floatData[i] = sample.toFloat() / 32768f
        }
        
        // Apply Hanning Window
        for (i in floatData.indices) {
            val multiplier = 0.5f * (1f - cos(2.0 * PI * i / (floatData.size - 1)).toFloat())
            floatData[i] = floatData[i] * multiplier
        }
        
        return computeFFT(floatData)
    }

    private fun fft(real: FloatArray, imag: FloatArray) {
        val n = real.size
        if (n <= 1) return

        // Bit-reversal permutation
        var j = 0
        for (i in 0 until n - 1) {
            if (i < j) {
                var temp = real[i]
                real[i] = real[j]
                real[j] = temp
            }
            var m = n / 2
            while (m <= j) {
                j -= m
                m /= 2
            }
            j += m
        }

        // Cooley-Tukey decimation-in-time radix-2 FFT
        var step = 2
        while (step <= n) {
            val halfStep = step / 2
            val angle = -2.0 * PI / step
            val wReal = cos(angle).toFloat()
            val wImag = sin(angle).toFloat()

            for (i in 0 until n step step) {
                var currentWReal = 1f
                var currentWImag = 0f

                for (k in 0 until halfStep) {
                    val evenIndex = i + k
                    val oddIndex = i + k + halfStep

                    val tReal = currentWReal * real[oddIndex] - currentWImag * imag[oddIndex]
                    val tImag = currentWReal * imag[oddIndex] + currentWImag * real[oddIndex]

                    real[oddIndex] = real[evenIndex] - tReal
                    imag[oddIndex] = imag[evenIndex] - tImag
                    real[evenIndex] += tReal
                    imag[evenIndex] += tImag

                    val nextWReal = currentWReal * wReal - currentWImag * wImag
                    val nextWImag = currentWReal * wImag + currentWImag * wReal
                    currentWReal = nextWReal
                    currentWImag = nextWImag
                }
            }
            step *= 2
        }
    }
}
