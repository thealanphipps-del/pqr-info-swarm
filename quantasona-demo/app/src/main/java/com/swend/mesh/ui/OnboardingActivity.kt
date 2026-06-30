package com.swend.mesh.ui

import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.widget.Button
import android.widget.EditText
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.swend.mesh.net.MeshClient
import com.swend.mesh.worker.MobileWorkerService
import org.json.JSONObject

class OnboardingActivity : AppCompatActivity() {

    private lateinit var nameInput: EditText
    private lateinit var joinButton: Button

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContentView(R.layout.activity_onboarding)

        nameInput = findViewById(R.id.deviceNameInput)
        joinButton = findViewById(R.id.joinButton)

        joinButton.setOnClickListener {
            val name = nameInput.text.toString().trim()

            if (name.isEmpty()) {
                Toast.makeText(this, "Enter a name", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }

            saveDeviceName(name)
            registerWithMesh(name)
        }
    }

    private fun saveDeviceName(name: String) {
        val prefs = getSharedPreferences("mesh", Context.MODE_PRIVATE)
        prefs.edit().putString("device_name", name).apply()
    }

    private fun registerWithMesh(name: String) {
        // Build registration payload
        val registration = JSONObject().apply {
            put("type", "register")
            put("device", "android")
            put("name", name)
            put("model", android.os.Build.MODEL)
            put("manufacturer", android.os.Build.MANUFACTURER)

            put("capabilities", JSONObject().apply {
                put("local_llm", true)
                put("local_embeddings", true)
                put("snapdragon_npu", true)
            })
        }.toString()

        // Fire-and-forget registration via MeshClient
        val client = MeshClient(
            context = this,
            onTaskReceived = {}
        )

        client.connect()
        client.sendResult(registration)

        // Start worker service
        startService(Intent(this, MobileWorkerService::class.java))

        // Move to connected screen
        startActivity(Intent(this, ConnectedActivity::class.java))
        finish()
    }
}
