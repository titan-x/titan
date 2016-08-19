package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDB implementation for DB interface.
type DynamoDB struct {
	config  *aws.Config
	session *session.Session
	db      *dynamodb.DynamoDB
}

// NewDynamoDB creates a new AWS DynamoDB instance.
// region = Optional region setting. Will overwrite AWS_REGION env var if available.
// endpoint = Optional endpoint URL setting. Useful for specifying local/development service URL.
func NewDynamoDB(region string, endpoint string) *DynamoDB {
	db := DynamoDB{}

	// carefully crafting config elements not to mess with the defaults
	if region != "" || endpoint != "" {
		db.config = aws.NewConfig()

		if region != "" {
			db.config.WithRegion(region)
		}
		if endpoint != "" {
			db.config.WithEndpoint(endpoint)
		}

		db.session = session.New(db.config)
	} else {
		db.session = session.New()
	}

	db.db = dynamodb.New(db.session)

	return &db
}

func (db *DynamoDB) listTables() ([]string, error) {
	res, err := db.db.ListTables(&dynamodb.ListTablesInput{Limit: aws.Int64(100)})
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, name := range res.TableNames {
		names = append(names, *name)
	}

	return names, nil
}

func (db *DynamoDB) deleteTables() error {
	tables, err := db.listTables()
	if err != nil {
		return err
	}

	for _, tbl := range tables {
		if _, err := db.db.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	for _, tbl := range tables {
		if err := db.db.WaitUntilTableNotExists(&dynamodb.DescribeTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	return nil
}

// Seed creates and populates the database, overwriting existing data if specified.
func (db *DynamoDB) Seed(overwrite bool) error {
	tables := []string{"users"}

	if overwrite {
		if err := db.deleteTables(); err != nil {
			return err
		}
	}

	for _, tbl := range tables {
		tableParams := &dynamodb.CreateTableInput{
			TableName: aws.String(tbl),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
			// GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			// 	{
			// 		IndexName: aws.String("IndexName"),
			// 		KeySchema: []*dynamodb.KeySchemaElement{
			// 			{
			// 				AttributeName: aws.String("KeySchemaAttributeName"),
			// 				KeyType:       aws.String("KeyType"),
			// 			},
			//
			// 		},
			// 		Projection: &dynamodb.Projection{
			// 			NonKeyAttributes: []*string{
			// 				aws.String("NonKeyAttributeName"),
			//
			// 			},
			// 			ProjectionType: aws.String("ProjectionType"),
			// 		},
			// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			// 			ReadCapacityUnits:  aws.Int64(1),
			// 			WriteCapacityUnits: aws.Int64(1),
			// 		},
			// 	},
			// 	// More values...
			// },
			// LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex{
			// 	{
			// 		IndexName: aws.String("IndexName"),
			// 		KeySchema: []*dynamodb.KeySchemaElement{
			// 			{
			// 				AttributeName: aws.String("KeySchemaAttributeName"),
			// 				KeyType:       aws.String("KeyType"),
			// 			},
			//
			// 		},
			// 		Projection: &dynamodb.Projection{
			// 			NonKeyAttributes: []*string{
			// 				aws.String("NonKeyAttributeName"),
			//
			// 			},
			// 			ProjectionType: aws.String("ProjectionType"),
			// 		},
			// 	},
			// 	// More values...
			// },
			// StreamSpecification: &dynamodb.StreamSpecification{
			// 	StreamEnabled:  aws.Bool(true),
			// 	StreamViewType: aws.String("StreamViewType"),
			// },
		}

		_, err := db.db.CreateTable(tableParams)
		if err != nil {
			return err
		}

		// tables with secondary indexes need to be created sequentially so wait till table is ready
		if err := db.db.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	return nil
}
