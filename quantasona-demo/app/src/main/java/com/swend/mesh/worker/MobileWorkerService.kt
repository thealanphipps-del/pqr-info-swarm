package com.swend.mesh.worker

import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.IBinder
import android.util.Log
import com.swend.mesh.ui.MeshNotifications
import com.swend.mesh.protocol.TaskProtocol
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import okio.ByteString
import org.json.JSONObject

class MobileWorkerService : Service() {

    companion object {
        private const val TAG = "MobileWorkerService"
        private const val WS_URL = "ws://201.mh:8080/ws"
    }

    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())
    private val client by lazy { OkHttpClient() }

    private var webSocket: WebSocket? = null
    private lateinit var embeddingHandler: EmbeddingTaskHandler
    private var deviceName: String = "android-node"

    override fun onCreate() {
        super.onCreate()

        embeddingHandler = EmbeddingTaskHandler(this)
        deviceName = getSharedPreferences("mesh", Context.MODE_PRIVATE)
            .getString("device_name", "android-node") ?: "android-node"

        val notification = MeshNotifications.createForegroundNotification(this)
        startForeground(1, notification)

        connectWebSocket()
        startTelemetryLoop()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        // Keep service running
        return START_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        Log.i(TAG, "Service destroyed, closing WebSocket")
        webSocket?.close(1000, "Service destroyed")
        scope.cancel()
    }

    override fun onBind(intent: Intent?): IBinder? = null

    // -----------------------------
    // WebSocket connection
    // -----------------------------
    private fun connectWebSocket() {
        Log.i(TAG, "Connecting to SWEN at $WS_URL")

        val request = Request.Builder()
            .url(WS_URL)
            .build()

        webSocket = client.newWebSocket(request, object : WebSocketListener() {

            override fun onOpen(ws: WebSocket, response: okhttp3.Response) {
                Log.i(TAG, "WebSocket opened")
                sendRegistration(ws)
            }

            override fun onMessage(ws: WebSocket, text: String) {
                scope.launch {
                    handleMessage(text)
                }
            }

            override fun onMessage(ws: WebSocket, bytes: ByteString) {
                // Not used; mesh uses text frames
            }

            override fun onClosing(ws: WebSocket, code: Int, reason: String) {
                Log.i(TAG, "WebSocket closing: $code / $reason")
                ws.close(code, reason)
            }

            override fun onFailure(ws: WebSocket, t: Throwable, response: okhttp3.Response?) {
                Log.e(TAG, "WebSocket failure: ${t.message}", t)
                // Simple reconnect strategy
                scope.launch {
                    delay(5000)
                    connectWebSocket()
                }
            }
        })
    }

    private fun sendRegistration(ws: WebSocket) {
        val payload = JSONObject().apply {
            put("type", "register")
            put("device", deviceName)
            put("platform", "android")
        }
        ws.send(payload.toString())
        Log.i(TAG, "Sent registration for device: $deviceName")
    }

    // -----------------------------
    // Telemetry loop
    // -----------------------------
    private fun startTelemetryLoop() {
        scope.launch {
            while (true) {
                sendTelemetry()
                delay(30_000L)
            }
        }
    }

    private fun sendTelemetry() {
        val ws = webSocket ?: return

        val payload = JSONObject().apply {
            put("type", "telemetry")
            put("device", deviceName)
            put("status", "ok")
            put("timestamp", System.currentTimeMillis())
        }

        ws.send(payload.toString())
    }

    // -----------------------------
    // Task handling
    // -----------------------------
    private suspend fun handleMessage(text: String) {
        try {
            val task = TaskProtocol.parseTask(text)

            when (task.type) {
                "embedding" -> {
                    val resultJson = embeddingHandler.handle(task)
                    webSocket?.send(resultJson)
                }
                else -> {
                    val resultJson = TaskProtocol.buildResult(
                        id = task.id,
                        result = "Unsupported task type: ${task.type}",
                        status = "error"
                    )
                    webSocket?.send(resultJson)
                }
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error handling message: ${e.message}", e)
        }
    }
}
