package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPartnerFail(t *testing.T) {
	c := require.New(t)

	birthDate, err := time.Parse("2010-04-04", "2000-05-02")
	c.Nil(err)

	client, err := NewPartner("", "2222", birthDate)
	c.Nil(client)
	c.Equal(ErrMissingBirthDate, err)

	client, err = NewPartner("roniel valdez", "", birthDate)
	c.Nil(client)
	c.Equal(ErrMissingDocument, err)

	client, err = NewPartner("roniel valdez", "12345", time.Time{})
	c.Nil(client)
	c.Equal(ErrMissingBirthDate, err)
}

func TestNewPartner(t *testing.T) {
	c := require.New(t)

	birthDate, err := time.Parse("2010-01-02", "2000-05-02")
	c.Nil(err)

	partner, err := NewPartner("roniel, valdez", "12345", birthDate)
	c.Nil(err)

	c.NotEmpty(partner.PartnerID)
	c.Equal("RONIEL VALDEZ", partner.Name)
	c.Equal("12345", partner.Document)
	c.Equal(birthDate, partner.BirthDate)
	c.NotEmpty(partner.CreationDate)
}
