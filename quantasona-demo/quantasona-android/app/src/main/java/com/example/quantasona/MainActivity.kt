package com.example.quantasona

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.lifecycle.lifecycleScope
import com.example.quantasona.audio.AudioProcessor
import com.example.quantasona.audio.MfccExtractor
import com.example.quantasona.inference.BiomarkerModel
import com.example.quantasona.mesh.AndroidMeshClient
import com.example.quantasona.mesh.AndroidMeshStream
import com.example.quantasona.pipeline.BiomarkerPipeline
import android.util.Log
import kotlinx.coroutines.launch

class MainActivity : ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val audioProcessor = AudioProcessor()
        val mfccExtractor = MfccExtractor()
        val biomarkerModel = BiomarkerModel(this)
        val meshClient = AndroidMeshClient()
        val pipeline = BiomarkerPipeline(
            lifecycleScope,
            audioProcessor,
            mfccExtractor,
            biomarkerModel,
            meshClient
        )

        pipeline.start()

        val stream = AndroidMeshStream()

        lifecycleScope.launch {
            stream.stream { msg ->
                Log.d("MeshStream", msg)
            }
        }

        setContent {
            QuantasonaTheme {
                PatentDashboardScreen(
                    onProcessCycleComplete = {
                        pipeline.forceFlush()
                    }
                )
            }
        }
    }
}

@Composable
fun PatentDashboardScreen(onProcessCycleComplete: () -> Unit) {
    var patentText by remember { mutableStateOf("") }
    var parsedData by remember { mutableStateOf("") }
    var llmQuery by remember { mutableStateOf("") }
    var llmResponse by remember { mutableStateOf("") }

    MaterialTheme {
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(16.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            item {
                Text(
                    text = "Quantasona Patent Dashboard",
                    style = MaterialTheme.typography.h5
                )
            }

            // Patent Upload/Input Section
            item {
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text("Upload / Paste Patent Documentation", style = MaterialTheme.typography.h6)
                        Spacer(Modifier.height(8.dp))
                        OutlinedTextField(
                            value = patentText,
                            onValueChange = { patentText = it },
                            modifier = Modifier
                                .fillMaxWidth()
                                .height(150.dp),
                            label = { Text("Patent Text") }
                        )
                        Spacer(Modifier.height(8.dp))
                        Button(
                            onClick = {
                                // TODO: Integrate with Backend/LLM for parsing
                                parsedData = "Parsed output will appear here..."
                            },
                            modifier = Modifier.align(Alignment.End)
                        ) {
                            Text("Parse Patent")
                        }
                    }
                }
            }

            // Parsed Result Section
            if (parsedData.isNotEmpty()) {
                item {
                    Card(modifier = Modifier.fillMaxWidth()) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Parsed Data", style = MaterialTheme.typography.h6)
                            Spacer(Modifier.height(8.dp))
                            Text(parsedData, style = MaterialTheme.typography.body1)
                        }
                    }
                }
            }

            // LLM Interaction Section
            item {
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text("Interact with LLM", style = MaterialTheme.typography.h6)
                        Spacer(Modifier.height(8.dp))
                        OutlinedTextField(
                            value = llmQuery,
                            onValueChange = { llmQuery = it },
                            modifier = Modifier.fillMaxWidth(),
                            label = { Text("Ask a question about the patent") }
                        )
                        Spacer(Modifier.height(8.dp))
                        Button(
                            onClick = {
                                // TODO: Integrate with LLM API
                                llmResponse = "LLM response will appear here..."
                            },
                            modifier = Modifier.align(Alignment.End)
                        ) {
                            Text("Ask")
                        }
                        if (llmResponse.isNotEmpty()) {
                            Spacer(Modifier.height(16.dp))
                            Text("Response:", style = MaterialTheme.typography.subtitle1)
                            Text(llmResponse, style = MaterialTheme.typography.body1)
                        }
                    }
                }
            }

            // Core System Status Block (Maintained from previous screen)
            item {
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text("Audio Pipeline Status: ACTIVE", style = MaterialTheme.typography.h6)
                        Spacer(Modifier.height(8.dp))
                        Text("Sensory Input: Biomarker Stream Active", style = MaterialTheme.typography.body1)
                        Text("Network Link: gRPC Mesh Connected", style = MaterialTheme.typography.body1)
                    }
                }
            }

            item {
                Button(onClick = onProcessCycleComplete, modifier = Modifier.fillMaxWidth()) {
                    Text("Execute Biomarker Cycle & Publish")
                }
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun PatentDashboardPreview() {
    QuantasonaTheme {
        PatentDashboardScreen(onProcessCycleComplete = {})
    }
}

@Composable
fun QuantasonaTheme(content: @Composable () -> Unit) {
    MaterialTheme(content = content)
}
