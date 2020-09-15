package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRegionFromProvince(t *testing.T) {
	c := require.New(t)

	region, err := GetRegionFromProvince(ProvinceLaRomana)
	c.Nil(err)

	c.Equal(YumaRegion, region)
}

func TestGetRegionFromProvinceFail(t *testing.T) {
	c := require.New(t)

	_, err := GetRegionFromProvince("ohana")
	c.Equal(errNoRegionForProvince, err)
}
