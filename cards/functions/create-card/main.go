package main

import (
	"context"
	models "kargo-back/models/cards"
	clientModel "kargo-back/models/clients"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/cards"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	card, err := models.NewCard(random.GetSHA256WithPrefix(clientModel.ClientIDPrefix, username), body.Get("number"), body.Get("name"), body.Get("csv"), body.Get("year"), body.Get("month"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutCard(ctx, card)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, card), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
