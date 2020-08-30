package main

import (
	"context"
	"encoding/base64"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/partners-models"
	storage "kargo-back/storage/partners"
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

	partner, err := models.NewPartner(body.Get("name"), body.Get("document"), birthDate)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutPartner(ctx, partner)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, partner), nil

}

func main() {
	lambda.Start(apiGatewayHandler)
}
