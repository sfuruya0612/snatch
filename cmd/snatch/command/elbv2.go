package command

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListElbv2(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	err := aws.DescribeLoadBalancersv2(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
