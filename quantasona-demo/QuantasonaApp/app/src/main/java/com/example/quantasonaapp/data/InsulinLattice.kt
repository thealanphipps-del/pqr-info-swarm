package com.example.quantasonaapp.data

object InsulinLattice {
    // Axis A: Proteins
    const val INS = 0L
    const val IAPP = 1L
    const val PCSK1 = 2L
    const val CPE = 3L
    const val INSR = 4L

    // Axis B: Domains
    const val INS_A_CHAIN = 0L
    const val INS_B_CHAIN = 1L
    const val INS_C_PEPTIDE = 2L
    const val INS_BINDING_SURFACE = 3L

    const val PCSK1_CATALYTIC = 0L
    const val PCSK1_P_DOMAIN = 1L

    const val INSR_LIGAND_BINDING = 0L
    const val INSR_KINASE = 1L

    // Axis D: Tissues
    const val TISSUE_B_CELL = 0L
    const val TISSUE_A_CELL = 1L
    const val TISSUE_HEPATOCYTE = 2L

    // Axis E: States
    const val STATE_HEALTHY = 0L
    const val STATE_T1D = 1L
    const val STATE_T2D = 2L

    fun bootstrap(store: InMemoryFiveDStore) {
        val scope = kotlinx.coroutines.CoroutineScope(kotlinx.coroutines.Dispatchers.Default)
        kotlinx.coroutines.runBlocking {
            // 1. Bootstrap Healthy State (e = 0)
            val vHealthyIns = VertexRecord(
                addr = Addr5D(timeIndex = 0, spaceId = "pancreas", lineageHash = "healthy", contentHash = "INS_B", oracleContextId = "biology"),
                payloadData = "Insulin:B-chain:PCSK1".toByteArray(),
                metadata = mapOf("protein" to "INS", "domain" to "B-chain", "state" to "Healthy", "level" to "High"),
                lastUpdatedLineage = LineageRecord("healthy_lineage", 1, System.currentTimeMillis() / 1000)
            )
            store.putVertex(vHealthyIns)

            // 2. Bootstrap T1D State (e = 1) - Deletion of beta cell vertices (low weight / missing)
            val vT1dIns = VertexRecord(
                addr = Addr5D(timeIndex = 0, spaceId = "pancreas", lineageHash = "t1d", contentHash = "INS_B", oracleContextId = "biology"),
                payloadData = "Insulin:B-chain:Absent".toByteArray(),
                metadata = mapOf("protein" to "INS", "domain" to "B-chain", "state" to "T1D", "level" to "Absent"),
                lastUpdatedLineage = LineageRecord("t1d_lineage", 1, System.currentTimeMillis() / 1000)
            )
            store.putVertex(vT1dIns)

            // 3. Bootstrap T2D State (e = 2) - Warped weights (insulin resistance)
            val vT2dIns = VertexRecord(
                addr = Addr5D(timeIndex = 0, spaceId = "pancreas", lineageHash = "t2d", contentHash = "INS_B", oracleContextId = "biology"),
                payloadData = "Insulin:B-chain:Stressed".toByteArray(),
                metadata = mapOf("protein" to "INS", "domain" to "B-chain", "state" to "T2D", "level" to "Stressed"),
                lastUpdatedLineage = LineageRecord("t2d_lineage", 1, System.currentTimeMillis() / 1000)
            )
            store.putVertex(vT2dIns)

            // Add Edges representing signaling transitions
            store.putEdge(
                EdgeRecord(
                    source = vHealthyIns.addr,
                    target = vT2dIns.addr,
                    edgeType = "state_transition",
                    connectionStrength = 0.5f,
                    edgeMetadata = mapOf("transition" to "Healthy to T2D", "description" to "Deformation of insulin signaling cascade")
                )
            )
        }
    }
}
