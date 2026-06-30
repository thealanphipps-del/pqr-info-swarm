package com.example.quantasonaapp.ui.main

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.example.quantasonaapp.data.HpaGene

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HpaAtlasScreen(
    genes: List<HpaGene>,
    modifier: Modifier = Modifier
) {
    var searchQuery by remember { mutableStateOf("") }
    var selectedGene by remember { mutableStateOf<HpaGene?>(null) }

    val filteredGenes = genes.filter {
        it.symbol.contains(searchQuery, ignoreCase = true) ||
                it.ensemblId.contains(searchQuery, ignoreCase = true)
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(Color(0xFF0F172A)) // Slate 900
    ) {
        Text(
            text = "Human Protein Atlas Baselines",
            fontSize = 22.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFF38BDF8), // Sky 400
            modifier = Modifier.padding(bottom = 8.dp)
        )

        TextField(
            value = searchQuery,
            onValueChange = { searchQuery = it },
            placeholder = { Text("Search by Gene Symbol or Ensembl ID", color = Color.Gray) },
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 8.dp),
            colors = TextFieldDefaults.colors(
                focusedContainerColor = Color(0xFF1E293B),
                unfocusedContainerColor = Color(0xFF1E293B),
                focusedTextColor = Color.White,
                unfocusedTextColor = Color.White,
                cursorColor = Color(0xFF38BDF8)
            ),
            shape = RoundedCornerShape(8.dp)
        )

        if (selectedGene != null) {
            GeneDetailsCard(
                gene = selectedGene!!,
                onDismiss = { selectedGene = null }
            )
        } else {
            LazyColumn(
                verticalArrangement = Arrangement.spacedBy(8.dp),
                modifier = Modifier.weight(1f)
            ) {
                items(filteredGenes) { gene ->
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable { selectedGene = gene },
                        colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B)),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Row(
                            modifier = Modifier.padding(16.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Column(modifier = Modifier.weight(1f)) {
                                Text(
                                    text = gene.symbol,
                                    fontSize = 18.sp,
                                    fontWeight = FontWeight.Bold,
                                    color = Color.White
                                )
                                Text(
                                    text = gene.ensemblId,
                                    fontSize = 12.sp,
                                    color = Color(0xFF94A3B8)
                                )
                            }
                            Text(
                                text = "View Details",
                                color = Color(0xFF10B981), // Green
                                fontSize = 14.sp,
                                fontWeight = FontWeight.Bold
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun GeneDetailsCard(gene: HpaGene, onDismiss: () -> Unit) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 8.dp),
        colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B)),
        shape = RoundedCornerShape(12.dp)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Row(
                horizontalArrangement = Arrangement.SpaceBetween,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    text = gene.symbol,
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.White
                )
                Text(
                    text = "Close",
                    color = Color(0xFFEF4444),
                    fontWeight = FontWeight.Bold,
                    modifier = Modifier.clickable { onDismiss() }
                )
            }
            Text(
                text = "Ensembl ID: ${gene.ensemblId}",
                fontSize = 12.sp,
                color = Color(0xFF94A3B8),
                modifier = Modifier.padding(bottom = 8.dp)
            )
            Text(
                text = gene.description,
                fontSize = 14.sp,
                color = Color.White,
                modifier = Modifier.padding(bottom = 12.dp)
            )

            HorizontalDivider(color = Color(0xFF334155))

            Text(
                text = "Subcellular Locations",
                fontSize = 15.sp,
                fontWeight = FontWeight.SemiBold,
                color = Color(0xFF38BDF8),
                modifier = Modifier.padding(vertical = 8.dp)
            )
            gene.subcellularLocations.forEach { loc ->
                Text(text = "• $loc", color = Color.White, fontSize = 14.sp)
            }

            HorizontalDivider(color = Color(0xFF334155), modifier = Modifier.padding(vertical = 8.dp))

            Text(
                text = "Tissue Expression Levels (IHC)",
                fontSize = 15.sp,
                fontWeight = FontWeight.SemiBold,
                color = Color(0xFF38BDF8),
                modifier = Modifier.padding(bottom = 8.dp)
            )
            gene.tissueExpression.forEach { (tissue, level) ->
                Row(
                    horizontalArrangement = Arrangement.SpaceBetween,
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 2.dp)
                ) {
                    Text(text = tissue, color = Color.White, fontSize = 14.sp)
                    val badgeColor = when (level.lowercase()) {
                        "high" -> Color(0xFF10B981)
                        "medium" -> Color(0xFFF59E0B)
                        else -> Color(0xFF64748B)
                    }
                    Text(text = level, color = badgeColor, fontWeight = FontWeight.Bold, fontSize = 14.sp)
                }
            }
        }
    }
}
