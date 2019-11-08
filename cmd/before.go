package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func Before(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	return nil
}
