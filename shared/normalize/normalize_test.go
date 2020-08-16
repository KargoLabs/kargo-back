package normalize

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeName(t *testing.T) {
	c := require.New(t)

	c.Equal("CÉSAR GONZÁLEZ", NormalizeName("césar     gonzález    "))
}
