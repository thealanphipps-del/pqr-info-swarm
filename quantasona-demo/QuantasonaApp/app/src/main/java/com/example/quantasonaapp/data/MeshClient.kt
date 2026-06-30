package com.example.quantasonaapp.data

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.asSharedFlow

interface MeshNodeClient {
    val nodeId: String
    suspend fun sendPacket(packet: SilkRoadPacket)
    fun incomingPackets(): Flow<SilkRoadPacket>
    suspend fun reportSignal(signal: SignalVertex)
    fun rewardEvents(): Flow<RewardEvent>
    fun currentAddr5D(): Addr5D
}

class MockMeshClient(override val nodeId: String = "node-android-local") : MeshNodeClient {
    private val _packets = MutableSharedFlow<SilkRoadPacket>(extraBufferCapacity = 16)
    private val _rewards = MutableSharedFlow<RewardEvent>(extraBufferCapacity = 16)

    override suspend fun sendPacket(packet: SilkRoadPacket) {
        _packets.emit(packet)
    }

    override fun incomingPackets(): Flow<SilkRoadPacket> = _packets.asSharedFlow()

    override suspend fun reportSignal(signal: SignalVertex) {
        val reward = RewardEvent(
            nodeId = nodeId,
            type = RewardType.SCAN_SIGNAL,
            amountTrits = (10..50).random().toLong(),
            meta = mapOf("signal_id" to signal.id, "type" to signal.type.name)
        )
        _rewards.emit(reward)
    }

    suspend fun triggerReward(type: RewardType, amount: Long, meta: Map<String, String> = emptyMap()) {
        _rewards.emit(RewardEvent(nodeId, type, amount, meta))
    }

    override fun rewardEvents(): Flow<RewardEvent> = _rewards.asSharedFlow()

    override fun currentAddr5D(): Addr5D = Addr5D(
        timeIndex = System.currentTimeMillis() / 1000,
        spaceId = "pancreas",
        lineageHash = "local_node_lineage",
        contentHash = "empty_payload",
        oracleContextId = "local_mesh"
    )
}
