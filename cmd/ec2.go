package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Ec2 = &cli.Command{
	Name:      "ec2",
	Usage:     "Get a list of EC2 instance",
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

		input.Filters = append(input.Filters, types.Filter{
			Name:   aws.String("tag:" + spl[0]),
			Values: []string{spl[1]},
		})
	}

	c := saws.NewEc2Client(profile, region)
	instances, err := c.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintInstances(os.Stdout, instances); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
