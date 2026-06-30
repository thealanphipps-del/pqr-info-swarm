package com.example.quantasona.audio

import android.media.AudioFormat
import android.media.AudioRecord
import android.media.MediaRecorder
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.channels.BufferOverflow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.launch
import kotlin.math.cos
import kotlin.math.PI

class AudioProcessor {

    private val sampleRate = 16_000
    private val frameSize = 1024

    private var job: Job? = null

    private val _frames =
        MutableSharedFlow<FloatArray>(replay = 0, extraBufferCapacity = 8, onBufferOverflow = BufferOverflow.DROP_OLDEST)
    val frames: SharedFlow<FloatArray> = _frames

    private val recorder: AudioRecord by lazy {
        val minBuf = AudioRecord.getMinBufferSize(
            sampleRate,
            AudioFormat.CHANNEL_IN_MONO,
            AudioFormat.ENCODING_PCM_16BIT
        )
        AudioRecord(
            MediaRecorder.AudioSource.MIC,
            sampleRate,
            AudioFormat.CHANNEL_IN_MONO,
            AudioFormat.ENCODING_PCM_16BIT,
            minBuf
        )
    }

    fun start(scope: CoroutineScope) {
        if (job != null) return
        recorder.startRecording()
        job = scope.launch(Dispatchers.Default) {
            val buffer = ShortArray(frameSize)
            while (true) {
                val read = recorder.read(buffer, 0, frameSize)
                if (read > 0) {
                    val preprocessed = preprocess(buffer, read)
                    _frames.emit(preprocessed)
                }
            }
        }
    }

    fun stop() {
        job?.cancel()
        job = null
        recorder.stop()
        recorder.release()
    }

    private fun preprocess(data: ShortArray, len: Int): FloatArray {
        val out = FloatArray(len)

        // int16 -> float [-1,1]
        for (i in 0 until len) {
            out[i] = data[i] / 32768f
        }

        // pre-emphasis
        for (i in len - 1 downTo 1) {
            out[i] = out[i] - 0.97f * out[i - 1]
        }

        // Hamming window
        for (i in 0 until len) {
            val w = (0.54 - 0.46 * cos(2.0 * PI * i / (len - 1))).toFloat()
            out[i] *= w
        }

        return out
    }
}
