package main

import (
	"context"
	"errors"
	usersModels "kargo-back/models/users"
	"kargo-back/shared/apigateway"
	phoneValidation "kargo-back/shared/phone-validation"
	"kargo-back/shared/random"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errInvalidUserType = errors.New("invalid user type parameter")
)

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	userType := usersModels.UserType(body.Get("user_type"))
	if !usersModels.ValidUserTypes[userType] {
		return apigateway.NewErrorResponse(400, errInvalidUserType), nil
	}

	username, err := apigateway.GetUsername(request)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	// truoraAccountID will be the same as our parter and client IDs
	truoraAccountID := random.GetSHA256WithPrefix(usersModels.UserTypeToPrefix[userType], username)

	createResponse, err := phoneValidation.CreatePhoneValidation(truoraAccountID)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	if createResponse.TruoraError != nil {
		return apigateway.NewJSONResponse(createResponse.HTTPCode, createResponse.TruoraError), nil
	}

	return apigateway.NewJSONResponse(201, createResponse.TruoraValidation), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
