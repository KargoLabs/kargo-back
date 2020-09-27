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

// TransactionStatus is the type handler for payment TransactionStatus.
type TransactionStatus string

const (
	// TransactionStatusOnHold represents when a transaction is on hold
	TransactionStatusOnHold TransactionStatus = "on hold"
	// TransactionStatusReversed represents when a transaction is reversed
	TransactionStatusReversed TransactionStatus = "reversed"
	// TransactionStatusCompleted represents when a transaction is completed
	TransactionStatusCompleted TransactionStatus = "completed"
)

var (
	// ValidTransactionStatuses represents a map of all the valid transaction status
	ValidTransactionStatuses = map[TransactionStatus]bool{
		TransactionStatusOnHold:    true,
		TransactionStatusReversed:  true,
		TransactionStatusCompleted: true,
	}
)

// Transaction is the struct handler for transaction
type Transaction struct {
	TransactionID string            `json:"transaction_id"`
	PartnerID     string            `json:"partner_id"`
	ClientID      string            `json:"client_id"`
	Amount        float64           `json:"amount"`
	Status        TransactionStatus `json:"status"`
	Date          time.Time         `json:"date"`
}

// NewTransaction returns a Transaction structure with given values
func NewTransaction(clientID, partnerID string, amount float64) (*Transaction, error) {
	if clientID == "" {
		return nil, ErrMissingClientID
	}

	if partnerID == "" {
		return nil, ErrMissingPartnerID
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	return &Transaction{
		ClientID:  clientID,
		PartnerID: partnerID,
		Amount:    amount,
		Status:    TransactionStatusOnHold,
		Date:      time.Now(),
		TransactionID: random.GenerateID(TransactionIDPrefix,
			random.StandardBitSize),
	}, nil
}
