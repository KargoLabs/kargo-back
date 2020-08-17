package main

import (
	"context"
	"errors"
	"fmt"
	"kargo-back/clients/storage"
	"kargo-back/shared/apigateway"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func logAndReturnError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())

	return apigateway.NewJSONResponse(500, err)
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return logAndReturnError(err)
	}

	client, err := storage.LoadClient(ctx, body.Get("client_id"))
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewJSONResponse(404, err)
	}

	if err != nil {
		return logAndReturnError(err)
	}

	return apigateway.NewJSONResponse(200, client)
}

func main() {
	lambda.Start(apiGatewayHandler)
}
