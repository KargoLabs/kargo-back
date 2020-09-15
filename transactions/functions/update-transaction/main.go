package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/transaction-models"
	storage "kargo-back/storage/transactions"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	transaction, err := storage.LoadTransaction(ctx,
		body.Get("transaction_id"))
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	status := body.Get("status")

	if status == "" {
		return apigateway.LogAndReturnError(
			errors.New("missing status parameter")), nil
	}

	transaction.Status = models.TransactionStatus(status)
	err = storage.PutTransaction(ctx, transaction)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, transaction), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
