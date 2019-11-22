package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetEc2List(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	tag := c.String("tag")

	input := &ec2.DescribeInstancesInput{}

	if len(tag) > 0 {
		if !strings.Contains(tag, ":") {
			return fmt.Errorf("%v", "tag is different (e.g. Name:hogehoge)")
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

	ec2 := saws.NewEc2Sess(profile, region)
	if err := ec2.DescribeInstances(input); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func GetEc2SystemLog(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	id := c.String("instanceid")
	if len(id) == 0 {
		return fmt.Errorf("--instanceid or -i option is required")
	}

	input := &ec2.GetConsoleOutputInput{
		InstanceId: aws.String(id),
	}

	ec2 := saws.NewEc2Sess(profile, region)
	if err := ec2.GetConsoleOutput(input); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
