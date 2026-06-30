package com.example.patentdemo.ui

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.nativeCanvas
import androidx.compose.ui.unit.dp
import com.example.patentdemo.dsp.FftEngine
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.nio.ByteBuffer
import java.nio.ByteOrder
import kotlin.math.PI
import kotlin.math.cos

@Composable
fun FftScreen(
    pcmData: ByteArray?,
    onFftComputed: (FloatArray) -> Unit
) {
    if (pcmData == null) {
        Text("No PCM data available")
        return
    }

    val fftResult = remember { mutableStateOf<FloatArray?>(null) }

    LaunchedEffect(Unit) {
        withContext(Dispatchers.Default) {
            // 1. Convert PCM → FloatArray
            val floats = ByteBuffer.wrap(pcmData)
                .order(ByteOrder.LITTLE_ENDIAN)
                .asShortBuffer()
                .let { buf ->
                    val arr = ShortArray(buf.limit())
                    buf.get(arr)
                    arr.map { it / 32768f }.toFloatArray()
                }

            // 2. Apply Hanning window
            val windowed = FloatArray(floats.size) { i ->
                val w = 0.5f - 0.5f * cos((2f * PI * i) / (floats.size - 1)).toFloat()
                floats[i] * w
            }

            // 3. Run FFT
            // Assuming FftEngine().process(ByteArray) or similar. Adapting slightly if needed.
            // But we will use the user's provided structure.
            val fft = FftEngine().computeFFT(windowed)

            fftResult.value = fft
            // The original logic calls onFftComputed from the button, but we can also set the state.
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            "Frequency Spectrum",
            style = MaterialTheme.typography.headlineMedium
        )

        Spacer(Modifier.height(24.dp))

        fftResult.value?.let { spectrum ->
            SpectrumGraph(spectrum)
        } ?: Text("Computing FFT…")

        Spacer(Modifier.height(24.dp))

        Button(
            onClick = { fftResult.value?.let(onFftComputed) },
            modifier = Modifier.fillMaxWidth()
        ) {
            Text("Continue")
        }
    }
}

@Composable
fun SpectrumGraph(magnitudes: FloatArray) {
    Canvas(
        modifier = Modifier
            .fillMaxWidth()
            .height(240.dp)
            .background(Color(0xFF111111))
    ) {
        val maxMag = magnitudes.maxOrNull() ?: 1f
        val logBins = magnitudes.size
        val widthStep = size.width / logBins

        // Draw magnitude bars
        magnitudes.forEachIndexed { i, mag ->
            val norm = mag / maxMag
            val x = i * widthStep
            drawLine(
                color = Color.Cyan,
                start = Offset(x, size.height),
                end = Offset(x, size.height - norm * size.height),
                strokeWidth = 2f
            )
        }

        // Draw octave markers
        val sampleRate = 44100f
        val freqs = listOf(125, 250, 500, 1000, 2000, 4000, 8000)

        freqs.forEach { f ->
            val bin = ((f / sampleRate) * logBins).toInt().coerceIn(0, logBins - 1)
            val x = bin * widthStep

            drawLine(
                color = Color(0x55FFFFFF),
                start = Offset(x, 0f),
                end = Offset(x, size.height),
                strokeWidth = 1f
            )

            drawContext.canvas.nativeCanvas.drawText(
                "${f}Hz",
                x + 4,
                20f,
                android.graphics.Paint().apply {
                    color = android.graphics.Color.WHITE
                    textSize = 24f
                }
            )
        }
    }
}
