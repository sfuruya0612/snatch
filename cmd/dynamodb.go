package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetTablesList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewDynamoDbSess(profile, region)
	tables, err := client.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Println(tables)

	return nil
}
