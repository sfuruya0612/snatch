package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Ec2 = &cli.Command{
	Name:      "ec2",
	Usage:     "Get a list of EC2 resources",
	ArgsUsage: "[ --tag | -t ] <Key:Value>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "The Key-Value of the tag to filter",
		},
	},
	Action: func(c *cli.Context) error {
		return getEc2List(c.String("profile"), c.String("region"), c.String("tag"))
	},
}

func getEc2List(profile, region, tag string) error {
	input := &ec2.DescribeInstancesInput{}
	if len(tag) > 0 {
		if !strings.Contains(tag, ":") {
			return fmt.Errorf("tag is different (e.g. Name:hogehoge)")
		}

		spl := strings.Split(tag, ":")
		if len(spl) == 0 {
			return fmt.Errorf("parse tag=%s", tag)
		}

		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("tag:" + spl[0]),
			Values: []*string{aws.String(spl[1])},
		})
	}

	client := saws.NewEc2Sess(profile, region)
	resources, err := client.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintInstances(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
