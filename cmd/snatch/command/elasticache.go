package command

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

// ListElasticache returns error
func ListElasticache(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	err := aws.DescribeCacheClusters(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

// ListReplicationGroup returns error
func ListReplicationGroups(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	err := aws.DescribeReplicationGroups(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
