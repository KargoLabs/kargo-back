package models

import (
	"errors"
	"kargo-back/shared/normalize"
	"kargo-back/shared/random"
	"time"
)

const (
	clientIDPrefix  = "CLI"
	clientIDBitSize = 16
)

var (
	// ErrMissingName error when name is missing
	ErrMissingName = errors.New("missing name parameter")
	// ErrMissingDocument error when document is missing
	ErrMissingDocument = errors.New("missing document parameter")
	// ErrMissingBirthDate error when birth date is missing
	ErrMissingBirthDate = errors.New("missing birth date parameter")
)

// Client is the struct handler for client
type Client struct {
	ClientID     string    `json:"client_id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	BirthDate    time.Time `json:"birth_date"`
	CreationDate time.Time `json:"creation_date"`
}

// NewClient returns Client structure with given values
func NewClient(name, document string, birthDate time.Time) (*Client, error) {
	if name == "" {
		return nil, ErrMissingName
	}

	if document == "" {
		return nil, ErrMissingDocument
	}

	if (birthDate == time.Time{}) {
		return nil, ErrMissingBirthDate
	}

	return &Client{
		ClientID:     random.GenerateID(clientIDPrefix, clientIDBitSize),
		Name:         normalize.NormalizeName(name),
		Document:     document,
		BirthDate:    birthDate,
		CreationDate: time.Now(),
	}, nil
}
