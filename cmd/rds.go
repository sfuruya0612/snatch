package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/rds"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Rds = &cli.Command{
	Name:  "rds",
	Usage: "Get a list of RDS instance",
	Action: func(c *cli.Context) error {
		return getRdsList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:    "cluster",
			Aliases: []string{"c"},
			Usage:   "Get a list of RDS cluster",
			Action: func(c *cli.Context) error {
				return getRdsClusterList(c.String("profile"), c.String("region"))
			},
			Subcommands: []*cli.Command{
				{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Usage:   "Get a list of RDS cluster endpoint",
					Action: func(c *cli.Context) error {
						return getRdsClusterEndpoints(c.String("profile"), c.String("region"))
					},
				},
			},
		},
		{
			Name:    "s3export",
			Aliases: []string{"e"},
			Usage:   "Get a list of RDS S3 export",
			Action: func(c *cli.Context) error {
				return getRdsS3ExportList(c.String("profile"), c.String("region"))
			},
		},
	},
}

func getRdsList(profile, region string) error {
	c := saws.NewRdsClient(profile, region)

	instances, err := c.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintDBInstances(os.Stdout, instances); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func getRdsClusterList(profile, region string) error {
	c := saws.NewRdsClient(profile, region)

	clusters, err := c.DescribeDBClusters(&rds.DescribeDBClustersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintDBClusters(os.Stdout, clusters); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func getRdsClusterEndpoints(profile, region string) error {
	c := saws.NewRdsClient(profile, region)

	endpoints, err := c.DescribeDBClusterEndpoints(&rds.DescribeDBClusterEndpointsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintDBClusterEndpoints(os.Stdout, endpoints); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func getRdsS3ExportList(profile, region string) error {
	c := saws.NewRdsClient(profile, region)

	exports, err := c.DescribeExportTasks(&rds.DescribeExportTasksInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintExportTasks(os.Stdout, exports); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
