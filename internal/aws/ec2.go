package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Ec2Opts struct {
	Tag        string
	Running    bool
	Stopped    bool
	InstanceId string
	File       string
}

// EC2Instance は選択可能なEC2インスタンスを表す
type EC2Instance struct {
	Name       string
	InstanceID string
}

// Title は表示名を返す
func (i EC2Instance) Title() string {
	return fmt.Sprintf("%s (%s)", i.Name, i.InstanceID)
}

// ID はインスタンスIDを返す
func (i EC2Instance) ID() string {
	return i.InstanceID
}

type ec2Api interface {
	DescribeInstances(ctx context.Context, input *ec2.DescribeInstancesInput, opts ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func NewEC2Client(profile, region string) ec2Api {
	return ec2.NewFromConfig(GetSession(profile, region))
}

func getEc2TagValue(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func GenerateDescribeInstancesInput(opts *Ec2Opts) (*ec2.DescribeInstancesInput, error) {
	i := &ec2.DescribeInstancesInput{}

	// t := opts.Tag
	// if len(t) > 0 {
	// 	if !strings.Contains(t, ":") {
	// 		return nil, fmt.Errorf("tag format is invalid (e.g. Name:MyInstance)")
	// 	}

	// 	s := strings.Split(t, ":")
	// 	i.Filters = append(i.Filters, types.Filter{
	// 		Name:   aws.String("tag:" + s[0]),
	// 		Values: []string{s[1]},
	// 	})
	// }

	if opts.Running {
		i.Filters = append(i.Filters, types.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []string{"running"},
		})
	}

	if opts.InstanceId != "" {
		i.Filters = append(i.Filters, types.Filter{
			Name:   aws.String("instance-id"),
			Values: []string{opts.InstanceId},
		})
	}

	return i, nil
}

func DescribeInstances(api ec2Api, input *ec2.DescribeInstancesInput) ([][]string, error) {
	o, err := api.DescribeInstances(context.Background(), input)
	if err != nil {
		return nil, err
	}

	var instances [][]string
	for _, r := range o.Reservations {
		for _, i := range r.Instances {
			name := getEc2TagValue(i.Tags, "Name")

			lifecycle := "OnDemand"
			if i.InstanceLifecycle != "" {
				lifecycle = string(i.InstanceLifecycle)
			}

			privateIP := "None"
			if i.PrivateIpAddress != nil {
				privateIP = *i.PrivateIpAddress
			}

			publicIP := "None"
			if i.PublicIpAddress != nil {
				publicIP = *i.PublicIpAddress
			}

			keyName := "None"
			if i.KeyName != nil {
				keyName = *i.KeyName
			}

			instance := []string{
				name,
				aws.ToString(i.InstanceId),
				string(i.InstanceType),
				lifecycle,
				privateIP,
				publicIP,
				string(i.State.Name),
				keyName,
				string(*i.Placement.AvailabilityZone),
				i.LaunchTime.String(),
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}
