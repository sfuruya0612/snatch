package aws

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Instance struct {
	Name             string
	InstanceId       string
	InstanceType     string
	PrivateIpAddress string
	PublicIpAddress  string
	State            string
	KeyName          string
	AvailabilityZone string
}

type Instances []Instance

func newEc2Sess(profile string, region string) *ec2.EC2 {
	sess := getSession(profile, region)
	return ec2.New(sess)
}

func DescribeInstances(profile, region, tag string) error {
	client := newEc2Sess(profile, region)

	input := &ec2.DescribeInstancesInput{}

	if tag != "" {
		spl := strings.Split(tag, ":")
		if len(spl) == 0 {
			return fmt.Errorf("parse tag=%s", tag)
		}

		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("tag:" + spl[0]),
			Values: []*string{aws.String(spl[1])},
		})
	}

	res, err := client.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("Describe running instances: %v", err)
	}

	list := Instances{}
	for _, r := range res.Reservations {
		for _, i := range r.Instances {

			var tag_name string
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					tag_name = *t.Value
				}
			}

			if i.PrivateIpAddress == nil {
				i.PrivateIpAddress = aws.String("NULL")
			}

			if i.PublicIpAddress == nil {
				i.PublicIpAddress = aws.String("NULL")
			}

			if i.KeyName == nil {
				i.KeyName = aws.String("NULL")
			}

			list = append(list, Instance{
				Name:             tag_name,
				InstanceId:       *i.InstanceId,
				InstanceType:     *i.InstanceType,
				PrivateIpAddress: *i.PrivateIpAddress,
				PublicIpAddress:  *i.PublicIpAddress,
				State:            *i.State.Name,
				KeyName:          *i.KeyName,
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

func getInstancesByInstanceIds(profile, region string, ids []string) (Instances, error) {
	client := newEc2Sess(profile, region)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(ids),
	}

	res, err := client.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Describe instances by instance ids: %v", err)
	}

	list := Instances{}
	for _, r := range res.Reservations {
		for _, i := range r.Instances {

			var tag_name string
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					tag_name = *t.Value
				}
			}

			list = append(list, Instance{
				Name:       tag_name,
				InstanceId: *i.InstanceId,
			})
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
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
