package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	trips "kargo-back/shared/trips-models"
	models "kargo-back/shared/truck-models"
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

	mileague, err := strconv.Atoi(body.Get("mileague"))
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidMileague), nil
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
		Mileague:          mileague,
		Type:              models.TruckType(body.Get("type")),
		Location:          trips.Region(body.Get("location")),
	}

	truck, err := models.NewTruck(*truckParam)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
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
