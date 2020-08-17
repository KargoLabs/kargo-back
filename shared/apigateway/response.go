package apigateway

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// NewJSONResponse creates a new JSON response given a serializable val
func NewJSONResponse(statusCode int, val interface{}) *events.APIGatewayProxyResponse {
	data, _ := json.Marshal(val)
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
