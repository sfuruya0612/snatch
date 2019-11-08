package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func DescribeLogGroups(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	flag := c.GlobalBool("f")

	err := aws.DescribeLogGroups(profile, region, flag)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
