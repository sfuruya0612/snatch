package cmd

import (
	"fmt"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func ListBuckets(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	flag := c.Bool("l")

	s3 := aws.NewS3Sess(profile, region)
	if err := s3.ListBuckets(flag); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
