package main

import (
	"context"
	"errors"
	trips "kargo-back/models/trips"
	models "kargo-back/models/trucks"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/trucks"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
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
			return apigateway.NewErrorResponse(400, models.ErrInvalidYear), nil
		}

		truck.Year = year
	}

	if body.Get("max_volume") != "" {
		maxVolume, err := strconv.ParseFloat(body.Get("max_volume"), 64)
		if err != nil {
			return apigateway.NewErrorResponse(400, models.ErrInvalidMaxVolume), nil
		}

		truck.MaxVolume = maxVolume
	}

	if body.Get("max_weight") != "" {
		maxWeight, err := strconv.ParseFloat(body.Get("max_weight"), 64)
		if err != nil {
			return apigateway.NewErrorResponse(400, models.ErrInvalidMaxWeight), nil
		}

		truck.MaxWeight = maxWeight
	}

	if body.Get("base_price") != "" {
		basePrice, err := strconv.ParseFloat(body.Get("base_price"), 64)
		if err != nil {
			return apigateway.NewErrorResponse(400, models.ErrInvalidBasePrice), nil
		}

		truck.BasePrice = basePrice
	}

	if body.Get("per_region_price") != "" {
		perRegionPrice, err := strconv.ParseFloat(body.Get("per_region_price"), 64)
		if err != nil {
			return apigateway.NewErrorResponse(400, models.ErrInvalidPerRegionPrice), nil
		}

		truck.PerRegionPrice = perRegionPrice
	}

	if body.Get("available") != "" {
		available, err := strconv.ParseBool(body.Get("available"))
		if err != nil {
			return apigateway.NewErrorResponse(400, models.ErrInvalidAvailable), nil
		}

		truck.Available = available
	}

	if body.Get("type") != "" {
		truckType := models.TruckType(body.Get("type"))
		_, ok := models.ValidTruckTypes[truckType]
		if !ok {
			return apigateway.NewErrorResponse(400, models.ErrInvalidTruckType), nil
		}
		truck.Type = truckType
	}

	regionsParam, ok := body["regions"]
	if ok {
		regions := make([]trips.Region, len(regionsParam))
		for i := range regionsParam {
			regions[i] = trips.Region(regionsParam[i])
		}
		truck.Regions = regions
	}

	var uploadProfilePhotoURL string

	if body.Get("profile_photo") != "" {
		uploadProfilePhotoURL, err = s3.GetPutPreSignedURL(ctx, profilePhotosBucket, truck.ProfilePhotoS3Path)
		if err != nil {
			return apigateway.LogAndReturnError(err), nil
		}
	}

	truck.UpdateDate = time.Now()

	err = storage.PutTruck(ctx, truck)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, &s3.StructWithUploadURL{
		Struct:                truck,
		UploadProfilePhotoURL: uploadProfilePhotoURL,
	}), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
