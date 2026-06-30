package com.example.quantasonaapp.ui.main

import android.Manifest
import android.content.pm.PackageManager
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.content.ContextCompat
import com.example.quantasonaapp.audio.VoiceRecorderManager
import com.example.quantasonaapp.data.InMemoryFilecoinStore
import com.example.quantasonaapp.data.NeighborView
import com.example.quantasonaapp.domain.TesseractGenerator
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

@Composable
fun HudTelemetryScreen(
    neighbors: List<NeighborView>,
    modifier: Modifier = Modifier
) {
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    val scrollState = rememberScrollState()

    var packetsPerSecond by remember { mutableStateOf(5) }
    var avgLatencyMs by remember { mutableStateOf(45) }

    val voiceRecorder = remember { VoiceRecorderManager(context) }
    val filecoinStore = remember { InMemoryFilecoinStore() }

    var isRecording by remember { mutableStateOf(false) }
    var hasPermission by remember {
        mutableStateOf(
            ContextCompat.checkSelfPermission(context, Manifest.permission.RECORD_AUDIO) == PackageManager.PERMISSION_GRANTED
        )
    }

    var currentTesseract by remember { mutableStateOf<String?>(null) }
    var txId by remember { mutableStateOf<String?>(null) }
    var isUploading by remember { mutableStateOf(false) }

    val permissionLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        hasPermission = isGranted
    }

    LaunchedEffect(Unit) {
        while (true) {
            delay(2000)
            packetsPerSecond = (2..8).random()
            avgLatencyMs = (38..48).random()
        }
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(Color(0xFF0F172A)) // Slate 900
            .verticalScroll(scrollState)
            .padding(16.dp)
    ) {
        Text(
            text = "Sovereign Mesh Node HUD",
            fontSize = 22.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFF38BDF8), // Sky 400
            modifier = Modifier.padding(bottom = 8.dp)
        )

        Card(
            colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B)),
            shape = RoundedCornerShape(12.dp),
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 8.dp)
        ) {
            Row(
                horizontalArrangement = Arrangement.SpaceAround,
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp)
            ) {
                Column {
                    Text(text = "PACKETS/SEC", fontSize = 11.sp, color = Color(0xFF94A3B8))
                    Text(text = "$packetsPerSecond Hz", fontSize = 18.sp, fontWeight = FontWeight.Bold, color = Color.White)
                }
                Column {
                    Text(text = "PEERS", fontSize = 11.sp, color = Color(0xFF94A3B8))
                    Text(text = "${neighbors.size}", fontSize = 18.sp, fontWeight = FontWeight.Bold, color = Color.White)
                }
                Column {
                    Text(text = "LATENCY", fontSize = 11.sp, color = Color(0xFF94A3B8))
                    Text(text = "$avgLatencyMs ms", fontSize = 18.sp, fontWeight = FontWeight.Bold, color = Color.White)
                }
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        Text(
            text = "Lineage Canvas (5-D State Graph)",
            fontSize = 16.sp,
            fontWeight = FontWeight.Bold,
            color = Color.White,
            modifier = Modifier.padding(bottom = 8.dp)
        )

        Canvas(
            modifier = Modifier
                .fillMaxWidth()
                .height(250.dp)
                .background(Color(0xFF1E293B), RoundedCornerShape(12.dp))
        ) {
            val center = Offset(size.width / 2, size.height / 2)
            drawCircle(color = Color(0xFF38BDF8), radius = 10f, center = center)

            neighbors.forEachIndexed { index, neighbor ->
                val angle = (2 * Math.PI * index / neighbors.size)
                val distance = 100f
                val peerX = (center.x + distance * Math.cos(angle)).toFloat()
                val peerY = (center.y + distance * Math.sin(angle)).toFloat()
                val peer = Offset(peerX, peerY)

                val edgeColor = Color(0xFF64748B).copy(alpha = neighbor.strength.coerceIn(0.1f, 1.0f))
                val strokeWidth = 2f + neighbor.strength * 8f

                drawLine(
                    color = edgeColor,
                    start = center,
                    end = peer,
                    strokeWidth = strokeWidth
                )

                val nodeColor = when (neighbor.edgeType.lowercase()) {
                    "iot" -> Color(0xFFEF4444)
                    "ble" -> Color(0xFF3B82F6)
                    "gps" -> Color(0xFF10B981)
                    else -> Color(0xFFF59E0B)
                }
                drawCircle(color = nodeColor, radius = 8f, center = peer)
            }
        }

        Spacer(modifier = Modifier.height(24.dp))

        Text(
            text = "Voice Anomaly Analyzer",
            fontSize = 18.sp,
            fontWeight = FontWeight.Bold,
            color = Color.White,
            modifier = Modifier.padding(bottom = 8.dp)
        )

        if (!hasPermission) {
            Button(
                onClick = { permissionLauncher.launch(Manifest.permission.RECORD_AUDIO) },
                colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFF59E0B)),
                modifier = Modifier.fillMaxWidth()
            ) {
                Text("Grant Microphone Permission")
            }
        } else {
            Button(
                onClick = {
                    if (isRecording) {
                        isRecording = false
                        val audioData = voiceRecorder.stopRecording()
                        val dataToHash = audioData ?: "mock_voice_anomaly_data".toByteArray()
                        currentTesseract = TesseractGenerator.generateTesseractHash(dataToHash)
                        txId = null
                    } else {
                        currentTesseract = null
                        txId = null
                        isRecording = true
                        voiceRecorder.startRecording()
                    }
                },
                colors = ButtonDefaults.buttonColors(
                    containerColor = if (isRecording) Color(0xFFEF4444) else Color(0xFF10B981)
                ),
                modifier = Modifier
                    .fillMaxWidth()
                    .height(50.dp)
            ) {
                Text(
                    text = if (isRecording) "Stop Recording" else "Record Voice Anomaly",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Bold
                )
            }

            if (isRecording) {
                Text(
                    text = "Listening to vocal frequencies...",
                    color = Color(0xFF38BDF8),
                    modifier = Modifier.padding(top = 8.dp, bottom = 8.dp).align(Alignment.CenterHorizontally)
                )
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        currentTesseract?.let { tesseract ->
            Card(
                colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B)),
                modifier = Modifier.fillMaxWidth()
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("5D Cascading Complexity Tesseract", color = Color(0xFF94A3B8), fontSize = 14.sp)
                    Text(
                        text = tesseract,
                        color = Color(0xFF34D399),
                        fontSize = 12.sp,
                        fontFamily = androidx.compose.ui.text.font.FontFamily.Monospace,
                        modifier = Modifier.padding(top = 8.dp),
                        textAlign = TextAlign.Center
                    )

                    Spacer(modifier = Modifier.height(16.dp))

                    if (txId != null) {
                        Text("Filecoin Proof-of-Replication TxID:", color = Color(0xFF94A3B8), fontSize = 14.sp)
                        Text(
                            text = txId!!,
                            color = Color(0xFF38BDF8),
                            fontSize = 12.sp,
                            fontWeight = FontWeight.Bold,
                            modifier = Modifier.padding(top = 4.dp)
                        )
                    } else if (isUploading) {
                        Box(modifier = Modifier.fillMaxWidth(), contentAlignment = Alignment.Center) {
                            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                                CircularProgressIndicator(
                                    color = Color(0xFF38BDF8),
                                    modifier = Modifier.size(24.dp)
                                )
                                Text(
                                    text = "Uploading to Filecoin...",
                                    color = Color(0xFF94A3B8),
                                    fontSize = 12.sp,
                                    modifier = Modifier.padding(top = 8.dp)
                                )
                            }
                        }
                    } else {
                        Button(
                            onClick = {
                                isUploading = true
                                scope.launch {
                                    val newTxId = filecoinStore.storeTesseract(tesseract, emptyMap())
                                    txId = newTxId
                                    isUploading = false
                                }
                            },
                            modifier = Modifier.align(Alignment.CenterHorizontally),
                            colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF6366F1))
                        ) {
                            Text("Upload to Filecoin Storage")
                        }
                    }
                }
            }
        }
    }
}
