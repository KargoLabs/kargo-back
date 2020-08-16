package clientModels

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClientFail(t *testing.T) {
	c := require.New(t)

	birthDate, err := time.Parse("2006-01-02", "1999-07-20")
	c.Nil(err)

	client, err := NewClient("", "12345", birthDate)
	c.Nil(client)
	c.Equal(ErrMissingName, err)

	client, err = NewClient("pablo ocampo", "", birthDate)
	c.Nil(client)
	c.Equal(ErrMissingDocument, err)

	client, err = NewClient("pablo ocampo", "12345", time.Time{})
	c.Nil(client)
	c.Equal(ErrMissingBirthDate, err)
}

func TestNewClient(t *testing.T) {
	c := require.New(t)

	birthDate, err := time.Parse("2006-01-02", "1999-07-20")
	c.Nil(err)

	client, err := NewClient("pablo ocampo", "12345", birthDate)
	c.Nil(err)

	c.NotEmpty(client.ClientID)
	c.Equal("PABLO OCAMPO", client.Name)
	c.Equal("12345", client.Document)
	c.Equal(birthDate, client.BirthDate)
	c.NotEmpty(client.CreationDate)
}
