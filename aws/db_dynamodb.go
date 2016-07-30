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
func NewDynamoDB(region string) *DynamoDB {
	db := DynamoDB{}

	if region != "" {
		db.config = aws.NewConfig().WithRegion(region)
		db.session = session.New(db.config)
	} else {
		db.session = session.New()
	}

	db.db = dynamodb.New(db.session)

	return &db
}
