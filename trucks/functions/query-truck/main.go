package main

import (
	"context"
	"errors"
	tripsModels "kargo-back/models/trips"
	models "kargo-back/models/trucks"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/trucks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingOrigin      = errors.New("missing origin parameter")
	errMissingDestination = errors.New("missing destination parameter")
	errMissingWeight      = errors.New("missing weight parameter")
	errMissingVolume      = errors.New("missing volume parameter")

	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	truckType, ok := request.QueryStringParameters["truck_type"]
	if !ok || truckType == "" {
		return apigateway.NewErrorResponse(400, models.ErrMissingTruckType), nil
	}

	originProvince, ok := request.QueryStringParameters["origin"]
	if !ok || originProvince == "" {
		return apigateway.NewErrorResponse(400, errMissingOrigin), nil
	}

	destinationProvince, ok := request.QueryStringParameters["destination"]
	if !ok || destinationProvince == "" {
		return apigateway.NewErrorResponse(400, errMissingDestination), nil
	}

	weight, ok := request.QueryStringParameters["weight"]
	if !ok || weight == "" {
		return apigateway.NewErrorResponse(400, errMissingWeight), nil
	}

	volume, ok := request.QueryStringParameters["volume"]
	if !ok || volume == "" {
		return apigateway.NewErrorResponse(400, errMissingVolume), nil
	}

	origin, err := tripsModels.GetRegionFromProvince(tripsModels.Province(originProvince))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	destination, err := tripsModels.GetRegionFromProvince(tripsModels.Province(destinationProvince))
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	trucksQuery, err := models.NewTrucksQuery(truckType, weight, volume, origin, destination)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	trucks, err := storage.QueryTrucks(ctx, trucksQuery)
	if errors.Is(err, storage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	trucksWithTripPriceAndGetURL := []*s3.StructWithGetURL{}

	for _, truck := range trucks {
		getProfilePhotoURL, err := s3.GetGetPreSignedURL(ctx, profilePhotosBucket, truck.ProfilePhotoS3Path)
		if err != nil {
			return apigateway.LogAndReturnError(err), nil
		}

		truckWithTripPrice := models.NewTruckWithTripPrice(origin, destination, truck)

		trucksWithTripPriceAndGetURL = append(trucksWithTripPriceAndGetURL, &s3.StructWithGetURL{
			Struct:             truckWithTripPrice,
			GetProfilePhotoURL: getProfilePhotoURL,
		})
	}

	return apigateway.NewJSONResponse(200, trucksWithTripPriceAndGetURL), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
