package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	models "kargo-back/shared/transaction-models"
	clientStorage "kargo-back/storage/clients"
	partnerStorage "kargo-back/storage/partners"
	storage "kargo-back/storage/transactions"

	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	amount, err := strconv.Atoi(body.Get("amount"))
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	partner, err := partnerStorage.LoadPartner(ctx, body.Get("partner_id"))
	if errors.Is(err, partnerStorage.ErrPartnerNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	client, err := clientStorage.LoadClient(ctx, body.Get("client_id"))
	if errors.Is(err, clientStorage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	transaction, err := models.NewTransaction(client.ClientID,
		partner.PartnerID, amount)

	err = storage.PutTransaction(ctx, transaction)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, transaction), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
