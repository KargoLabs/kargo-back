package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamoClient dynamodbiface.DynamoDBAPI

	tableConfigurations = []TableConfiguration{
		TableConfiguration{
			TableName:  "clients",
			TableIndex: "client_id",
		},
		TableConfiguration{
			TableName:  "partners",
			TableIndex: "partner_id",
		},
		TableConfiguration{
			TableName:  "orders",
			TableIndex: "order_id",
		},
		TableConfiguration{
			TableName:  "trips",
			TableIndex: "trip_id",
		},
		TableConfiguration{
			TableName:  "commodities",
			TableIndex: "commodity_id",
		},
		TableConfiguration{
			TableName:  "payment-accounts",
			TableIndex: "payment_account_id",
		},
		TableConfiguration{
			TableName:  "payment-transactions",
			TableIndex: "payment_transaction_id",
		},
		TableConfiguration{
			TableName:  "payment-methods",
			TableIndex: "payment_method_id",
		},
	}
)

// TableConfiguration is struct handler for values necessary to create dynamo table
type TableConfiguration struct {
	TableName  string
	TableIndex string
}

func init() {
	sess := session.Must(session.NewSession())
	dynamoClient = dynamodb.New(sess)
}

func createTable(tableConfiguration TableConfiguration) {
	input := &dynamodb.CreateTableInput{
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(tableConfiguration.TableIndex),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableConfiguration.TableName),
	}

	result, err := dynamoClient.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}

func main() {
	for _, tableConfiguration := range tableConfigurations {
		createTable(tableConfiguration)
	}
}
