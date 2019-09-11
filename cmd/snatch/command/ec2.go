package command

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

// ListEc2 returns nil
func ListEc2(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	err := aws.DescribeInstances(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
