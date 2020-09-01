package models

import (
	"errors"
	"kargo-back/shared/normalize"
	"kargo-back/shared/random"
	"time"
)

const (
	// ClientIDPrefix is the prefix for identifying client IDs
	ClientIDPrefix = "CLI"
)

var (
	// ErrMissingName error when name is missing
	ErrMissingName = errors.New("missing name parameter")
	// ErrMissingDocument error when document is missing
	ErrMissingDocument = errors.New("missing document parameter")
	// ErrMissingBirthDate error when birth date is missing
	ErrMissingBirthDate = errors.New("missing birth date parameter")
	// ErrMissingUsername error when username is missing
	ErrMissingUsername = errors.New("missing username parameter")
)

// Client is the struct handler for client
type Client struct {
	ClientID     string    `json:"client_id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	BirthDate    time.Time `json:"birthdate"`
	CreationDate time.Time `json:"creation_date"`
}

// NewClient returns Client structure with given values
func NewClient(username, name, document string, birthDate time.Time) (*Client, error) {
	if username == "" {
		return nil, ErrMissingName
	}

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
		ClientID:     random.GetSHA256WithPrefix(ClientIDPrefix, username),
		Name:         normalize.NormalizeName(name),
		Document:     document,
		BirthDate:    birthDate,
		CreationDate: time.Now(),
	}, nil
}
