package com.example.quantasonaapp.data

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class HpaGene(
    val symbol: String,
    val ensemblId: String,
    val description: String,
    val subcellularLocations: List<String>,
    val tissueExpression: Map<String, String> // Tissue -> Expression Level
)

enum class GeologyScannerState {
    IDLE, SCANNING, COMPLETED, TIMEOUT
}

data class MineralScan(
    val name: String,
    val hardness: Double,
    val crystalStructure: String,
    val tritMultiplier: Double
)

interface DataRepository {
    val data: Flow<List<String>>
    val tritBalance: StateFlow<Long>
    val hpaGenes: StateFlow<List<HpaGene>>
    val client: MockMeshClient
    val neighbors: StateFlow<List<NeighborView>>

    // NEW: Expose Gem Match Score
    val gemScore: StateFlow<Int>

    // NEW: Expose Geology Scanner state and scans
    val scannerState: StateFlow<GeologyScannerState>
    val availableScans: List<MineralScan>
    val identifiedScans: StateFlow<List<MineralScan>>

    suspend fun addTrits(amount: Long)
    suspend fun recordSignal(signal: SignalVertex)
    
    // NEW: Expose functions to update scores and scanner states
    suspend fun incrementGemScore(points: Int)
    fun setScannerState(state: GeologyScannerState)
    fun recordMineralScan(scan: MineralScan)
}

class DefaultDataRepository(
    override val client: MockMeshClient = MockMeshClient(),
    private val heliumClient: HeliumClient = HeliumClientImpl("hotspot-android-local")
) : DataRepository {

    override val data: Flow<List<String>> = flow {
        emit(listOf("Sovereign Mesh Core Active", "Tesseract Genome Syncing"))
    }

    private val _tritBalance = MutableStateFlow(1000L)
    override val tritBalance: StateFlow<Long> = _tritBalance.asStateFlow()

    // 5-D CRDT Graph Engine integration
    val store = InMemoryFiveDStore()
    val graphEngine = DefaultGraphQueryEngine(store, store)

    private val _neighbors = MutableStateFlow<List<NeighborView>>(emptyList())
    override val neighbors: StateFlow<List<NeighborView>> = _neighbors.asStateFlow()

    private val _gemScore = MutableStateFlow(0)
    override val gemScore: StateFlow<Int> = _gemScore.asStateFlow()

    private val _scannerState = MutableStateFlow(GeologyScannerState.IDLE)
    override val scannerState: StateFlow<GeologyScannerState> = _scannerState.asStateFlow()

    override val availableScans: List<MineralScan> = listOf(
        MineralScan("Basalt", 6.0, "Aphanitic / Igneous", 1.2),
        MineralScan("Granite", 6.5, "Phaneritic / Plutonic", 1.5),
        MineralScan("Quartzite", 7.0, "Granoblastic / Metamorphic", 2.0)
    )

    private val _identifiedScans = MutableStateFlow<List<MineralScan>>(emptyList())
    override val identifiedScans: StateFlow<List<MineralScan>> = _identifiedScans.asStateFlow()

    private val _hpaGenes = MutableStateFlow(
        listOf(
            HpaGene(
                symbol = "ZNF268",
                ensemblId = "ENSG00000090612",
                description = "Zinc finger protein 268",
                subcellularLocations = listOf("Nucleoplasm", "Nuclear bodies"),
                tissueExpression = mapOf(
                    "Heart muscle" to "Medium",
                    "Cerebrum" to "Medium",
                    "Liver" to "Low",
                    "Appendix" to "High",
                    "Bone marrow" to "Low"
                )
            ),
            HpaGene(
                symbol = "TP53",
                ensemblId = "ENSG00000141510",
                description = "Tumor protein p53",
                subcellularLocations = listOf("Nucleoplasm", "Mitochondria", "Cytosol"),
                tissueExpression = mapOf(
                    "Heart muscle" to "Low",
                    "Cerebrum" to "Low",
                    "Liver" to "Medium",
                    "Spleen" to "High"
                )
            ),
            HpaGene(
                symbol = "ERBB2",
                ensemblId = "ENSG00000141736",
                description = "Erb-b2 receptor tyrosine kinase 2",
                subcellularLocations = listOf("Plasma membrane"),
                tissueExpression = mapOf(
                    "Heart muscle" to "Medium",
                    "Cerebrum" to "Not Detected",
                    "Gastrointestinal tract" to "High"
                )
            )
        )
    )
    override val hpaGenes: StateFlow<List<HpaGene>> = _hpaGenes.asStateFlow()

    private val repositoryScope = CoroutineScope(Dispatchers.Default)
    private val heliumMeshBridge = HeliumMeshBridge(heliumClient, client, this, repositoryScope)

    init {
        // Bootstrap the 5-D CRDT graph engine with metabolic state configurations
        InsulinLattice.bootstrap(store)

        // Start the Helium telemetry pipeline to stream beacons dynamically
        heliumMeshBridge.start()

        repositoryScope.launch {
            client.rewardEvents().collect { reward ->
                _tritBalance.update { it + reward.amountTrits }
            }
        }
    }

    override suspend fun addTrits(amount: Long) {
        if (amount > 0) {
            client.triggerReward(RewardType.MATCH_GAME, amount)
        }
    }

    override suspend fun incrementGemScore(points: Int) {
        _gemScore.update { it + points }
        addTrits(points.toLong())
    }

    override fun setScannerState(state: GeologyScannerState) {
        _scannerState.value = state
    }

    override fun recordMineralScan(scan: MineralScan) {
        _identifiedScans.update { (it + scan).distinctBy { s -> s.name } }
    }

    override suspend fun recordSignal(signal: SignalVertex) {
        client.reportSignal(signal)

        // R3: Process signal inside 5-D CRDT Graph Engine to dynamically calculate connection strength
        val currentAddr = client.currentAddr5D()
        val signalAddr = signal.addr5d

        // Normalize RSSI using extracted normalizer in HeliumSignalAdapter
        val normalizedStrength = HeliumSignalAdapter.normalizeRssi(signal.strength)

        val edge = EdgeRecord(
            source = currentAddr,
            target = signalAddr,
            edgeType = signal.type.name,
            connectionStrength = normalizedStrength,
            edgeMetadata = mapOf(
                "signal_id" to signal.id,
                "rssi" to signal.strength.toString(),
                "observed_at" to signal.observedAtEpochSeconds.toString()
            )
        )

        // Update CRDT store
        store.putEdge(edge)

        // Update active neighbors flow for UI binding
        val updatedNeighbors = graphEngine.queryNeighbors(currentAddr).map { neighbor ->
            NeighborView(
                addr = neighbor.addr,
                edgeType = neighbor.edge.edgeType,
                strength = neighbor.edge.connectionStrength
            )
        }
        _neighbors.value = updatedNeighbors
    }
}
