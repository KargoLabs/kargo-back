package storage

import (
	"context"
	"errors"
	models "kargo-back/models/accounts"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	// ErrAccountNotFound when no account was found
	ErrAccountNotFound = errors.New("account not found")
	// ErrAccountNotBelongPartner when the account does not belong to the requesting partner
	ErrAccountNotBelongPartner = errors.New("account does not belong to partner")

	accountsTableName          = environment.GetString("ACCOUNTS_TABLE_NAME", "accounts")
	accountsPartnerIDIndexName = environment.GetString("ACCOUNT_PARTNER_ID_INDEX_NAME", "account_partner-id-index")
	dynamoClient               dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutAccount saves an Account in DynamoDB
func PutAccount(ctx context.Context, account *models.Account) error {
	item, err := dynamodbattribute.MarshalMap(account)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(accountsTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadAccount loads an account from DynamoDB based on its partner
func LoadAccount(ctx context.Context, partnerID string) (*models.Account, error) {
	keyCondition := expression.KeyEqual(expression.Key("partner_id"), expression.Value(partnerID))
	dynamoExpression := expression.NewBuilder().WithKeyCondition(keyCondition)

	dynamoquery, err := dynamoExpression.Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  dynamoquery.Names(),
		ExpressionAttributeValues: dynamoquery.Values(),
		KeyConditionExpression:    dynamoquery.KeyCondition(),
		FilterExpression:          dynamoquery.Filter(),
		IndexName:                 aws.String(accountsPartnerIDIndexName),
		TableName:                 aws.String(accountsTableName),
	}

	response, err := dynamoClient.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, ErrAccountNotFound
	}

	var account []*models.Account

	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &account)

	return account[0], nil
}

// DeleteAccount deletes a partner account
func DeleteAccount(ctx context.Context, partnerID string) (*models.Account, error) {
	account, err := LoadAccount(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"account_id": {
				S: aws.String(account.AccountID),
			},
		},
		TableName: aws.String(accountsTableName),
	}

	_, err = dynamoClient.DeleteItem(input)
	if err != nil {
		return nil, err
	}

	return account, nil
}
