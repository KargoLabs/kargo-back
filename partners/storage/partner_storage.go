package storage

import (
	"context"
	"errors"
	"kargo-back/shared/environment"
	models "kargo-back/shared/partners-models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// ErrPartnerNotFound when no Partner was found
	ErrPartnerNotFound = errors.New("partnert not found")

	partnersTableName = environment.GetString("PARTNERS_TABLE_NAME", "partners")
	dynamoClient      dynamodbiface.DynamoDBAPI
)

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

// PutPartner saves a Partnert in DynamoDB
func PutPartner(ctx context.Context, partner *models.Partner) error {
	item, err := dynamodbattribute.MarshalMap(partner)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(partnersTableName),
	}

	_, err = dynamoClient.PutItemWithContext(ctx, params)

	return err
}

// LoadPartner loads a Partner from DynamoDB
func LoadPartner(ctx context.Context, partnerID string) (*models.Partner, error) {
	key := map[string]*dynamodb.AttributeValue{
		"partner_id": {S: aws.String(partnerID)},
	}

	params := &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(partnersTableName),
		ConsistentRead: aws.Bool(true),
	}

	response, err := dynamoClient.GetItemWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, ErrPartnerNotFound
	}

	var partner models.Partner
	err = dynamodbattribute.UnmarshalMap(response.Item, &partner)

	return &partner, nil
}
