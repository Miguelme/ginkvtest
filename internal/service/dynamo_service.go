package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoService struct {
	client *dynamodb.DynamoDB
}

// NewDynamoService initializes the DynamoDB client and ensures the table and data exist.
func NewDynamoService() *DynamoService {

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""), // Dummy credentials

	}))

	client := dynamodb.New(sess)
	service := &DynamoService{client: client}

	// Ensure the table and initial data exist
	service.ensureTableAndDataExists()

	return service
}

// ensureTableAndDataExists ensures the table exists and inserts initial data if needed.
func (s *DynamoService) ensureTableAndDataExists() {
	tableName := "KeyValueTable"

	// 1. Check if the table exists

	_, err := s.client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		// Table does not exist, create it
		fmt.Println("Creating DynamoDB table:", tableName)

		_, err := s.client.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("key"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("key"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		})

		if err != nil {
			fmt.Println("Failed to create DynamoDB table:", err)
			return
		}
		fmt.Println("DynamoDB table created successfully:", tableName)
	} else {
		fmt.Println("DynamoDB table already exists:", tableName)
	}

	// 2. Insert initial data if the table is empty
	fmt.Println("Checking for existing data...")

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {S: aws.String("test_key")},
		},
	}

	result, err := s.client.GetItem(input)
	if err != nil {
		fmt.Println("Error checking existing data:", err)
		return
	}

	if result.Item == nil {
		// Insert the initial data
		fmt.Println("Inserting initial data into the table...")

		_, err := s.client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"key":   {S: aws.String("test_key")},
				"value": {S: aws.String("test_value")},
			},
		})

		if err != nil {
			fmt.Println("Failed to insert initial data:", err)
			return
		}
		fmt.Println("Initial data inserted successfully.")
	} else {
		fmt.Println("Initial data already exists.")
	}
}

// GetValueByKey retrieves the value associated with a key from DynamoDB.
func (s *DynamoService) GetValueByKey(key string) (string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("KeyValueTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {S: aws.String(key)},
		},
	}

	result, err := s.client.GetItem(input)
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", fmt.Errorf("key not found")
	}

	return *result.Item["value"].S, nil
}
