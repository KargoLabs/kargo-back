package models

import (
	"errors"
	"kargo-back/shared/random"
)

const (
	// AccountIDPrefix is the prefix for identifying accounts
	AccountIDPrefix = "ACC"
)

var (
	// ErrMissingName error when name is missing
	ErrMissingName = errors.New("missing name parameter")
	// ErrMissingPartnerID error when partner_id is missing
	ErrMissingPartnerID = errors.New("missing partner_id parameter")
	// ErrMissingDocument error when document is missing
	ErrMissingDocument = errors.New("missing document parameter")
	// ErrMissingNumber error when account number is missing
	ErrMissingNumber = errors.New("missing number parameter")
)

// Account is the struct handler for an account
type Account struct {
	AccountID string `json:"account_id"`
	PartnerID string `json:"partner_id"`
	Name      string `json:"name"`
	// Could be national ID or RNC
	Document string `json:"document"`
	Number   string `json:"number"`
}

// NewAccount validates and returns an Account structure with given values
func NewAccount(partnerID, name, document, number string) (*Account, error) {
	if partnerID == "" {
		return nil, ErrMissingPartnerID
	}

	if name == "" {
		return nil, ErrMissingName
	}

	if document == "" {
		return nil, ErrMissingDocument
	}

	if number == "" {
		return nil, ErrMissingNumber
	}

	return &Account{
		AccountID: random.GenerateID(AccountIDPrefix, random.StandardBitSize),
		PartnerID: partnerID,
		Name:      name,
		Document:  document,
		Number:    number,
	}, nil
}
