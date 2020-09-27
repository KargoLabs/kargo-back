package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	storage "kargo-back/storage/trips"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingTripID = errors.New("missing trip id in query parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	tripID, ok := request.QueryStringParameters["trip_id"]
	if !ok || tripID == "" {
		return apigateway.NewErrorResponse(400, errMissingTripID), nil
	}

	trip, err := storage.LoadTrip(ctx, tripID)
	if errors.Is(err, storage.ErrTripNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, trip), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
