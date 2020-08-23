package main

import (
	"context"
	"encoding/base64"
	"kargo-back/clients/storage"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/clients-models"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	body, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	birthDate, err := time.Parse("2006-01-02", body.Get("birth_date"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	client, err := models.NewClient(body.Get("name"), body.Get("document"), birthDate)
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
