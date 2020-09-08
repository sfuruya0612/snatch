package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

func GetBucketList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewS3Sess(profile, region)
	buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, i := range buckets {
		fmt.Printf("%v\n", i)
	}

	return nil
}

func GetObjectList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	bucket := c.String("bucket")

	client := saws.NewS3Sess(profile, region)

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

func CatObject(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	bucket := c.String("bucket")
	key := c.String("key")

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("--bucket and --key is required")
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	client := saws.NewS3Sess(profile, region)
	body, err := client.GetObject(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var buf bytes.Buffer
	if _, rerr := buf.ReadFrom(body); rerr != nil {
		return fmt.Errorf("read body from s3://%s/%s: %v", bucket, key, err)
	}

	if !strings.HasSuffix(key, ".gz") {
		fmt.Println(buf.String())
		return nil
	}

	reader, err := gzip.NewReader(&buf)
	if err != nil {
		return fmt.Errorf("new gzip reader: %v", err)
	}

	bytea, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read gzip s3://%s/%s: %v", bucket, key, err)
	}

	if err := reader.Close(); err != nil {
		return fmt.Errorf("close reader s3://%s/%s: %v", bucket, key, err)
	}

	fmt.Println(string(bytea))

	return nil
}

func DownloadObject(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	bucket := c.String("bucket")
	key := c.String("key")

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("--bucket and --key is required")
	}

	f, err := os.Create(filepath.Base(key))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	client := saws.NewS3DownloaderSess(profile, region)
	bytea, err := client.Download(f, input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("\nDownloadedSize: %d byte", bytea)
	return nil
}
