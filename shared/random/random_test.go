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
