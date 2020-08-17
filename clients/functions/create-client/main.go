package main

import (
	"context"
	"fmt"
	"kargo-back/clients/storage"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/clients-models"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func logAndReturnError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())

	return apigateway.NewErrorResponse(500, err)
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return logAndReturnError(err), nil
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
		return logAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, client), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
