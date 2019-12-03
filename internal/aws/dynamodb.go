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

// ListTables return []string (dynamodb.ListTablesOutput.TableNames)
// input dynamodb.ListTablesInput
func (c *DynamoDB) ListTables(input *dynamodb.ListTablesInput) ([]string, error) {
	output, err := c.Client.ListTables(input)
	if err != nil {
		return nil, fmt.Errorf("List tables: %v", err)
	}

	tables := []string{}
	for _, t := range output.TableNames {
		tables = append(tables, *t)
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("No tables")
	}

	return tables, nil
}
