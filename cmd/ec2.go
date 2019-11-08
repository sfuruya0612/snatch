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

	ec2 := aws.NewEc2Sess(profile, region)
	if err := ec2.DescribeInstances(tag); err != nil {
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

	ec2 := aws.NewEc2Sess(profile, region)
	if err := ec2.GetConsoleOutput(id); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
