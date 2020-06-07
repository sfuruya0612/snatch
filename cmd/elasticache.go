package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/elasticache"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetEcClusterList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewElastiCacheSess(profile, region)
	resources, err := client.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintCacheClusters(os.Stdout, resources); err != nil {
		return fmt.Errorf("Failed to print resources")
	}

	return nil
}

func GetEcGroupsList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewElastiCacheSess(profile, region)
	resources, err := client.DescribeReplicationGroups(&elasticache.DescribeReplicationGroupsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintRepricationGroups(os.Stdout, resources); err != nil {
		return fmt.Errorf("Failed to print resources")
	}

	return nil
}
