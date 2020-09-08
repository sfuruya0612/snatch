package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CloudWatch client struct
type CloudWatch struct {
	Client *cloudwatch.CloudWatch
}

// NewMetricsSess return CloudWatch struct initialized
func NewMetricsSess(profile, region string) *CloudWatch {
	return &CloudWatch{
		Client: cloudwatch.New(GetSession(profile, region)),
	}
}

// GetMetricData return *cloudwatch.GetMetricDataOutput
// input cloudwatch.GetMetricDataInput
func (c *CloudWatch) GetMetricData(input *cloudwatch.GetMetricDataInput) error {
	output, err := c.Client.GetMetricData(input)
	if err != nil {
		return fmt.Errorf("get metric data: %v", err)
	}

	fmt.Println(output)

	return nil
}
