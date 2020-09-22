package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	storage "kargo-back/storage/trucks"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	truck, err := storage.LoadTruck(ctx, body.Get("truck_id"))
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
