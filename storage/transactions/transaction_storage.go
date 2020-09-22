package storage

import (
	"context"
	"errors"
	models "kargo-back/models/transactions"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrTransactionNotFound when no transaction was found
	ErrTransactionNotFound = errors.New("transaction not found")

	transactionsTableName = environment.GetString("TRANSACTIONS_TABLE_NAME", "transactions")
	dynamoClient          dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutTransaction saves a transaction in DynamoDB
func PutTransaction(ctx context.Context, transaction *models.Transaction) error {
	item, err := dynamodbattribute.MarshalMap(transaction)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(transactionsTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadTransaction loads a Transaction from DynamoDB
func LoadTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	key := map[string]*dynamodb.AttributeValue{
		"transaction_id": {S: aws.String(transactionID)},
	}

	params := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(transactionsTableName),
		ConsistentRead: aws.Bool(true),
	}

	response, err := dynamoClient.GetItemWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, ErrTransactionNotFound
	}

	var transaction models.Transaction
	err = dynamodbattribute.UnmarshalMap(response.Item, &transaction)

	return &transaction, err
}
