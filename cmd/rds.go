package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListRds(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	rds := aws.NewRdsSess(profile, region)
	if err := rds.DescribeDBInstances(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
