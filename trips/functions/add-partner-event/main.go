package main

import (
	"context"
	"errors"
	tripEvents "kargo-back/models/events"
	transactionModels "kargo-back/models/transactions"
	models "kargo-back/models/trips"
	usersModels "kargo-back/models/users"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	lambdaLibrary "kargo-back/shared/lambda"
	storage "kargo-back/storage/trips"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	updateTransactionLambdaName = environment.GetString("UPDATE_TRANSACTION_LAMBDA_NAME", "transactions_update-transaction")

	errMissingTripID       = errors.New("missing trip id in body parameter")
	errInvalidEventRoute   = errors.New("invalid event route parameter")
	errForbiddenEventRoute = errors.New("given event route is forbidden")
	errTripFinished        = errors.New("events cannot be added to a finished trip")
)

type partnerEventRequest struct {
	*events.APIGatewayProxyRequest
	body    url.Values
	message string
	*models.Trip
}

func (partnerEventReq *partnerEventRequest) getTrip(ctx context.Context) *events.APIGatewayProxyResponse {
	tripID := partnerEventReq.body.Get("trip_id")
	if tripID == "" {
		return apigateway.NewErrorResponse(400, errMissingTripID)
	}

	trip, err := storage.LoadTrip(ctx, tripID)
	if errors.Is(err, storage.ErrTripNotFound) {
		return apigateway.NewErrorResponse(404, err)
	}

	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	if trip.Finished {
		return apigateway.NewErrorResponse(403, errTripFinished)
	}

	partnerEventReq.Trip = trip

	return nil
}

func (partnerEventReq *partnerEventRequest) setDenialEvent(ctx context.Context) *events.APIGatewayProxyResponse {
	err := partnerEventReq.Trip.AddTripDenialEvent(partnerEventReq.message)
	if err != nil {
		return apigateway.NewErrorResponse(403, err)
	}

	// If the event is authorized it means the trip had not started therefore transaction needs to be reversed
	errResponse := partnerEventReq.reverseTransaction(ctx)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (partnerEventReq *partnerEventRequest) setCancellationEvent(ctx context.Context) *events.APIGatewayProxyResponse {
	err := partnerEventReq.Trip.AddTripCancellationEvent(partnerEventReq.message, usersModels.UserTypePartner)
	if err != nil {
		return apigateway.NewErrorResponse(403, err)
	}

	// If the event is authorized it means the trip had not started therefore transaction needs to be reversed
	errResponse := partnerEventReq.reverseTransaction(ctx)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (partnerEventReq *partnerEventRequest) reverseTransaction(ctx context.Context) *events.APIGatewayProxyResponse {
	params := url.Values{}

	params.Set("transaction_id", partnerEventReq.TransactionID)
	params.Set("status", string(transactionModels.TransactionStatusReversed))

	lambdaResponse, err := lambdaLibrary.InvokeAPIGatewayWithURLEncodedParams(ctx, updateTransactionLambdaName, params)
	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	if lambdaResponse.StatusCode != 200 {
		// It is a 500 because supposedly all inputs were correct before triggering the transaction lambda
		lambdaResponse.StatusCode = 500

		return lambdaResponse
	}

	return nil
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	partnerEventReq := &partnerEventRequest{
		APIGatewayProxyRequest: &request,
	}

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	partnerEventReq.body = body
	partnerEventReq.message = body.Get("message")

	errResponse := partnerEventReq.getTrip(ctx)
	if errResponse != nil {
		return errResponse, nil
	}

	eventRoute := tripEvents.EventRoute(body.Get("event_route"))
	if !tripEvents.ValidEventRoutes[eventRoute] {
		return apigateway.NewErrorResponse(400, errInvalidEventRoute), nil
	}

	switch eventRoute {
	case tripEvents.EventRouteNatural:
		err := partnerEventReq.Trip.AddNaturalFlowPartnerEvent()
		if err != nil {
			return apigateway.NewErrorResponse(403, err), nil
		}

	case tripEvents.EventRouteDenial:
		errResponse = partnerEventReq.setDenialEvent(ctx)
		if errResponse != nil {
			return errResponse, nil
		}

	case tripEvents.EventRouteCancellation:
		errResponse = partnerEventReq.setCancellationEvent(ctx)
		if errResponse != nil {
			return errResponse, nil
		}

	case tripEvents.EventRouteReport:
		err := partnerEventReq.Trip.AddReportEvent(partnerEventReq.message, usersModels.UserTypePartner)
		if err != nil {
			return apigateway.NewErrorResponse(400, err), nil
		}

	default:
		// If none of the case statements is accomplished the event route is forbidden for client
		return apigateway.NewErrorResponse(403, errForbiddenEventRoute), nil
	}

	err = storage.PutTrip(ctx, partnerEventReq.Trip)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, partnerEventReq.Trip), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
