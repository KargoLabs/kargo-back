package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"kargo-back/partners/storage"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/normalize"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingPartnerID = errors.New("missing partner id in body parameter")
)

func logAndReturnError(err error) *events.APIGatewayProxyResponse {
	fmt.Println(err.Error())

	return apigateway.NewErrorResponse(500, err)
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return logAndReturnError(err), nil
	}

	body, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return logAndReturnError(err), nil
	}

	partnerID := body.Get("partner_id")
	if partnerID == "" {
		return apigateway.NewErrorResponse(400, errMissingPartnerID), nil
	}

	partner, err := storage.LoadPartner(ctx, partnerID)
	if errors.Is(err, storage.ErrPartnerNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return logAndReturnError(err), nil
	}

	if body.Get("name") != "" {
		partner.Name = normalize.NormalizeName(body.Get("name"))
	}

	if body.Get("document") != "" {
		partner.Document = body.Get("document")
	}

	if body.Get("birth_date") != "" {
		birthDate, err := time.Parse("2006-01-02", body.Get("birth_date"))
		if err != nil {
			return apigateway.NewErrorResponse(400, err), nil
		}

		partner.BirthDate = birthDate
	}

	err = storage.PutPartner(ctx, partner)
	if err != nil {
		return logAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, partner), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
