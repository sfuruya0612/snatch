package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Before(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	return nil
}
