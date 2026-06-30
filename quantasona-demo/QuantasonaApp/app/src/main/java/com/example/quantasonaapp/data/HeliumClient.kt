package com.example.quantasonaapp.data

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.launch

data class HeliumBeacon(
    val id: String,
    val rssi: Int,
    val frequency: Double?,
    val lat: Double?,
    val lon: Double?,
    val timestampEpochSeconds: Long,
    val metadata: Map<String, String>
)

data class HeliumReward(
    val hotspotId: String,
    val amountMobile: Double,
    val reason: String,
    val timestampEpochSeconds: Long
)

interface HeliumClient {
    fun beacons(): Flow<HeliumBeacon>
    fun rewards(): Flow<HeliumReward>
}

class HeliumClientImpl(
    private val hotspotId: String
) : HeliumClient {

    private val peerHotspots = listOf("hotspot-alpha", "hotspot-beta", "hotspot-gamma")

    override fun beacons(): Flow<HeliumBeacon> = flow {
        var index = 0
        while (true) {
            val peerId = peerHotspots[index % peerHotspots.size]
            emit(
                HeliumBeacon(
                    id = peerId,
                    rssi = (-100..-40).random(),
                    frequency = 915.2,
                    lat = 37.7749,
                    lon = -122.4194,
                    timestampEpochSeconds = System.currentTimeMillis() / 1000,
                    metadata = mapOf("channel" to (index % 3).toString(), "snr" to "12.5")
                )
            )
            index++
            delay(1500) // Alternate beacons every 1.5 seconds to build dynamic updates
        }
    }

    override fun rewards(): Flow<HeliumReward> = flow {
        emit(
            HeliumReward(
                hotspotId = hotspotId,
                amountMobile = 15.25,
                reason = "poc_challenger",
                timestampEpochSeconds = System.currentTimeMillis() / 1000
            )
        )
    }
}

object HeliumSignalAdapter {
    fun normalizeRssi(rssi: Double): Float {
        return ((rssi - (-100.0)) / (-40.0 - (-100.0))).toFloat().coerceIn(0.0f, 1.0f)
    }

    fun toSignalVertex(beacon: HeliumBeacon, localAddr: Addr5D): SignalVertex {
        val peerAddr = localAddr.copy(
            spaceId = beacon.id,
            lineageHash = "lineage-" + beacon.id.hashCode()
        )
        return SignalVertex(
            id = beacon.id,
            type = when (beacon.id) {
                "hotspot-alpha" -> SignalType.IOT
                "hotspot-beta" -> SignalType.BLE
                "hotspot-gamma" -> SignalType.GPS
                else -> SignalType.IOT
            },
            addr5d = peerAddr,
            strength = beacon.rssi.toDouble(),
            observedAtEpochSeconds = beacon.timestampEpochSeconds,
            metadata = beacon.metadata
        )
    }
}

object HeliumRewardMapper {
    fun toRewardEvent(nodeId: String, hr: HeliumReward): RewardEvent {
        return RewardEvent(
            nodeId = nodeId,
            type = RewardType.SCAN_SIGNAL,
            amountTrits = (hr.amountMobile * 10).toLong(),
            meta = mapOf(
                "helium_hotspot" to hr.hotspotId,
                "reason" to hr.reason,
                "helium_mobile_amount" to hr.amountMobile.toString()
            )
        )
    }
}

class HeliumMeshBridge(
    private val helium: HeliumClient,
    private val mesh: MeshNodeClient,
    private val repository: DataRepository,
    private val scope: CoroutineScope = CoroutineScope(Dispatchers.Default)
) {
    fun start() {
        // Stream beacons -> Sovereign Mesh signals
        scope.launch {
            helium.beacons().collect { beacon ->
                val addr = mesh.currentAddr5D()
                val signal = HeliumSignalAdapter.toSignalVertex(beacon, addr)
                mesh.reportSignal(signal)
                repository.recordSignal(signal)
            }
        }
    }
}
