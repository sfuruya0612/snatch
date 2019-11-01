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

	err := aws.DescribeCacheClusters(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func ListReplicationGroups(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	err := aws.DescribeReplicationGroups(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
