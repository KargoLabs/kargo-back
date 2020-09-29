package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	storage "kargo-back/storage/trucks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrMissingTruckID is when the partner_id query parameter is missing
	ErrMissingTruckID = errors.New("missing truck_id query param")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	params := request.QueryStringParameters

	truckID, ok := params["truck_id"]

	if !ok {
		return apigateway.NewErrorResponse(403, ErrMissingTruckID), nil
	}

	truck, err := storage.LoadTruck(ctx, truckID)

	if errors.Is(err, storage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, truck), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
