package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPartnerTruckQueryFail(t *testing.T) {
	c := require.New(t)

	partnerTruckQuery, err := NewPartnerTruckQuery("", "true")
	c.Equal(ErrMissingPartnerID, err)
	c.Empty(partnerTruckQuery)

	partnerTruckQuery, err = NewPartnerTruckQuery("pablito123", "non-bool")
	c.Equal(ErrInvalidAvailable, err)
}

func TestNewPartnerTruckQuery(t *testing.T) {
	c := require.New(t)

	partnerTruckQuery, err := NewPartnerTruckQuery("enriquito123", "")
	c.Nil(err)
	c.Equal("enriquito123", partnerTruckQuery.PartnerID)
	c.False(partnerTruckQuery.FilterAvailable)
	c.False(partnerTruckQuery.Available)

	partnerTruckQuery, err = NewPartnerTruckQuery("albinito123", "true")
	c.Nil(err)
	c.Equal("albinito123", partnerTruckQuery.PartnerID)
	c.True(partnerTruckQuery.FilterAvailable)
	c.True(partnerTruckQuery.Available)

	partnerTruckQuery, err = NewPartnerTruckQuery("jasoncito123", "false")
	c.Nil(err)
	c.Equal("jasoncito123", partnerTruckQuery.PartnerID)
	c.True(partnerTruckQuery.FilterAvailable)
	c.False(partnerTruckQuery.Available)
}
