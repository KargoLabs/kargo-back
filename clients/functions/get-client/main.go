package main

import (
	"context"
	"errors"
	"kargo-back/clients/storage"
	"kargo-back/shared/apigateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingClientID = errors.New("missing client id in query parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	clientID, ok := request.QueryStringParameters["client_id"]
	if !ok || clientID == "" {
		return apigateway.NewErrorResponse(400, errMissingClientID), nil
	}

	client, err := storage.LoadClient(ctx, clientID)
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, client), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
