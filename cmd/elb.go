package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Elb = &cli.Command{
	Name:  "elb",
	Usage: "Get a list of ELB(Classic) resources.",
	Action: func(c *cli.Context) error {
		return getElbList(c.String("profile"), c.String("region"))
	},
}

var ElbV2 = &cli.Command{
	Name:  "elbv2",
	Usage: "Get a list of ELB(Application & Network) resources",
	Action: func(c *cli.Context) error {
		return getElbV2List(c.String("profile"), c.String("region"))
	},
}

func getElbList(profile, region string) error {
	client := saws.NewElbSess(profile, region)

	resources, err := client.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintBalancers(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}

func getElbV2List(profile, region string) error {
	client := saws.NewElbV2Sess(profile, region)

	resources, err := client.DescribeLoadBalancersV2(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintBalancersV2(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
