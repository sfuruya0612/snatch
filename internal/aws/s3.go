package aws

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sfuruya0612/snatch/internal/util"
)

// S3 client struct
type S3 struct {
	Client *s3.S3
}

// NewS3Sess return S3 struct initialized
func NewS3Sess(profile, region string) *S3 {
	return &S3{
		Client: s3.New(getSession(profile, region)),
	}
}

// Object s3 object struct
type Object struct {
	Key          string
	Size         string
	LastModified string
}

// Objects Object struct slice
type Objects []Object

func (c *S3) ListBuckets(flag bool) error {
	output, err := c.Client.ListBuckets(nil)
	if err != nil {
		return fmt.Errorf("List s3 buckets: %v", err)
	}

	elements := []string{}
	for _, r := range output.Buckets {
		item := *r.Name

		elements = append(elements, item)
	}

	// flagがなければ出力してreturn
	if !flag {
		for _, i := range elements {
			fmt.Printf("%v\n", i)
		}

		return nil
	}

	bucket, err := util.Prompt(elements, "Select Bucket")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err = c.listObjects(bucket); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (c *S3) listObjects(bucket string) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	output, err := c.Client.ListObjectsV2(input)
	if err != nil {
		return fmt.Errorf("List s3 objects: %v", err)
	}

	list := Objects{}
	for _, r := range output.Contents {

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
