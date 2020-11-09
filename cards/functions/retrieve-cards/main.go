package main

import (
	"context"
	clientModel "kargo-back/models/clients"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/cards"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	cards, err := storage.LoadUserCards(ctx, random.GetSHA256WithPrefix(clientModel.ClientIDPrefix, username))
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, cards), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
