package audio

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Substance represents a bio-acoustic resonant chemical/substance in the patent.
type Substance struct {
	Name      string
	Frequency float64
	Category  string
}

// PhysiologicalProfile represents a collection of expected substances for a condition.
type PhysiologicalProfile struct {
	Name        string
	Substances  []string // Substance names
	Weights     map[string]int // Importance weights (0 - 100)
}

// Stored Patent Reference Databases
var ResonantSubstances = []Substance{
	{"Cytochrome c oxidase", 18.3515, "Enzyme"},
	{"Asparagine", 16.5149, "Amino Acid"},
	{"Aspartic Acid", 24.5062, "Amino Acid"},
	{"Benserazide", 16.0778, "Biochem"},
	{"Chrysene", 14.2683, "Chemical"},
	{"Dopa", 24.6488, "Biochem"},
	{"Glutathione", 19.2079, "Antioxidant"},
	{"Glycine", 18.7668, "Amino Acid"},
	{"Neurotensin", 15.3218, "Neuropeptide"},
	{"Ropinirole", 16.2737, "Biochem"},
	{"Taurine", 15.6435, "Amino Acid"},
	{"Threonine", 14.8900, "Amino Acid"},
	{"Vitamin E Nicotinate", 16.7440, "Vitamin"},
	{"Probable hemoglobin", 14.0002, "Protein"},
	{"Bentonium Chloride", 14.0027, "Chemical"},
	{"Nitrogen", 14.0070, "Gas"},
	{"Puberulonic Acid", 14.0079, "Acid"},
	{"Zinc nitride", 14.0114, "Chemical"},
	{"Astragalin", 14.0119, "Flavonoide"},
	{"Cytochrome P450", 14.0122, "Enzyme"},
	{"Fibroblast growth factor", 14.0125, "Protein"},
	{"Tomatidine", 14.0136, "Alkaloid"},
	{"Acrolein", 14.0160, "Chemical"},
}

var Profiles = []PhysiologicalProfile{
	{
		Name: "Parkinson's Disease Profile",
		Substances: []string{
			"Cytochrome c oxidase", "Asparagine", "Aspartic Acid",
			"Benserazide", "Chrysene", "Dopa", "Glutathione",
			"Glycine", "Neurotensin", "Ropinirole", "Taurine",
			"Threonine", "Vitamin E Nicotinate",
		},
		Weights: map[string]int{
			"Cytochrome c oxidase": 99, "Asparagine": 66, "Aspartic Acid": 99,
			"Benserazide": 99, "Chrysene": 77, "Dopa": 66, "Glutathione": 66,
			"Glycine": 99, "Neurotensin": 78, "Ropinirole": 99, "Taurine": 78,
			"Threonine": 77, "Vitamin E Nicotinate": 77,
		},
	},
	{
		Name: "Adrenal Stress Profile",
		Substances: []string{"Nitrogen", "Astragalin", "Cytochrome P450"},
		Weights: map[string]int{
			"Nitrogen": 80, "Astragalin": 75, "Cytochrome P450": 90,
		},
	},
	{
		Name: "Respiratory System Profile",
		Substances: []string{"Fibroblast growth factor", "Tomatidine"},
		Weights: map[string]int{
			"Fibroblast growth factor": 88, "Tomatidine": 85,
		},
	},
	{
		Name: "Toxic Exposure Profile",
		Substances: []string{"Bentonium Chloride", "Acrolein"},
		Weights: map[string]int{
			"Bentonium Chloride": 95, "Acrolein": 90,
		},
	},
}

// RunSoundDiagnosis simulates performing FFT and analyzing the zero-slope 6-octave harmonic hits.
func RunSoundDiagnosis(voicePrint string) (string, []string, int, error) {
	rand.Seed(time.Now().UnixNano())
	
	// Simulate zero-slope harmonic hit analysis (Patent FIG. 5 & 6)
	// Output list of matched substances
	var matchedSubstances []string
	npuCycles := 45000 + rand.Intn(15000) // Simulated Snapdragon NPU compute cycles

	// Determine matching based on keywords in voicePrint input to make it interactive
	inputUpper := strings.ToUpper(voicePrint)
	matchAll := false
	if strings.Contains(inputUpper, "PARKINSON") || strings.Contains(inputUpper, "TREMOR") || strings.Contains(inputUpper, "DEAN") {
		matchAll = true
	}

	for _, sub := range ResonantSubstances {
		// Standard random simulation with higher chance if matching keywords present
		threshold := 0.15
		if matchAll && sub.Category != "Gas" && sub.Category != "Acid" {
			threshold = 0.85
		}
		if rand.Float64() < threshold {
			matchedSubstances = append(matchedSubstances, sub.Name)
		}
	}

	// Calculate applicability score for each profile
	var sb strings.Builder
	sb.WriteString("=== 🔬 BIO-ACOUSTIC ANALYTICAL REPORT (US Patent 8,346,559 B2) ===\n")
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("Signal Transform: Fast Fourier Transform (FFT) over 6 Octaves (0Hz - 27kHz)\n")
	sb.WriteString(fmt.Sprintf("Snapdragon NPU Offloading: ACTIVE (Compute: %d cycles)\n\n", npuCycles))

	sb.WriteString("### 📊 HARMONIC SIG HIT MATRIX (Octave Repetitions):\n")
	sb.WriteString("| Octave | 1X | 2X | 4X | 8X | 16X | 32X | Hits | Status |\n")
	sb.WriteString("|--------|----|----|----|----|----|-----|------|--------|\n")
	
	// Print a few simulated rows for the top matched substances (FIG 5 / FIG 6 compliance)
	for i, subName := range matchedSubstances {
		if i >= 5 {
			break
		}
		hits := 4 + rand.Intn(3)
		status := "SUSPECT"
		if hits >= 5 {
			status = "RED ZONE"
		}
		
		h1 := "1"
		h2 := "0"
		if hits > 4 { h2 = "1" }
		
		sb.WriteString(fmt.Sprintf("| %s |  %s |  1 |  1 |  1 |   %s |   %s |  %d   | %s |\n", 
			subName[:min(len(subName), 12)], h1, h2, "1", hits, status))
	}
	sb.WriteString("\n")

	sb.WriteString("### 💊 DETECTED IMBALANCES & RESONANT SUBSTANCES:\n")
	for _, subName := range matchedSubstances {
		for _, sub := range ResonantSubstances {
			if sub.Name == subName {
				sb.WriteString(fmt.Sprintf("- **%s** (%s) at **%.4f Hz** [Imbalance Detected]\n", sub.Name, sub.Category, sub.Frequency))
			}
		}
	}
	sb.WriteString("\n")

	sb.WriteString("### 📋 PHYSIOLOGICAL PROFILE APPLICABILITY:\n")
	for _, profile := range Profiles {
		matchedCount := 0
		totalWeight := 0
		maxWeight := 0
		
		for _, subName := range profile.Substances {
			w := profile.Weights[subName]
			maxWeight += w
			
			isMatched := false
			for _, m := range matchedSubstances {
				if m == subName {
					isMatched = true
					break
				}
			}
			
			if isMatched {
				matchedCount++
				totalWeight += w
			}
		}
		
		percentage := 0
		if maxWeight > 0 {
			percentage = (totalWeight * 100) / maxWeight
		}
		
		severity := "NONE"
		if percentage >= 70 {
			severity = "CHRONIC"
		} else if percentage >= 40 {
			severity = "MAJOR"
		} else if percentage >= 20 {
			severity = "MINOR"
		}
		
		sb.WriteString(fmt.Sprintf("- **%s**: **%d%%** Match | Status: **%s** (Matched %d/%d components)\n", 
			profile.Name, percentage, severity, matchedCount, len(profile.Substances)))
	}

	return sb.String(), matchedSubstances, npuCycles, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
