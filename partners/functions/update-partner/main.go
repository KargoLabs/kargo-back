package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/normalize"
	models "kargo-back/shared/partners-models"
	"kargo-back/shared/random"
	storage "kargo-back/storage/partners"

	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingPartnerID = errors.New("missing partner id in body parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	body, err := url.ParseQuery(string(request.Body))
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

	if body.Get("name") != "" {
		partner.Name = normalize.NormalizeName(body.Get("name"))
	}

	if body.Get("document") != "" {
		partner.Document = body.Get("document")
	}

	if body.Get("birthdate") != "" {
		birthdate, err := time.Parse("2006-01-02", body.Get("birthdate"))
		if err != nil {
			return apigateway.NewErrorResponse(400, err), nil
		}

		partner.Birthdate = birthdate
	}

	partner.UpdateDate = time.Now()

	err = storage.PutPartner(ctx, partner)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, partner), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
