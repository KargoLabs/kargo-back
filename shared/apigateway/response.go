package apigateway

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// ErrorResponse represents an API error
type ErrorResponse struct {
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

// NewErrorResponse returns an error response
func NewErrorResponse(statusCode int, err error) *events.APIGatewayProxyResponse {
	return NewJSONResponse(statusCode, ErrorResponse{
		HTTPCode: statusCode,
		Message:  err.Error(),
	})
}

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
