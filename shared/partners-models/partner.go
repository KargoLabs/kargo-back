package models

import (
	"errors"
	"kargo-back/shared/normalize"
	"kargo-back/shared/random"
	"time"
)

const (
	partnerIDPrefix  = "PAR"
	partnerIDBitSize = 16
)

var (
	// ErrMissingName error when name is missing
	ErrMissingName = errors.New("missing name parameter")
	// ErrMissingDocument error when document is missing
	ErrMissingDocument = errors.New("missing document parameter")
	// ErrMissingBirthDate error when birth date is missing
	ErrMissingBirthDate = errors.New("missing birth date parameter")
)

// Partner is the struct handler for partner
type Partner struct {
	PartnerID    string    `json:"partner_id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	BirthDate    time.Time `json:"birth_date"`
	CreationDate time.Time `json:"creation_date"`
}

// NewPartner returns Partner structure with given values
func NewPartner(name, document string, birthDate time.Time) (*Partner, error) {
	if name == "" {
		return nil, ErrMissingName
	}

	if document == "" {
		return nil, ErrMissingDocument
	}

	if (birthDate == time.Time{}) {
		return nil, ErrMissingBirthDate
	}

	return &Partner{
		PartnerID:    random.GenerateID(partnerIDPrefix, partnerIDBitSize),
		Name:         normalize.NormalizeName(name),
		Document:     document,
		CreationDate: time.Now(),
	}, nil
}
