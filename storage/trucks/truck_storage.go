package storage

import (
	"context"
	"errors"
	models "kargo-back/models/trucks"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrTruckNotFound when no truck was found
	ErrTruckNotFound = errors.New("truck not found")

	trucksTableName     = environment.GetString("TRUCKS_TABLE_NAME", "trucks")
	trucksTypeIndexName = environment.GetString("TRUCK_TYPE_INDEX_NAME", "truck_type-index")
	dynamoClient        dynamodbiface.DynamoDBAPI
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

// QueryTrucks queries trucks from DynamoDB with given input
func QueryTrucks(ctx context.Context, trucksQuery *models.TrucksQuery) ([]*models.Truck, error) {
	keyCondition := expression.KeyEqual(expression.Key("truck_type"), expression.Value(trucksQuery.TruckType))
	avaibleFilter := expression.Name("available").Equal(expression.Value(true))
	weightFilter := expression.Name("max_weight").GreaterThanEqual(expression.Value(trucksQuery.Weight))
	volumeFilter := expression.Name("max_volume").GreaterThanEqual(expression.Value(trucksQuery.Volume))
	originFilter := expression.Name("regions").Contains(trucksQuery.Origin)
	destinationFilter := expression.Name("regions").Contains(trucksQuery.Destination)
	filterExpression := expression.And(avaibleFilter, originFilter, weightFilter, volumeFilter, destinationFilter)

	dynamoExpression, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithFilter(filterExpression).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  dynamoExpression.Names(),
		ExpressionAttributeValues: dynamoExpression.Values(),
		KeyConditionExpression:    dynamoExpression.KeyCondition(),
		FilterExpression:          dynamoExpression.Filter(),
		IndexName:                 aws.String(trucksTypeIndexName),
		TableName:                 aws.String(trucksTableName),
	}

	response, err := dynamoClient.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, ErrTruckNotFound
	}

	trucks := []*models.Truck{}

	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &trucks)
	if err != nil {
		return nil, err
	}

	return trucks, nil
}
