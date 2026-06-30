package com.example.patentdemo.ui

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.*
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.example.patentdemo.dsp.Band
import com.example.patentdemo.signature.ZeroSlopeDetector
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

@Composable
fun SignatureScreen(
    fftData: FloatArray?,
    onZeroSlopeDetected: (List<Band>) -> Unit
) {
    if (fftData == null) {
        Text("No FFT data available")
        return
    }

    val bands = remember { mutableStateOf<List<Band>?>(null) }

    LaunchedEffect(Unit) {
        withContext(Dispatchers.Default) {
            val detector = ZeroSlopeDetector()
            val detected = detector.detect(fftData)
            bands.value = detected
            onZeroSlopeDetected(detected)
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            "Zero‑Slope Signatures",
            style = MaterialTheme.typography.headlineMedium
        )

        Spacer(Modifier.height(24.dp))

        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(240.dp)
        ) {
            SpectrumGraph(fftData)

            bands.value?.let { list ->
                ZeroSlopeOverlay(
                    magnitudes = fftData,
                    bands = list
                )
            }
        }

        Spacer(Modifier.height(24.dp))

        Button(
            onClick = { bands.value?.let(onZeroSlopeDetected) },
            modifier = Modifier.fillMaxWidth()
        ) {
            Text("Continue")
        }
    }
}

@Composable
fun ZeroSlopeOverlay(
    magnitudes: FloatArray,
    bands: List<Band>
) {
    Canvas(
        modifier = Modifier
            .fillMaxSize()
    ) {
        val width = size.width
        val height = size.height
        val binWidth = width / magnitudes.size

        bands.forEach { band ->
            val startX = band.startIndex * binWidth
            val endX = band.endIndex * binWidth

            drawRect(
                color = Color(0x5527A8E0),
                topLeft = Offset(startX, 0f),
                size = Size(endX - startX, height)
            )
        }
    }
}
