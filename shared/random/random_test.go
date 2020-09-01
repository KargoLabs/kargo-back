package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomData(t *testing.T) {
	c := require.New(t)

	c.Equal(8, len(GenerateRandomData(8)))
}

func TestGenerateRandomHexString(t *testing.T) {
	c := require.New(t)

	c.Equal(16, len(GenerateRandomHexString(8)))
}

func TestGenerateID(t *testing.T) {
	c := require.New(t)

	c.Equal(34, len(GenerateID("PR", 16)))
}

func TestGetSHA256WithPrefix(t *testing.T) {
	c := require.New(t)
	apiKey := "dummy"

	enc := GetSHA256WithPrefix("DUM", apiKey)
	c.Equal("DUMb5a2c96250612366ea272ffac6d9744aaf4b45aacd96aa7cfcb931ee3b558259", enc)
}
