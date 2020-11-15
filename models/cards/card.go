/*Package models WARNING: This is only for development purposes, never save
 * credit card's data directly on the database, better offload the requirements
 * of that by integrating an external payment service such as stripe
 */
package models

import (
	"errors"
	"kargo-back/shared/random"
	"strconv"
	"time"

	creditcard "github.com/durango/go-credit-card"
)

const (
	// CardIDPrefix is the prefix for identifying cards
	CardIDPrefix = "CRD"
)

var (
	// ErrMissingUserID error when user id is missing
	ErrMissingUserID = errors.New("missing user id parameter")
	// ErrMissingUserType error when user type is missing
	ErrMissingUserType = errors.New("missing user type parameter")
	// ErrMissingNumber error when the number param is mssing
	ErrMissingNumber = errors.New("missing number parameter")
	// ErrMissingName error when the card holder name is missing
	ErrMissingName = errors.New("missing name parameter")
	// ErrMissingCSV error when the csv param is mssing
	ErrMissingCSV = errors.New("missing csv parameter")
	// ErrMissingMonth error when the month  is mssing
	ErrMissingMonth = errors.New("missing month parameter")
	// ErrInvalidMonth error when the month  is invalid
	ErrInvalidMonth = errors.New("invalid month parameter")
	// ErrMissingYear error when the year  is mssing
	ErrMissingYear = errors.New("missing year parameter")
	// ErrInvalidYear error when the year param is invalid
	ErrInvalidYear = errors.New("invalid year parameter")
)

// Card is the struct handler for a card
type Card struct {
	CardID         string `json:"card_id"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	LastFourDigits string `json:"last_four_digits"`
	Company        string `json:"company"`
}

// NewCard validates and returns a Card structure with given values
func NewCard(userID, number, name, csv, year, month string) (*Card, error) {
	currYear, currMonth, _ := time.Now().Date()

	if userID == "" {
		return nil, ErrMissingUserID
	}

	if name == "" {
		return nil, ErrMissingName
	}

	if number == "" {
		return nil, ErrMissingNumber
	}

	if csv == "" {
		return nil, ErrMissingCSV
	}

	if year == "" {
		return nil, ErrMissingYear
	}

	if intYear, err := strconv.Atoi(year); err != nil || intYear < currYear {
		return nil, ErrInvalidYear
	}

	if month == "" {
		return nil, ErrMissingMonth
	}

	if intMonth, err := strconv.Atoi(month); err != nil || intMonth < int(currMonth) {
		return nil, ErrInvalidMonth
	}

	card := creditcard.Card{Number: number, Cvv: csv, Month: month, Year: year}

	// Allow test cards
	err := card.Validate(true)
	if err != nil {
		return nil, err
	}

	lastFour, err := card.LastFourDigits()
	if err != nil {
		return nil, err
	}

	err = card.Method()
	if err != nil {
		return nil, err
	}

	return &Card{
		CardID:         random.GenerateID(CardIDPrefix, random.StandardBitSize),
		UserID:         userID,
		Name:           name,
		LastFourDigits: lastFour,
		Company:        card.Company.Long,
	}, nil
}
