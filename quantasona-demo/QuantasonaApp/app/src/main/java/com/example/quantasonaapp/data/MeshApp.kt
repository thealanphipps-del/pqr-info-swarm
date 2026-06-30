package com.example.quantasonaapp.data

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.launch

fun meshApp(
    nodeClient: MeshNodeClient,
    graph: GraphQueryEngine,
    block: MeshAppScope.() -> Unit
): MeshAppRuntime {
    val scope = MeshAppScopeImpl(nodeClient, graph)
    scope.block()
    return MeshAppRuntime(scope).also { it.start() }
}

interface MeshAppScope {
    fun onStart(block: suspend MeshContext.() -> Unit)
    fun onMessage(block: suspend MeshContext.(packet: SilkRoadPacket) -> Unit)
    fun every(millis: Long, block: suspend MeshContext.() -> Unit)
    fun ifNeighborsMatch(
        minStrength: Float = 0.5f,
        block: suspend MeshContext.(neighbors: List<NeighborView>) -> Unit
    )
}

data class NeighborView(
    val addr: Addr5D,
    val edgeType: String,
    val strength: Float
) {
    val isConnected: Boolean
        get() = strength > 0.1f
}

interface MeshContext {
    val nodeId: String
    val nodeClient: MeshNodeClient
    suspend fun send(packet: SilkRoadPacket)
    suspend fun queryNeighbors(minStrength: Float = 0.0f): List<NeighborView>
}

internal class MeshAppScopeImpl(
    val nodeClient: MeshNodeClient,
    val graph: GraphQueryEngine
) : MeshAppScope {

    var onStartBlock: (suspend MeshContext.() -> Unit)? = null
    val onMessageBlocks = mutableListOf<suspend MeshContext.(SilkRoadPacket) -> Unit>()
    val periodicBlocks = mutableListOf<Pair<Long, suspend MeshContext.() -> Unit>>()
    val neighborBlocks = mutableListOf<Pair<Float, suspend MeshContext.(List<NeighborView>) -> Unit>>()

    override fun onStart(block: suspend MeshContext.() -> Unit) {
        onStartBlock = block
    }

    override fun onMessage(block: suspend MeshContext.(packet: SilkRoadPacket) -> Unit) {
        onMessageBlocks += block
    }

    override fun every(millis: Long, block: suspend MeshContext.() -> Unit) {
        periodicBlocks += millis to block
    }

    override fun ifNeighborsMatch(
        minStrength: Float,
        block: suspend MeshContext.(neighbors: List<NeighborView>) -> Unit
    ) {
        neighborBlocks += minStrength to block
    }
}

class MeshAppRuntime internal constructor(
    private val scopeImpl: MeshAppScopeImpl
) {
    private val coroutineScope = CoroutineScope(Dispatchers.Default)

    fun start() {
        val ctx = MeshContextImpl(scopeImpl)

        scopeImpl.onStartBlock?.let { block ->
            coroutineScope.launch { ctx.block() }
        }

        coroutineScope.launch {
            ctx.nodeClient.incomingPackets().collect { pkt ->
                scopeImpl.onMessageBlocks.forEach { block ->
                    launch { ctx.block(pkt) }
                }
            }
        }

        scopeImpl.periodicBlocks.forEach { (interval, block) ->
            coroutineScope.launch {
                while (true) {
                    ctx.block()
                    delay(interval)
                }
            }
        }

        scopeImpl.neighborBlocks.forEach { (minStrength, block) ->
            coroutineScope.launch {
                while (true) {
                    val neighbors = ctx.queryNeighbors(minStrength)
                    if (neighbors.isNotEmpty()) {
                        ctx.block(neighbors)
                    }
                    delay(1000L)
                }
            }
        }
    }
}

private class MeshContextImpl(
    private val scopeImpl: MeshAppScopeImpl
) : MeshContext {

    override val nodeClient: MeshNodeClient get() = scopeImpl.nodeClient
    private val graph: GraphQueryEngine get() = scopeImpl.graph

    override val nodeId: String
        get() = nodeClient.nodeId

    override suspend fun send(packet: SilkRoadPacket) {
        nodeClient.sendPacket(packet)
    }

    override suspend fun queryNeighbors(minStrength: Float): List<NeighborView> {
        val addr = nodeClient.currentAddr5D()
        val neighbors = graph.queryNeighbors(addr, minStrength)
        return neighbors.map {
            NeighborView(
                addr = it.addr,
                edgeType = it.edge.edgeType,
                strength = it.edge.connectionStrength
            )
        }
    }
}
