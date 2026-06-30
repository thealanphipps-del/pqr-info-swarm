package com.example.quantasonaapp.ui.main

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

enum class GemType(val color: Color, val symbol: String) {
	BLE(Color(0xFF3B82F6), "💎"),       // Blue
	GPS(Color(0xFF10B981), "🟢"),       // Green
	IOT(Color(0xFFEF4444), "🔥"),       // Red
	GEOCACHE(Color(0xFFF59E0B), "✨")    // Amber
}

@Composable
fun GemMatchScreen(
	score: Int,
	onMatchTriggered: (Long) -> Unit,
	modifier: Modifier = Modifier
) {
	val gridWidth = 4
	val gridHeight = 4
	val gridState = remember {
		mutableStateListOf<GemType>().apply {
			for (i in 0 until gridWidth * gridHeight) {
				add(GemType.values().random())
			}
		}
	}

	Column(
		modifier = modifier
			.fillMaxSize()
			.background(Color(0xFF0F172A)), // Slate 900
		horizontalAlignment = Alignment.CenterHorizontally
	) {
		Text(
			text = "Genome Tesseract Gem Matcher",
			fontSize = 22.sp,
			fontWeight = FontWeight.Bold,
			color = Color(0xFF10B981), // Validation Green
			modifier = Modifier.padding(bottom = 8.dp)
		)

		Text(
			text = "Match adjacent gem waveforms to align genome nodes and earn Protein Trits.",
			fontSize = 14.sp,
			color = Color(0xFF94A3B8),
			modifier = Modifier.padding(start = 16.dp, end = 16.dp, bottom = 16.dp)
		)

		Text(
			text = "Session Score: $score",
			fontSize = 18.sp,
			fontWeight = FontWeight.Bold,
			color = Color.White,
			modifier = Modifier.padding(bottom = 16.dp)
		)

		Column(
			verticalArrangement = Arrangement.spacedBy(8.dp),
			modifier = Modifier
				.wrapContentSize()
				.background(Color(0xFF1E293B), RoundedCornerShape(12.dp))
				.padding(16.dp)
		) {
			for (r in 0 until gridHeight) {
				Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
					for (c in 0 until gridWidth) {
						val index = r * gridWidth + c
						val gemType = gridState[index]
						Box(
							modifier = Modifier
								.size(60.dp)
								.clip(RoundedCornerShape(8.dp))
								.background(gemType.color.copy(alpha = 0.2f))
								.clickable {
									// Simple tap matching gameplay: clear and regenerate gem, reward Trits!
									gridState[index] = GemType.values().random()
									onMatchTriggered(15L)
								},
							contentAlignment = Alignment.Center
						) {
							Text(text = gemType.symbol, fontSize = 24.sp)
						}
					}
				}
			}
		}

		Spacer(modifier = Modifier.height(24.dp))

		Button(
			onClick = {
				for (i in 0 until gridWidth * gridHeight) {
					gridState[i] = GemType.values().random()
				}
			},
			colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF10B981))
		) {
			Text("Reshuffle Genome Grid")
		}
	}
}
