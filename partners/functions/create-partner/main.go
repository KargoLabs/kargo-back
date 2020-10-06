package main

import (
	"context"
	models "kargo-back/models/partners"
	"kargo-back/shared/apigateway"
	storage "kargo-back/storage/partners"
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

	email, err := apigateway.GetEmail(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	birthdate, err := time.Parse("2006-01-02", body.Get("birthdate"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	partner, err := models.NewPartner(username, body.Get("name"), body.Get("document"), body.Get("phone_number"), email, birthdate)
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
