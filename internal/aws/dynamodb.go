package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDB client struct
type DynamoDB struct {
	Client *dynamodb.DynamoDB
}

// NewDynamoDbSess return DynamoDB struct initialized
func NewDynamoDbSess(profile, region string) *DynamoDB {
	return &DynamoDB{
		Client: dynamodb.New(getSession(profile, region)),
	}
}

// ListTables return &dynamodb.ListTablesOutput.TableNames
func (c *DynamoDB) ListTables(input *dynamodb.ListTablesInput) ([]*string, error) {
	output, err := c.Client.ListTables(input)
	if err != nil {
		return nil, fmt.Errorf("List tables: %v", err)
	}

	tname := output.TableNames

	if len(tname) == 0 {
		return nil, fmt.Errorf("No tables")
	}

	return tname, nil
}
