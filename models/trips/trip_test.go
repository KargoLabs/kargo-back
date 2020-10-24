package models

import (
	events "kargo-back/models/events"
	users "kargo-back/models/users"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTripFail(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Equal(ErrMissingClientID, err)
	c.Empty(trip)

	trip, err = NewTrip("CLI123", "", "TRU123", "PAY123", 1234, startTime)
	c.Equal(ErrMissingPartnerID, err)
	c.Empty(trip)

	trip, err = NewTrip("CLI123", "PAR123", "", "PAY123", 1234, startTime)
	c.Equal(ErrMissingTruckID, err)
	c.Empty(trip)

	trip, err = NewTrip("CLI123", "PAR123", "TRU123", "", 1234, startTime)
	c.Equal(ErrMissingTransactionID, err)
	c.Empty(trip)

	trip, err = NewTrip("CLI123", "PAR123", "TRU123", "PAY123", -1, startTime)
	c.Equal(ErrInvalidTripPrice, err)
	c.Empty(trip)
}

func TestNewTrip(t *testing.T) {
	c := require.New(t)

	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	c.NotEmpty(trip.TripID)
	c.Equal("CLI123", trip.ClientID)
	c.Equal("PAR123", trip.PartnerID)
	c.Equal("TRU123", trip.TruckID)
	c.Equal("PAY123", trip.TransactionID)
	c.Equal(float64(1234), trip.TripPrice)
	c.Equal(0, trip.NaturalFlowStep)
	c.False(trip.Finished)
	c.Len(trip.EventsHistory, 1)
	c.Equal(events.EventTruckSelection, trip.EventsHistory[0].Event)
	c.Empty(trip.EventsHistory[0].Message)
	c.Equal(users.UserTypeClient, trip.EventsHistory[0].UserType)
	c.NotEmpty(trip.EventsHistory[0].Date)
	c.NotEmpty(trip.UpdateDate)
	c.NotEmpty(trip.CreationDate)
	c.Equal(trip.StartTime, startTime)
}

func TestTrip_AddNaturalFlowPartnerEventFail(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	// This is necessary to arrive to an event that cannot be triggered by partner
	for i := 0; i < 7; i++ {
		err = trip.AddNaturalFlowPartnerEvent()
		c.Nil(err)
	}

	err = trip.AddNaturalFlowPartnerEvent()
	c.Equal(ErrEventNotAuthorized, err)
}

func TestTrip_AddNaturalFlowPartnerEvent(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	err = trip.AddNaturalFlowPartnerEvent()
	c.Nil(err)

	c.Len(trip.EventsHistory, 2)
	c.Equal(events.EventTripAcceptance, trip.EventsHistory[1].Event)
	c.Empty(trip.EventsHistory[1].Message)
	c.Equal(users.UserTypePartner, trip.EventsHistory[1].UserType)
	c.NotEmpty(trip.EventsHistory[1].Date)
	c.Equal(1, trip.NaturalFlowStep)
}

func TestTrip_AddNaturalFlowClientEventFail(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	err = trip.AddNaturalFlowClientEvent()
	c.Equal(ErrEventNotAuthorized, err)
}

func TestTrip_AddNaturalFlowClientEvent(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	// This is necessary to arrive to an event that can be triggered by client
	for i := 0; i < 7; i++ {
		err = trip.AddNaturalFlowPartnerEvent()
		c.Nil(err)
	}

	err = trip.AddNaturalFlowClientEvent()
	c.Nil(err)

	c.Len(trip.EventsHistory, 9)
	c.Equal(events.EventReceiptConfirmation, trip.EventsHistory[8].Event)
	c.Empty(trip.EventsHistory[8].Message)
	c.Equal(users.UserTypeClient, trip.EventsHistory[8].UserType)
	c.NotEmpty(trip.EventsHistory[8].Date)
	c.Equal(8, trip.NaturalFlowStep)
	c.True(trip.Finished)
}

func TestTrip_AddTripDenialEventFail(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	// This is necessary so the denial event cannot be triggered
	err = trip.AddNaturalFlowPartnerEvent()
	c.Nil(err)

	err = trip.AddTripDenialEvent("sorry bro")
	c.Equal(ErrEventNotAuthorized, err)
}

func TestTrip_AddTripDenialEvent(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	err = trip.AddTripDenialEvent("sorry bro")
	c.Nil(err)

	c.Len(trip.EventsHistory, 2)
	c.Equal(events.EventTripDenial, trip.EventsHistory[1].Event)
	c.Equal("sorry bro", trip.EventsHistory[1].Message)
	c.Equal(users.UserTypePartner, trip.EventsHistory[1].UserType)
	c.NotEmpty(trip.EventsHistory[1].Date)
	c.Equal(0, trip.NaturalFlowStep)
	c.True(trip.Finished)
}

func TestTrip_AddTripCancellationEventFail(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	// This is necessary so the cancellation event cannot be triggered
	for i := 0; i < 2; i++ {
		err = trip.AddNaturalFlowPartnerEvent()
		c.Nil(err)
	}

	err = trip.AddTripCancellationEvent("sorry bro", users.UserTypeClient)
	c.Equal(ErrEventNotAuthorized, err)
}

func TestTrip_AddTripCancellationEvent(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	err = trip.AddTripCancellationEvent("sorry bro", users.UserTypeClient)
	c.Nil(err)

	c.Len(trip.EventsHistory, 2)
	c.Equal(events.EventTripCancellation, trip.EventsHistory[1].Event)
	c.Equal("sorry bro", trip.EventsHistory[1].Message)
	c.Equal(users.UserTypeClient, trip.EventsHistory[1].UserType)
	c.NotEmpty(trip.EventsHistory[1].Date)
	c.Equal(0, trip.NaturalFlowStep)
	c.True(trip.Finished)
}

func TestTrip_AddReportEventFail(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	err = trip.AddReportEvent("", users.UserTypeClient)
	c.Equal(ErrMissingMessage, err)
}

func TestTrip_AddReportEvent(t *testing.T) {
	c := require.New(t)
	startTime := time.Now()

	trip, err := NewTrip("CLI123", "PAR123", "TRU123", "PAY123", 1234, startTime)
	c.Nil(err)

	err = trip.AddReportEvent("some minions attacked me", users.UserTypePartner)
	c.Nil(err)

	c.Len(trip.EventsHistory, 2)
	c.Equal(events.EventReport, trip.EventsHistory[1].Event)
	c.Equal("some minions attacked me", trip.EventsHistory[1].Message)
	c.Equal(users.UserTypePartner, trip.EventsHistory[1].UserType)
	c.NotEmpty(trip.EventsHistory[1].Date)
	c.Equal(0, trip.NaturalFlowStep)
	c.False(trip.Finished)
}
