package com.example.patentdemo.profiles

import kotlin.math.abs
import com.example.patentdemo.dsp.Band
import com.example.patentdemo.report.ConditionReportBuilder

data class ProfileMatchResult(
    val profileName: String,
    val confidence: Float,
    val matchedBands: List<Band>,
    val conditionName: String? = null,
    val summary: String? = null
)

data class MatchedProfile(
    val profile: PhysiologicalProfile,
    val matchConfidence: Float,
    val matchedFrequencies: List<Float>
)


class ProfileMatcher {
    
    fun match(zeroSlopeBands: List<Band>): ProfileMatchResult {
        // Simple demo implementation to satisfy UI
        val confidence = if (zeroSlopeBands.isNotEmpty()) 0.85f else 0.0f
        
        val baseResult = ProfileMatchResult(
            profileName = "Respiratory Pattern A",
            confidence = confidence,
            matchedBands = zeroSlopeBands
        )
        
        val builder = ConditionReportBuilder()
        val report = builder.build(baseResult)

        return baseResult.copy(
            conditionName = report.conditionName,
            summary = report.summary
        )
    }

    /**
     * Cross-references detected signatures against known physiological profiles.
     * Evaluates base frequencies and their harmonics (octaves) as per Patent 8,346,559 B2.
     */
    fun crossReference(
        signatures: List<Float>, 
        profiles: List<PhysiologicalProfile>,
        toleranceHz: Float = 0.04f,
        octavesToExamine: Int = 6
    ): List<MatchedProfile> {
        val results = mutableListOf<MatchedProfile>()
        
        for (profile in profiles) {
            val matchedFreqs = mutableListOf<Float>()
            var totalPossibleHits = 0
            var actualHits = 0
            
            for (baseFreq in profile.frequencies) {
                // Check base frequency and its octaves
                for (octave in 1..octavesToExamine) {
                    val targetFreq = baseFreq * octave
                    totalPossibleHits++
                    
                    // Check if any signature matches the target frequency within tolerance
                    val match = signatures.find { abs(it - targetFreq) <= toleranceHz }
                    if (match != null) {
                        actualHits++
                        matchedFreqs.add(match)
                    }
                }
            }
            
            val confidence = if (totalPossibleHits > 0) {
                actualHits.toFloat() / totalPossibleHits.toFloat()
            } else {
                0f
            }
            
            if (confidence > 0) {
                results.add(MatchedProfile(profile, confidence, matchedFreqs))
            }
        }
        
        // Sort by highest confidence first
        return results.sortedByDescending { it.matchConfidence }
    }
}
