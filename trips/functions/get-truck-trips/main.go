package main

import (
	"context"
	"errors"
	models "kargo-back/models/trips"
	"kargo-back/shared/apigateway"
	storage "kargo-back/storage/trips"
	truckStorage "kargo-back/storage/trucks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrInvalidFinishedValue is when a non boolean value is given to the finished param
	ErrInvalidFinishedValue = errors.New("invalid value for finished param")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	params := request.QueryStringParameters

	truckID, ok := request.QueryStringParameters["truck_id"]
	if !ok || truckID == "" {
		return apigateway.NewErrorResponse(400, models.ErrMissingTruckID), nil
	}

	truck, err := truckStorage.LoadTruck(ctx, truckID)
	if errors.Is(err, truckStorage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	// Optional parameter, filter by wether a trip has finished or not
	finished := params["finished"]

	query, err := models.NewTripsQuery(truck.TruckID, finished)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	trips, err := storage.QueryTruckTrips(ctx, query)
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
