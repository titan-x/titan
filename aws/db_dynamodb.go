// Package aws provides AWS implementation of data interfaces.
//
// AWS Go SDK: https://github.com/aws/aws-sdk-go
// Go SDK API Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/
// Dev Guide: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.Modifying.html
// REST API Ref: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/titan-x/titan"
)

// DynamoDB implementation for DB interface.
type DynamoDB struct {
	DB      *dynamodb.DynamoDB
	Tables  []string
	Config  *aws.Config
	Session *session.Session
}

// NewDynamoDB creates a new AWS DynamoDB instance.
// region = Optional region setting. Will overwrite AWS_REGION env var if available.
// endpoint = Optional endpoint URL setting. Useful for specifying local/development service URL.
func NewDynamoDB(region string, endpoint string) *DynamoDB {
	db := DynamoDB{}
	db.Tables = []string{"users"}

	// carefully crafting config elements not to mess with the defaults
	if region != "" || endpoint != "" {
		db.Config = aws.NewConfig()

		if region != "" {
			db.Config.WithRegion(region)
		}
		if endpoint != "" {
			db.Config.WithEndpoint(endpoint)
		}

		db.Session = session.New(db.Config)
	} else {
		db.Session = session.New()
	}

	db.DB = dynamodb.New(db.Session)

	return &db
}

func (db *DynamoDB) listTables() ([]string, error) {
	res, err := db.DB.ListTables(&dynamodb.ListTablesInput{Limit: aws.Int64(100)})
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
		if _, err := db.DB.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	for _, tbl := range tables {
		if err := db.DB.WaitUntilTableNotExists(&dynamodb.DescribeTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	return nil
}

// Seed creates and populates the database, overwriting existing data if specified.
func (db *DynamoDB) Seed(overwrite bool) error {
	if !overwrite {
		if tbls, err := db.listTables(); err != nil {
			return err
		} else if len(tbls) != 0 {
			return nil
		}
	}

	if err := db.deleteTables(); err != nil {
		return err
	}

	for _, tbl := range db.Tables {
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
			// },
			// StreamSpecification: &dynamodb.StreamSpecification{
			// 	StreamEnabled:  aws.Bool(true),
			// 	StreamViewType: aws.String("StreamViewType"),
			// },
		}

		_, err := db.DB.CreateTable(tableParams)
		if err != nil {
			return err
		}

		// tables with secondary indexes need to be created sequentially so wait till table is ready
		if err := db.DB.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves a user by ID with OK indicator.
func (db *DynamoDB) GetByID(id string) (u *titan.User, ok bool) {
	return nil, false
}

// GetByMail retrieves a user by e-mail with OK indicator.
func (db *DynamoDB) GetByMail(email string) (u *titan.User, ok bool) {
	return nil, false
}

// SaveUser creates or updates a user. Upon creation, users are assigned a unique ID.
func (db *DynamoDB) SaveUser(u *titan.User) error {
	return nil
}
