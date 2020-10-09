package main

import (
	"context"
	"errors"
	trips "kargo-back/models/trips"
	models "kargo-back/models/trucks"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	"kargo-back/shared/s3"
	partnerStorage "kargo-back/storage/partners"
	storage "kargo-back/storage/trucks"
	"net/url"
	"strconv"

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

	year, err := strconv.Atoi(body.Get("year"))
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidYear), nil
	}

	maxVolume, err := strconv.ParseFloat(body.Get("max_volume"), 64)
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidMaxVolume), nil
	}

	maxWeight, err := strconv.ParseFloat(body.Get("max_weight"), 64)
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidMaxWeight), nil
	}

	basePrice, err := strconv.ParseFloat(body.Get("base_price"), 64)
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidBasePrice), nil
	}

	perRegionPrice, err := strconv.ParseFloat(body.Get("per_region_price"), 64)
	if err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidPerRegionPrice), nil
	}

	regionsParam, ok := body["regions"]
	if !ok || regionsParam[0] == "" {
		return apigateway.NewErrorResponse(400, models.ErrMissingRegion), nil
	}

	regions := make([]trips.Region, len(regionsParam))
	for i := range regionsParam {
		regions[i] = trips.Region(regionsParam[i])
	}

	partner, err := partnerStorage.LoadPartner(ctx, body.Get("partner_id"))
	if errors.Is(err, partnerStorage.ErrPartnerNotFound) {
		return apigateway.NewErrorResponse(404, err), nil
	}

	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	truck := &models.Truck{
		PartnerID:         partner.PartnerID,
		RegistrationPlate: body.Get("registration_plate"),
		Brand:             body.Get("brand"),
		Model:             body.Get("model"),
		Year:              year,
		Type:              models.TruckType(body.Get("type")),
		Regions:           regions,
		MaxVolume:         maxVolume,
		MaxWeight:         maxWeight,
		BasePrice:         basePrice,
		PerRegionPrice:    perRegionPrice,
	}

	err = truck.ValidateTruck()
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutTruck(ctx, truck)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	uploadProfilePhotoURL, err := s3.GetPutPreSignedURL(ctx, profilePhotosBucket, truck.ProfilePhotoS3Path)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, &s3.StructWithUploadURL{
		Struct:                truck,
		UploadProfilePhotoURL: uploadProfilePhotoURL,
	}), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
