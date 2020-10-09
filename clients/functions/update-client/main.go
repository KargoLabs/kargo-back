package main

import (
	"context"
	"errors"
	models "kargo-back/models/clients"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/normalize"
	"kargo-back/shared/random"
	"kargo-back/shared/s3"
	storage "kargo-back/storage/clients"
	"net/url"
	"time"

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

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	client, err := storage.LoadClient(ctx, random.GetSHA256WithPrefix(models.ClientIDPrefix, username))
	if errors.Is(err, storage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	if body.Get("name") != "" {
		client.Name = normalize.NormalizeName(body.Get("name"))
	}

	if body.Get("document") != "" {
		client.Document = body.Get("document")
	}

	if body.Get("birthdate") != "" {
		birthdate, err := time.Parse("2006-01-02", body.Get("birthdate"))
		if err != nil {
			return apigateway.NewErrorResponse(400, err), nil
		}

		client.Birthdate = birthdate
	}

	var uploadProfilePhotoURL string

	if body.Get("profile_photo") != "" {
		uploadProfilePhotoURL, err = s3.GetPutPreSignedURL(ctx, profilePhotosBucket, client.ProfilePhotoS3Path)
		if err != nil {
			return apigateway.LogAndReturnError(err), nil
		}
	}

	client.UpdateDate = time.Now()

	err = storage.PutClient(ctx, client)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(200, &s3.StructWithUploadURL{
		Struct:                client,
		UploadProfilePhotoURL: uploadProfilePhotoURL,
	}), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
