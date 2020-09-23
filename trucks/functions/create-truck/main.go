package main

import (
	"context"
	"errors"
	trips "kargo-back/models/trips"
	models "kargo-back/models/trucks"
	"kargo-back/shared/apigateway"
	partnerStorage "kargo-back/storage/partners"
	storage "kargo-back/storage/trucks"
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

	year, err := strconv.Atoi(body.Get("year"))
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidYear), nil
	}

	regionsParam, ok := body["regions"]
	if !ok || regionsParam[0] == "" {
		return apigateway.NewErrorResponse(400, models.ErrMissingRegion), nil
	}

	regions := make([]trips.Region, len(regionsParam))
	for i := range regionsParam {
		regions[i] = trips.Region(regionsParam[i])
	}

	partner, err := partnerStorage.LoadPartner(ctx, body.Get("partner_id"))
	if errors.Is(err, partnerStorage.ErrPartnerNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	truckParam := &models.Truck{
		PartnerID:         partner.PartnerID,
		RegistrationPlate: body.Get("registration_plate"),
		Brand:             body.Get("brand"),
		Model:             body.Get("model"),
		Year:              year,
		Type:              models.TruckType(body.Get("type")),
		Location:          trips.Region(body.Get("location")),
		Regions:           regions,
	}

	truck, err := models.NewTruck(*truckParam)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutTruck(ctx, truck)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, truck), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
