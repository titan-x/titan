// Package aws provides AWS implementation of data interfaces.
//
// AWS Go SDK: https://github.com/aws/aws-sdk-go
// Go SDK API Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/
// Dev Guide: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.Modifying.html
// REST API Ref: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Operations.html
package aws

import (
	"fmt"
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
func (db *DynamoDB) Seed(overwrite bool, jwtPass string) error {
	cred, err := db.DB.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("dynamodb: failed to initialize: %v", err)
	}

	log.Printf("dynamodb: initialized with region: %v, access key ID: %v, endpoint: %v", *(db.DB.Config.Region), cred.AccessKeyID, db.DB.Config.Endpoint)

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
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("ID"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("Email"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("ID"),
					KeyType:       aws.String("HASH"),
				},
			},
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String("Email"),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("Email"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &dynamodb.Projection{
						NonKeyAttributes: []*string{
							aws.String("Email"),
						},
						ProjectionType: aws.String("INCLUDE"),
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
				},
			},
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
	if err := data.SeedInit(jwtPass); err != nil {
		return err
	}

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
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Printf("dynamodb: error: %v", err)
		return nil, false
	}
	if len(res.Item) == 0 {
		return nil, false
	}

	var user models.User
	if err := dynamodbattribute.UnmarshalMap(res.Item, &user); err != nil {
		log.Printf("dynamodb: getbyid error: %v", err)
		return nil, false
	}

	return &user, true
}

// GetByEmail retrieves a user by e-mail with OK indicator.
func (db *DynamoDB) GetByEmail(email string) (u *models.User, ok bool) {
	res, err := db.DB.Query(&dynamodb.QueryInput{
		// ConsistentRead: aws.Bool(true),
		TableName:              aws.String("users"),
		IndexName:              aws.String("Email"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		KeyConditionExpression: aws.String("Email = :Email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		log.Printf("dynamodb: getbymail error: %v", err)
		return nil, false
	}

	var user models.User
	if err := dynamodbattribute.UnmarshalMap(res.Items[0], &user); err != nil {
		log.Printf("dynamodb: getbyid error: %v", err)
		return nil, false
	}

	return &user, true
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

	item, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return err
	}

	_, err = db.DB.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}
