package com.swend.mesh.protocol

import org.json.JSONObject

object TaskProtocol {

    // -----------------------------
    // Incoming Task (SWEN → Android)
    // -----------------------------
    fun parseTask(json: String): Task {
        val obj = JSONObject(json)

        return Task(
            id = obj.getString("id"),
            type = obj.getString("type"),
            input = obj.opt("input"),
            timestamp = obj.optLong("timestamp", System.currentTimeMillis())
        )
    }

    data class Task(
        val id: String,
        val type: String,
        val input: Any?,
        val timestamp: Long
    )

    // -----------------------------
    // Outgoing Result (Android → SWEN)
    // -----------------------------
    fun buildResult(
        id: String,
        result: Any,
        status: String = "ok"
    ): String {
        val obj = JSONObject()

        obj.put("id", id)
        obj.put("device", "android")
        obj.put("status", status)
        obj.put("result", result)

        return obj.toString()
    }

    // -----------------------------
    // Registration Payload
    // -----------------------------
    fun buildRegistration(
        name: String,
        model: String,
        manufacturer: String
    ): String {
        val caps = JSONObject().apply {
            put("local_llm", true)
            put("local_embeddings", true)
            put("snapdragon_npu", true)
        }

        val obj = JSONObject().apply {
            put("type", "register")
            put("device", "android")
            put("name", name)
            put("model", model)
            put("manufacturer", manufacturer)
            put("capabilities", caps)
        }

        return obj.toString()
    }

    // -----------------------------
    // Telemetry Packet
    // -----------------------------
    fun buildTelemetry(
        batteryPercent: Int,
        batteryTempC: Float,
        cpuLoad: Float,
        signalStrength: Int,
        model: String,
        manufacturer: String
    ): String {
        val caps = JSONObject().apply {
            put("local_llm", true)
            put("local_embeddings", true)
            put("snapdragon_npu", true)
        }

        val obj = JSONObject().apply {
            put("type", "telemetry")
            put("device", "android")
            put("model", model)
            put("manufacturer", manufacturer)

            put("battery_percent", batteryPercent)
            put("battery_temp_c", batteryTempC)
            put("cpu_load", cpuLoad)
            put("signal_strength", signalStrength)

            put("capabilities", caps)
        }

        return obj.toString()
    }
}
