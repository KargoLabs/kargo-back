package models

import (
	"errors"
	"fmt"
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
	// ErrMissingBirthdate error when birthdate is missing
	ErrMissingBirthdate = errors.New("missing birthdate parameter")
	// ErrMissingPhoneNumber error when phone number is missing
	ErrMissingPhoneNumber = errors.New("missing phone number parameter")
)

// Partner is the struct handler for partner
type Partner struct {
	PartnerID          string    `json:"partner_id"`
	Name               string    `json:"name"`
	Document           string    `json:"document"`
	PhoneNumber        string    `json:"phone_number"`
	Email              string    `json:"email"`
	ProfilePhotoS3Path string    `json:"profile_photo_s3_path"`
	Birthdate          time.Time `json:"birthdate"`
	CreationDate       time.Time `json:"creation_date"`
	UpdateDate         time.Time `json:"update_date"`
}

// NewPartner returns Partner structure with given values
func NewPartner(username, name, document, phoneNumber, email string, birthdate time.Time) (*Partner, error) {
	if name == "" {
		return nil, ErrMissingName
	}

	if document == "" {
		return nil, ErrMissingDocument
	}

	if phoneNumber == "" {
		return nil, ErrMissingPhoneNumber
	}

	if (birthdate == time.Time{}) {
		return nil, ErrMissingBirthdate
	}

	partnerID := random.GetSHA256WithPrefix(PartnerIDPrefix, username)

	return &Partner{
		PartnerID:          partnerID,
		Name:               normalize.NormalizeName(name),
		Document:           document,
		PhoneNumber:        phoneNumber,
		Email:              email,
		ProfilePhotoS3Path: fmt.Sprintf("partners/%s.png", partnerID),
		Birthdate:          birthdate,
		CreationDate:       time.Now(),
		UpdateDate:         time.Now(),
	}, nil
}
