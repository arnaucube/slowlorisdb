package core

import (
	"bytes"
	"crypto/sha256"
)

// Hash is the type for a hash data packet
type Hash [32]byte

// IsZero returns true if the Hash is empty (all zeroes)
func (h *Hash) IsZero() bool {
	z := Hash{}
	if bytes.Equal(z[:], h[:]) {
		return true
	}
	return false
}

// HashBytes performs a hash over a given byte array
func HashBytes(b []byte) Hash {
	h := sha256.Sum256(b)
	return h
}

// PoWData is the interface for the data that have the Nonce parameter to calculate the Proof-of-Work
type PoWData interface {
	Bytes() []byte
	GetNonce() uint64
	IncrementNonce()
}

// CheckPoW verifies the PoW difficulty of a Hash
func CheckPoW(hash Hash, difficulty int) bool {
	var empty [32]byte
	if !bytes.Equal(hash[:][0:difficulty], empty[0:difficulty]) {
		return false
	}
	return true
}

// CalculatePoW calculates the nonce for the given data in order to fit in the current Proof of Work difficulty
func CalculatePoW(data PoWData, difficulty int) (uint64, error) {
	hash := HashBytes(data.Bytes())
	for !CheckPoW(hash, difficulty) {
		data.IncrementNonce()

		hash = HashBytes(data.Bytes())
	}
	return data.GetNonce(), nil
}
