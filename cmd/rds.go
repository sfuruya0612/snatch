package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/rds"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetRdsList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	input := &rds.DescribeDBInstancesInput{}

	rds := saws.NewRdsSess(profile, region)
	resources, err := rds.DescribeDBInstances(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintDBInstances(os.Stdout, resources); err != nil {
		return fmt.Errorf("Failed to print resources")
	}

	return nil
}
