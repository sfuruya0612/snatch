package aws

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sfuruya0612/snatch/internal/util"
)

// EC2 client struct
type EC2 struct {
	Client *ec2.EC2
}

// NewEc2Sess return EC2 struct initialized
func NewEc2Sess(profile, region string) *EC2 {
	return &EC2{
		Client: ec2.New(getSession(profile, region)),
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

// DescribeInstances return error. Print ec2.DescribeInstances
func (c *EC2) DescribeInstances(input *ec2.DescribeInstancesInput) error {
	output, err := c.Client.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("Describe instances: %v", err)
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

			list = append(list, Instance{
				Name:             name,
				InstanceId:       *i.InstanceId,
				InstanceType:     *i.InstanceType,
				PrivateIpAddress: priip,
				PublicIpAddress:  pubip,
				State:            *i.State.Name,
				KeyName:          key,
				AvailabilityZone: *i.Placement.AvailabilityZone,
			})
		}
	}
	f := util.Formatln(
		list.Name(),
		list.InstanceId(),
		list.InstanceType(),
		list.PrivateIpAddress(),
		list.PublicIpAddress(),
		list.State(),
		list.KeyName(),
		list.AvailabilityZone(),
	)

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.InstanceId,
			i.InstanceType,
			i.PrivateIpAddress,
			i.PublicIpAddress,
			i.State,
			i.KeyName,
			i.AvailabilityZone,
		)
	}

	return nil
}

func (c *EC2) getInstancesByInstanceIds(ids []string) (Instances, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(ids),
	}

	output, err := c.Client.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Describe instances by instance ids: %v", err)
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

			time := i.LaunchTime.String()

			list = append(list, Instance{
				Name:       name,
				InstanceId: *i.InstanceId,
				LaunchTime: time,
			})
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

// GetConsoleOutput return error. Print ec2.GetConsoleOutput.Output
func (c *EC2) GetConsoleOutput(input *ec2.GetConsoleOutputInput) error {
	output, err := c.Client.GetConsoleOutput(input)
	if err != nil {
		return fmt.Errorf("Get console output: %v", err)
	}

	d, err := util.DecodeString(*output.Output)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Println(d)

	return nil
}

func (ins Instances) Name() []string {
	name := []string{}
	for _, i := range ins {
		name = append(name, i.Name)
	}
	return name
}

func (ins Instances) InstanceId() []string {
	id := []string{}
	for _, i := range ins {
		id = append(id, i.InstanceId)
	}
	return id
}

func (ins Instances) InstanceType() []string {
	ty := []string{}
	for _, i := range ins {
		ty = append(ty, i.InstanceType)
	}
	return ty
}

func (ins Instances) PrivateIpAddress() []string {
	pip := []string{}
	for _, i := range ins {
		pip = append(pip, i.PrivateIpAddress)
	}
	return pip
}

func (ins Instances) PublicIpAddress() []string {
	gip := []string{}
	for _, i := range ins {
		gip = append(gip, i.PublicIpAddress)
	}
	return gip
}

func (ins Instances) State() []string {
	st := []string{}
	for _, i := range ins {
		st = append(st, i.State)
	}
	return st
}

func (ins Instances) KeyName() []string {
	key := []string{}
	for _, i := range ins {
		key = append(key, i.KeyName)
	}
	return key
}

func (ins Instances) AvailabilityZone() []string {
	az := []string{}
	for _, i := range ins {
		az = append(az, i.AvailabilityZone)
	}
	return az
}
