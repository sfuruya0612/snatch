package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

func GetEc2List(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	tag := c.String("tag")
	short := c.Bool("short")

	input := &ec2.DescribeInstancesInput{}
	if len(tag) > 0 {
		if !strings.Contains(tag, ":") {
			return fmt.Errorf("Tag is different (e.g. Name:hogehoge)")
		}

		spl := strings.Split(tag, ":")
		if len(spl) == 0 {
			return fmt.Errorf("Parse tag=%s", tag)
		}

		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("tag:" + spl[0]),
			Values: []*string{aws.String(spl[1])},
		})
	}

	if short {
		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("running")},
		})
	}

	client := saws.NewEc2Sess(profile, region)
	resources, err := client.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintInstances(os.Stdout, resources); err != nil {
		return fmt.Errorf("Failed to print resources")
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

	client := saws.NewEc2Sess(profile, region)
	output, err := client.GetConsoleOutput(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	d, err := util.DecodeString(*output.Output)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Println(d)

	return nil
}

func TerminateEc2(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	id := c.String("instanceid")
	if len(id) == 0 {
		return fmt.Errorf("--instanceid or -i option is required")
	}

	if !util.Confirm(id) {
		fmt.Printf("\nCancel terminate: %v\n", id)
		return nil
	}

	input := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{id}),
	}

	client := saws.NewEc2Sess(profile, region)
	output, err := client.TerminateInstances(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, o := range output.TerminatingInstances {
		fmt.Printf("\nInstanceId %v is terminated\n", *o.InstanceId)
	}

	return nil
}
