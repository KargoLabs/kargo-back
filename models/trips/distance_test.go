package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDistance(t *testing.T) {
	c := require.New(t)

	c.Equal(float64(7), GetDistance(EnriquilloRegion, CibaoNoroesteRegion))
}
