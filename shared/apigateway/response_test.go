package apigateway

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJSONResponse(t *testing.T) {
	c := require.New(t)
	response := NewJSONResponse(200, map[string]string{"hello": "world"})

	c.Equal(200, response.StatusCode)
	c.Equal("application/json", response.Headers["Content-Type"])
	c.JSONEq(`{"hello":"world"}`, response.Body)
}
