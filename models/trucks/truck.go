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
	ErrMissingRegistrationPlate = errors.New("missing registration plate parameter")
	// ErrMissingBrand error when truck brand is missing
	ErrMissingBrand = errors.New("missing brand parameter")
	// ErrMissingModel error when truck_model is missing
	ErrMissingModel = errors.New("missing model parameter")
	// ErrMissingLocation error when truck location is missing
	ErrMissingLocation = errors.New("missing location parameter")
	// ErrInvalidAvailable error when an invalid truck available param is given
	ErrInvalidAvailable = errors.New("invalid available parameter")
	// ErrMissingTruckType error when an invalid truck type is given
	ErrMissingTruckType = errors.New("missing truck type parameter")
	// ErrInvalidBasePrice error when base price is invalid
	ErrInvalidBasePrice = errors.New("invalid base price type parameter")
	// ErrInvalidPerRegionPrice error when per region price is invalid
	ErrInvalidPerRegionPrice = errors.New("invalid per region price type parameter")
	// ErrInvalidYear error when an invalid truck year is given
	ErrInvalidYear = errors.New("invalid year parameter")
	// ErrInvalidTruckType error when an invalid truck type is given
	ErrInvalidTruckType = errors.New("invalid truck type parameter")
	// ErrInvalidRegion error when an invalid truck region is given
	ErrInvalidRegion = errors.New("invalid truck region")
	// ErrMissingRegion error when truck regions is missing
	ErrMissingRegion = errors.New("invalid truck region")
	// ErrInvalidMaxVolume error when an invalid truck max volume is given
	ErrInvalidMaxVolume = errors.New("invalid truck max volume")
	// ErrInvalidMaxWeight error when an invalid truck max weight is given
	ErrInvalidMaxWeight = errors.New("invalid truck max weight")
	// ErrInvalidCompletedTrips error when an invalid truck completed_trips is given
	ErrInvalidCompletedTrips = errors.New("invalid completed_trips")
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
	MaxVolume         float64        `json:"max_volume"`
	MaxWeight         float64        `json:"max_weight"`
	BasePrice         float64        `json:"base_price"`
	PerRegionPrice    float64        `json:"per_region_price"`
	Type              TruckType      `json:"truck_type"`
	CreationDate      time.Time      `json:"creation_date"`
	UpdateDate        time.Time      `json:"update_date"`
	Regions           []trips.Region `json:"regions"`
}

// ValidateTruck returns a cleaned Truck structure with given values
func (truck *Truck) ValidateTruck() error {
	if truck.PartnerID == "" {
		return ErrMissingPartnerID
	}

	if truck.RegistrationPlate == "" {
		return ErrMissingRegistrationPlate
	}

	if truck.Brand == "" {
		return ErrMissingBrand
	}

	if truck.Model == "" {
		return ErrMissingModel
	}

	if truck.Year <= 0 {
		return ErrInvalidYear
	}

	if truck.Type == "" {
		return ErrMissingTruckType
	}

	if _, ok := ValidTruckTypes[truck.Type]; !ok {
		return ErrInvalidTruckType
	}

	if len(truck.Regions) == 0 {
		return ErrMissingRegion
	}

	for _, region := range truck.Regions {
		_, ok := trips.ValidRegions[region]
		if !ok {
			return ErrInvalidRegion
		}
	}

	if truck.MaxVolume <= 0 {
		return ErrInvalidMaxVolume
	}

	if truck.MaxWeight <= 0 {
		return ErrInvalidMaxWeight
	}

	if truck.BasePrice <= 0 {
		return ErrInvalidBasePrice
	}

	if truck.PerRegionPrice <= 0 {
		return ErrInvalidPerRegionPrice
	}

	truck.Available = true
	truck.TruckID = random.GenerateID(TruckIDPrefix, random.StandardBitSize)
	truck.CreationDate = time.Now()
	truck.UpdateDate = time.Now()

	return nil
}
