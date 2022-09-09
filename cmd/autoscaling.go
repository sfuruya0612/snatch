package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli/v2"
)

var AutoScaling = &cli.Command{
	Name:    "autoscaling",
	Aliases: []string{"as"},
	Usage:   "Get a list of EC2 AutoScalingGroups",
	Action: func(c *cli.Context) error {
		return getASGList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:      "capacity",
			Aliases:   []string{"cap"},
			Usage:     "Update autoscaling group capacity",
			ArgsUsage: "[ --name | -n ] <AutoScalingGroupName> [ --desired ] <CapacityNum> [ --min ]  <CapacityNum> [ --max ] <CapacityNum> ",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Set target autoscaling group name",
					Required: true,
				},
				&cli.Int64Flag{
					Name:  "desired",
					Usage: "Set desired capacity",
				},
				&cli.Int64Flag{
					Name:  "min",
					Usage: "Set minimum capacity",
				},
				&cli.Int64Flag{
					Name:  "max",
					Usage: "Set maximum capacity",
				},
			},
			Action: func(c *cli.Context) error {
				return updateCapacity(c.String("profile"), c.String("region"), c.String("name"), c.Int64("desired"), c.Int64("min"), c.Int64("max"))
			},
		},
	},
}

func getASGList(profile, region string) error {
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

func updateCapacity(profile, region, name string, desired, min, max int64) error {
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
