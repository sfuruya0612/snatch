package aws

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2 client struct
type EC2 struct {
	Client *ec2.EC2
}

// NewEc2Sess return EC2 struct initialized
func NewEc2Sess(profile, region string) *EC2 {
	return &EC2{
		Client: ec2.New(GetSession(profile, region)),
	}
}

// Instance ec2 instance struct
type Instance struct {
	Name             string
	InstanceId       string
	InstanceType     string
	PrivateIpAddress string
	PublicIpAddress  string
	State            string
	KeyName          string
	AvailabilityZone string
	LaunchTime       string
}

// Instances Instance struct slice
type Instances []Instance

// DescribeInstances return Instances
// input ec2.DescribeInstancesInput
func (c *EC2) DescribeInstances(input *ec2.DescribeInstancesInput) (Instances, error) {
	output, err := c.Client.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("describe instances: %v", err)
	}

	list := Instances{}
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
				InstanceType:     *i.InstanceType,
				PrivateIpAddress: priip,
				PublicIpAddress:  pubip,
				State:            *i.State.Name,
				KeyName:          key,
				AvailabilityZone: az,
				LaunchTime:       i.LaunchTime.String(),
			})
		}
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func PrintInstances(wrt io.Writer, resources Instances) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"InstanceID",
		"InstanceType",
		"PrivateIP",
		"PublicIP",
		"State",
		"KeyName",
		"AZ",
		"LaunchTime",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.Ec2TabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Instance) Ec2TabString() string {
	fields := []string{
		i.Name,
		i.InstanceId,
		i.InstanceType,
		i.PrivateIpAddress,
		i.PublicIpAddress,
		i.State,
		i.KeyName,
		i.AvailabilityZone,
		i.LaunchTime,
	}

	return strings.Join(fields, "\t")
}
