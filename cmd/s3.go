package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

func GetS3List(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	flag := c.Bool("l")

	client := saws.NewS3Sess(profile, region)
	buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// flagがなければ出力してreturn
	if !flag {
		for _, i := range buckets {
			fmt.Printf("%v\n", i)
		}
		return nil
	}

	bucket, err := util.Prompt(buckets, "Select Bucket")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	resources, err := client.ListObjects(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintObjects(os.Stdout, resources); err != nil {
		return fmt.Errorf("Failed to print resources")
	}

	return nil
}
