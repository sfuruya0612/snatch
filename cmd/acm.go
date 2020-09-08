package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/acm"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetCertificatesList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewAcmSess(profile, region)
	resources, err := client.ListCertificates(&acm.ListCertificatesInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintCertificates(os.Stdout, resources); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
