package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTrucksQueryFail(t *testing.T) {
	c := require.New(t)

	trucksQuery, err := NewTrucksQuery("camioncito", "500", "600", "DO-33", "DO-34")
	c.Equal(ErrInvalidTruckType, err)
	c.Empty(trucksQuery)

	trucksQuery, err = NewTrucksQuery("rigid", "500", "600", "ocampo region", "DO-34")
	c.Equal(ErrInvalidOrigin, err)
	c.Empty(trucksQuery)

	trucksQuery, err = NewTrucksQuery("rigid", "500", "600", "DO-34", "ocampo region")
	c.Equal(ErrInvalidDestination, err)
	c.Empty(trucksQuery)

	trucksQuery, err = NewTrucksQuery("rigid", "-2", "600", "DO-33", "DO-34")
	c.Equal(ErrInvalidWeight, err)
	c.Empty(trucksQuery)

	trucksQuery, err = NewTrucksQuery("rigid", "600", "hh", "DO-33", "DO-34")
	c.Equal(ErrInvalidVolume, err)
	c.Empty(trucksQuery)
}

func TestNewTrucksQuery(t *testing.T) {
	c := require.New(t)

	trucksQuery, err := NewTrucksQuery("rigid", "500", "600", "DO-33", "DO-34")
	c.Nil(err)

	c.Equal("rigid", trucksQuery.TruckType)
	c.Equal("DO-33", trucksQuery.Origin)
	c.Equal("DO-34", trucksQuery.Destination)
	c.Equal(float64(500), trucksQuery.Weight)
	c.Equal(float64(600), trucksQuery.Volume)
}
