package com.example.quantasonaapp.data

import kotlinx.serialization.Serializable

@Serializable
data class Addr5D(
    val timeIndex: Long,
    val spaceId: String,
    val lineageHash: String,
    val contentHash: String,
    val oracleContextId: String
)

enum class PriorityClass { FAST, NORMAL, BULK, LOST_DEVICE }

@Serializable
data class SilkRoadHeader(
    val src: Addr5D,
    val dst: Addr5D,
    val priority: PriorityClass,
    val maxLatencySeconds: Long,
    val expiryEpochSeconds: Long,
    val feeBudgetTritSubunits: Long
)

data class SilkRoadPacket(
    val header: SilkRoadHeader,
    val body: ByteArray
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false
        other as SilkRoadPacket
        if (header != other.header) return false
        if (!body.contentEquals(other.body)) return false
        return true
    }

    override fun hashCode(): Int {
        var result = header.hashCode()
        result = 31 * result + body.contentHashCode()
        return result
    }
}

enum class SignalType { BLE, GPS, IOT, GEO_CACHE, LOST_DEVICE }

@Serializable
data class SignalVertex(
    val id: String,
    val type: SignalType,
    val addr5d: Addr5D,
    val strength: Double,
    val observedAtEpochSeconds: Long,
    val metadata: Map<String, String>
)

enum class RewardType {
    SCAN_SIGNAL, ROUTE_PACKET, STORE_PACKET, VALIDATE_STATE,
    TREASURE_DISCOVERY, GEO_SCAN, GEO_CACHE, MATCH_GAME
}

@Serializable
data class RewardEvent(
    val nodeId: String,
    val type: RewardType,
    val amountTrits: Long,
    val meta: Map<String, String>
)
