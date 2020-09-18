package random

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"kargo-back/shared/environment"
)

const (
	attemptsToReadRandomData = 3
)

var (
	passphrase = environment.GetString("PASSPHRASE", "")

	// StandardBitSize is the standard bit size for models
	StandardBitSize = int(environment.GetInt64("BITSIZE", 1))
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

// GetSHA256WithPrefix returns hexadecimal representation of a string with SHA with prefix
func GetSHA256WithPrefix(prefix, val string) string {
	algorithm := sha256.New()
	algorithm.Write([]byte(val))

	return prefix + hex.EncodeToString(algorithm.Sum([]byte(passphrase)))
}
