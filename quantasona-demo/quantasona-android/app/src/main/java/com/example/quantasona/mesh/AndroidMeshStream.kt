package com.example.quantasona.mesh

import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import quantasona.epistemic.MeshDeltaServiceGrpc
import quantasona.epistemic.StreamRequest

class AndroidMeshStream(
    host: String = "10.0.2.2",
    port: Int = 50051
) {
    private val channel = ManagedChannelBuilder
        .forAddress(host, port)
        .usePlaintext()
        .build()

    private val stub = MeshDeltaServiceGrpc.newBlockingStub(channel)

    suspend fun stream(onDelta: (String) -> Unit) = withContext(Dispatchers.IO) {
        while (true) {
            try {
                val req = StreamRequest.newBuilder().build()
                val stream = stub.streamDeltas(req)

                while (true) {
                    val delta = stream.next()
                    
                    onDelta("Collapse: ${delta.sigmaId} (conf=${delta.confidence})")
                }
            } catch (e: kotlinx.coroutines.CancellationException) {
                throw e
            } catch (e: Exception) {
                onDelta("Stream fault intercepted: ${e.message}. Reconnecting gracefully...")
                kotlinx.coroutines.delay(2000)
            }
        }
    }
}
