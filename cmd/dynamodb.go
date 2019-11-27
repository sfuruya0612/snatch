package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ScanDb(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	input := &dynamodb.ListTablesInput{}

	dynamodb := saws.NewDynamoDbSess(profile, region)
	tname, err := dynamodb.ListTables(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Println(tname)

	return nil
}
