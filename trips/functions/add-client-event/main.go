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

type clientEventRequest struct {
	*events.APIGatewayProxyRequest
	body    url.Values
	message string
	*models.Trip
}

func (clientEventReq *clientEventRequest) getTrip(ctx context.Context) *events.APIGatewayProxyResponse {
	tripID := clientEventReq.body.Get("trip_id")
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

	clientEventReq.Trip = trip

	return nil
}

func (clientEventReq *clientEventRequest) setNaturalFlowEvent(ctx context.Context) *events.APIGatewayProxyResponse {
	err := clientEventReq.Trip.AddNaturalFlowClientEvent()
	if err != nil {
		return apigateway.NewErrorResponse(403, err)
	}

	// If the event is authorized it means the trip finished and payment needs to be done
	errResponse := clientEventReq.updateTransaction(ctx, transactionModels.TransactionStatusCompleted)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (clientEventReq *clientEventRequest) setCancellationEvent(ctx context.Context) *events.APIGatewayProxyResponse {
	err := clientEventReq.Trip.AddTripCancellationEvent(clientEventReq.message, usersModels.UserTypeClient)
	if err != nil {
		return apigateway.NewErrorResponse(403, err)
	}

	// If the event is authorized it means the trip had not started therefore transaction needs to be reversed
	errResponse := clientEventReq.updateTransaction(ctx, transactionModels.TransactionStatusReversed)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (clientEventReq *clientEventRequest) updateTransaction(ctx context.Context, transactionStatus transactionModels.TransactionStatus) *events.APIGatewayProxyResponse {
	params := url.Values{}

	params.Set("transaction_id", clientEventReq.TransactionID)
	params.Set("status", string(transactionStatus))

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
	clientEventReq := &clientEventRequest{
		APIGatewayProxyRequest: &request,
	}

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	clientEventReq.body = body
	clientEventReq.message = body.Get("message")

	errResponse := clientEventReq.getTrip(ctx)
	if errResponse != nil {
		return errResponse, nil
	}

	eventRoute := tripEvents.EventRoute(body.Get("event_route"))
	if !tripEvents.ValidEventRoutes[eventRoute] {
		return apigateway.NewErrorResponse(400, errInvalidEventRoute), nil
	}

	switch eventRoute {
	case tripEvents.EventRouteNatural:
		errResponse = clientEventReq.setNaturalFlowEvent(ctx)
		if errResponse != nil {
			return errResponse, nil
		}

	case tripEvents.EventRouteCancellation:
		errResponse = clientEventReq.setCancellationEvent(ctx)
		if errResponse != nil {
			return errResponse, nil
		}

	case tripEvents.EventRouteReport:
		err := clientEventReq.Trip.AddReportEvent(clientEventReq.message, usersModels.UserTypeClient)
		if err != nil {
			return apigateway.NewErrorResponse(400, err), nil
		}

	default:
		// If none of the case statements is accomplished the event route is forbidden for client
		return apigateway.NewErrorResponse(403, errForbiddenEventRoute), nil
	}

	err = storage.PutTrip(ctx, clientEventReq.Trip)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, clientEventReq.Trip), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
