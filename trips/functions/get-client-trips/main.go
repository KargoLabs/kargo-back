package main

import (
	"context"
	"errors"
	clientModels "kargo-back/models/clients"
	models "kargo-back/models/trips"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	clientStorage "kargo-back/storage/clients"
	storage "kargo-back/storage/trips"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrInvalidFinishedValue is when a non boolean value is given to the finished param
	ErrInvalidFinishedValue = errors.New("invalid value for finished param")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	client, err := clientStorage.LoadClient(ctx, random.GetSHA256WithPrefix(clientModels.ClientIDPrefix, username))
	if errors.Is(err, clientStorage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	params := request.QueryStringParameters

	// Optional parameter, filter by wether a trip has finished or not
	finished := params["finished"]

	query, err := models.NewTripsQuery(client.ClientID, finished)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	trips, err := storage.QueryClientTrips(ctx, query)
	if errors.Is(err, storage.ErrTripNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, trips), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
