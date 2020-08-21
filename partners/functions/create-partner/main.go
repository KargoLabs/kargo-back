package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"kargo-back/partners/storage"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/partners-models"
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
	bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return logAndReturnError(err), nil
	}

	body, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return logAndReturnError(err), nil
	}

	birthDate, err := time.Parse("2006-01-02", body.Get("birth_date"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	partner, err := models.NewPartner(body.Get("name"), body.Get("document"), birthDate)
	if err != nil {
		return logAndReturnError(err), nil
	}

	err = storage.PutPartner(ctx, partner)
	if err != nil {
		return logAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, partner), nil

}

func main() {
	lambda.Start(apiGatewayHandler)
}
