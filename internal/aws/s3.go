package aws

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

// S3Downloader client struct
type S3Downloader struct {
	Client *s3manager.Downloader
}

// NewS3DownloaderSess return S3Manager Downloader struct initialized
func NewS3DownloaderSess(profile, region string) *S3Downloader {
	return &S3Downloader{
		Client: s3manager.NewDownloader(getSession(profile, region)),
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

// ListBuckets return []string (s3.ListBuckets.Buckets)
// input s3.ListBucketsInput
func (c *S3) ListBuckets(input *s3.ListBucketsInput) ([]string, error) {
	output, err := c.Client.ListBuckets(input)
	if err != nil {
		return nil, fmt.Errorf("List buckets: %v", err)
	}

	buckets := []string{}
	for _, l := range output.Buckets {
		name := *l.Name
		buckets = append(buckets, name)
	}

	if len(buckets) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	return buckets, nil
}

// ListObjects return Objects
// input s3.ListObjectsV2Input
func (c *S3) ListObjects(input *s3.ListObjectsV2Input) (Objects, error) {
	output, err := c.Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("List objects: %v", err)
	}

	list := Objects{}
	for _, l := range output.Contents {

		size := strconv.FormatInt(*l.Size, 10)

		list = append(list, Object{
			Key:          *l.Key,
			Size:         size,
			LastModified: l.LastModified.String(),
		})
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].LastModified < list[j].LastModified
	})

	return list, nil
}

// GetObject return io.ReadCloser
// input s3.GetObjectInput
func (c *S3) GetObject(input *s3.GetObjectInput) (io.ReadCloser, error) {
	output, err := c.Client.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("Get object: %v", err)
	}

	return output.Body, nil
}

// Download return int64
// input io.WriterAt, s3.GetObjectInput
func (c *S3Downloader) Download(w io.WriterAt, input *s3.GetObjectInput) (int64, error) {
	output, err := c.Client.Download(w, input)
	if err != nil {
		return 0, fmt.Errorf("Get object: %v", err)
	}

	return output, nil
}

func PrintObjects(wrt io.Writer, resources Objects) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Key",
		"Size",
		"LastModified",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.S3TabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Object) S3TabString() string {
	fields := []string{
		i.Key,
		i.Size,
		i.LastModified,
	}

	return strings.Join(fields, "\t")
}
