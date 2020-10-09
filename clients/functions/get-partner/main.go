package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/partners"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingPartnerID = errors.New("missing partner id in query parameter")

	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
)

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
		return apigateway.LogAndReturnError(err), nil
	}

	getProfilePhotoURL, err := s3.GetGetPreSignedURL(ctx, profilePhotosBucket, partner.ProfilePhotoS3Path)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, &s3.StructWithGetURL{
		Struct:             partner,
		GetProfilePhotoURL: getProfilePhotoURL,
	}), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
