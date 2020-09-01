package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/clients-models"
	"kargo-back/shared/random"
	storage "kargo-back/storage/clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	client, err := storage.LoadClient(ctx, random.GetSHA256WithPrefix(models.ClientIDPrefix, username))
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
