package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
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
	desired := c.Int64("desired")
	min := c.Int64("min")
	max := c.Int64("max")

	if len(name) == 0 {
		return fmt.Errorf("--name or -n option is required")
	}

	// desired, min, maxの値の関係性を確認
	if max < min || desired < min || max < desired {
		return fmt.Errorf("capacity options number have incorrect relationship")
	}

	client := saws.NewAsgSess(profile, region)

	// 実施前のパラメータを取得
	before, err := client.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: aws.StringSlice([]string{name}),
	})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("Tatget: %v, Terminate Policies: %v\n\n", name, before[0].TerminationPolicies)
	fmt.Printf("Before Parameters:\n\tDesired: %v, MinSize: %v, MaxSize: %v\n", before[0].Desired, before[0].Min, before[0].Max)
	fmt.Printf("After  Parameters:\n\tDesired: %v, MinSize: %v, MaxSize: %v\n", desired, min, max)

	// Capacityに0が指定された場合、警告文を出しておく
	if desired == 0 || min == 0 || max == 0 {
		fmt.Printf("\n\x1b[35mAutoScaling Group capacity is 0\x1b[0m\n")
	}

	if !util.Confirm(name) {
		fmt.Printf("Cancel update autoscaling group: %v\n", name)
		return nil
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(desired),
		MinSize:              aws.Int64(min),
		MaxSize:              aws.Int64(max),
	}

	if err := client.UpdateAutoScalingGroup(input); err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("\n%v groups is updated\n", name)

	return nil
}
