package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/rds"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Rds = &cli.Command{
	Name:  "rds",
	Usage: "Get a list of RDS resources",
	Action: func(c *cli.Context) error {
		return getRdsList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:  "cluster",
			Usage: "Get a list of RDS Cluster resources",
			Action: func(c *cli.Context) error {
				return getRdsClusterList(c.String("profile"), c.String("region"))
			},
		},
	},
}

func getRdsList(profile, region string) error {
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

func getRdsClusterList(profile, region string) error {
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
