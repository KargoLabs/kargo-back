package lambda

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsLambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var (
	lambdaClient lambdaiface.LambdaAPI
)

func init() {
	sess := session.Must(session.NewSession())
	lambdaClient = awsLambda.New(sess)
}

// InvokeAPIGatewayWithURLEncodedParams invokes apigateway lambda with url encoded params
func InvokeAPIGatewayWithURLEncodedParams(ctx context.Context, lambdaName string, params url.Values) (*events.APIGatewayProxyResponse, error) {
	request := &events.APIGatewayProxyRequest{}

	request.Body = params.Encode()

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	invokeInput := &awsLambda.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      payload,
	}

	invokeResponse, err := lambdaClient.InvokeWithContext(ctx, invokeInput)
	if err != nil {
		return nil, err
	}

	response := events.APIGatewayProxyResponse{}

	err = json.Unmarshal(invokeResponse.Payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
