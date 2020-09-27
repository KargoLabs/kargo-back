package models

type Event string

const (
	// Natural flow

	// EventTruckSelection event
	EventTruckSelection Event = "truck selection"
	// EventTripAcceptance event
	EventTripAcceptance Event = "trip acceptance"
	// EventTruckDeparture event
	EventTruckDeparture Event = "truck departure"
	// EventArrivalAtOrigin event
	EventArrivalAtOrigin Event = "arrival at origin"
	// EventLoadingProcess event
	EventLoadingProcess Event = "loading process"
	// EventOnWay event
	EventOnWay Event = "on way"
	// EventArrivalAtDestination event
	EventArrivalAtDestination Event = "arrival at destination"
	// EvenetUnloadingProcess event
	EvenetUnloadingProcess Event = "unloading process"
	// EventReceiptConfirmation event
	EventReceiptConfirmation Event = "receipt confirmation"

	// Interruption

	// EventTripDenial event
	EventTripDenial Event = "trip denial"
	// EventTripCancellation event
	EventTripCancellation Event = "trip cancellation"

	// Report

	// EventReport event
	EventReport Event = "report"
)

var (
	// NaturalFlowSteps represent the natural flow of a trip
	NaturalFlowSteps = []Event{
		EventTruckSelection,
		EventTripAcceptance,
		EventTruckDeparture,
		EventArrivalAtOrigin,
		EventLoadingProcess,
		EventOnWay,
		EventArrivalAtDestination,
		EvenetUnloadingProcess,
		EventReceiptConfirmation,
	}

	// NaturalFlowPartnerEvents represent the natural flow events triggered by partner
	NaturalFlowPartnerEvents = map[Event]bool{
		EventTripAcceptance:       true,
		EventTruckDeparture:       true,
		EventArrivalAtOrigin:      true,
		EventLoadingProcess:       true,
		EventOnWay:                true,
		EventArrivalAtDestination: true,
		EvenetUnloadingProcess:    true,
	}

	// NaturalFlowClientEvents represent the natural flow events triggered by client
	NaturalFlowClientEvents = map[Event]bool{
		EventTruckSelection:      true,
		EventReceiptConfirmation: true,
	}
)
