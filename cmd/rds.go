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

	client := saws.NewRdsSess(profile, region)
	resources, err := client.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintDBInstances(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}

func GetRdsClusterList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewRdsSess(profile, region)

	clusters, err := client.DescribeDBClusters(&rds.DescribeDBClustersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintDBClusters(os.Stdout, clusters); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
