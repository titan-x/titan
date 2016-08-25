// Package aws provides AWS implementation of data interfaces.
//
// AWS Go SDK: https://github.com/aws/aws-sdk-go
// Go SDK API Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/
// Dev Guide: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.Modifying.html
// REST API Ref: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html
package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/neptulon/shortid"
	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/models"
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

	// create the tables
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

		if _, err := db.DB.CreateTable(tableParams); err != nil {
			return err
		}

		// tables with secondary indexes need to be created sequentially so wait till table is ready
		if err := db.DB.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: aws.String(tbl)}); err != nil {
			return err
		}
	}

	// insert the seed data
	for _, u := range data.SeedUsers {
		if err := db.SaveUser(&u); err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves a user by ID with OK indicator.
func (db *DynamoDB) GetByID(id string) (u *models.User, ok bool) {
	res, err := db.DB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("1"),
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		log.Printf("dynamodb: error: %v", err)
		return nil, false
	}
	if len(res.Item) == 0 {
		return nil, false
	}

	log.Printf("dynamodb: getbyid: consumed capacity: %v", res.ConsumedCapacity)

	var user models.User
	if err := dynamodbattribute.UnmarshalMap(res.Item, &user); err != nil {
		log.Printf("dynamodb: error: %v", err)
		return nil, false
	}

	return &user, true
}

// GetByEmail retrieves a user by e-mail with OK indicator.
func (db *DynamoDB) GetByEmail(email string) (u *models.User, ok bool) {
	// db.DB.Query(nil)
	return nil, false
}

// SaveUser creates or updates a user. Upon creation, users are assigned a unique ID.
func (db *DynamoDB) SaveUser(u *models.User) error {
	if u.ID == "" {
		id, err := shortid.ID(64)
		if err != nil {
			return err
		}

		u.ID = id
	}

	res, err := db.DB.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(u.ID),
			},
		},
		// UpdateExpression: aws.String("SET a=:value1, b=:value2"),
	})
	if err != nil {
		return err
	}

	log.Printf("dynamodb: saveuser: consumed capacity: %v", res.ConsumedCapacity)

	return nil
}
