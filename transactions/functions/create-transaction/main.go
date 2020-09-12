package main

import (
	"context"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/transaction-models"
	storage "kargo-back/storage/transactions"
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	_, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	amount, err := strconv.Atoi(body.Get("amount"))
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	transaction, err := models.NewTransaction(body.Get("client_id"),
		body.Get("partner_id"), amount)

	err = storage.PutTransaction(ctx, transaction)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, transaction), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
