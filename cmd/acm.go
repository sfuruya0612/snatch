package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListCertificates(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	acm := aws.NewAcmSess(profile, region)
	if err := acm.ListCertificates(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
