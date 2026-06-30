package protocol

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Spec: protocol/VERSION.yaml
const ProtocolVersion = "v0.4.0"

// SimulationLeaf represents a single workload simulation case output.
type SimulationLeaf struct {
	CaseID    string
	OldHash   string
	NewHash   string
	HashMatch bool
}

// SimulationConsensusProof represents a Merkle-rooted cryptographic commitment to simulation outcomes.
type SimulationConsensusProof struct {
	RFCID           string
	ProtocolVersion string
	WorkloadSetHash string
	SimulationRoot  string
	NodeSignatures  []string
}

// SimulationVote represents a node's Slingshot consensus vote containing the simulation root.
type SimulationVote struct {
	NodeID         string
	SimulationRoot string
}

// SerializeSimulationLeaf serializes a leaf deterministically into a binary structure.
func SerializeSimulationLeaf(l SimulationLeaf) []byte {
	buf := &bytes.Buffer{}
	
	// Write length-prefixed strings in big-endian
	writeString := func(s string) {
		length := uint16(len(s))
		var lenBuf [2]byte
		binary.BigEndian.PutUint16(lenBuf[:], length)
		buf.Write(lenBuf[:])
		buf.Write([]byte(s))
	}

	writeString(l.CaseID)
	writeString(l.OldHash)
	writeString(l.NewHash)

	match := byte(0)
	if l.HashMatch {
		match = 1
	}
	buf.WriteByte(match)

	return buf.Bytes()
}

// BuildSimulationLeafHashes calculates leaf hashes for Merkle tree builder
func BuildSimulationLeafHashes(leaves []SimulationLeaf) [][]byte {
	var hashes [][]byte
	for _, l := range leaves {
		data := SerializeSimulationLeaf(l)
		h := sha256.Sum256(data)
		hashes = append(hashes, h[:])
	}
	return hashes
}

// ComputeSimulationRoot constructs a Merkle tree of leaf hashes and returns the hex Merkle Root
func ComputeSimulationRoot(leaves []SimulationLeaf) string {
	hashes := BuildSimulationLeafHashes(leaves)
	if len(hashes) == 0 {
		return fmt.Sprintf("%x", sha256.Sum256([]byte("empty")))
	}

	for len(hashes) > 1 {
		var nextLevel [][]byte
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				combined := append(hashes[i], hashes[i+1]...)
				h := sha256.Sum256(combined)
				nextLevel = append(nextLevel, h[:])
			} else {
				nextLevel = append(nextLevel, hashes[i])
			}
		}
		hashes = nextLevel
	}

	return fmt.Sprintf("%x", hashes[0])
}

// VerifySimulationConsensus checks if a 2/3 majority of nodes agree on the SimulationRoot.
func VerifySimulationConsensus(votes []SimulationVote) (string, bool) {
	counts := make(map[string]int)
	for _, v := range votes {
		counts[v.SimulationRoot]++
	}
	total := len(votes)
	if total == 0 {
		return "", false
	}
	for root, c := range counts {
		if float64(c) >= (2.0/3.0)*float64(total) {
			return root, true
		}
	}
	return "", false
}
