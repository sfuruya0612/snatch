package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetASGList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewAsgSess(profile, region)
	groups, err := client.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintGroups(os.Stdout, groups); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func UpdateCapacity(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	name := c.String("name")
	disire := c.Int64("disire")
	min := c.Int64("min")
	max := c.Int64("max")

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(disire),
		MinSize:              aws.Int64(min),
		MaxSize:              aws.Int64(max),
	}

	client := saws.NewAsgSess(profile, region)
	if err := client.UpdateAutoScalingGroup(input); err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("\n\x1b[35m%v groups is updated\x1b[0m", name)

	return nil
}
