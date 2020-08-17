package environment

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStringDefaultValue(t *testing.T) {
	t.Parallel()
	c := require.New(t)

	defaultValue := "val"
	c.Equal(defaultValue, GetString("GET_STRING", defaultValue))
}

func TestGetStringCustomValue(t *testing.T) {
	t.Parallel()
	c := require.New(t)

	withTestEnv("custom", func(varName string) {
		c.Equal("custom", GetString(varName, ""))
	})
}

func withTestEnv(val interface{}, codeBlock func(varName string)) {
	varName := fmt.Sprintf("TEST_%d_%d", rand.Intn(math.MaxInt32), rand.Intn(math.MaxInt32))
	_ = os.Setenv(varName, fmt.Sprintf("%v", val))
	defer func() {
		_ = os.Unsetenv(varName)
	}()

	codeBlock(varName)
}
