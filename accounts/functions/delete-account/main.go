package main

import (
	"context"
	partnerModel "kargo-back/models/partners"

	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/accounts"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	partnerID := random.GetSHA256WithPrefix(partnerModel.PartnerIDPrefix, username)

	account, err := storage.DeleteAccount(ctx, partnerID)
	if err == storage.ErrAccountNotFound {
		return apigateway.NewErrorResponse(404, storage.ErrAccountNotFound), nil
	}
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, account), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
