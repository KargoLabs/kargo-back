package main

import (
	"context"
	"errors"
	"fmt"
	"kargo-back/clients/storage"
	"kargo-back/shared/apigateway"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingClientID = errors.New("missing client id in query parameter")
)

func logAndReturnError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())

	return apigateway.NewJSONResponse(500, err)
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return logAndReturnError(err)
	}

	clientID := body.Get("client_id")
	if clientID == "" {
		return apigateway.NewJSONResponse(400, errMissingClientID)
	}

	client, err := storage.LoadClient(ctx, clientID)
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewJSONResponse(404, err)
	}

	if err != nil {
		return logAndReturnError(err)
	}

	if body.Get("name") != "" {
		client.Name = body.Get("name")
	}

	if body.Get("document") != "" {
		client.Document = body.Get("document")
	}

	if body.Get("birth_date") != "" {
		birthDate, err := time.Parse("2006-01-02", body.Get("birth_date"))
		if err != nil {
			return apigateway.NewJSONResponse(400, err)
		}

		client.BirthDate = birthDate
	}

	err = storage.PutClient(ctx, client)
	if err != nil {
		return logAndReturnError(err)
	}

	return apigateway.NewJSONResponse(200, client)
}

func main() {
	lambda.Start(apiGatewayHandler)
}
