package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli/v2"
)

var S3 = &cli.Command{
	Name:  "s3",
	Usage: "Get a list of S3 Buckets",
	Action: func(c *cli.Context) error {
		return getBucketList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:      "object",
			Usage:     "Get S3 object list",
			ArgsUsage: "[ --bucket | -b ] <BucketName>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "bucket",
					Aliases: []string{"b"},
					Usage:   "Set bucket name",
				},
			},
			Action: func(c *cli.Context) error {
				return getObjectList(c.String("profile"), c.String("region"), c.String("bucket"))
			},
		},
		{
			Name:      "cat",
			Usage:     "Desplay S3 object file",
			ArgsUsage: "[ --bucket | -b ] <BucketName> [ --key | -k ] <ObjectKey>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "bucket",
					Aliases:  []string{"b"},
					Usage:    "Set bucket name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "key",
					Aliases:  []string{"k"},
					Usage:    "Set object key",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "download",
					Aliases: []string{"d"},
					Usage:   "Download object file",
				},
			},
			Action: func(c *cli.Context) error {
				return catObject(c.String("profile"), c.String("region"), c.String("bucket"), c.String("key"), c.Bool("download"))
			},
		},
	},
}

func getBucketList(profile, region string) error {
	client := saws.NewS3Client(profile, region)

	buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, i := range buckets {
		fmt.Printf("%v\n", i)
	}

	return nil
}

func getObjectList(profile, region, bucket string) error {
	client := saws.NewS3Client(profile, region)

	if len(bucket) == 0 {
		buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		bucket, err = util.Prompt(buckets, "Select Bucket")
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	objects, err := client.ListObjects(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintObjects(os.Stdout, objects); err != nil {
		return fmt.Errorf("failed to print objects")
	}

	return nil
}

func catObject(profile, region, bucket, key string, download bool) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	client := saws.NewS3Client(profile, region)
	body, err := client.GetObject(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	b, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read file s3://%s/%s: %v", bucket, key, err)
	}

	if download {
		f, err := os.Create(filepath.Base(key))
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer f.Close()

		_, err = f.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}

		fmt.Printf("Downloaded s3://%s/%s\n", bucket, key)

		return nil
	}

	fmt.Println(string(b))

	return nil
}
