package com.example.quantasonaapp.ui.main

import com.example.quantasonaapp.data.DataRepository
import com.example.quantasonaapp.data.HpaGene
import com.example.quantasonaapp.data.MockMeshClient
import com.example.quantasonaapp.data.SignalVertex
import com.example.quantasonaapp.data.NeighborView
import com.example.quantasonaapp.data.GeologyScannerState
import com.example.quantasonaapp.data.MineralScan
import junit.framework.TestCase.assertEquals
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.test.runTest
import org.junit.Test

class MainScreenViewModelTest {
  @Test
  fun uiState_initiallyLoading() = runTest {
    val viewModel = MainScreenViewModel(FakeMyModelRepository())
    assertEquals(viewModel.uiState.value, MainScreenUiState.Loading)
  }

  @Test
  fun uiState_loadsSuccess() = runTest {
    val viewModel = MainScreenViewModel(FakeMyModelRepository())
    val successState = viewModel.uiState.first { it is MainScreenUiState.Success } as MainScreenUiState.Success
    assertEquals(listOf("Sample"), successState.data)
  }
}

private class FakeMyModelRepository : DataRepository {
  override val data: Flow<List<String>> = flow { emit(listOf("Sample")) }
  override val tritBalance: StateFlow<Long> = MutableStateFlow(100L)
  override val hpaGenes: StateFlow<List<HpaGene>> = MutableStateFlow(emptyList())
  override val client = MockMeshClient()
  override val neighbors: StateFlow<List<NeighborView>> = MutableStateFlow(emptyList()) // Added override

  override val gemScore: StateFlow<Int> = MutableStateFlow(0)
  override val scannerState: StateFlow<GeologyScannerState> = MutableStateFlow(GeologyScannerState.IDLE)
  override val availableScans: List<MineralScan> = emptyList()
  override val identifiedScans: StateFlow<List<MineralScan>> = MutableStateFlow(emptyList())

  override suspend fun addTrits(amount: Long) {}
  override suspend fun recordSignal(signal: SignalVertex) {}
  override suspend fun incrementGemScore(points: Int) {}
  override fun setScannerState(state: GeologyScannerState) {}
  override fun recordMineralScan(scan: MineralScan) {}
}
