package com.example.quantasona.audio

import kotlin.math.ln
import kotlin.math.cos
import kotlin.math.PI
import kotlin.math.pow
import kotlin.math.sqrt

class MfccExtractor(
    private val sampleRate: Int = 16_000,
    private val numCoeffs: Int = 13,
    private val numFilters: Int = 26,
) {

    fun extract(frame: FloatArray): FloatArray {
        val n = frame.size
        val fftReal = frame.copyOf()
        val fftImag = FloatArray(n)

        fft(fftReal, fftImag)

        val power = FloatArray(n / 2)
        for (i in power.indices) {
            val re = fftReal[i]
            val im = fftImag[i]
            power[i] = (re * re + im * im) / n
        }

        val melEnergies = melFilterbank(power)
        for (i in melEnergies.indices) {
            melEnergies[i] = ln(melEnergies[i] + 1e-10f)
        }

        return dct(melEnergies, numCoeffs)
    }

    // Cooley–Tukey radix-2 FFT (in-place)
    private fun fft(real: FloatArray, imag: FloatArray) {
        val n = real.size
        var j = 0
        for (i in 1 until n - 1) {
            var bit = n shr 1
            while (j >= bit) {
                j -= bit
                bit = bit shr 1
            }
            j += bit
            if (i < j) {
                val tr = real[i]; real[i] = real[j]; real[j] = tr
                val ti = imag[i]; imag[i] = imag[j]; imag[j] = ti
            }
        }
        var len = 2
        while (len <= n) {
            val ang = -2.0 * PI / len
            val wlenCos = cos(ang).toFloat()
            val wlenSin = kotlin.math.sin(ang).toFloat()
            var i = 0
            while (i < n) {
                var wCos = 1f
                var wSin = 0f
                for (k in 0 until len / 2) {
                    val uRe = real[i + k]
                    val uIm = imag[i + k]
                    val vRe = real[i + k + len / 2] * wCos - imag[i + k + len / 2] * wSin
                    val vIm = real[i + k + len / 2] * wSin + imag[i + k + len / 2] * wCos
                    real[i + k] = uRe + vRe
                    imag[i + k] = uIm + vIm
                    real[i + k + len / 2] = uRe - vRe
                    imag[i + k + len / 2] = uIm - vIm

                    val tmpCos = wCos * wlenCos - wSin * wlenSin
                    wSin = wCos * wlenSin + wSin * wlenCos
                    wCos = tmpCos
                }
                i += len
            }
            len = len shl 1
        }
    }

    private fun hzToMel(hz: Double) = 2595.0 * ln(1 + hz / 700.0)
    private fun melToHz(mel: Double) = 700.0 * (kotlin.math.exp(mel / 2595.0) - 1)

    private fun melFilterbank(power: FloatArray): FloatArray {
        val nFft = power.size * 2
        val lowMel = hzToMel(0.0)
        val highMel = hzToMel(sampleRate / 2.0)
        val melPoints = DoubleArray(numFilters + 2) { i ->
            lowMel + (highMel - lowMel) * i / (numFilters + 1)
        }
        val hzPoints = melPoints.map { melToHz(it) }
        val bins = hzPoints.map { (it / sampleRate * nFft).toInt() }

        val fb = Array(numFilters) { FloatArray(power.size) }
        for (m in 1..numFilters) {
            val fMMinus = bins[m - 1]
            val fM = bins[m]
            val fMPlus = bins[m + 1]

            for (k in fMMinus until fM) {
                if (k in power.indices) {
                    fb[m - 1][k] = ((k - fMMinus).toFloat() / (fM - fMMinus).toFloat())
                }
            }
            for (k in fM until fMPlus) {
                if (k in power.indices) {
                    fb[m - 1][k] = ((fMPlus - k).toFloat() / (fMPlus - fM).toFloat())
                }
            }
        }

        val melEnergies = FloatArray(numFilters)
        for (m in 0 until numFilters) {
            var sum = 0f
            for (k in power.indices) {
                sum += fb[m][k] * power[k]
            }
            melEnergies[m] = sum
        }
        return melEnergies
    }

    private fun dct(input: FloatArray, numCoeffs: Int): FloatArray {
        val n = input.size
        val out = FloatArray(numCoeffs)
        for (k in 0 until numCoeffs) {
            var sum = 0.0
            for (nIdx in 0 until n) {
                sum += input[nIdx] * cos(PI * k * (2 * nIdx + 1) / (2.0 * n))
            }
            out[k] = (sum * sqrt(2.0 / n)).toFloat()
        }
        return out
    }
}
