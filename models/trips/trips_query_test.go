package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTripsQueryFail(t *testing.T) {
	c := require.New(t)

	tripsQuery, err := NewTripsQuery("", "true")
	c.Equal(ErrMissingValue, err)
	c.Empty(tripsQuery)

	tripsQuery, err = NewTripsQuery("pablito123", "non-bool")
	c.Equal(ErrInvalidFinished, err)
}

func TestNewTripsQuery(t *testing.T) {
	c := require.New(t)

	tripsQuery, err := NewTripsQuery("enriquito123", "")
	c.Nil(err)
	c.Equal("enriquito123", tripsQuery.Value)
	c.False(tripsQuery.FilterFinished)
	c.False(tripsQuery.Finished)

	tripsQuery, err = NewTripsQuery("albinito123", "true")
	c.Nil(err)
	c.Equal("albinito123", tripsQuery.Value)
	c.True(tripsQuery.FilterFinished)
	c.True(tripsQuery.Finished)

	tripsQuery, err = NewTripsQuery("jasoncito123", "false")
	c.Nil(err)
	c.Equal("jasoncito123", tripsQuery.Value)
	c.True(tripsQuery.FilterFinished)
	c.False(tripsQuery.Finished)
}
