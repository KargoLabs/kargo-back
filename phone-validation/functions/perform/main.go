package main

import (
	"context"
	"errors"
	"kargo-back/shared/apigateway"
	phoneValidation "kargo-back/shared/phone-validation"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errMissingCode         = errors.New("missing code parameter")
	errMissingValidationID = errors.New("missing validation ID parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	code := body.Get("code")
	if code == "" {
		return apigateway.NewErrorResponse(400, errMissingCode), nil
	}

	validationID := body.Get("validation_id")
	if validationID == "" {
		return apigateway.NewErrorResponse(400, errMissingValidationID), nil
	}

	performResponse, err := phoneValidation.PerformPhoneValidation(code, validationID)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	if performResponse.TruoraError != nil {
		return apigateway.NewJSONResponse(performResponse.HTTPCode, performResponse.TruoraError), nil
	}

	return apigateway.NewJSONResponse(200, performResponse.TruoraValidation), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
