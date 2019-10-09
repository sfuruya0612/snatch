package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListEc2(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	tag := c.String("tag")

	err := aws.DescribeInstances(profile, region, tag)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
