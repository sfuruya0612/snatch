package command

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

// ListEc2 returns ec2.DescribeInstances
func ListEc2(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	list := aws.DescribeInstances(profile, region)

	fmt.Println(list)

	return nil
}
