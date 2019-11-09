package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListElbV2(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	elbv2 := aws.NewElbV2Sess(profile, region)
	if err := elbv2.DescribeLoadBalancersV2(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
