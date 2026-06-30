package com.swend.mesh.worker

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log

class BootReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action != Intent.ACTION_BOOT_COMPLETED) return

        val prefs = context.getSharedPreferences("mesh", Context.MODE_PRIVATE)
        val name = prefs.getString("device_name", null)

        // Only auto-start if onboarding was completed
        if (name != null) {
            Log.i("BootReceiver", "Starting MobileWorkerService after reboot")

            val serviceIntent = Intent(context, MobileWorkerService::class.java)
            context.startForegroundService(serviceIntent)
        } else {
            Log.i("BootReceiver", "Device not onboarded; skipping auto-start")
        }
    }
}
