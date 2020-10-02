package models

// PartnerTrucksQuery is the struct for making dynamodb queries for partner trucks
type PartnerTrucksQuery struct {
	PartnerID       string
	FilterAvailable bool
	Available       bool
}
