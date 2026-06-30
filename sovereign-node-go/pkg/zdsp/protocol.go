package zdsp

import (
	"crypto/sha256"
	"encoding/binary"
)

type ZDSPPacket struct {
	Timestamp  uint64
	ParentHash [32]byte
	Payload    []byte
	Signature  []byte
}

func (p *ZDSPPacket) Hash() [32]byte {
	data := make([]byte, 8+32+len(p.Payload))
	binary.BigEndian.PutUint64(data[0:8], p.Timestamp)
	copy(data[8:40], p.ParentHash[:])
	copy(data[40:], p.Payload)
	return sha256.Sum256(data)
}
