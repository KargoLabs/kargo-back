package storage

import (
	"context"
	"errors"
	models "kargo-back/models/cards"
	"kargo-back/shared/environment"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrCardNotFound when no card was found
	ErrCardNotFound     = errors.New("card not found")
	cardsTableName      = environment.GetString("CARDS_TABLE_NAME", "cards")
	cardsUserIDndexName = environment.GetString("CARD_USER_ID_INDEX_NAME", "card_user-id-index")
	dynamoClient        dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutCard saves a Card in DynamoDB
func PutCard(ctx context.Context, card *models.Card) error {
	item, err := dynamodbattribute.MarshalMap(card)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(cardsTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadUserCards load all cards related to an user from Dynamodb
func LoadUserCards(ctx context.Context, userID string) ([]*models.Card, error) {
	keyCondition := expression.KeyEqual(expression.Key("user_id"), expression.Value(userID))
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
		IndexName:                 aws.String(cardsUserIDndexName),
		TableName:                 aws.String(cardsTableName),
	}

	response, err := dynamoClient.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	cards := []*models.Card{}

	// Return empty list
	if len(response.Items) == 0 {
		return cards, nil
	}

	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &cards)
	if err != nil {
		return nil, err
	}

	return cards, nil
}
