package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListElb(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	elb := aws.NewElbSess(profile, region)
	if err := elb.DescribeLoadBalancers(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
