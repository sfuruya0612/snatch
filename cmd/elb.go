package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetElbList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

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

func GetElbV2List(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

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
