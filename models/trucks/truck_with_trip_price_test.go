package models

import (
	models "kargo-back/models/trips"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTruckWithTripPrice(t *testing.T) {
	c := require.New(t)

	truck := &Truck{
		PartnerID:         "PAR123",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234,
		BasePrice:         1234,
		PerRegionPrice:    1234,
		Regions:           []models.Region{"DO-33", "DO-35"}}

	truckWithTripPrice := NewTruckWithTripPrice("DO-33", "DO-35", truck)

	c.Equal("PAR123", truckWithTripPrice.PartnerID)
	c.Equal("12345678901", truckWithTripPrice.RegistrationPlate)
	c.Equal("Freightliner", truckWithTripPrice.Brand)
	c.Equal("Cascadia", truckWithTripPrice.Model)
	c.Equal(2012, truckWithTripPrice.Year)
	c.Equal(TruckType("articulated"), truckWithTripPrice.Type)
	c.Equal(float64(1234), truckWithTripPrice.MaxWeight)
	c.Equal(float64(1234), truckWithTripPrice.MaxVolume)
	c.Equal(float64(1234), truckWithTripPrice.BasePrice)
	c.Equal(float64(1234), truckWithTripPrice.PerRegionPrice)
	c.Equal([]models.Region{"DO-33", "DO-35"}, truckWithTripPrice.Regions)
	c.Equal(float64(3702), truckWithTripPrice.TripPrice)
}
