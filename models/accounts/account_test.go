package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAccountFails(t *testing.T) {
	c := require.New(t)

	account, err := NewAccount("", "Osvaldo", "12345678901", "5678")
	c.Nil(account)
	c.Equal(ErrMissingPartnerID, err)

	account, err = NewAccount("PAR123", "", "12345678901", "5678")
	c.Nil(account)
	c.Equal(ErrMissingName, err)

	account, err = NewAccount("PAR123", "Osvaldo", "", "5678")
	c.Nil(account)
	c.Equal(ErrMissingDocument, err)

	account, err = NewAccount("PAR123", "Osvaldo", "12345678901", "")
	c.Nil(account)
	c.Equal(ErrMissingNumber, err)
}

func TestNewCard(t *testing.T) {
	c := require.New(t)

	account, err := NewAccount("PAR123", "Osvaldo", "12345678901", "5678")
	c.Nil(err)
	c.Equal("PAR123", account.PartnerID)
	c.Equal("Osvaldo", account.Name)
	c.Equal("12345678901", account.Document)
	c.Equal("5678", account.Number)
}
