package com.example.quantasona.mesh

class AffinityClient(private val substrate: SubstrateRpcClient) {

    suspend fun getAffinity(a: String, b: String): Int {
        return substrate.getAffinityScore(a, b)
    }
}
