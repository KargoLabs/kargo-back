package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/normalize"
	storage "kargo-back/storage/clients"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingClientID = errors.New("missing client id in body parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	clientID := body.Get("client_id")
	if clientID == "" {
		return apigateway.NewErrorResponse(400, errMissingClientID), nil
	}

	client, err := storage.LoadClient(ctx, clientID)
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	if body.Get("name") != "" {
		client.Name = normalize.NormalizeName(body.Get("name"))
	}

	if body.Get("document") != "" {
		client.Document = body.Get("document")
	}

	if body.Get("birth_date") != "" {
		birthDate, err := time.Parse("2006-01-02", body.Get("birth_date"))
		if err != nil {
			return apigateway.NewErrorResponse(400, err), nil
		}

		client.BirthDate = birthDate
	}

	err = storage.PutClient(ctx, client)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, client), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
