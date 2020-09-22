package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	trips "kargo-back/shared/trips-models"
	models "kargo-back/shared/truck-models"
	storage "kargo-back/storage/trucks"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	truck, err := storage.LoadTruck(ctx, body.Get("truck_id"))
	if errors.Is(err, storage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	if body.Get("registration_plate") != "" {
		truck.RegistrationPlate = body.Get("registration_plate")
	}

	if body.Get("brand") != "" {
		truck.Brand = body.Get("brand")
	}

	if body.Get("model") != "" {
		truck.Model = body.Get("model")
	}

	if body.Get("year") != "" {
		year, err := strconv.Atoi(body.Get("year"))
		if err != nil {
			return apigateway.LogAndReturnError(models.ErrInvalidYear), nil
		}

		truck.Year = year
	}

	if body.Get("mileague") != "" {
		mileague, err := strconv.Atoi(body.Get("mileague"))
		if err != nil {
			return apigateway.LogAndReturnError(models.ErrInvalidMileague), nil
		}

		truck.Mileague = mileague
	}

	if body.Get("available") != "" {
		available, err := strconv.ParseBool(body.Get("available"))
		if err != nil {
			return apigateway.LogAndReturnError(models.ErrInvalidAvailable), nil
		}

		truck.Available = available
	}

	if body.Get("type") != "" {
		truckType := models.TruckType(body.Get("type"))
		_, ok := models.ValidTruckTypes[truckType]
		if !ok {
			return apigateway.LogAndReturnError(models.ErrInvalidTruckType), nil
		}
		truck.Type = truckType
	}

	if body.Get("location") != "" {
		location := trips.Region(body.Get("location"))
		_, ok := trips.ValidRegions[location]
		if !ok {
			return apigateway.LogAndReturnError(models.ErrInvalidLocation), nil
		}
		truck.Location = location
	}

	truck.UpdateDate = time.Now()

	err = storage.PutTruck(ctx, truck)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, truck), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}