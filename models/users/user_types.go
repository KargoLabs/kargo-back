package models

// UserType is enum of all Kargo User Types
type UserType string

const (
	// UserTypeClient represents client user
	UserTypeClient UserType = "client"
	// UserTypePartner represents partner user
	UserTypePartner UserType = "partner"
)
