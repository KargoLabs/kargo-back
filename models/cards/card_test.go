package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCardFails(t *testing.T) {
	c := require.New(t)

	card, err := NewCard("", "4242424242424242", "test", "800", "2025", "12")
	c.Nil(card)
	c.Equal(ErrMissingUserID, err)

	card, err = NewCard("USER123", "", "test", "800", "2025", "12")
	c.Nil(card)
	c.Equal(ErrMissingNumber, err)

	card, err = NewCard("USER123", "4242424242424242", "", "800", "2025", "12")
	c.Nil(card)
	c.Equal(ErrMissingName, err)

	card, err = NewCard("USER123", "4242424242424242", "test", "", "2025", "12")
	c.Nil(card)
	c.Equal(ErrMissingCSV, err)

	card, err = NewCard("USER123", "4242424242424242", "test", "100", "", "12")
	c.Nil(card)
	c.Equal(ErrMissingYear, err)

	card, err = NewCard("USER123", "4242424242424242", "test", "100", "2000", "12")
	c.Nil(card)
	c.Equal(ErrInvalidYear, err)

	card, err = NewCard("USER123", "4242424242424242", "test", "100", "2025", "")
	c.Nil(card)
	c.Equal(ErrMissingMonth, err)

	card, err = NewCard("USER123", "4242424242424242", "test", "100", "2020", "01")
	c.Nil(card)
	c.Equal(ErrInvalidMonth, err)

	card, err = NewCard("USER123", "434234232", "test", "100", "2025", "12")
	c.Nil(card)
}

func TestNewCard(t *testing.T) {
	c := require.New(t)

	card, err := NewCard("USR123", "4242424242424242", "test", "800", "2025", "12")
	c.Nil(err)
	c.Equal("USR123", card.UserID)
	c.Equal("test", card.Name)
	c.Equal("4242", card.LastFourDigits)
	c.Equal("Visa", card.Company)
}
