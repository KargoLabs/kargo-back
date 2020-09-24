package models

import (
	models "kargo-back/models/trips"
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
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err := truckParam.ValidateTruck()
	c.Equal(ErrMissingPartnerID, err)

	truckParam = &Truck{
		PartnerID: "partner01",
		Brand:     "Freightliner",
		Model:     "Cascadia",
		Year:      2012,
		Type:      "articulated",
		MaxWeight: 1234,
		MaxVolume: 1234,
		Regions:   []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrMissingRegistrationPlate, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrMissingBrand, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrMissingModel, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrInvalidYear, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrMissingTruckType, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "invalid",
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrInvalidTruckType, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234}
	err = truckParam.ValidateTruck()
	c.Equal(ErrMissingRegion, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"invalid"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrInvalidRegion, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrInvalidMaxWeight, err)

	truckParam = &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		Regions:           []models.Region{"DO-33"}}
	err = truckParam.ValidateTruck()
	c.Equal(ErrInvalidMaxVolume, err)
}

func TestNewTruck(t *testing.T) {
	c := require.New(t)

	truck := &Truck{
		PartnerID:         "partner01",
		RegistrationPlate: "12345678901",
		Brand:             "Freightliner",
		Model:             "Cascadia",
		Year:              2012,
		Type:              "articulated",
		MaxWeight:         1234,
		MaxVolume:         1234,
		Regions:           []models.Region{"DO-33"}}

	err := truck.ValidateTruck()
	c.Nil(err)
	c.True(truck.Available)
	c.NotEmpty(truck.TruckID)
	c.NotEmpty(truck.CreationDate)
	c.NotEmpty(truck.UpdateDate)
}
