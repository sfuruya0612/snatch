package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func StartSession(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	err := aws.StartSession(profile, region)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func SendCommand(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	args := c.Args()
	file := c.String("file")
	if len(args) == 0 && len(file) == 0 {
		return fmt.Errorf("Args or file is required")
	}

	id := c.String("instanceid")
	tag := c.String("tag")
	if len(id) == 0 && len(tag) == 0 {
		return fmt.Errorf("Instance id or tag is required")
	}

	err := aws.SendCommand(profile, region, file, id, tag, args)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
