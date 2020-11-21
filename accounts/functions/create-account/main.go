package main

import (
	"context"
	"errors"
	models "kargo-back/models/accounts"
	partnerModel "kargo-back/models/partners"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/accounts"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrAccountAlreadyExists when there's already an account for the partner
	ErrAccountAlreadyExists = errors.New("an account is already associated with this partner")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	partnerID := random.GetSHA256WithPrefix(partnerModel.PartnerIDPrefix, username)

	// Check if there's already an account
	account, err := storage.LoadAccount(ctx, partnerID)
	if account != nil {
		return apigateway.NewErrorResponse(403, ErrAccountAlreadyExists), nil
	}

	account, err = models.NewAccount(partnerID, body.Get("name"), body.Get("document"), body.Get("number"))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutAccount(ctx, account)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, account), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
