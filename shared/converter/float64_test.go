package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToGreaterThanZeroFloat(t *testing.T) {
	c := require.New(t)

	float, ok := ConvertToGreaterThanZeroFloat("10")
	c.True(ok)
	c.Equal(float64(10), float)

	float, ok = ConvertToGreaterThanZeroFloat("-10")
	c.False(ok)
	c.Equal(float64(-10), float)

	float, ok = ConvertToGreaterThanZeroFloat("ajua")
	c.False(ok)
	c.Empty(float)
}
