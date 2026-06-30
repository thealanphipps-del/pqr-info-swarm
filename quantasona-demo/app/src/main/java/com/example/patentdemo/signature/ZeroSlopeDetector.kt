package com.example.patentdemo.signature

import kotlin.math.abs
import com.example.patentdemo.dsp.Band

class ZeroSlopeDetector {
    
    /**
     * Identifies frequencies where the amplitude shift is less than the threshold
     * over a consecutive number of bins (representing a frequency band).
     * Based on Patent 8,346,559 B2: "change in amplitude is no greater than 0.15 dB
     * between any two measurement points for a given frequency band"
     */
    fun identifySignatures(
        spectrum: FloatArray, 
        frequencyResolution: Float, 
        amplitudeThreshold: Float = 0.15f,
        bandWidthHz: Float = 0.33f
    ): List<Float> {
        val signatures = mutableListOf<Float>()
        val binsInBand = Math.max(2, (bandWidthHz / frequencyResolution).toInt())
        
        var currentBandStart = 0
        var isZeroSlope = true
        
        for (i in 1 until spectrum.size) {
            val amplitudeShift = abs(spectrum[i] - spectrum[i - 1])
            
            if (amplitudeShift > amplitudeThreshold) {
                // Slope exceeded threshold, check if the previous run was long enough
                if (isZeroSlope && (i - currentBandStart) >= binsInBand) {
                    val centerBin = currentBandStart + (i - currentBandStart) / 2
                    signatures.add(centerBin * frequencyResolution)
                }
                currentBandStart = i
                isZeroSlope = true
            }
        }
        
        // Check the last run
        if (isZeroSlope && (spectrum.size - currentBandStart) >= binsInBand) {
            val centerBin = currentBandStart + (spectrum.size - currentBandStart) / 2
            signatures.add(centerBin * frequencyResolution)
        }
        
        return signatures
    }
    fun detect(spectrum: FloatArray): List<Band> {
        val bands = mutableListOf<Band>()
        val amplitudeThreshold = 0.15f
        val binsInBand = 2 // minimal bins to be considered a band for demo

        var currentBandStart = 0
        var isZeroSlope = true

        for (i in 1 until spectrum.size) {
            val amplitudeShift = abs(spectrum[i] - spectrum[i - 1])

            if (amplitudeShift > amplitudeThreshold) {
                if (isZeroSlope && (i - currentBandStart) >= binsInBand) {
                    bands.add(Band(currentBandStart, i, 0f))
                }
                currentBandStart = i
                isZeroSlope = true
            }
        }

        if (isZeroSlope && (spectrum.size - currentBandStart) >= binsInBand) {
            bands.add(Band(currentBandStart, spectrum.size - 1, 0f))
        }

        return bands
    }
}
