package com.example.patentdemo.ui

import android.Manifest
import android.content.pm.PackageManager
import android.media.AudioFormat
import android.media.AudioRecord
import android.media.MediaRecorder
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
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
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.content.ContextCompat
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import java.nio.ByteBuffer
import java.nio.ByteOrder

@Composable
fun RecordScreen(
    onRecorded: (ByteArray) -> Unit
) {
    val context = LocalContext.current
    val isRecording = remember { mutableStateOf(false) }
    val waveform = remember { mutableStateListOf<Float>() }
    val pcmBuffer = remember { mutableStateListOf<Byte>() }

    val hasPermission = remember {
        mutableStateOf(
            ContextCompat.checkSelfPermission(
                context,
                Manifest.permission.RECORD_AUDIO
            ) == PackageManager.PERMISSION_GRANTED
        )
    }

    val permissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        hasPermission.value = isGranted
    }

    LaunchedEffect(Unit) {
        if (!hasPermission.value) {
            permissionLauncher.launch(Manifest.permission.RECORD_AUDIO)
        }
    }

    val recorder = remember {
        if (hasPermission.value) {
            AudioRecord(
                MediaRecorder.AudioSource.MIC,
                44100,
                AudioFormat.CHANNEL_IN_MONO,
                AudioFormat.ENCODING_PCM_16BIT,
                44100 * 2
            )
        } else null
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {

        Text(
            "Record Sound Print",
            style = MaterialTheme.typography.headlineMedium
        )

        Spacer(Modifier.height(32.dp))

        // Waveform preview
        WaveformView(samples = waveform)

        Spacer(Modifier.height(32.dp))

        Button(
            onClick = {
                if (!hasPermission.value) {
                    permissionLauncher.launch(Manifest.permission.RECORD_AUDIO)
                    return@Button
                }

                val safeRecorder = recorder ?: return@Button

                if (!isRecording.value) {
                    isRecording.value = true
                    pcmBuffer.clear()
                    waveform.clear()

                    safeRecorder.startRecording()

                    // Background audio capture
                    CoroutineScope(Dispatchers.IO).launch {
                        val buffer = ByteArray(2048)
                        while (isRecording.value) {
                            val read = safeRecorder.read(buffer, 0, buffer.size)
                            if (read > 0) {
                                pcmBuffer.addAll(buffer.take(read))
                                // Convert to floats for waveform
                                val floats = ShortArray(read / 2)
                                ByteBuffer.wrap(buffer).order(ByteOrder.LITTLE_ENDIAN).asShortBuffer().get(floats)
                                waveform.addAll(floats.map { it / 32768f })
                            }
                        }
                    }
                } else {
                    isRecording.value = false
                    safeRecorder.stop()
                    onRecorded(pcmBuffer.toByteArray())
                }
            },
            modifier = Modifier.fillMaxWidth()
        ) {
            Text(if (isRecording.value) "Stop & Analyze" else "Start Recording")
        }
    }
}

@Composable
fun WaveformView(samples: List<Float>) {
    Canvas(
        modifier = Modifier
            .fillMaxWidth()
            .height(120.dp)
            .background(Color(0xFF111111))
    ) {
        val mid = size.height / 2f
        val step = size.width / samples.size.coerceAtLeast(1)

        samples.forEachIndexed { i, v ->
            val x = i * step
            drawLine(
                color = Color.Cyan,
                start = Offset(x, mid - v * mid),
                end = Offset(x, mid + v * mid),
                strokeWidth = 2f
            )
        }
    }
}
