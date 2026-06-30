package com.example.patentdemo.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.example.patentdemo.dsp.Band
import com.example.patentdemo.profiles.ProfileMatchResult
import com.example.patentdemo.profiles.ProfileMatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

@Composable
fun ProfileScreen(
    zeroSlopeBands: List<Band>?,
    onProfileMatched: (ProfileMatchResult) -> Unit
) {
    if (zeroSlopeBands == null) {
        Text("No signature bands available")
        return
    }

    val matchResult = remember { mutableStateOf<ProfileMatchResult?>(null) }

    LaunchedEffect(Unit) {
        withContext(Dispatchers.Default) {
            val matcher = ProfileMatcher()
            val result = matcher.match(zeroSlopeBands)
            matchResult.value = result
            onProfileMatched(result)
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            "Physiological Profile",
            style = MaterialTheme.typography.headlineMedium
        )

        Spacer(Modifier.height(24.dp))

        matchResult.value?.let { result ->
            ProfileSummaryCard(result)
        } ?: Text("Analyzing…")

        Spacer(Modifier.height(24.dp))

        Button(
            onClick = { matchResult.value?.let(onProfileMatched) },
            modifier = Modifier.fillMaxWidth()
        ) {
            Text("Continue")
        }
    }
}

@Composable
fun ProfileSummaryCard(result: ProfileMatchResult) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = Color(0xFF1A1A1A))
    ) {
        Column(modifier = Modifier.padding(20.dp)) {

            Text(
                text = result.profileName,
                style = MaterialTheme.typography.headlineSmall,
                color = Color.Cyan
            )

            Spacer(Modifier.height(12.dp))

            ConfidenceBar(result.confidence)

            Spacer(Modifier.height(16.dp))

            Text(
                "Matched Bands:",
                style = MaterialTheme.typography.titleMedium
            )

            Spacer(Modifier.height(8.dp))

            result.matchedBands.forEach { band ->
                Text(
                    "- ${band.startIndex} to ${band.endIndex} (slope=${"%.3f".format(band.slope)})",
                    style = MaterialTheme.typography.bodyMedium
                )
            }
        }
    }
}

@Composable
fun ConfidenceBar(confidence: Float) {
    val pct = (confidence * 100).toInt()

    Column {
        Text("Confidence: $pct%", style = MaterialTheme.typography.bodyLarge)

        Spacer(Modifier.height(6.dp))

        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(12.dp)
                .background(Color.DarkGray)
        ) {
            Box(
                modifier = Modifier
                    .fillMaxHeight()
                    .fillMaxWidth(confidence.coerceIn(0f, 1f))
                    .background(Color(0xFF27A8E0))
            )
        }
    }
}
