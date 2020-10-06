package storage

import (
	"context"
	"errors"
	models "kargo-back/models/trips"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrTripNotFound when no trip was found
	ErrTripNotFound = errors.New("trip not found")

	tripsTableName         = environment.GetString("TRIPS_TABLE_NAME", "trips")
	tripClienIDFieldName   = environment.GetString("TRIP_CLIENT_ID_FIELD_NAME", "client_id")
	tripClientIDIndexName  = environment.GetString("TRIP_CLIENT_ID_INDEX_NAME", "trip-client-id-index")
	tripPartnerIDFieldName = environment.GetString("TRIP_PARTNER_ID_FIELD_NAME", "partner_id")
	tripPartnerIDIndexName = environment.GetString("TRIP_CLIENT_ID_INDEX_NAME", "trip-partner-id-index")
	dynamoClient           dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutTrip saves a trip in DynamoDB
func PutTrip(ctx context.Context, trip *models.Trip) error {
	item, err := dynamodbattribute.MarshalMap(trip)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tripsTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadTrip loads a Trip from DynamoDB
func LoadTrip(ctx context.Context, tripID string) (*models.Trip, error) {
	key := map[string]*dynamodb.AttributeValue{
		"trip_id": {S: aws.String(tripID)},
	}

	params := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(tripsTableName),
		ConsistentRead: aws.Bool(true),
	}

	response, err := dynamoClient.GetItemWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, ErrTripNotFound
	}

	var trip models.Trip
	err = dynamodbattribute.UnmarshalMap(response.Item, &trip)

	return &trip, err
}

// QueryTrips queries trips from dynamodb based on the given indexed field
func queryTrips(ctx context.Context, query models.TripQuery, fieldName, indexName string) ([]*models.Trip, error) {
	keyCondition := expression.KeyEqual(expression.Key(fieldName), expression.Value(query.Value))
	dynamoExpression := expression.NewBuilder().WithKeyCondition(keyCondition)

	if query.FilterFinished {
		finishedFilter := expression.Name("finished").Equal(expression.Value(query.Finished))
		dynamoExpression.WithFilter(finishedFilter)
	}

	dynamoQuery, err := dynamoExpression.Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  dynamoQuery.Names(),
		ExpressionAttributeValues: dynamoQuery.Values(),
		KeyConditionExpression:    dynamoQuery.KeyCondition(),
		FilterExpression:          dynamoQuery.Filter(),
		IndexName:                 aws.String(indexName),
		TableName:                 aws.String(tripsTableName),
	}

	response, err := dynamoClient.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, ErrTripNotFound
	}

	trips := []*models.Trip{}

	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &trips)
	if err != nil {
		return nil, err
	}

	return trips, nil
}

// QueryClientTrips retrieves all the trips that belongs to a client, optionally filtering by finished
func QueryClientTrips(ctx context.Context, clientQuery models.TripQuery) ([]*models.Trip, error) {
	trips, err := queryTrips(ctx, clientQuery, tripClienIDFieldName, tripClientIDIndexName)
	if err != nil {
		return nil, err
	}

	return trips, nil
}

// QueryPartnerTrips retrieves all the trips that belongs to a partner, optionally filtering by finished
func QueryPartnerTrips(ctx context.Context, partnerQuery models.TripQuery) ([]*models.Trip, error) {
	trips, err := queryTrips(ctx, partnerQuery, tripPartnerIDFieldName, tripPartnerIDIndexName)
	if err != nil {
		return nil, err
	}

	return trips, nil
}
