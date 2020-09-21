package storage

import (
	"errors"
	"kargo-back/shared/environment"

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
