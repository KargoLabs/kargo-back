package main

import (
	"context"
	"errors"
	"fmt"
	"kargo-back/partners/storage"
	"kargo-back/shared/apigateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingPartnerID = errors.New("missing partner id in query parameter")
)

func logAndReturnError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())

	return apigateway.NewErrorResponse(500, err)
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	partnerID, ok := request.QueryStringParameters["partner_id"]

	if !ok || partnerID == "" {
		return apigateway.NewErrorResponse(400, errMissingPartnerID), nil
	}

	partner, err := storage.LoadPartner(ctx, partnerID)
	if errors.Is(err, storage.ErrPartnerNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return logAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, partner), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
