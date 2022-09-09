package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/route53"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Route53 = &cli.Command{
	Name:  "route53",
	Usage: "Get a list of Rotue53 Record resources",
	Action: func(c *cli.Context) error {
		return getRecordsList(c.String("profile"), c.String("region"))
	},
}

func getRecordsList(profile, region string) error {
	client := saws.NewRoute53Sess(profile, region)

	resources, err := client.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintRecords(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
