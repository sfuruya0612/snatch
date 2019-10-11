package aws

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Object struct {
	Key          string
	Size         string
	LastModified string
}

type Objects []Object

func newS3Sess(profile string, region string) *s3.S3 {
	sess := getSession(profile, region)
	return s3.New(sess)
}

func ListBuckets(profile string, region string) error {
	client := newS3Sess(profile, region)

	res, err := client.ListBuckets(nil)
	if err != nil {
		return fmt.Errorf("List s3 buckets: %v", err)
	}

	elements := []string{}
	for _, r := range res.Buckets {
		buckets := *r.Name

		elements = append(elements, buckets)
	}

	bucket, err := util.Prompt(elements, "Select Bucket")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	err = ListObjects(client, bucket)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func ListObjects(client *s3.S3, bucket string) error {
	fmt.Println(bucket)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	res, err := client.ListObjectsV2(input)
	if err != nil {
		return fmt.Errorf("List s3 objects: %v", err)
	}

	list := Objects{}
	for _, r := range res.Contents {

		size := strconv.FormatInt(*r.Size, 10)
		lastmod := r.LastModified.String()

		list = append(list, Object{
			Key:          *r.Key,
			Size:         size,
			LastModified: lastmod,
		})
	}

	f := util.Formatln(
		list.key(),
		list.size(),
		list.lastModified(),
	)

	sort.Slice(list, func(i, j int) bool {
		return list[i].LastModified < list[j].LastModified
	})

	for _, i := range list {
		fmt.Printf(
			f,
			i.Key,
			i.Size,
			i.LastModified,
		)
	}

	return nil
}

func (obj Objects) key() []string {
	key := []string{}
	for _, i := range obj {
		key = append(key, i.Key)
	}
	return key
}

func (obj Objects) size() []string {
	size := []string{}
	for _, i := range obj {
		size = append(size, i.Size)
	}
	return size
}

func (obj Objects) lastModified() []string {
	lastmod := []string{}
	for _, i := range obj {
		lastmod = append(lastmod, i.LastModified)
	}
	return lastmod
}
