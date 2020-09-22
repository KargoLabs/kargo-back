package storage

import (
	"context"
	"errors"
	models "kargo-back/models/clients"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrClientNotFound when no Client was found
	ErrClientNotFound = errors.New("client not found")

	clientsTableName = environment.GetString("CLIENTS_TABLE_NAME", "clients")
	dynamoClient     dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutClient saves a Client in DynamoDB
func PutClient(ctx context.Context, client *models.Client) error {
	item, err := dynamodbattribute.MarshalMap(client)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(clientsTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadClient loads a Client from DynamoDB
func LoadClient(ctx context.Context, clientID string) (*models.Client, error) {
	key := map[string]*dynamodb.AttributeValue{
		"client_id": {S: aws.String(clientID)},
	}

	params := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(clientsTableName),
		ConsistentRead: aws.Bool(true),
	}

	response, err := dynamoClient.GetItemWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, ErrClientNotFound
	}

	var client models.Client
	err = dynamodbattribute.UnmarshalMap(response.Item, &client)

	return &client, err
}
