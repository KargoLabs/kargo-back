package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTruckFail(t *testing.T) {
	c := require.New(t)

	truckParam := &Truck{
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		Location:          "DO-33",
	}
	truck, err := NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrMissingPartnerID, err)

	truckParam = &Truck{
		PartnerID: "partner01",
		Brand:     "Freightliner",
		Model:     "Cascadia",
		Year:      2012,
		Type:      "articulated",
		Location:  "DO-33"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrMissingRegistrationPlate, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		Location:          "DO-33"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrMissingBrand, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Year:              2012,
		Type:              "articulated",
		Location:          "DO-33"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrMissingModel, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Type:              "articulated",
		Location:          "DO-33"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrInvalidYear, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Location:          "DO-33"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrMissingTruckType, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "invalid",
		Location:          "DO-33"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrInvalidTruckType, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrMissingLocation, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		Location:          "invalid"}
	truck, err = NewTruck(*truckParam)
	c.Nil(truck)
	c.Equal(ErrInvalidLocation, err)
}

func TestNewTruck(t *testing.T) {
	c := require.New(t)

	truckParam := &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		Location:          "DO-33"}
	truck, err := NewTruck(*truckParam)
	c.Nil(err)
	c.Equal("partner01", truck.PartnerID)
	c.Equal("12345678901", truck.RegistrationPlate)
	c.Equal("Freightliner", truck.Brand)
	c.Equal("Cascadia", truck.Model)
	c.Equal(2012, truck.Year)
	c.Equal("articulated", string(truck.Type))
	c.Equal("DO-33", string(truck.Location))
	c.Equal(0, truck.CompletedTrips)
	c.Equal(true, truck.Available)
	c.NotEmpty(truck.CreationDate)
	c.NotEmpty(truck.UpdateDate)
	c.Equal(truck.CreationDate, truck.UpdateDate)
}
