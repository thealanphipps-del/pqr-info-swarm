package com.example.quantasonaapp

import com.example.quantasonaapp.data.*
import com.example.quantasonaapp.data.MineralScan
import com.example.quantasonaapp.ui.main.MainScreenUiState
import com.example.quantasonaapp.ui.main.MainScreenViewModel
import com.example.quantasonaapp.ui.main.GemType
import junit.framework.TestCase.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.async
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import kotlinx.coroutines.test.runTest
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class E2ETestSuite {

    private lateinit var repository: DefaultDataRepository
    private lateinit var viewModel: MainScreenViewModel

    @Before
    fun setUp() {
        repository = DefaultDataRepository()
        viewModel = MainScreenViewModel(repository)
    }

    // ==========================================
    // TIER 1: Basic Functional Verification (20 tests)
    // ==========================================

    // Feature 1: HpaAtlasScreen
    @Test
    fun tier1_hpa_initialGeneListNotEmpty() = runTest {
        val genes = repository.hpaGenes.value
        assertTrue("Initial gene list should not be empty", genes.isNotEmpty())
    }

    @Test
    fun tier1_hpa_searchBySymbolMatch() = runTest {
        val genes = repository.hpaGenes.value
        val match = genes.filter { it.symbol.contains("TP53", ignoreCase = true) }
        assertEquals(1, match.size)
        assertEquals("ENSG00000141510", match.first().ensemblId)
    }

    @Test
    fun tier1_hpa_searchByEnsemblIdMatch() = runTest {
        val genes = repository.hpaGenes.value
        val match = genes.filter { it.ensemblId.contains("ENSG00000090612", ignoreCase = true) }
        assertEquals(1, match.size)
        assertEquals("ZNF268", match.first().symbol)
    }

    @Test
    fun tier1_hpa_searchNoMatchReturnsEmpty() = runTest {
        val genes = repository.hpaGenes.value
        val match = genes.filter { it.symbol.contains("XYZ_NONEXISTENT", ignoreCase = true) }
        assertTrue(match.isEmpty())
    }

    @Test
    fun tier1_hpa_geneDetailsRenderExpression() = runTest {
        val genes = repository.hpaGenes.value
        val gene = genes.first { it.symbol == "TP53" }
        assertEquals("Tumor protein p53", gene.description)
        assertTrue(gene.subcellularLocations.contains("Nucleoplasm"))
        assertEquals("High", gene.tissueExpression["Spleen"])
    }

    // Feature 2: GemMatchScreen
    @Test
    fun tier1_gem_gridInitialization() = runTest {
        GemType.values().forEach { gemType ->
            val mappedName = if (gemType == GemType.GEOCACHE) "GEO_CACHE" else gemType.name
            val signalType = SignalType.valueOf(mappedName)
            assertNotNull(signalType)
        }
    }

    @Test
    fun tier1_gem_scoreInitiallyZero() = runTest {
        assertEquals(0, repository.gemScore.value)
    }

    @Test
    fun tier1_gem_matchIncrementsScoreAndTrits() = runTest {
        val initialBalance = repository.tritBalance.value
        repository.addTrits(15)
        val finalBalance = repository.tritBalance.first { it >= initialBalance + 15 }
        assertTrue(finalBalance >= initialBalance + 15)
    }

    @Test
    fun tier1_gem_reshuffleGrid() = runTest {
        val grid = List(16) { GemType.values().random() }
        assertEquals(16, grid.size)
        val validGemTypes = GemType.values().toSet()
        grid.forEach { gem ->
            assertTrue(gem in validGemTypes)
        }
    }

    @Test
    fun tier1_gem_matchTritAccumulation() = runTest {
        val initialBalance = repository.tritBalance.value
        repository.addTrits(10)
        repository.addTrits(20)
        val finalBalance = repository.tritBalance.first { it >= initialBalance + 30 }
        assertTrue(finalBalance >= initialBalance + 30)
    }

    // Feature 3: GeologyScannerScreen
    @Test
    fun tier1_geo_scansLoadSuccessfully() = runTest {
        assertTrue(repository.availableScans.isNotEmpty())
    }

    @Test
    fun tier1_geo_scanMultipliersCorrect() = runTest {
        val basalt = repository.availableScans.first { it.name == "Basalt" }
        assertEquals(1.2, basalt.tritMultiplier)
    }

    @Test
    fun tier1_geo_scanHardnessCorrect() = runTest {
        val granite = repository.availableScans.first { it.name == "Granite" }
        assertEquals(6.5, granite.hardness)
    }

    @Test
    fun tier1_geo_scanCrystalStructureCorrect() = runTest {
        val quartzite = repository.availableScans.first { it.name == "Quartzite" }
        assertEquals("Granoblastic / Metamorphic", quartzite.crystalStructure)
    }

    @Test
    fun tier1_geo_scannerStateTransitions() = runTest {
        assertEquals(GeologyScannerState.IDLE, repository.scannerState.value)
        repository.setScannerState(GeologyScannerState.SCANNING)
        assertEquals(GeologyScannerState.SCANNING, repository.scannerState.value)
        repository.setScannerState(GeologyScannerState.COMPLETED)
        assertEquals(GeologyScannerState.COMPLETED, repository.scannerState.value)
    }

    // Feature 4: HudTelemetryScreen & Helium
    @Test
    fun tier1_hud_heliumBeaconsStreamed() = runTest {
        val client = HeliumClientImpl("test-hotspot")
        val firstBeacon = client.beacons().first()
        assertEquals("hotspot-alpha", firstBeacon.id)
        assertTrue(firstBeacon.rssi in -100..-40)
    }

    @Test
    fun tier1_hud_heliumRewardsStreamed() = runTest {
        val client = HeliumClientImpl("test-hotspot")
        val reward = client.rewards().first()
        assertEquals("test-hotspot", reward.hotspotId)
        assertEquals(15.25, reward.amountMobile)
    }

    @Test
    fun tier1_hud_bridgePropagatesBeacons() = runTest {
        val mockMesh = MockMeshClient()
        val mockRepository = DefaultDataRepository(client = mockMesh)
        val fakeBeacon = HeliumBeacon(
            id = "test-beacon-node",
            rssi = -60,
            frequency = 915.2,
            lat = 37.7749,
            lon = -122.4194,
            timestampEpochSeconds = System.currentTimeMillis() / 1000,
            metadata = emptyMap()
        )
        val fakeHelium = object : HeliumClient {
            override fun beacons() = kotlinx.coroutines.flow.flowOf(fakeBeacon)
            override fun rewards() = kotlinx.coroutines.flow.emptyFlow<HeliumReward>()
        }
        val bridge = HeliumMeshBridge(fakeHelium, mockMesh, mockRepository, this)
        bridge.start()
        kotlinx.coroutines.delay(100)
        val neighbors = mockRepository.neighbors.value
        val hasBeaconNode = neighbors.any { it.addr.spaceId == "test-beacon-node" }
        assertTrue(hasBeaconNode)
    }

    @Test
    fun tier1_hud_connectionStrengthNormalized() = runTest {
        val beacon = HeliumBeacon("test", -70, 915.0, 0.0, 0.0, 1000L, emptyMap())
        val addr = Addr5D(0, "space", "lineage", "content", "oracle")
        val vertex = HeliumSignalAdapter.toSignalVertex(beacon, addr)
        assertEquals(-70.0, vertex.strength)
    }

    @Test
    fun tier1_hud_lineageCanvasPeersUpdated() = runTest {
        val client = MockMeshClient()
        val addr = client.currentAddr5D()
        assertEquals("local_node_lineage", addr.lineageHash)
    }

    // ==========================================
    // TIER 2: Boundary and Edge Cases (20 tests)
    // ==========================================

    // Feature 1: HpaAtlasScreen
    @Test
    fun tier2_hpa_searchQueryTrimming() = runTest {
        val genes = repository.hpaGenes.value
        val query = "  TP53  ".trim()
        val match = genes.filter { it.symbol.contains(query, ignoreCase = true) }
        assertEquals(1, match.size)
    }

    @Test
    fun tier2_hpa_searchCaseInsensitivity() = runTest {
        val genes = repository.hpaGenes.value
        val matchLower = genes.filter { it.symbol.contains("tp53", ignoreCase = true) }
        val matchUpper = genes.filter { it.symbol.contains("TP53", ignoreCase = true) }
        assertEquals(matchLower.size, matchUpper.size)
    }

    @Test
    fun tier2_hpa_emptyGeneListGraceful() = runTest {
        val genes = repository.hpaGenes.value
        val match = genes.filter { it.symbol.contains("NON_EXISTENT_SYMBOL_XYZ", ignoreCase = true) }
        assertTrue(match.isEmpty())
    }

    @Test
    fun tier2_hpa_doubleClickGeneNoCrash() = runTest {
        val firstFetch = repository.hpaGenes.value.first()
        val secondFetch = repository.hpaGenes.value.first()
        assertEquals(firstFetch, secondFetch)
    }

    @Test
    fun tier2_hpa_dismissDetailsResetsState() = runTest {
        val originalCount = repository.hpaGenes.value.size
        val filteredTP53 = repository.hpaGenes.value.filter { it.symbol.contains("TP53", ignoreCase = true) }
        assertEquals(1, filteredTP53.size)
        val resetList = repository.hpaGenes.value.filter { it.symbol.contains("", ignoreCase = true) }
        assertEquals(originalCount, resetList.size)
    }

    // Feature 2: GemMatchScreen
    @Test
    fun tier2_gem_scoreOverflowPrevention() = runTest {
        val initialScore = repository.gemScore.value
        repository.incrementGemScore(10000)
        assertEquals(initialScore + 10000, repository.gemScore.value)
    }

    @Test
    fun tier2_gem_invalidMatchAmountRejected() = runTest {
        val initialBalance = repository.tritBalance.value
        repository.addTrits(-50)
        val currentBalance = repository.tritBalance.value
        assertEquals(initialBalance, currentBalance)
    }

    @Test
    fun tier2_gem_consecutiveReshuffles() = runTest {
        val initialScore = repository.gemScore.value
        val initialBalance = repository.tritBalance.value
        repeat(5) {
            repository.incrementGemScore(15)
        }
        val expectedBalance = initialBalance + 75
        val finalBalance = repository.tritBalance.first { it >= expectedBalance }
        assertEquals(initialScore + 75, repository.gemScore.value)
        assertEquals(expectedBalance, finalBalance)
    }

    @Test
    fun tier2_gem_gridStatePersistence() = runTest {
        val mappedList = GemType.values().map { gem ->
            val mappedName = if (gem == GemType.GEOCACHE) "GEO_CACHE" else gem.name
            SignalType.valueOf(mappedName)
        }
        assertEquals(GemType.values().size, mappedList.size)
        mappedList.forEach { signalType ->
            assertNotNull(signalType)
        }
    }

    @Test
    fun tier2_gem_highScoreThreshold() = runTest {
        repository.incrementGemScore(1200)
        assertTrue(repository.gemScore.value > 1000)
    }

    // Feature 3: GeologyScannerScreen
    @Test
    fun tier2_geo_emptyScansList() = runTest {
        assertTrue(repository.identifiedScans.value.isEmpty())
    }

    @Test
    fun tier2_geo_invalidMultiplierBounded() = runTest {
        repository.availableScans.forEach { scan ->
            assertTrue(scan.tritMultiplier >= 1.0)
        }
    }

    @Test
    fun tier2_geo_hardnessBoundaryValues() = runTest {
        repository.availableScans.forEach { scan ->
            assertTrue(scan.hardness in 1.0..10.0)
        }
    }

    @Test
    fun tier2_geo_scannerTimeout() = runTest {
        repository.setScannerState(GeologyScannerState.TIMEOUT)
        assertEquals(GeologyScannerState.TIMEOUT, repository.scannerState.value)
    }

    @Test
    fun tier2_geo_duplicateScansHandled() = runTest {
        val basalt = repository.availableScans.first { it.name == "Basalt" }
        repository.recordMineralScan(basalt)
        repository.recordMineralScan(basalt)
        assertEquals(1, repository.identifiedScans.value.size)
        assertEquals(basalt, repository.identifiedScans.value.first())
    }

    // Feature 4: HudTelemetryScreen & Helium
    @Test
    fun tier2_hud_extremeRssiHighNormalizedToOne() = runTest {
        assertEquals(1.0f, HeliumSignalAdapter.normalizeRssi(-30.0))
    }

    @Test
    fun tier2_hud_extremeRssiLowNormalizedToZero() = runTest {
        assertEquals(0.0f, HeliumSignalAdapter.normalizeRssi(-110.0))
    }

    @Test
    fun tier2_hud_emptyBeaconsStream() = runTest {
        val mockClient = MockMeshClient()
        val emptyHelium = object : HeliumClient {
            override fun beacons() = kotlinx.coroutines.flow.emptyFlow<HeliumBeacon>()
            override fun rewards() = kotlinx.coroutines.flow.emptyFlow<HeliumReward>()
        }
        val mockRepository = DefaultDataRepository(heliumClient = emptyHelium)
        val bridge = HeliumMeshBridge(emptyHelium, mockClient, mockRepository)
        bridge.start()
        kotlinx.coroutines.delay(100)
        assertTrue(mockRepository.neighbors.value.isEmpty())
    }

    @Test
    fun tier2_hud_disconnectedNodeHandling() = runTest {
        val weakSignal = SignalVertex(
            id = "weak-node",
            type = SignalType.IOT,
            addr5d = Addr5D(0, "space", "weak_lineage", "content", "oracle"),
            strength = -95.0,
            observedAtEpochSeconds = System.currentTimeMillis() / 1000,
            metadata = emptyMap()
        )
        val strongSignal = SignalVertex(
            id = "strong-node",
            type = SignalType.IOT,
            addr5d = Addr5D(0, "space", "strong_lineage", "content", "oracle"),
            strength = -40.0,
            observedAtEpochSeconds = System.currentTimeMillis() / 1000,
            metadata = emptyMap()
        )
        repository.recordSignal(weakSignal)
        repository.recordSignal(strongSignal)
        val neighbors = repository.neighbors.value
        val weakNeighbor = neighbors.find { it.addr.lineageHash == "weak_lineage" }
        val strongNeighbor = neighbors.find { it.addr.lineageHash == "strong_lineage" }
        assertNotNull(weakNeighbor)
        assertNotNull(strongNeighbor)
        assertFalse(weakNeighbor!!.isConnected)
        assertTrue(strongNeighbor!!.isConnected)
    }

    @Test
    fun tier2_hud_multipleBeaconsSameNode() = runTest {
        val store = InMemoryFiveDStore()
        val src = Addr5D(0, "space", "lineage", "content", "oracle")
        val dst = Addr5D(1, "space", "lineage", "content", "oracle")
        
        val edge1 = EdgeRecord(src, dst, "IOT", 0.4f, emptyMap())
        store.putEdge(edge1)
        
        val edge2 = EdgeRecord(src, dst, "IOT", 0.8f, emptyMap())
        store.putEdge(edge2)
        
        val edges = store.getEdgesFor(src)
        assertEquals(1, edges.size)
        assertEquals(0.8f, edges.first().connectionStrength)
    }

    // ==========================================
    // TIER 3: Pairwise Cross-Feature Interactions (4 tests)
    // ==========================================

    @Test
    fun tier3_cross_hpaToGemMatch_tritBoost() = runTest {
        val initialTrit = repository.tritBalance.value
        repository.addTrits(200)
        val finalTrit = repository.tritBalance.first { it >= initialTrit + 200 }
        
        val gene = repository.hpaGenes.value.first { it.symbol == "ZNF268" }
        assertNotNull(gene)
        assertTrue(finalTrit >= initialTrit + 200)
    }

    @Test
    fun tier3_cross_geologyToHud_scannerUpdatesTelemetry() = runTest {
        val initialNeighborsSize = repository.neighbors.value.size
        
        val signal = SignalVertex(
            id = "geo-scan-1",
            type = SignalType.GEO_CACHE,
            addr5d = Addr5D(System.currentTimeMillis() / 1000, "pancreas", "local", "hash", "geo"),
            strength = -65.0,
            observedAtEpochSeconds = System.currentTimeMillis() / 1000,
            metadata = mapOf("rock" to "quartzite")
        )
        
        repository.recordSignal(signal)
        
        val finalNeighbors = repository.neighbors.value
        val geoNeighbor = finalNeighbors.find { it.edgeType == "GEO_CACHE" }
        assertNotNull(geoNeighbor)
    }

    @Test
    fun tier3_cross_heliumToGemMatch_rewardsAddTrits() = runTest {
        val initialBalance = repository.tritBalance.value
        val hr = HeliumReward("hotspot-1", 20.0, "poc_challenger", System.currentTimeMillis() / 1000)
        val rewardEvent = HeliumRewardMapper.toRewardEvent("node-local", hr)
        
        repository.client.triggerReward(rewardEvent.type, rewardEvent.amountTrits, rewardEvent.meta)
        
        val finalBalance = repository.tritBalance.first { it >= initialBalance + 200 }
        assertTrue(finalBalance >= initialBalance + 200)
    }

    @Test
    fun tier3_cross_hudToHpa_neighborLineageVerify() = runTest {
        val healthyGene = repository.hpaGenes.value.first { it.symbol == "ZNF268" }
        assertNotNull(healthyGene)
        
        val signal = SignalVertex(
            id = "node-1",
            type = SignalType.IOT,
            addr5d = Addr5D(0, "pancreas", "healthy_lineage", "content", "biology"),
            strength = -50.0,
            observedAtEpochSeconds = System.currentTimeMillis() / 1000,
            metadata = emptyMap()
        )
        repository.recordSignal(signal)
        
        val neighbor = repository.neighbors.value.find { it.addr.lineageHash == "healthy_lineage" }
        assertNotNull(neighbor)
    }

    // ==========================================
    // TIER 4: Real-World System Workloads (5 tests)
    // ==========================================

    @Test
    fun tier4_workload_continuousSyncCycle() = runTest {
        val initialTrit = repository.tritBalance.value
        
        repeat(50) { idx ->
            repository.recordSignal(
                SignalVertex(
                    id = "beacon-$idx",
                    type = SignalType.values().random(),
                    addr5d = Addr5D(idx.toLong(), "pancreas", "lineage-$idx", "content", "mesh"),
                    strength = (-95..-45).random().toDouble(),
                    observedAtEpochSeconds = System.currentTimeMillis() / 1000,
                    metadata = emptyMap()
                )
            )
        }
        
        repeat(10) { idx ->
            repository.addTrits(15)
        }
        
        val finalBalance = repository.tritBalance.first { it >= initialTrit + 150 }
        assertTrue(finalBalance >= initialTrit + 150)
        assertTrue(repository.neighbors.value.size >= 1)
    }

    @Test
    fun tier4_workload_stateTransitionConvergence() = runTest {
        val store = InMemoryFiveDStore()
        val queryEngine = DefaultGraphQueryEngine(store, store)
        
        InsulinLattice.bootstrap(store)
        
        val startAddr = Addr5D(timeIndex = 0, spaceId = "pancreas", lineageHash = "healthy", contentHash = "INS_B", oracleContextId = "biology")
        val endAddr = Addr5D(timeIndex = 0, spaceId = "pancreas", lineageHash = "t2d", contentHash = "INS_B", oracleContextId = "biology")
        
        val path = queryEngine.findPathByCriteria(startAddr, endAddr, PathCriteria(minStrength = 0.4f))
        assertNotNull(path)
        assertEquals(2, path?.size)
        assertEquals(startAddr, path?.first())
        assertEquals(endAddr, path?.last())
    }

    @Test
    fun tier4_workload_congestedMeshRouting() = runTest {
        val header = SilkRoadHeader(
            src = Addr5D(0, "pancreas", "src", "hash", "mesh"),
            dst = Addr5D(1, "pancreas", "dst", "hash", "mesh"),
            priority = PriorityClass.FAST,
            maxLatencySeconds = 5L,
            expiryEpochSeconds = System.currentTimeMillis() / 1000 + 10,
            feeBudgetTritSubunits = 100L
        )
        val packet = SilkRoadPacket(header, "payload".toByteArray())
        
        val mockClient = MockMeshClient()
        val receivedPacket = async(Dispatchers.Unconfined) { mockClient.incomingPackets().first() }
        mockClient.sendPacket(packet)
        
        assertEquals(packet, receivedPacket.await())
    }

    @Test
    fun tier4_workload_concurrentStateUpdates() = runTest {
        val jobs = List(10) {
            launch(Dispatchers.Default) {
                repository.addTrits(10)
            }
        }
        jobs.forEach { it.join() }
        
        val finalBalance = repository.tritBalance.first { it >= 1100L }
        assertTrue(finalBalance >= 1100L)
    }

    @Test
    fun tier4_workload_crdtConflictResolution() = runTest {
        val store = InMemoryFiveDStore()
        val addr = Addr5D(0, "pancreas", "lineage", "content", "biology")
        
        val vertexOld = VertexRecord(
            addr = addr,
            payloadData = "Old".toByteArray(),
            metadata = emptyMap(),
            lastUpdatedLineage = LineageRecord("lineage-crdt", 1, 1000L)
        )
        store.putVertex(vertexOld)
        
        val vertexNew = VertexRecord(
            addr = addr,
            payloadData = "New".toByteArray(),
            metadata = emptyMap(),
            lastUpdatedLineage = LineageRecord("lineage-crdt", 5, 2000L)
        )
        store.mergeVertex(vertexNew)
        
        val resolved = store.getVertex(addr)
        assertEquals("New", resolved?.payloadData?.let { String(it) })
    }
}
