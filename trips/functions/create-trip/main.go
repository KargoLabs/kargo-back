package main

import (
	"context"
	"encoding/json"
	"errors"
	cardModels "kargo-back/models/cards"
	clientModels "kargo-back/models/clients"
	transactionModels "kargo-back/models/transactions"
	models "kargo-back/models/trips"
	"kargo-back/shared/apigateway"
	"kargo-back/shared/environment"
	lambdaLibrary "kargo-back/shared/lambda"
	"kargo-back/shared/random"
	cardStorage "kargo-back/storage/cards"
	clientStorage "kargo-back/storage/clients"
	storage "kargo-back/storage/trips"
	truckStorage "kargo-back/storage/trucks"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	createTransactionLambdaName = environment.GetString("CREATE_TRANSACTION_LAMBDA_NAME", "transactions_create-transaction")
)

type tripRequest struct {
	*events.APIGatewayProxyRequest
	body            url.Values
	clientID        string
	partnerID       string
	truckID         string
	cardID          string
	transactionID   string
	tripPriceString string
}

func (tripReq *tripRequest) setCLientID(ctx context.Context) *events.APIGatewayProxyResponse {
	username, err := apigateway.GetUsername(*tripReq.APIGatewayProxyRequest)
	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	client, err := clientStorage.LoadClient(ctx, random.GetSHA256WithPrefix(clientModels.ClientIDPrefix, username))
	if errors.Is(err, clientStorage.ErrClientNotFound) {
		return apigateway.NewErrorResponse(404, err)
	}

	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	tripReq.clientID = client.ClientID

	return nil
}

func (tripReq *tripRequest) setTruckIDAndPartnerID(ctx context.Context) *events.APIGatewayProxyResponse {
	truckID := tripReq.body.Get("truck_id")
	if truckID == "" {
		return apigateway.NewErrorResponse(400, models.ErrMissingTruckID)
	}

	truck, err := truckStorage.LoadTruck(ctx, truckID)
	if errors.Is(err, truckStorage.ErrTruckNotFound) {
		return apigateway.NewErrorResponse(404, err)
	}

	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	tripReq.truckID = truck.TruckID
	tripReq.partnerID = truck.PartnerID

	return nil
}

func (tripReq *tripRequest) setCardID(ctx context.Context) *events.APIGatewayProxyResponse {
	cardID := tripReq.body.Get("card_id")
	if cardID == "" {
		return apigateway.NewErrorResponse(400, cardModels.ErrMissingCardID)
	}

	card, err := cardStorage.LoadCard(ctx, cardID)
	if errors.Is(err, cardStorage.ErrCardNotFound) {
		return apigateway.NewErrorResponse(404, err)
	}

	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	tripReq.cardID = card.CardID

	return nil
}

func (tripReq *tripRequest) createTransaction(ctx context.Context) *events.APIGatewayProxyResponse {
	params := url.Values{}

	params.Set("client_id", tripReq.clientID)
	params.Set("partner_id", tripReq.partnerID)
	params.Set("amount", tripReq.tripPriceString)
	params.Set("card_id", tripReq.cardID)

	lambdaResponse, err := lambdaLibrary.InvokeAPIGatewayWithURLEncodedParams(ctx, createTransactionLambdaName, params)
	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	if lambdaResponse.StatusCode != 201 {
		// It is a 500 because supposedly all inputs were correct before triggering the transaction lambda
		lambdaResponse.StatusCode = 500

		return lambdaResponse
	}

	transaction := transactionModels.Transaction{}

	err = json.Unmarshal([]byte(lambdaResponse.Body), &transaction)
	if err != nil {
		return apigateway.LogAndReturnError(err)
	}

	tripReq.transactionID = transaction.TransactionID

	return nil
}

func apiGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	tripReq := &tripRequest{
		APIGatewayProxyRequest: &request,
	}

	errResponse := tripReq.setCLientID(ctx)
	if errResponse != nil {
		return errResponse, nil
	}

	body, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	tripReq.body = body

	tripReq.tripPriceString = body.Get("trip_price")

	tripPrice, err := strconv.ParseFloat(tripReq.tripPriceString, 64)
	if err != nil || tripPrice <= 0 {
		return apigateway.NewErrorResponse(400, models.ErrInvalidTripPrice), nil
	}

	startTimeString := body.Get("start_time")
	startTimeInt, err := strconv.ParseInt(startTimeString, 10, 64)

	if startTimeString == "" || startTimeInt < 0 || err != nil {
		return apigateway.NewErrorResponse(400, models.ErrInvalidStartTime), nil
	}

	startTime := time.Unix(startTimeInt, 0)

	errResponse = tripReq.setTruckIDAndPartnerID(ctx)
	if errResponse != nil {
		return errResponse, nil
	}

	errResponse = tripReq.setCardID(ctx)
	if errResponse != nil {
		return errResponse, nil
	}

	errResponse = tripReq.createTransaction(ctx)
	if errResponse != nil {
		return errResponse, nil
	}

	trip, err := models.NewTrip(tripReq.clientID, tripReq.partnerID, tripReq.truckID, tripReq.transactionID, body.Get("departure_url"), body.Get("departure_direction"), body.Get("arrival_url"), body.Get("arrival_direction"), tripPrice, startTime)
	if err != nil {
		return apigateway.NewErrorResponse(400, err), nil
	}

	err = storage.PutTrip(ctx, trip)
	if err != nil {
		return apigateway.LogAndReturnError(err), nil
	}

	return apigateway.NewJSONResponse(201, trip), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
