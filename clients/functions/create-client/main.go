package main

import (
	"context"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/clients-models"
	storage "kargo-back/storage/clients"
	"net/url"
	"time"

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

	birthDate, err := time.Parse("2006-01-02", body.Get("birthdate"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	client, err := models.NewClient(username, body.Get("name"), body.Get("document"), birthDate)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutClient(ctx, client)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, client), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
