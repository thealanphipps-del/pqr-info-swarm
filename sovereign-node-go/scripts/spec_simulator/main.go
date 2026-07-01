package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sovereign-node-go/pkg/audio"
	"sovereign-node-go/pkg/protocol"
)

// Simulation Workload defines historical outputs to compare against
type Workload struct {
	ID            string             `json:"id"`
	Z             []float64          `json:"z"`
	Nutrients     map[string]float64 `json:"nutrients"`
	ExpectedHash  string             `json:"expected_hash"`
	FailuresCount int                `json:"failures_count"`
}

func main() {
	fmt.Println("=== 🧪 Sovereign Mesh Protocol Upgrade Simulation Engine ===")

	// Dummy/default values for version binding verification
	dummySimRoot := "0000000000000000000000000000000000000000000000000000000000000000"

	// 1. Load historical workloads
	workloadsPath := "protocol/simulation/workloads.json"
	workloadBytes, err := ioutil.ReadFile(workloadsPath)
	if err != nil {
		fmt.Printf("[SIM] No workloads found at %s. Creating default seed workloads...\n", workloadsPath)
		// Seed default workloads
		defaultWorkloads := []Workload{
			{
				ID:            "workload-1",
				Z:             []float64{0.2, -0.1, 0.5, 0.7, -0.3},
				Nutrients:     map[string]float64{"B12": 0.2, "Iron": 0.5, "Mg": 0.4},
				ExpectedHash:  "",
				FailuresCount: 2,
			},
			{
				ID:            "workload-2",
				Z:             []float64{0.1, 0.1, 0.1, 0.1, 0.1},
				Nutrients:     map[string]float64{"B12": 0.8, "Iron": 0.9, "Mg": 0.95},
				ExpectedHash:  "",
				FailuresCount: 0,
			},
		}

		// Calculate initial actual hashes for the seed
		proteins := audio.GetDefaultProteins()
		diseaseMatrix := audio.GetDefaultDiseaseMatrix()
		for i, w := range defaultWorkloads {
			diag := audio.RunBioInference(w.Z, proteins, w.Nutrients, diseaseMatrix, protocol.ProtocolVersion, dummySimRoot)
			defaultWorkloads[i].ExpectedHash = diag.ExecutionHash
			defaultWorkloads[i].FailuresCount = len(diag.ProteinFailures)
		}

		os.MkdirAll(filepath.Dir(workloadsPath), 0755)
		data, _ := json.MarshalIndent(defaultWorkloads, "", "  ")
		ioutil.WriteFile(workloadsPath, data, 0644)
		workloadBytes = data
	}

	var workloads []Workload
	if err := json.Unmarshal(workloadBytes, &workloads); err != nil {
		panic("Failed to parse workloads: " + err.Error())
	}

	// 2. Execute bio-inference simulation under current local code engine
	proteins := audio.GetDefaultProteins()
	diseaseMatrix := audio.GetDefaultDiseaseMatrix()

	mismatches := 0
	total := len(workloads)
	
	fmt.Printf("[SIM] Loaded %d historical workloads for compliance checking.\n", total)

	var leaves []protocol.SimulationLeaf

	for i, w := range workloads {
		diag := audio.RunBioInference(w.Z, proteins, w.Nutrients, diseaseMatrix, protocol.ProtocolVersion, dummySimRoot)
		
		fmt.Printf("\n[WORKLOAD %d] Input Z: %v\n", i+1, w.Z)
		fmt.Printf("  -> Expected Hash: %s\n", w.ExpectedHash)
		fmt.Printf("  -> Actual Hash:   %s\n", diag.ExecutionHash)
		fmt.Printf("  -> Failures: %v\n", diag.ProteinFailures)

		match := diag.ExecutionHash == w.ExpectedHash
		if !match {
			fmt.Printf("  ⚠️ [DRIFT DETECTED] Hash mismatch!\n")
			mismatches++
		} else {
			fmt.Printf("  ✅ [PASS] Hash matches exactly.\n")
		}

		leaves = append(leaves, protocol.SimulationLeaf{
			CaseID:    w.ID,
			OldHash:   w.ExpectedHash,
			NewHash:   diag.ExecutionHash,
			HashMatch: match,
		})
	}

	// Calculate and commit Simulation Consensus Proof (SCP) Merkle Root
	simRoot := protocol.ComputeSimulationRoot(leaves)
	fmt.Printf("\n[SCP] Computed Merkle-Rooted Simulation Consensus Proof: %s\n", simRoot)

	driftScore := float64(mismatches) / float64(total)
	fmt.Printf("\n=== Drift Report ===\n")
	fmt.Printf("Total Workloads Run: %d\n", total)
	fmt.Printf("Mismatched Hashes:  %d\n", mismatches)
	fmt.Printf("Drift Score:        %.2f%%\n", driftScore*100.0)

	// Enforce thresholds
	threshold := 0.0 // Reject any drift in execution hash
	if driftScore > threshold {
		fmt.Println("❌ [REJECTED] Spec upgrade validation failed: Consensus execution hash drift detected!")
		os.Exit(2)
	}

	fmt.Println("✅ [APPROVED] Spec upgrade validation passed: Zero execution hash drift. Upgrades are consensus-safe.")
}
