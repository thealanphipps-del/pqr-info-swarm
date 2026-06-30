package com.example.quantasonaapp.data

data class LineageRecord(
    val lineageHash: String,
    val version: Long,
    val timestampEpochSeconds: Long
)

data class VertexRecord(
    val addr: Addr5D,
    var payloadData: ByteArray,
    var metadata: Map<String, Any>,
    var lastUpdatedLineage: LineageRecord? = null
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false
        other as VertexRecord
        if (addr != other.addr) return false
        if (!payloadData.contentEquals(other.payloadData)) return false
        if (metadata != other.metadata) return false
        if (lastUpdatedLineage != other.lastUpdatedLineage) return false
        return true
    }

    override fun hashCode(): Int {
        var result = addr.hashCode()
        result = 31 * result + payloadData.contentHashCode()
        result = 31 * result + metadata.hashCode()
        result = 31 * result + (lastUpdatedLineage?.hashCode() ?: 0)
        return result
    }
}

data class EdgeRecord(
    val source: Addr5D,
    val target: Addr5D,
    var edgeType: String,
    var connectionStrength: Float,
    val edgeMetadata: Map<String, String>,
    var genesisLineage: LineageRecord? = null
) {
    fun normalizedKey(): Pair<Addr5D, Addr5D> =
        if (source.hashCode() <= target.hashCode()) source to target else target to source
}

interface FiveDVertexStore {
    suspend fun getVertex(addr: Addr5D): VertexRecord?
    suspend fun putVertex(record: VertexRecord)
    suspend fun mergeVertex(remote: VertexRecord)
}

interface FiveDEdgeStore {
    suspend fun getEdgesFor(addr: Addr5D): List<EdgeRecord>
    suspend fun putEdge(edge: EdgeRecord)
    suspend fun mergeEdge(remote: EdgeRecord)
}

class InMemoryFiveDStore : FiveDVertexStore, FiveDEdgeStore {
    private val vertices = mutableMapOf<Addr5D, VertexRecord>()
    private val edges = mutableMapOf<Pair<Addr5D, Addr5D>, EdgeRecord>()

    override suspend fun getVertex(addr: Addr5D): VertexRecord? = vertices[addr]

    override suspend fun putVertex(record: VertexRecord) {
        vertices[record.addr] = record
    }

    override suspend fun mergeVertex(remote: VertexRecord) {
        val local = vertices[remote.addr]
        vertices[remote.addr] = mergeVertexCrdt(local, remote)
    }

    override suspend fun getEdgesFor(addr: Addr5D): List<EdgeRecord> =
        edges.values.filter { it.source == addr || it.target == addr }

    override suspend fun putEdge(edge: EdgeRecord) {
        edges[edge.normalizedKey()] = edge
    }

    override suspend fun mergeEdge(remote: EdgeRecord) {
        val key = remote.normalizedKey()
        val local = edges[key]
        edges[key] = mergeEdgeCrdt(local, remote)
    }

    private fun mergeVertexCrdt(local: VertexRecord?, remote: VertexRecord): VertexRecord =
        when {
            local == null -> remote
            (remote.lastUpdatedLineage?.version ?: 0L) >
                (local.lastUpdatedLineage?.version ?: 0L) -> remote
            else -> local
        }

    private fun mergeEdgeCrdt(local: EdgeRecord?, remote: EdgeRecord): EdgeRecord =
        when {
            local == null -> remote
            (remote.genesisLineage?.version ?: 0L) >
                (local.genesisLineage?.version ?: 0L) -> remote
            else -> local
        }
}

data class Neighbor(
    val addr: Addr5D,
    val edge: EdgeRecord
)

data class PathCriteria(
    val minStrength: Float = 0.0f,
    val allowedEdgeTypes: Set<String>? = null
)

data class ConvergenceDelta(
    val vertexDiff: Map<String, Any?>,
    val edgeDiff: Map<String, Any?>
)

interface GraphQueryEngine {
    suspend fun queryNeighbors(addr: Addr5D, minStrength: Float = 0.0f): List<Neighbor>
    suspend fun findPathByCriteria(
        start: Addr5D,
        end: Addr5D,
        criteria: PathCriteria
    ): List<Addr5D>?

    suspend fun calculateConvergence(addrA: Addr5D, addrB: Addr5D): ConvergenceDelta
}

class DefaultGraphQueryEngine(
    private val vertexStore: FiveDVertexStore,
    private val edgeStore: FiveDEdgeStore
) : GraphQueryEngine {

    override suspend fun queryNeighbors(addr: Addr5D, minStrength: Float): List<Neighbor> =
        edgeStore.getEdgesFor(addr)
            .filter { it.connectionStrength >= minStrength }
            .map { edge ->
                val neighborAddr = if (edge.source == addr) edge.target else edge.source
                Neighbor(neighborAddr, edge)
            }

    override suspend fun findPathByCriteria(
        start: Addr5D,
        end: Addr5D,
        criteria: PathCriteria
    ): List<Addr5D>? {
        val visited = mutableSetOf<Addr5D>()
        val queue = ArrayDeque<List<Addr5D>>()
        queue.add(listOf(start))

        while (queue.isNotEmpty()) {
            val path = queue.removeFirst()
            val current = path.last()
            if (current == end) return path
            if (!visited.add(current)) continue

            val neighbors = queryNeighbors(current, criteria.minStrength)
                .filter { n ->
                    criteria.allowedEdgeTypes?.contains(n.edge.edgeType) ?: true
                }

            neighbors.forEach { n ->
                if (!visited.contains(n.addr)) {
                    queue.add(path + n.addr)
                }
            }
        }
        return null
    }

    override suspend fun calculateConvergence(addrA: Addr5D, addrB: Addr5D): ConvergenceDelta {
        val vA = vertexStore.getVertex(addrA)
        val vB = vertexStore.getVertex(addrB)

        val vertexDiff = mutableMapOf<String, Any?>()
        if (vA != null && vB != null) {
            vertexDiff["payload_diff"] = (vA.payloadData.contentHashCode() != vB.payloadData.contentHashCode())
            vertexDiff["metadata_diff"] = (vA.metadata != vB.metadata)
        }

        val edgesA = edgeStore.getEdgesFor(addrA)
        val edgesB = edgeStore.getEdgesFor(addrB)

        val edgeDiff = mutableMapOf<String, Any?>()
        edgeDiff["edge_count_diff"] = edgesA.size != edgesB.size

        return ConvergenceDelta(
            vertexDiff = vertexDiff,
            edgeDiff = edgeDiff
        )
    }
}
