package com.example.quantasonaapp.ui.main

import android.Manifest
import android.annotation.SuppressLint
import android.content.Context
import android.content.pm.PackageManager
import android.graphics.SurfaceTexture
import android.hardware.camera2.*
import android.os.Handler
import android.os.HandlerThread
import android.view.Surface
import android.view.TextureView
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.viewinterop.AndroidView
import androidx.core.content.ContextCompat
import com.example.quantasonaapp.data.GeologyScannerState
import com.example.quantasonaapp.data.MineralScan

@Composable
fun GeologyScannerScreen(
    scans: List<MineralScan>,
    identifiedScans: List<MineralScan>,
    scannerState: GeologyScannerState,
    onScanFinished: (MineralScan) -> Unit,
    modifier: Modifier = Modifier
) {
    val context = LocalContext.current

    var hasCameraPermission by remember {
        mutableStateOf(
            ContextCompat.checkSelfPermission(context, Manifest.permission.CAMERA) == PackageManager.PERMISSION_GRANTED
        )
    }

    val permissionLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        hasCameraPermission = isGranted
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(Color(0xFF0F172A)), // Slate 900
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "Geology Spectral Scanner",
            fontSize = 22.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFFF59E0B), // Amber 400
            modifier = Modifier.padding(bottom = 8.dp)
        )

        Text(
            text = "Point device camera at rocks to parse mineral signatures and earn multipliers.",
            fontSize = 14.sp,
            color = Color(0xFF94A3B8),
            modifier = Modifier.padding(start = 16.dp, end = 16.dp, bottom = 16.dp)
        )

        if (!hasCameraPermission) {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(200.dp)
                    .background(Color(0xFF1E293B), RoundedCornerShape(12.dp)),
                contentAlignment = Alignment.Center
            ) {
                Button(
                    onClick = { permissionLauncher.launch(Manifest.permission.CAMERA) },
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFF59E0B))
                ) {
                    Text("Grant Camera Permission")
                }
            }
        } else {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(200.dp)
                        .background(Color(0xFF1E293B), RoundedCornerShape(12.dp)),
                    contentAlignment = Alignment.Center
                ) {
                    Camera2Preview(modifier = Modifier.fillMaxSize())
                }
                Spacer(modifier = Modifier.height(8.dp))
                Button(
                    onClick = {
                        val randomScan = scans.randomOrNull() ?: MineralScan("Basalt", 6.0, "Aphanitic / Igneous", 1.2)
                        onScanFinished(randomScan)
                    },
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFF59E0B))
                ) {
                    Text("Trigger Spectral Scan (State: ${scannerState.name})")
                }
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        Text(
            text = "Identified Mineral Signatures",
            fontSize = 16.sp,
            fontWeight = FontWeight.Bold,
            color = Color.White,
            modifier = Modifier
                .align(Alignment.Start)
                .padding(bottom = 8.dp)
        )

        LazyColumn(
            verticalArrangement = Arrangement.spacedBy(8.dp),
            modifier = Modifier.fillMaxWidth().weight(1f)
        ) {
            items(identifiedScans) { scan ->
                Card(
                    colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B)),
                    shape = RoundedCornerShape(8.dp)
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Row(
                            horizontalArrangement = Arrangement.SpaceBetween,
                            modifier = Modifier.fillMaxWidth()
                        ) {
                            Text(
                                text = scan.name,
                                fontSize = 18.sp,
                                fontWeight = FontWeight.Bold,
                                color = Color.White
                            )
                            Text(
                                text = "Multiplier: x${scan.tritMultiplier}",
                                color = Color(0xFFF59E0B),
                                fontWeight = FontWeight.Bold
                            )
                        }
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(text = "Hardness: ${scan.hardness} Mohs", color = Color(0xFF94A3B8), fontSize = 14.sp)
                        Text(text = "Structure: ${scan.crystalStructure}", color = Color(0xFF94A3B8), fontSize = 14.sp)
                    }
                }
            }
        }
    }
}

@Composable
fun Camera2Preview(modifier: Modifier = Modifier) {
    val context = LocalContext.current
    var cameraDevice by remember { mutableStateOf<CameraDevice?>(null) }
    var captureSession by remember { mutableStateOf<CameraCaptureSession?>(null) }

    val backgroundThread = remember { HandlerThread("CameraBackground").apply { start() } }
    val backgroundHandler = remember { Handler(backgroundThread.looper) }

    DisposableEffect(Unit) {
        onDispose {
            captureSession?.close()
            cameraDevice?.close()
            backgroundThread.quitSafely()
        }
    }

    AndroidView(
        factory = { ctx ->
            TextureView(ctx).apply {
                surfaceTextureListener = object : TextureView.SurfaceTextureListener {
                    @SuppressLint("MissingPermission")
                    override fun onSurfaceTextureAvailable(surfaceTexture: SurfaceTexture, width: Int, height: Int) {
                        val cameraManager = ctx.getSystemService(Context.CAMERA_SERVICE) as CameraManager
                        try {
                            val cameraId = cameraManager.cameraIdList.firstOrNull { id ->
                                val chars = cameraManager.getCameraCharacteristics(id)
                                chars.get(CameraCharacteristics.LENS_FACING) == CameraCharacteristics.LENS_FACING_BACK
                            } ?: cameraManager.cameraIdList.first()

                            cameraManager.openCamera(cameraId, object : CameraDevice.StateCallback() {
                                override fun onOpened(camera: CameraDevice) {
                                    cameraDevice = camera
                                    val surface = Surface(surfaceTexture)
                                    val builder = camera.createCaptureRequest(CameraDevice.TEMPLATE_PREVIEW).apply {
                                        addTarget(surface)
                                    }
                                    camera.createCaptureSession(listOf(surface), object : CameraCaptureSession.StateCallback() {
                                        override fun onConfigured(session: CameraCaptureSession) {
                                            captureSession = session
                                            builder.set(CaptureRequest.CONTROL_MODE, CameraMetadata.CONTROL_MODE_AUTO)
                                            session.setRepeatingRequest(builder.build(), null, backgroundHandler)
                                        }

                                        override fun onConfigureFailed(session: CameraCaptureSession) {}
                                    }, backgroundHandler)
                                }

                                override fun onDisconnected(camera: CameraDevice) {
                                    camera.close()
                                    cameraDevice = null
                                }

                                override fun onError(camera: CameraDevice, error: Int) {
                                    camera.close()
                                    cameraDevice = null
                                }
                            }, backgroundHandler)
                        } catch (e: CameraAccessException) {
                            e.printStackTrace()
                        }
                    }

                    override fun onSurfaceTextureSizeChanged(surface: SurfaceTexture, width: Int, height: Int) {}
                    override fun onSurfaceTextureDestroyed(surface: SurfaceTexture): Boolean = true
                    override fun onSurfaceTextureUpdated(surface: SurfaceTexture) {}
                }
            }
        },
        modifier = modifier
    )
}
