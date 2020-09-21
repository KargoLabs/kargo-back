package models

import (
	"errors"
	"kargo-back/shared/random"
	trips "kargo-back/shared/trips-models"
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
	// ErrMissingMileague error when truck mileague is missing
	ErrMissingMileague = errors.New("missing mileague parameter")
	// ErrMissingLocation error when truck location is missing
	ErrMissingLocation = errors.New("missing location parameter")
	// ErrMissingAvailable error when truck avaliability is missing
	ErrMissingAvailable = errors.New("missing available parameter")
	// ErrMissingTruckType error when an invalid truck type is given
	ErrMissingTruckType = errors.New("missing truck type parameter")
	// ErrInvalidYear error when an invalid truck year is given
	ErrInvalidYear = errors.New("invalid year parameter")
	// ErrInvalidMileague error when an invalid truck mileague is given
	ErrInvalidMileague = errors.New("invalid mileague parameter")
	// ErrInvalidTruckType error when an invalid truck type is given
	ErrInvalidTruckType = errors.New("invalid truck type parameter")
	// ErrInvalidLocation error when an invalid truck location is given
	ErrInvalidLocation = errors.New("invalid truck location")
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
)

var (
	// ValidTruckTypes is the map of all the valid truck types
	ValidTruckTypes = map[TruckType]bool{
		TruckTypeRigid:       true,
		TruckTypeArticulated: true,
		TruckTypeTrailer:     true,
		TruckTypeHighway:     true,
	}
)

// Truck is the struct handler for a truck
type Truck struct {
	TruckID           string       `json:"truck_id"`
	PartnerID         string       `json:"partner_id"`
	RegistrationPlate string       `json:"registration_id"`
	Brand             string       `json:"brand"`
	Model             string       `json:"model"`
	Year              int          `json:"year"`
	Mileague          int          `json:"mileague"`
	CompletedTrips    int          `json:"completed_trips"`
	Available         bool         `json:"available"`
	Type              TruckType    `json:"type"`
	Location          trips.Region `json:"location"`
	CreationDate      time.Time    `json:"creation_date"`
	UpdateDate        time.Time    `json:"update_date"`
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

	if truck.Mileague <= 0 {
		return nil, ErrInvalidMileague
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

	now := time.Now()

	return &Truck{
		TruckID:           random.GenerateID(TruckIDPrefix, random.StandardBitSize),
		PartnerID:         truck.PartnerID,
		RegistrationPlate: truck.RegistrationPlate,
		Brand:             truck.Brand,
		Model:             truck.Model,
		Year:              truck.Year,
		Mileague:          truck.Mileague,
		Type:              truck.Type,
		Location:          truck.Location,
		Available:         true,
		CompletedTrips:    0,
		CreationDate:      now,
		UpdateDate:        now,
	}, nil
}
