package models

import (
	"errors"
	"strconv"
)

var (
	// ErrMissingValue error when key value is missing
	ErrMissingValue = errors.New("missing key parameter")
	// ErrInvalidFinished error when the finished parameter is invalid
	ErrInvalidFinished = errors.New("invalid finished value")
)

// TripsQuery is the struct for making trip-related queries based on the indexed key field
type TripsQuery struct {
	Value          string
	FilterFinished bool
	Finished       bool
}

// NewTripsQuery returns PartnerTrucksQuery structure with given values
func NewTripsQuery(value, finished string) (*TripsQuery, error) {
	if value == "" {
		return nil, ErrMissingValue
	}

	filterFinished := false
	finishedOption := false

	if finished != "" {
		finishedBool, err := strconv.ParseBool(finished)
		if err != nil {
			return nil, ErrInvalidFinished
		}

		finishedOption = finishedBool
		filterFinished = true
	}

	return &TripsQuery{
		Value:          value,
		FilterFinished: filterFinished,
		Finished:       finishedOption,
	}, nil
}
