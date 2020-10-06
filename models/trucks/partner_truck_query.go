package models

import (
	"strconv"
)

// PartnerTrucksQuery is the struct for making dynamodb queries for partner trucks
type PartnerTrucksQuery struct {
	PartnerID       string
	FilterAvailable bool
	Available       bool
}

// NewPartnerTruckQuery returns PartnerTrucksQuery structure with given values
func NewPartnerTruckQuery(partnerID, available string) (*PartnerTrucksQuery, error) {
	if partnerID == "" {
		return nil, ErrMissingPartnerID
	}

	filterAvailable := false
	availableOption := false

	if available != "" {
		availableBool, err := strconv.ParseBool(available)
		if err != nil {
			return nil, ErrInvalidAvailable
		}

		availableOption = availableBool
		filterAvailable = true
	}

	return &PartnerTrucksQuery{
		PartnerID:       partnerID,
		FilterAvailable: filterAvailable,
		Available:       availableOption,
	}, nil
}
