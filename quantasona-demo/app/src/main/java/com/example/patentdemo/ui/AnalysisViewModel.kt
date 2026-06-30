package com.example.patentdemo.ui

import androidx.lifecycle.ViewModel
import com.example.patentdemo.dsp.Band
import com.example.patentdemo.profiles.ProfileMatchResult

class AnalysisViewModel : ViewModel() {
    var pcmData: ByteArray? = null
    var fftData: FloatArray? = null
    var zeroSlopeBands: List<Band>? = null
    var profileMatch: ProfileMatchResult? = null
}
