package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTransactionFail(t *testing.T) {
	c := require.New(t)

	transaction, err := NewTransaction("", "partner01", "CARD123", 100)
	c.Nil(transaction)
	c.Equal(ErrMissingClientID, err)

	transaction, err = NewTransaction("client01", "", "CARD123", 100)
	c.Nil(transaction)
	c.Equal(ErrMissingPartnerID, err)

	transaction, err = NewTransaction("client01", "partner01", "", 100)
	c.Nil(transaction)
	c.Equal(ErrMissingCardID, err)

	transaction, err = NewTransaction("client01", "partner01", "CARD123", -100)
	c.Nil(transaction)
	c.Equal(ErrInvalidAmount, err)
}

func TestNewTrancation(t *testing.T) {
	c := require.New(t)

	transaction, err := NewTransaction("client01", "partner01", "CARD123", 100)
	c.Nil(err)
	c.Equal("client01", transaction.ClientID)
	c.Equal("partner01", transaction.PartnerID)
	c.Equal(float64(100), transaction.Amount)
	c.Equal(TransactionStatusOnHold, transaction.Status)
	c.NotEmpty(transaction.TransactionID)
	c.NotEmpty(transaction.Date)
}
