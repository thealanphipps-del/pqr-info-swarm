package com.example.patentdemo.report

import com.example.patentdemo.profiles.ProfileMatchResult

data class ConditionReport(
    val conditionName: String,
    val summary: String
)

class ConditionReportBuilder {

    fun build(match: ProfileMatchResult): ConditionReport {
        val condition = determineCondition(match.profileName)
        val summary = buildSummary(match, condition)

        return ConditionReport(
            conditionName = condition,
            summary = summary
        )
    }

    private fun determineCondition(profileName: String): String =
        when (profileName) {
            "Respiratory Pattern A" -> "Respiratory Stress Indicator"
            "Vocal Fold Pattern B" -> "Vocal Fold Tension Signature"
            "Sinus Resonance C" -> "Sinus Congestion Pattern"
            else -> "General Acoustic Signature"
        }

    private fun buildSummary(
        match: ProfileMatchResult,
        condition: String
    ): String {
        val confidencePct = (match.confidence * 100).toInt()

        val bandDescriptions = match.matchedBands.joinToString("\n") { band ->
            "- Flat band from bin ${band.startIndex} to ${band.endIndex} (slope=${"%.3f".format(band.slope)})"
        }

        return """
            Condition: $condition
            Confidence: $confidencePct%

            The acoustic analysis identified ${match.matchedBands.size} substantial zero‑slope regions within the frequency spectrum. These regions are characteristic of the "$condition" profile.

            Matched Signature Bands:
            $bandDescriptions

            This pattern reflects the presence of stable harmonic energy across specific frequency ranges, consistent with the physiological characteristics associated with this condition.
        """.trimIndent()
    }
}
