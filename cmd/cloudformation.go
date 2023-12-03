package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var CloudFormation = &cli.Command{
	Name:    "cloudformation",
	Aliases: []string{"cfn"},
	Usage:   "Get a list of stacks",
	Action: func(c *cli.Context) error {
		return getStackList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:      "events",
			Usage:     "Get stack events",
			ArgsUsage: "[ --name | -n ] <CfnStackName>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Set stack name",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				return getStackEvents(c.String("profile"), c.String("region"), c.String("name"))
			},
		},
		{
			Name:      "outputs",
			Usage:     "Get stack outputs",
			ArgsUsage: "[ --name | -n ] <CfnStackName>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Set stack name",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				return getStackEvents(c.String("profile"), c.String("region"), c.String("name"))
			},
		},
	},
}

func getStackList(profile, region string) error {
	client := saws.NewCfnClient(profile, region)

	resources, err := client.DescribeStacks(&cloudformation.DescribeStacksInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintStacks(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}

func getStackEvents(profile, region, name string) error {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(name),
	}

	client := saws.NewCfnClient(profile, region)
	events, err := client.DescribeStackEvents(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintEvents(os.Stdout, events); err != nil {
		return fmt.Errorf("failed to print events")
	}

	return nil
}
