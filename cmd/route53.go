package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListHostedZones(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	route53 := aws.NewRoute53Sess(profile, region)
	if err := route53.ListHostedZones(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
