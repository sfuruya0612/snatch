package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func DescribeStacks(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	cloudformation := aws.NewCfnSess(profile, region)
	if err := cloudformation.DescribeStacks(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
