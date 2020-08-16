package random

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	attemptsToReadRandomData = 3
)

// GenerateRandomData returns secure random data with the given size
func GenerateRandomData(bitSize int) []byte {
	buffer := make([]byte, bitSize)

	for i := 0; i < attemptsToReadRandomData; i++ {
		_, err := rand.Read(buffer)
		if err == nil {
			break
		}
	}

	return buffer
}

// GenerateRandomHexString returns secure random hex string
func GenerateRandomHexString(bitSize int) string {
	return hex.EncodeToString(GenerateRandomData(bitSize))
}

// GenerateID generates a ID with the given prefix
func GenerateID(prefix string, bitSize int) string {
	return prefix + GenerateRandomHexString(bitSize)
}
