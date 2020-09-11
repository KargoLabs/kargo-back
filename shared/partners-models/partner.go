package models

import (
	"errors"
	"kargo-back/shared/normalize"
	"kargo-back/shared/random"
	"time"
)

const (
	// PartnerIDPrefix is the prefix for identifying partners
	PartnerIDPrefix = "PAR"
)

var (
	// ErrMissingName error when name is missing
	ErrMissingName = errors.New("missing name parameter")
	// ErrMissingDocument error when document is missing
	ErrMissingDocument = errors.New("missing document parameter")
	// ErrMissingBirthdate error when birth date is missing
	ErrMissingBirthdate = errors.New("missing birth date parameter")
)

// Partner is the struct handler for partner
type Partner struct {
	PartnerID    string    `json:"partner_id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	Birthdate    time.Time `json:"birth_date"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"update_date"`
}

// NewPartner returns Partner structure with given values
func NewPartner(username, name, document string, birthDate time.Time) (*Partner, error) {
	if name == "" {
		return nil, ErrMissingName
	}

	if document == "" {
		return nil, ErrMissingDocument
	}

	if (birthDate == time.Time{}) {
		return nil, ErrMissingBirthdate
	}

	return &Partner{
		PartnerID:    random.GetSHA256WithPrefix(PartnerIDPrefix, username),
		Name:         normalize.NormalizeName(name),
		Document:     document,
		CreationDate: time.Now(),
		UpdateDate:   time.Now(),
	}, nil
}
