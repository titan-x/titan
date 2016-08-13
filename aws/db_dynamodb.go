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
