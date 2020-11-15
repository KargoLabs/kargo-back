package main

import (
	"context"
	"errors"
	clientModel "kargo-back/models/clients"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/cards"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrMissingCardID when card_id parameter is missing
	ErrMissingCardID = errors.New("missing card_id parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	cardID := request.QueryStringParameters["card_id"]

	if cardID == "" {
		return apigateway.NewErrorResponse(400, ErrMissingCardID), nil
	}

	card, err := storage.DeleteCard(ctx, random.GetSHA256WithPrefix(clientModel.ClientIDPrefix, username), cardID)
	if err == storage.ErrCardNotFound {
		return apigateway.NewErrorResponse(404, storage.ErrCardNotFound), nil
	}
	if err != nil {
		return apigateway.NewErrorResponse(500, err), nil
	}

	return apigateway.NewJSONResponse(200, card), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
