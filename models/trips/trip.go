package models

import (
	"errors"
	events "kargo-back/models/events"
	"kargo-back/shared/random"
	"time"
)

const (
	// TripIDPrefix is the prefix for identifying trips
	TripIDPrefix = "TRI"
)

var (
	// ErrMissingClientID error when client_id is missing
	ErrMissingClientID = errors.New("missing client id parameter")
	// ErrMissingPartnerID error when the partner_id is missing
	ErrMissingPartnerID = errors.New("missing partner id parameter")
	// ErrMissingTruckID error when the truck_id is missing
	ErrMissingTruckID = errors.New("missing truck id parameter")
	// ErrMissingTransactionID error when the transaction_id is missing
	ErrMissingTransactionID = errors.New("missing transaction id parameter")
	// ErrMissingMessage error when the message is missing
	ErrMissingMessage = errors.New("missing message parameter")
	// ErrInvalidTripPrice error when an invalid trip price is given
	ErrInvalidTripPrice = errors.New("invalid trip price parameter")
	// ErrEventNotAuthorized error when user is not authorized for event
	ErrEventNotAuthorized = errors.New("user is not authorized for event")
)

// Trip is the struct handler for a trip
type Trip struct {
	TripID          string                `json:"trip_id"`
	ClientID        string                `json:"client_id"`
	PartnerID       string                `json:"partner_id"`
	TruckID         string                `json:"truck_id"`
	TransactionID   string                `json:"transaction_id"`
	TripPrice       float64               `json:"trip_price"`
	NaturalFlowStep int                   `json:"natural_flow_step"`
	Finished        bool                  `json:"finished"`
	EventsHistory   []events.EventHistory `json:"events_history"`
	CreationDate    time.Time             `json:"creation_date"`
	UpdateDate      time.Time             `json:"update_date"`
}

// NewTrip returns truck struct with given input
func NewTrip(clientID, partnerID, truckID, transactionID string, tripPrice float64) (*Trip, error) {
	if clientID == "" {
		return nil, ErrMissingClientID
	}

	if partnerID == "" {
		return nil, ErrMissingPartnerID
	}

	if truckID == "" {
		return nil, ErrMissingTruckID
	}

	if transactionID == "" {
		return nil, ErrMissingTransactionID
	}

	if tripPrice <= 0 {
		return nil, ErrInvalidTripPrice
	}

	// NaturalFlowStep starts with 0 by default which is the first index for the NaturalFlowSteps slice
	// This field will be used for knowing in which step the trip is in without needing to check the event history
	return &Trip{
		TripID:        random.GenerateID(TripIDPrefix, random.StandardBitSize),
		ClientID:      clientID,
		PartnerID:     partnerID,
		TruckID:       truckID,
		TransactionID: transactionID,
		TripPrice:     tripPrice,
		EventsHistory: []events.EventHistory{
			events.EventHistory{
				Event: events.EventTruckSelection,
				Date:  time.Now(),
			},
		},
		CreationDate: time.Now(),
		UpdateDate:   time.Now(),
	}, nil
}

// AddNaturalFlowPartnerEvent adds partner triggered natural flow event
func (trip *Trip) AddNaturalFlowPartnerEvent() error {
	event := events.NaturalFlowSteps[trip.NaturalFlowStep+1]
	if !events.NaturalFlowPartnerEvents[event] {
		return ErrEventNotAuthorized
	}

	trip.NaturalFlowStep++

	trip.EventsHistory = append(trip.EventsHistory, events.EventHistory{
		Event: event,
		Date:  time.Now(),
	})

	trip.UpdateDate = time.Now()

	return nil
}

// AddNaturalFlowClientEvent adds client triggered natural flow event
func (trip *Trip) AddNaturalFlowClientEvent() error {
	event := events.NaturalFlowSteps[trip.NaturalFlowStep+1]
	if !events.NaturalFlowClientEvents[event] {
		return ErrEventNotAuthorized
	}

	if event == events.EventReceiptConfirmation {
		trip.Finished = true
	}

	trip.NaturalFlowStep++

	trip.EventsHistory = append(trip.EventsHistory, events.EventHistory{
		Event: event,
		Date:  time.Now(),
	})

	trip.UpdateDate = time.Now()

	return nil
}

// AddTripDenialEvent adds trip denial event
// Can just be triggered by partner
func (trip *Trip) AddTripDenialEvent(message string) error {
	// Partner can just deny a trip that has not been accepted
	if trip.NaturalFlowStep != 0 {
		return ErrEventNotAuthorized
	}

	trip.EventsHistory = append(trip.EventsHistory, events.EventHistory{
		Event:   events.EventTripDenial,
		Message: message,
		Date:    time.Now(),
	})

	trip.Finished = true
	trip.UpdateDate = time.Now()

	return nil
}

// AddTripCancellationEvent adds trip cancellation event
// Can be triggered by partner and client
func (trip *Trip) AddTripCancellationEvent(message string) error {
	// Trip can just be cancelled if truck departure event has not happened
	if trip.NaturalFlowStep > 1 {
		return ErrEventNotAuthorized
	}

	trip.EventsHistory = append(trip.EventsHistory, events.EventHistory{
		Event:   events.EventTripCancellation,
		Message: message,
		Date:    time.Now(),
	})

	trip.Finished = true
	trip.UpdateDate = time.Now()

	return nil
}

// AddReportEvent adds report event
// Can be triggered by partner and client and message is mandatory
func (trip *Trip) AddReportEvent(message string) error {
	if message == "" {
		return ErrMissingMessage
	}

	trip.EventsHistory = append(trip.EventsHistory, events.EventHistory{
		Event:   events.EventReport,
		Message: message,
		Date:    time.Now(),
	})

	trip.UpdateDate = time.Now()

	return nil
}
