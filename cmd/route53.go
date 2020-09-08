package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/route53"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetRecordsList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

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
