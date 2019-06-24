package command

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

// ListRds returns nil
func ListRds(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	list := aws.DescribeDBInstances(profile, region)

	fmt.Println(list)

	return nil
}
