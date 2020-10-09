package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/trucks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrMissingTruckID is when the truck_id query parameter is missing
	ErrMissingTruckID = errors.New("missing truck_id query param")

	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	params := request.QueryStringParameters

	truckID, ok := params["truck_id"]

	if !ok {
		return apigateway.NewErrorResponse(403, ErrMissingTruckID), nil
	}

	truck, err := storage.LoadTruck(ctx, truckID)
	if errors.Is(err, storage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	getProfilePhotoURL, err := s3.GetGetPreSignedURL(ctx, profilePhotosBucket, truck.ProfilePhotoS3Path)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, &s3.StructWithGetURL{
		Struct:             truck,
		GetProfilePhotoURL: getProfilePhotoURL,
	}), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
