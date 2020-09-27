package storage

import (
	"context"
	"errors"
	models "kargo-back/models/trips"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrTripNotFound when no trip was found
	ErrTripNotFound = errors.New("trip not found")

	tripsTableName = environment.GetString("TRIPS_TABLE_NAME", "trips")
	dynamoClient   dynamodbiface.DynamoDBAPI
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
