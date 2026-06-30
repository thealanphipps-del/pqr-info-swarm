package com.swend.mesh.ui

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.os.Build
import androidx.core.app.NotificationCompat
import com.swend.mesh.R

object MeshNotifications {

    private const val CHANNEL_ID = "mesh_worker"
    private const val CHANNEL_NAME = "Sovereign Mesh Worker"

    fun createForegroundNotification(context: Context): Notification {
        createChannel(context)

        return NotificationCompat.Builder(context, CHANNEL_ID)
            .setContentTitle("Sovereign Mesh Worker")
            .setContentText("Running background inference tasks")
            .setSmallIcon(R.drawable.ic_mesh_worker) // Antigravity will supply this
            .setOngoing(true)
            .setPriority(NotificationCompat.PRIORITY_LOW)
            .build()
    }

    private fun createChannel(context: Context) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                CHANNEL_NAME,
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "Keeps the Sovereign Mesh worker running in the background"
            }

            val manager = context.getSystemService(Context.NOTIFICATION_SERVICE)
                    as NotificationManager

            manager.createNotificationChannel(channel)
        }
    }
}
