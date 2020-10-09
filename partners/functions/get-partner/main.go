package main

import (
	"context"
	"errors"
	models "kargo-back/models/partners"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/random"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/partners"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, err := apigateway.GetUsername(request)

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
