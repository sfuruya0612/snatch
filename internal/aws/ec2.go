package aws

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2 structure is ec2 client.
type EC2 struct {
	Client *ec2.Client
}

// NewEc2Client returns EC2 struct initialized.
func NewEc2Client(profile, region string) *EC2 {
	return &EC2{
		Client: ec2.NewFromConfig(GetSessionV2(profile, region)),
	}
}

// Instance structure is ec2 instance information.
type Instance struct {
	Name             string
	InstanceId       string
	InstanceType     string
	Lifecycle        string
	PrivateIpAddress string
	PublicIpAddress  string
	State            string
	KeyName          string
	AvailabilityZone string
	LaunchTime       string
}

// DescribeInstances returns slice Instance structure.
func (c *EC2) DescribeInstances(input *ec2.DescribeInstancesInput) ([]Instance, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe instances: %v", err)
	}

	if len(output.Reservations) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []Instance{}
	for _, r := range output.Reservations {
		for _, i := range r.Instances {
			name := ""
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					name = *t.Value
				}
			}

			priip := "None"
			if i.PrivateIpAddress != nil {
				priip = *i.PrivateIpAddress
			}

			pubip := "None"
			if i.PublicIpAddress != nil {
				pubip = *i.PublicIpAddress
			}

			key := "None"
			if i.KeyName != nil {
				key = *i.KeyName
			}

			// AvailabilityZoneは末尾(1a, 1c...)のみ取得する
			spl := strings.Split(*i.Placement.AvailabilityZone, "-")
			az := spl[2]

			list = append(list, Instance{
				Name:             name,
				InstanceId:       *i.InstanceId,
				InstanceType:     string(i.InstanceType),
				Lifecycle:        string(i.InstanceLifecycle),
				PrivateIpAddress: priip,
				PublicIpAddress:  pubip,
				State:            string(i.State.Name),
				KeyName:          key,
				AvailabilityZone: az,
				LaunchTime:       i.LaunchTime.String(),
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func PrintInstances(wrt io.Writer, resources []Instance) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"InstanceID",
		"InstanceType",
		"Lifecycle",
		"PrivateIP",
		"PublicIP",
		"State",
		"KeyName",
		"AZ",
		"LaunchTime",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("header join: %v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.Ec2TabString()); err != nil {
			return fmt.Errorf("resources join: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush: %v", err)
	}

	return nil
}

func (i *Instance) Ec2TabString() string {
	fields := []string{
		i.Name,
		i.InstanceId,
		i.InstanceType,
		i.Lifecycle,
		i.PrivateIpAddress,
		i.PublicIpAddress,
		i.State,
		i.KeyName,
		i.AvailabilityZone,
		i.LaunchTime,
	}

	return strings.Join(fields, "\t")
}
