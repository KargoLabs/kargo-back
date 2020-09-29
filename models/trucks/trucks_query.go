package models

import (
	"errors"
	trips "kargo-back/models/trips"
	"kargo-back/shared/converter"
)

var (
	// ErrInvalidOrigin error when an invalid origin is given
	ErrInvalidOrigin = errors.New("invalid origin parameter")
	// ErrInvalidDestination error when an invalid destination is given
	ErrInvalidDestination = errors.New("invalid destination parameter")
	// ErrInvalidVolume error when truck max volume is invalid
	ErrInvalidVolume = errors.New("invalid volume parameter")
	// ErrInvalidWeight error when truck max weight is invalid
	ErrInvalidWeight = errors.New("invalid weight parameter")
)

type TrucksQuery struct {
	TruckType   string
	Origin      string
	Destination string
	Weight      float64
	Volume      float64
}

// NewTrucksQuery is the struct for making dynamodb queries for trucks
func NewTrucksQuery(truckType, weightString, volumeString string, origin, destination trips.Region) (*TrucksQuery, error) {
	if !ValidTruckTypes[TruckType(truckType)] {
		return nil, ErrInvalidTruckType
	}

	if !trips.ValidRegions[origin] {
		return nil, ErrInvalidOrigin
	}

	if !trips.ValidRegions[destination] {
		return nil, ErrInvalidDestination
	}

	weight, ok := converter.ConvertToGreaterThanZeroFloat(weightString)
	if !ok {
		return nil, ErrInvalidWeight
	}

	volume, ok := converter.ConvertToGreaterThanZeroFloat(volumeString)
	if !ok {
		return nil, ErrInvalidVolume
	}

	return &TrucksQuery{
		TruckType:   truckType,
		Origin:      string(origin),
		Destination: string(destination),
		Weight:      weight,
		Volume:      volume,
	}, nil
}

// PartnerTrucksQuery is the struct for making dynamodb queries for partner trucks
type PartnerTrucksQuery struct {
	PartnerID       string
	FilterAvailable bool
	Available       bool
}
