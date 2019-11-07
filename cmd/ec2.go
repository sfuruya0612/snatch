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

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	if err := aws.DescribeInstances(profile, region, tag); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func GetEc2Log(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	id := c.String("instanceid")
	if len(id) == 0 {
		return fmt.Errorf("--instanceid or -i option is required")
	}

	fmt.Printf("\x1b[32mAWS_PROFILE: %s , REGION: %s\x1b[0m\n", profile, region)

	if err := aws.GetConsoleOutput(profile, region, id); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
