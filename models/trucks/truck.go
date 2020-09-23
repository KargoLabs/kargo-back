package models

import (
	"errors"
	trips "kargo-back/models/trips"
	"kargo-back/shared/random"
	"time"
)

const (
	// TruckIDPrefix is the prefix for identifying transactions
	TruckIDPrefix = "TRU"
)

var (
	// ErrMissingPartnerID error when partner_id is missing
	ErrMissingPartnerID = errors.New("missing partner id parameter")
	// ErrMissingRegistrationPlate error when registration plate is missing
	ErrMissingRegistrationPlate = errors.New("missing registration id parameter")
	// ErrMissingBrand error when truck brand is missing
	ErrMissingBrand = errors.New("missing brand parameter")
	// ErrMissingModel error when truck_model is missing
	ErrMissingModel = errors.New("missing model parameter")
	// ErrMissingLocation error when truck location is missing
	ErrMissingLocation = errors.New("missing location parameter")
	// ErrMissingAvailable error when truck avaliability is missing
	ErrMissingAvailable = errors.New("missing available parameter")
	// ErrInvalidAvailable error when an invalid truck available param is given
	ErrInvalidAvailable = errors.New("invalid available parameter")
	// ErrMissingTruckType error when an invalid truck type is given
	ErrMissingTruckType = errors.New("missing truck type parameter")
	// ErrInvalidYear error when an invalid truck year is given
	ErrInvalidYear = errors.New("invalid year parameter")
	// ErrInvalidTruckType error when an invalid truck type is given
	ErrInvalidTruckType = errors.New("invalid truck type parameter")
	// ErrInvalidLocation error when an invalid truck location is given
	ErrInvalidLocation = errors.New("invalid truck location")
	// ErrInvalidRegion error when an invalid truck region is given
	ErrInvalidRegion = errors.New("invalid truck region")
	// ErrMissingRegion error when truck regions is missing
	ErrMissingRegion = errors.New("invalid truck region")
	// ErrMissingMaxVolume error when truck max volume is missing
	ErrMissingMaxVolume = errors.New("invalid truck max volume")
	// ErrMissingMaxWeight error when truck max weight is missing
	ErrMissingMaxWeight = errors.New("invalid truck max weight")
	// ErrInvalidMaxVolume error when an invalid truck max volume is given
	ErrInvalidMaxVolume = errors.New("invalid truck max volume")
	// ErrInvalidMaxWeight error when an invalid truck max weight is given
	ErrInvalidMaxWeight = errors.New("invalid truck max weight")
)

// TruckType is the type handler for the different types of trucks
type TruckType string

const (
	// TruckTypeRigid is a truck of type rigid
	TruckTypeRigid TruckType = "rigid"
	// TruckTypeArticulated is a truck of type Articulated
	TruckTypeArticulated TruckType = "articulated"
	// TruckTypeTrailer is a truck of type Trailer
	TruckTypeTrailer TruckType = "trailer"
	// TruckTypeHighway is a truck of type Highway
	TruckTypeHighway TruckType = "highway"
	// TruckTypeRefrigerated is a truck of type refrigerated
	TruckTypeRefrigerated TruckType = "refrigerated"
	// TruckTypeAnimal is a truck which is able to carry animals
	TruckTypeAnimal TruckType = "animal"
)

var (
	// ValidTruckTypes is the map of all the valid truck types
	ValidTruckTypes = map[TruckType]bool{
		TruckTypeRigid:        true,
		TruckTypeArticulated:  true,
		TruckTypeTrailer:      true,
		TruckTypeHighway:      true,
		TruckTypeRefrigerated: true,
		TruckTypeAnimal:       true,
	}
)

// Truck is the struct handler for a truck
type Truck struct {
	TruckID           string         `json:"truck_id"`
	PartnerID         string         `json:"partner_id"`
	RegistrationPlate string         `json:"registration_id"`
	Brand             string         `json:"brand"`
	Model             string         `json:"model"`
	Year              int            `json:"year"`
	CompletedTrips    int            `json:"completed_trips"`
	Available         bool           `json:"available"`
	MaxVolume         float32        `json:"max_volume"`
	MaxWeight         float32        `json:"max_weight"`
	Type              TruckType      `json:"type"`
	Location          trips.Region   `json:"location"`
	CreationDate      time.Time      `json:"creation_date"`
	UpdateDate        time.Time      `json:"update_date"`
	Regions           []trips.Region `json:"regions"`
}

// NewTruck returns a cleaned Truck structure with given values
func NewTruck(truck Truck) (*Truck, error) {
	if truck.PartnerID == "" {
		return nil, ErrMissingPartnerID
	}

	if truck.RegistrationPlate == "" {
		return nil, ErrMissingRegistrationPlate
	}

	if truck.Brand == "" {
		return nil, ErrMissingBrand
	}

	if truck.Model == "" {
		return nil, ErrMissingModel
	}

	if truck.Year <= 0 {
		return nil, ErrInvalidYear
	}

	if truck.Type == "" {
		return nil, ErrMissingTruckType
	}

	if _, ok := ValidTruckTypes[truck.Type]; !ok {
		return nil, ErrInvalidTruckType
	}

	if truck.Location == "" {
		return nil, ErrMissingLocation
	}

	if _, ok := trips.ValidRegions[truck.Location]; !ok {
		return nil, ErrInvalidLocation
	}

	if len(truck.Regions) == 0 {
		return nil, ErrMissingRegion
	}

	for _, region := range truck.Regions {
		_, ok := trips.ValidRegions[region]
		if !ok {
			return nil, ErrInvalidRegion
		}
	}

	if truck.MaxVolume <= 0 {
		return nil, ErrInvalidMaxVolume
	}

	if truck.MaxWeight <= 0 {
		return nil, ErrInvalidMaxWeight
	}

	now := time.Now()

	return &Truck{
		TruckID:           random.GenerateID(TruckIDPrefix, random.StandardBitSize),
		PartnerID:         truck.PartnerID,
		RegistrationPlate: truck.RegistrationPlate,
		Brand:             truck.Brand,
		Model:             truck.Model,
		Year:              truck.Year,
		Type:              truck.Type,
		Location:          truck.Location,
		Regions:           truck.Regions,
		Available:         true,
		CompletedTrips:    0,
		MaxVolume:         truck.MaxVolume,
		MaxWeight:         truck.MaxWeight,
		CreationDate:      now,
		UpdateDate:        now,
	}, nil
}