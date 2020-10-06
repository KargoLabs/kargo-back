package models

import (
	clients "kargo-back/models/clients"
	partners "kargo-back/models/partners"
)

// UserType is enum of all Kargo User Types
type UserType string

const (
	// UserTypeClient represents client user
	UserTypeClient UserType = "client"
	// UserTypePartner represents partner user
	UserTypePartner UserType = "partner"
)

var (
	// ValidUserTypes is map that represents valid user types
	ValidUserTypes = map[UserType]bool{
		UserTypeClient:  true,
		UserTypePartner: true,
	}

	// UserTypeToPrefix maps user type to according prefix
	UserTypeToPrefix = map[UserType]string{
		UserTypeClient:  clients.ClientIDPrefix,
		UserTypePartner: partners.PartnerIDPrefix,
	}
)
