package main

import (
	"context"
	"errors"
	models "kargo-back/models/partners"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/random"
	storage "kargo-back/storage/partners"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingPartnerID = errors.New("missing partner id in query parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	partner, err := storage.LoadPartner(ctx, random.GetSHA256WithPrefix(models.PartnerIDPrefix, username))
	if errors.Is(err, storage.ErrPartnerNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, partner), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
