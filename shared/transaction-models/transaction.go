package models

import (
	"errors"
	"kargo-back/shared/random"
	"time"
)

const (
	// TransactionIDPrefix is the prefix for identifying transactions
	TransactionIDPrefix = "PAY"
)

var (
	// ErrMissingClientID error when client_id is missing
	ErrMissingClientID = errors.New("missing client id parameter")
	// ErrMissingPartnerID error when the partner_id is missing
	ErrMissingPartnerID = errors.New("missing partner id parameter")
	// ErrInvalidAmount error when an invalid amount is given
	ErrInvalidAmount = errors.New("invalid amount parameter")
)

// PaymentStatus is the type handler for payment statuses.
type PaymentStatus int

const (
	// Started is when a trip has started
	Started PaymentStatus = iota
	// InProgress is when a trip is in progress
	InProgress
	// Completed is when a trip has been completed
	Completed
)

// Transaction is the struct handler for transaction
type Transaction struct {
	TransactionID string        `json:"transaction_id"`
	PartnerID     string        `json:"partner_id"`
	ClientID      string        `json:"client_id"`
	Amount        int           `json:"amount"`
	Status        PaymentStatus `json:"status"`
	Date          time.Time     `json:"date"`
}

// NewTransaction returns a Transaction structure with given values
func NewTransaction(clientID, partnerID string, amount int) (*Transaction,
	error) {

	if clientID == "" {
		return nil, ErrMissingClientID
	}

	if partnerID == "" {
		return nil, ErrMissingPartnerID
	}

	if amount < 0 {
		return nil, ErrInvalidAmount
	}

	return &Transaction{
		TransactionID: random.GetSHA256WithPrefix(TransactionIDPrefix,
			time.Now().String()),
	}, nil
}
