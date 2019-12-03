package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetStacksList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewCfnSess(profile, region)
	resources, err := client.DescribeStacks(&cloudformation.DescribeStacksInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintStacks(os.Stdout, resources); err != nil {
		return fmt.Errorf("Failed to print resources")
	}

	return nil
}
