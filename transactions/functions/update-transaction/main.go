package main

import (
	"context"
	"errors"
	models "kargo-back/models/transactions"
	"kargo-back/shared/apigateway"
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

	transaction, err := storage.LoadTransaction(ctx, body.Get("transaction_id"))
	if errors.Is(err, storage.ErrTransactionNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	status := models.TransactionStatus(body.Get("status"))
	if status == "" {
		return apigateway.LogAndReturnError(errors.New("missing status parameter")), nil
	}

	if !models.ValidTransactionStatuses[status] {
		return apigateway.NewErrorResponse(400, errors.New("invalid status parameter")), nil
	}

	transaction.Status = status

	err = storage.PutTransaction(ctx, transaction)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, transaction), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
