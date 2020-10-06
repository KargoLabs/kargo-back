package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPartnerFail(t *testing.T) {
	c := require.New(t)

	birthdate, err := time.Parse("2006-01-02", "2000-05-02")
	c.Nil(err)

	partner, err := NewPartner("dummyUser", "", "2222", "+18291234567", "kike@kargo.com", birthdate)
	c.Nil(partner)
	c.Equal(ErrMissingName, err)

	partner, err = NewPartner("dummyUser", "roniel valdez", "", "+18291234567", "kike@kargo.com", birthdate)
	c.Nil(partner)
	c.Equal(ErrMissingDocument, err)

	partner, err = NewPartner("dummyUser", "roniel valdez", "12345", "", "kike@kargo.com", birthdate)
	c.Nil(partner)
	c.Equal(ErrMissingPhoneNumber, err)

	partner, err = NewPartner("dummyUser", "roniel valdez", "12345", "+18291234567", "kike@kargo.com", time.Time{})
	c.Nil(partner)
	c.Equal(ErrMissingBirthdate, err)
}

func TestNewPartner(t *testing.T) {
	c := require.New(t)

	birthdate, err := time.Parse("2006-01-02", "2000-05-02")
	c.Nil(err)

	partner, err := NewPartner("dummyUser", "roniel valdez", "12345", "+18291234567", "kike@kargo.com", birthdate)
	c.Nil(err)

	c.NotEmpty(partner.PartnerID)
	c.Equal("RONIEL VALDEZ", partner.Name)
	c.Equal("12345", partner.Document)
	c.Equal(birthdate, partner.Birthdate)
	c.NotEmpty(partner.CreationDate)
}
