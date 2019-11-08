package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListElasticache(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	elasticache := aws.NewElastiCacheSess(profile, region)
	if err := elasticache.DescribeCacheClusters(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func ListReplicationGroups(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	elasticache := aws.NewElastiCacheSess(profile, region)
	if err := elasticache.DescribeReplicationGroups(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
