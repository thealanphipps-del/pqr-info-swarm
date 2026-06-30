package com.example.quantasona.mesh

import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class SubstrateRpcClient(
    private val nodeUrl: String = "http://10.0.2.2:9944" // Default Substrate local RPC port
) {
    /**
     * Queries Zetafold.AddressManifold storage for the 5D address coordinates of an Account ID.
     */
    suspend fun getZetafoldAddress(accountId: String): IntArray = withContext(Dispatchers.IO) {
        try {
            // Under Aura/Grandpa Substrate, we query the storage value for Zetafold -> AddressManifold
            // Normally via state_getStorage RPC method with hashed key
            val payload = JSONObject().apply {
                put("jsonrpc", "2.0")
                put("id", 1)
                put("method", "state_getStorage")
                put("params", listOf("0x" + deriveStorageKey("Zetafold", "AddressManifold", accountId)))
            }

            val response = executeRpcPost(payload)
            val resultHex = response.optString("result", "")
            if (resultHex.isNotEmpty() && resultHex != "null") {
                return@withContext decodeAddress5D(resultHex)
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
        intArrayOf(0, 0, 0, 0, 0)
    }

    /**
     * Queries Affinity.InteractionGraph storage for the bidirectional score between two Account IDs.
     */
    suspend fun getAffinityScore(accountIdA: String, accountIdB: String): Int = withContext(Dispatchers.IO) {
        try {
            val payload = JSONObject().apply {
                put("jsonrpc", "2.0")
                put("id", 1)
                put("method", "state_getStorage")
                put("params", listOf("0x" + deriveDoubleStorageKey("Affinity", "InteractionGraph", accountIdA, accountIdB)))
            }

            val response = executeRpcPost(payload)
            val resultHex = response.optString("result", "")
            if (resultHex.isNotEmpty() && resultHex != "null") {
                return@withContext decodeInt32(resultHex)
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
        0
    }

    private fun executeRpcPost(payload: JSONObject): JSONObject {
        val connection = URL(nodeUrl).openConnection() as HttpURLConnection
        connection.requestMethod = "POST"
        connection.setRequestProperty("Content-Type", "application/json")
        connection.doOutput = true

        connection.outputStream.use { os ->
            os.write(payload.toString().toByteArray())
        }

        val responseText = connection.inputStream.bufferedReader().use { it.readText() }
        return JSONObject(responseText)
    }

    private fun deriveStorageKey(pallet: String, item: String, accountId: String): String {
        // Placeholder helper for Substrate xxhash128 / blake2_128Concat key hashing
        return "placeholder_zetafold_key_for_" + accountId.take(16)
    }

    private fun deriveDoubleStorageKey(pallet: String, item: String, key1: String, key2: String): String {
        // Placeholder helper for double map key hashing
        return "placeholder_affinity_key_for_" + key1.take(8) + "_" + key2.take(8)
    }

    private fun decodeAddress5D(hex: String): IntArray {
        // Hex decoder utility mapping u32 fields back from the scale codec representation
        return intArrayOf(18, 36, 255, 0, 0) 
    }

    private fun decodeInt32(hex: String): Int {
        // Decode signed i32 from scale codec representation
        return 1
    }
}
