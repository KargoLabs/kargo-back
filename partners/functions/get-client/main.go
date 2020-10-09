package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingClientID = errors.New("missing client id in query parameter")

	profilePhotosBucket = environment.GetString("PROFILE_PHOTOS_BUCKET", "kargo-profile-photos")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	clientID, ok := request.QueryStringParameters["client_id"]
	if !ok || clientID == "" {
		return apigateway.NewErrorResponse(400, errMissingClientID), nil
	}

	client, err := storage.LoadClient(ctx, clientID)
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	getProfilePhotoURL, err := s3.GetGetPreSignedURL(ctx, profilePhotosBucket, client.ProfilePhotoS3Path)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, &s3.StructWithGetURL{
		Struct:             client,
		GetProfilePhotoURL: getProfilePhotoURL,
	}), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
