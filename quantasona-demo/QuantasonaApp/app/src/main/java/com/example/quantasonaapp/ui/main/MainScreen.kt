package com.example.quantasonaapp.ui.main

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.material3.TabRowDefaults.tabIndicatorOffset
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation3.runtime.NavKey
import com.example.quantasonaapp.data.DefaultDataRepository
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScreen(
    onItemClick: (NavKey) -> Unit,
    modifier: Modifier = Modifier
) {
    val repository = remember { DefaultDataRepository() }
    val tritBalance by repository.tritBalance.collectAsStateWithLifecycle()
    val hpaGenes by repository.hpaGenes.collectAsStateWithLifecycle()
    val neighbors by repository.neighbors.collectAsStateWithLifecycle() // R3: Bind dynamic connection strengths
    val gemScore by repository.gemScore.collectAsStateWithLifecycle()
    val scannerState by repository.scannerState.collectAsStateWithLifecycle()
    val identifiedScans by repository.identifiedScans.collectAsStateWithLifecycle()
    val availableScans = repository.availableScans
    val scope = rememberCoroutineScope()

    var selectedTab by remember { mutableStateOf(0) }
    val tabs = listOf("HPA Atlas", "Gem Match", "Geology", "Node HUD")

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(Color(0xFF0F172A))
    ) {
        // App Bar / Top Header
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 12.dp, horizontal = 16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = "QUANTASONA",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold,
                color = Color.White
            )
            // Trit Balance Display
            Card(
                colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B)),
                shape = RoundedCornerShape(8.dp)
            ) {
                Row(
                    modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(text = "Protein Trits: ", fontSize = 12.sp, color = Color(0xFF94A3B8))
                    Text(text = "$tritBalance", fontSize = 14.sp, fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                }
            }
        }

        TabRow(
            selectedTabIndex = selectedTab,
            containerColor = Color(0xFF1E293B),
            contentColor = Color.White,
            indicator = { tabPositions ->
                TabRowDefaults.SecondaryIndicator(
                    modifier = Modifier.tabIndicatorOffset(tabPositions[selectedTab]),
                    color = Color(0xFF38BDF8)
                )
            }
        ) {
            tabs.forEachIndexed { index, title ->
                Tab(
                    selected = selectedTab == index,
                    onClick = { selectedTab = index },
                    text = { Text(title, fontWeight = FontWeight.Medium) }
                )
            }
        }

        Box(
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            when (selectedTab) {
                0 -> HpaAtlasScreen(genes = hpaGenes)
                1 -> GemMatchScreen(
                    score = gemScore,
                    onMatchTriggered = { amount ->
                        scope.launch {
                            repository.addTrits(amount)
                            repository.incrementGemScore(amount.toInt())
                        }
                    }
                )
                2 -> GeologyScannerScreen(
                    scans = availableScans,
                    identifiedScans = identifiedScans,
                    scannerState = scannerState,
                    onScanFinished = { scan ->
                        repository.recordMineralScan(scan)
                    }
                )
                3 -> HudTelemetryScreen(neighbors = neighbors) // R2: Load HudTelemetryScreen with dynamic neighbor nodes
            }
        }
    }
}
