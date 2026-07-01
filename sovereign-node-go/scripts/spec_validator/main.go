package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	specBytes, err := ioutil.ReadFile("protocol/spec/serialization.md")
	if err != nil {
		panic("Failed to read spec: " + err.Error())
	}
	codeBytes, err := ioutil.ReadFile("pkg/audio/bio_inference.go")
	if err != nil {
		panic("Failed to read bio_inference.go: " + err.Error())
	}

	spec := string(specBytes)
	code := string(codeBytes)

	// Rule 1: Big-Endian Enforcement
	if strings.Contains(spec, "BIG-ENDIAN") && !strings.Contains(code, "binary.BigEndian") {
		panic("❌ Spec requires BIG-ENDIAN encoding, but bio_inference.go does not utilize binary.BigEndian")
	}

	// Rule 2: Map key sorting for determinism
	if strings.Contains(spec, "Sorted by key") && !strings.Contains(code, "sort.Strings") {
		panic("❌ Spec requires sorted map keys for tissues, but sort.Strings is missing in bio_inference.go")
	}

	// Rule 3: Lexicographic ordering of conditions/failures
	if strings.Contains(spec, "lexicographically") && !strings.Contains(code, "sort.Slice") {
		panic("❌ Spec requires lexicographical sorting, but sort.Slice is missing in bio_inference.go")
	}

	fmt.Println("✅ Spec validation passed: Code complies with protocol/spec/serialization.md")
}
