package audio

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
)

type ProteinAgent struct {
	ID        string
	Gene      string
	Beta      []float64
	Nutrients []NutrientDependency
	Tissues   map[string]float64
}

type NutrientDependency struct {
	Nutrient         string
	Role             string
	ActivationEffect float64
	StabilityEffect  float64
	BindingEffect    float64
}

type ProteinState struct {
	ID              string
	Activation      float64
	Stability       float64
	BindingAffinity float64
}

type DiagnosticTicket struct {
	ProteinFailures []string
	TissueImpacts   map[string]float64
	Conditions      []ConditionScore
	Confidence      float64
	ExecutionHash   string
}

type ConditionScore struct {
	Name  string
	Score float64
}

func GetDefaultProteins() []ProteinAgent {
	return []ProteinAgent{
		{
			ID:   "Cytochrome_c_oxidase",
			Gene: "COX1",
			Beta: []float64{0.8, 0.1, -0.3, 0.4, 0.2},
			Nutrients: []NutrientDependency{
				{Nutrient: "Iron", Role: "cofactor", ActivationEffect: 0.5, StabilityEffect: 0.7, BindingEffect: 0.8},
			},
			Tissues: map[string]float64{"brain": 0.9, "liver": 0.4, "kidney": 0.2},
		},
		{
			ID:   "Dopa_decarboxylase",
			Gene: "DDC",
			Beta: []float64{0.2, 0.9, 0.4, -0.1, 0.6},
			Nutrients: []NutrientDependency{
				{Nutrient: "B12", Role: "methylation", ActivationEffect: 0.6, StabilityEffect: 0.8, BindingEffect: 0.9},
				{Nutrient: "Mg", Role: "cofactor", ActivationEffect: 0.7, StabilityEffect: 0.9, BindingEffect: 0.9},
			},
			Tissues: map[string]float64{"brain": 0.95, "liver": 0.3},
		},
		{
			ID:   "Glutathione_synthetase",
			Gene: "GSS",
			Beta: []float64{-0.1, 0.3, 0.8, 0.5, 0.2},
			Nutrients: []NutrientDependency{
				{Nutrient: "Mg", Role: "synthesis", ActivationEffect: 0.5, StabilityEffect: 0.8, BindingEffect: 0.8},
			},
			Tissues: map[string]float64{"liver": 0.9, "brain": 0.5, "kidney": 0.6},
		},
	}
}

func GetDefaultDiseaseMatrix() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"Parkinson's Disease Profile": {
			"Cytochrome_c_oxidase": 0.6,
			"Dopa_decarboxylase":   0.9,
		},
		"Adrenal Stress Profile": {
			"Glutathione_synthetase": 0.8,
			"Cytochrome_c_oxidase":   0.4,
		},
	}
}

func ComputeProteinStates(proteins []ProteinAgent, z []float64) []ProteinState {
	var states []ProteinState
	for _, p := range proteins {
		activation := dot(p.Beta, z)
		state := ProteinState{
			ID:              p.ID,
			Activation:      max(0.0, activation),
			Stability:       1.0,
			BindingAffinity: 1.0,
		}
		states = append(states, state)
	}
	return states
}

func ApplyNutrientConstraints(states []ProteinState, proteins []ProteinAgent, nutrients map[string]float64) {
	for i := range states {
		p := proteins[i]
		state := &states[i]
		for _, dep := range p.Nutrients {
			level, ok := nutrients[dep.Nutrient]
			if !ok {
				level = 1.0
			}
			if level < 0.5 {
				state.Activation *= dep.ActivationEffect
				state.Stability *= dep.StabilityEffect
				state.BindingAffinity *= dep.BindingEffect
			}
		}
	}
}

func ComputeTissueImpact(states []ProteinState, proteins []ProteinAgent) map[string]float64 {
	result := make(map[string]float64)
	for i, state := range states {
		p := proteins[i]
		for tissue, expr := range p.Tissues {
			result[tissue] += state.Activation * expr
		}
	}
	return result
}

// GenerateCanonicalBinaryHash serializes states, failures, tissue impacts, and condition scores
// into a canonical binary stream according to the serialization specification, then returns its SHA-256 hash.
func GenerateCanonicalBinaryHash(
	states []ProteinState,
	failures []string,
	tissues map[string]float64,
	conditions []ConditionScore,
	protocolVersion string,
	simulationRoot string,
) string {
	hasher := sha256.New()

	// Helper to write uint16 length and string bytes
	writeString := func(s string) {
		length := uint16(len(s))
		var buf [2]byte
		binary.BigEndian.PutUint16(buf[:], length)
		hasher.Write(buf[:])
		hasher.Write([]byte(s))
	}

	// Helper to write float64 directly via binary representation
	writeFloat := func(f float64) {
		bits := math.Float64bits(f)
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], bits)
		hasher.Write(buf[:])
	}

	// 1. Hash Protein States (Assumed sorted lexicographically by ID beforehand)
	for _, s := range states {
		writeString(s.ID)
		writeFloat(s.Activation)
		writeFloat(s.Stability)
		writeFloat(s.BindingAffinity)
	}

	// 2. Hash Failures (lexicographically sorted)
	sortedFailures := make([]string, len(failures))
	copy(sortedFailures, failures)
	sort.Strings(sortedFailures)
	var failCountBuf [2]byte
	binary.BigEndian.PutUint16(failCountBuf[:], uint16(len(sortedFailures)))
	hasher.Write(failCountBuf[:])
	for _, f := range sortedFailures {
		writeString(f)
	}

	// 3. Hash Tissues (Sorted by key names)
	tissueKeys := make([]string, 0, len(tissues))
	for k := range tissues {
		tissueKeys = append(tissueKeys, k)
	}
	sort.Strings(tissueKeys)
	var tissueCountBuf [2]byte
	binary.BigEndian.PutUint16(tissueCountBuf[:], uint16(len(tissueKeys)))
	hasher.Write(tissueCountBuf[:])
	for _, k := range tissueKeys {
		writeString(k)
		writeFloat(tissues[k])
	}

	// 4. Hash Conditions (Sorted primarily by Score descending, then Name ascending)
	sortedConditions := make([]ConditionScore, len(conditions))
	copy(sortedConditions, conditions)
	sort.Slice(sortedConditions, func(i, j int) bool {
		if sortedConditions[i].Score == sortedConditions[j].Score {
			return sortedConditions[i].Name < sortedConditions[j].Name
		}
		return sortedConditions[i].Score > sortedConditions[j].Score
	})
	var condCountBuf [2]byte
	binary.BigEndian.PutUint16(condCountBuf[:], uint16(len(sortedConditions)))
	hasher.Write(condCountBuf[:])
	for _, c := range sortedConditions {
		writeString(c.Name)
		writeFloat(c.Score)
	}

	// 5. Append ProtocolVersion and SimulationRoot to secure future changes into consensus
	writeString(protocolVersion)
	writeString(simulationRoot)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func GenerateDiagnosis(states []ProteinState, tissues map[string]float64, diseaseMatrix map[string]map[string]float64, protocolVersion, simulationRoot string) DiagnosticTicket {
	ticket := DiagnosticTicket{}
	ticket.TissueImpacts = tissues
	for _, s := range states {
		if s.Activation < 0.3 || s.Stability < 0.5 {
			ticket.ProteinFailures = append(ticket.ProteinFailures, s.ID)
		}
	}
	for disease, weights := range diseaseMatrix {
		score := 0.0
		for _, s := range states {
			score += weights[s.ID] * s.Activation
		}
		if score > 0.5 {
			ticket.Conditions = append(ticket.Conditions, ConditionScore{
				Name:  disease,
				Score: score,
			})
		}
	}
	sort.Slice(ticket.Conditions, func(i, j int) bool {
		return ticket.Conditions[i].Score > ticket.Conditions[j].Score
	})
	if len(ticket.Conditions) > 0 {
		ticket.Confidence = ticket.Conditions[0].Score
	}

	// Create cryptographic verification hash across all execution outputs
	ticket.ExecutionHash = GenerateCanonicalBinaryHash(states, ticket.ProteinFailures, tissues, ticket.Conditions, protocolVersion, simulationRoot)

	return ticket
}

func RunBioInference(
	z []float64,
	proteins []ProteinAgent,
	nutrientState map[string]float64,
	diseaseMatrix map[string]map[string]float64,
	protocolVersion string,
	simulationRoot string,
) DiagnosticTicket {
	// Sort protein agents deterministically by ID to prevent order mismatches
	sort.Slice(proteins, func(i, j int) bool {
		return proteins[i].ID < proteins[j].ID
	})
	states := ComputeProteinStates(proteins, z)
	ApplyNutrientConstraints(states, proteins, nutrientState)
	tissues := ComputeTissueImpact(states, proteins)
	return GenerateDiagnosis(states, tissues, diseaseMatrix, protocolVersion, simulationRoot)
}

func dot(a, b []float64) float64 {
	val := 0.0
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		val += a[i] * b[i]
	}
	return val
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
