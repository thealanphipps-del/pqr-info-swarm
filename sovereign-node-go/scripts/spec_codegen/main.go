package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	versionBytes, err := ioutil.ReadFile("protocol/VERSION.yaml")
	if err != nil {
		panic("Failed to read VERSION.yaml: " + err.Error())
	}
	version := string(versionBytes)
	
	// Quick extraction of the main version field
	var mainVersion string
	lines := strings.Split(version, "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "version:") {
			mainVersion = strings.TrimSpace(strings.TrimPrefix(l, "version:"))
			break
		}
	}
	if mainVersion == "" {
		mainVersion = "v0.0.0"
	}

	output := fmt.Sprintf(`package protocol

// Spec: protocol/VERSION.yaml
const ProtocolVersion = "%s"
`, mainVersion)

	err = ioutil.WriteFile("pkg/protocol/version.go", []byte(output), 0644)
	if err != nil {
		panic("Failed to write version.go: " + err.Error())
	}
	fmt.Println("✅ version.go generated from spec")
}
