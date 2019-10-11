package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListBuckets(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	err := aws.ListBuckets(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
