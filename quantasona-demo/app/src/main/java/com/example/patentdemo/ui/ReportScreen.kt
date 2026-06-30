package com.example.patentdemo.ui

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.example.patentdemo.profiles.ProfileMatchResult

@Composable
fun ReportScreen(
    match: ProfileMatchResult?
) {
    if (match == null) {
        Text("No profile match available")
        return
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            "Condition Report",
            style = MaterialTheme.typography.headlineMedium
        )

        Spacer(Modifier.height(24.dp))

        ConditionReportCard(match)

        Spacer(Modifier.height(24.dp))

        Text(
            text = "This report is a demonstration based on acoustic signatures and is not a medical diagnosis.",
            style = MaterialTheme.typography.bodySmall,
            color = Color.Gray
        )
    }
}

@Composable
fun ConditionReportCard(match: ProfileMatchResult) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = Color(0xFF1A1A1A))
    ) {
        Column(modifier = Modifier.padding(20.dp)) {

            // Profile / condition header
            Text(
                text = match.profileName,
                style = MaterialTheme.typography.headlineSmall,
                color = Color.Cyan
            )

            Spacer(Modifier.height(8.dp))

            match.conditionName?.let {
                Text(
                    text = it,
                    style = MaterialTheme.typography.titleMedium,
                    color = Color(0xFFE0E027)
                )
            }

            Spacer(Modifier.height(16.dp))

            // Confidence
            ConfidenceBar(match.confidence)

            Spacer(Modifier.height(16.dp))

            // Bands summary
            Text(
                "Signature Bands:",
                style = MaterialTheme.typography.titleMedium
            )

            Spacer(Modifier.height(8.dp))

            match.matchedBands.forEach { band ->
                Text(
                    "- Bin ${band.startIndex} to ${band.endIndex} (slope=${"%.3f".format(band.slope)})",
                    style = MaterialTheme.typography.bodyMedium
                )
            }

            Spacer(Modifier.height(16.dp))

            // Narrative summary
            match.summary?.let {
                Text(
                    text = it,
                    style = MaterialTheme.typography.bodyMedium
                )
            }
        }
    }
}
