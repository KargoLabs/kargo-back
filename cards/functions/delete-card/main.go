package main

import (
	"context"
	"errors"
	clientModel "kargo-back/models/clients"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/cards"
	"net/url"

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

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	if body.Get("card_id") == "" {
		return apigateway.NewErrorResponse(400, ErrMissingCardID), nil
	}

	err = storage.DeleteCard(ctx, random.GetSHA256WithPrefix(clientModel.ClientIDPrefix, username), body.Get("card_id"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	return apigateway.NewJSONResponse(204, nil), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
