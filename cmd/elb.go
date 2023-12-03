package cmd

import (
	"fmt"
	"os"

	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Elb = &cli.Command{
	Name:  "elb",
	Usage: "Get a list of ELB",
	Action: func(c *cli.Context) error {
		return getElbList(c.String("profile"), c.String("region"))
	},
}

func getElbList(profile, region string) error {
	v1c := saws.NewElbClient(profile, region)
	lb, err := v1c.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	v2c := saws.NewElbV2Client(profile, region)
	lbv2, err := v2c.DescribeLoadBalancersV2(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	lb = append(lb, lbv2...)

	if err := saws.PrintBalancers(os.Stdout, lb); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
