package storage

import (
	"context"
	"errors"
	"kargo-back/shared/environment"
	models "kargo-back/shared/truck-models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrTruckNotFound when no truck was found
	ErrTruckNotFound = errors.New("truck not found")

	trucksTableName = environment.GetString("TRUCKS_TABLE_NAME", "trucks")
	dynamoClient    dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutTruck saves a truck in DynamoDBB
func PutTruck(ctx context.Context, truck *models.Truck) error {
	item, err := dynamodbattribute.MarshalMap(truck)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(trucksTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadTruck loads a truck from DynamoDB
func LoadTruck(ctx context.Context, truckID string) (*models.Truck, error) {
	key := map[string]*dynamodb.AttributeValue{
		"truck_id": {S: aws.String(truckID)},
	}

	params := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(trucksTableName),
		ConsistentRead: aws.Bool(true),
	}

	response, err := dynamoClient.GetItemWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, ErrTruckNotFound
	}

	var truck models.Truck
	err = dynamodbattribute.UnmarshalMap(response.Item, &truck)

	return &truck, err
}
