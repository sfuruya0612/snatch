package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/elasticache"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var ElastiCache = &cli.Command{
	Name:    "elasticache",
	Aliases: []string{"ec"},
	Usage:   "Get a list of ElastiCache Cluster resources",
	Action: func(c *cli.Context) error {
		return getEcClusterList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:  "node",
			Usage: "Get a list of ElastiCache Node resources",
			Action: func(c *cli.Context) error {
				return getEcClusterList(c.String("profile"), c.String("region"))
			},
		},
	},
}

func getEcClusterList(profile, region string) error {
	client := saws.NewElastiCacheSess(profile, region)

	resources, err := client.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintCacheClusters(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}

func getEcGroupsList(profile, region string) error {
	client := saws.NewElastiCacheSess(profile, region)

	resources, err := client.DescribeReplicationGroups(&elasticache.DescribeReplicationGroupsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintRepricationGroups(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
