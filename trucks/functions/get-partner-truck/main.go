package main

import (
	"context"
	"errors"
	models "kargo-back/models/trucks"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/trucks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrMissingPartnerID is when the partner_id query parameter is missing
	ErrMissingPartnerID = errors.New("missing partner_id query param")
	// ErrInvalidAvailableValue is when a non boolean value is given to the available param
	ErrInvalidAvailableValue = errors.New("invalid value for available param")

	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	params := request.QueryStringParameters

	partner, ok := params["partner_id"]
	if !ok || partner == "" {
		return apigateway.NewErrorResponse(400, ErrMissingPartnerID), nil
	}

	// Optional parameter, filter by truck avaliability
	available := params["available"]

	query, err := models.NewPartnerTruckQuery(partner, available)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	trucks, err := storage.QueryPartnerTrucks(ctx, query)
	if errors.Is(err, storage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	trucksWithGetURL := []*s3.StructWithGetURL{}

	for _, truck := range trucks {
		getProfilePhotoURL, err := s3.GetGetPreSignedURL(ctx, profilePhotosBucket, truck.ProfilePhotoS3Path)
		if err != nil {
			return apigateway.LogAndReturnError(err), nil
		}

		trucksWithGetURL = append(trucksWithGetURL, &s3.StructWithGetURL{
			Struct:             truck,
			GetProfilePhotoURL: getProfilePhotoURL,
		})
	}

	return apigateway.NewJSONResponse(200, trucksWithGetURL), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
