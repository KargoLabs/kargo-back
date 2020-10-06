package apigateway

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

var (
	errCastingClaimsMap = errors.New("error casting authorizer claims map")
	errCastingUsername  = errors.New("error casting username")
	errCastingEmail     = errors.New("error casting email")
)

// GetUsername extracts username from lambda proxy request model
func GetUsername(request events.APIGatewayProxyRequest) (string, error) {
	claimsMap, ok := request.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok || claimsMap == nil {
		return "", errCastingClaimsMap
	}

	username, ok := claimsMap["cognito:username"].(string)
	if !ok || username == "" {
		return "", errCastingUsername
	}

	return username, nil
}

// GetEmail extracts email from lambda proxy request model
func GetEmail(request events.APIGatewayProxyRequest) (string, error) {
	claimsMap, ok := request.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok || claimsMap == nil {
		return "", errCastingClaimsMap
	}

	email, ok := claimsMap["email"].(string)
	if !ok || email == "" {
		return "", errCastingEmail
	}

	return email, nil
}
