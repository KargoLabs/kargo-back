package main

import (
	"context"
	"errors"
	"fmt"
	"kargo-back/clients/storage"
	"kargo-back/shared/apigateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingClientID = errors.New("missing client id in query parameter")
)

func logAndReturnError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())

	return apigateway.NewJSONResponse(500, err)
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	clientID, ok := request.PathParameters["client_id"]
	if !ok || clientID == "" {
		return apigateway.NewJSONResponse(400, errMissingClientID), nil
	}

	client, err := storage.LoadClient(ctx, clientID)
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewJSONResponse(404, err), nil
	}

	if err != nil {
		return logAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, client), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
