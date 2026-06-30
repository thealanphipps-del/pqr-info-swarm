package com.swend.mesh.net

import android.content.Context
import android.util.Log
import okhttp3.*
import okio.ByteString
import java.util.concurrent.TimeUnit

class MeshClient(
    private val context: Context,
    private val onTaskReceived: (String) -> Unit
) {

    private val client = OkHttpClient.Builder()
        .pingInterval(20, TimeUnit.SECONDS)
        .retryOnConnectionFailure(true)
        .build()

    private var webSocket: WebSocket? = null

    private val request = Request.Builder()
        .url("wss://your-swen-endpoint/ws") // Antigravity will replace this
        .build()

    fun connect() {
        webSocket = client.newWebSocket(request, MeshSocketListener())
    }

    fun sendResult(json: String) {
        webSocket?.send(json)
    }

    fun close() {
        webSocket?.close(1000, "Client closed")
    }

    private inner class MeshSocketListener : WebSocketListener() {

        override fun onOpen(ws: WebSocket, response: Response) {
            Log.i("MeshClient", "Connected to SWEN")

            // Send registration payload
            val registration = """
                {
                    "type": "register",
                    "device": "android",
                    "model": "${android.os.Build.MODEL}",
                    "manufacturer": "${android.os.Build.MANUFACTURER}",
                    "capabilities": {
                        "local_llm": true,
                        "local_embeddings": true,
                        "snapdragon_npu": true
                    }
                }
            """.trimIndent()

            ws.send(registration)
        }

        override fun onMessage(ws: WebSocket, text: String) {
            onTaskReceived(text)
        }

        override fun onMessage(ws: WebSocket, bytes: ByteString) {
            onTaskReceived(bytes.utf8())
        }

        override fun onClosing(ws: WebSocket, code: Int, reason: String) {
            ws.close(1000, null)
        }

        override fun onFailure(ws: WebSocket, t: Throwable, response: Response?) {
            Log.e("MeshClient", "WebSocket error", t)

            // Attempt reconnect after delay
            Thread.sleep(2000)
            connect()
        }
    }
}
