package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

type Instance struct {
	Name             string
	InstanceId       string
	InstanceType     string
	PrivateIpAddress string
	PublicIpAddress  string
	State            string
	KeyName          string
}

type Instances []Instance

func NewEc2Sess(profile string, region string) *ec2.EC2 {
	sess := getSession(profile, region)
	return ec2.New(sess)
}

func DescribeEc2(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	svc := NewEc2Sess(profile, region)

	// 1k以上のホスト情報を取得する場合は
	// DescribeInstancesPages を使用する必要がある
	res, err := svc.DescribeInstances(nil)
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

			list = append(list, Instance{
				Name:             tag_name,
				InstanceId:       *i.InstanceId,
				InstanceType:     *i.InstanceType,
				PrivateIpAddress: *i.PrivateIpAddress,
				PublicIpAddress:  *i.PublicIpAddress,
				State:            *i.State.Name,
				KeyName:          *i.KeyName,
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
		[]string{""},
	)

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
		)
	}

	return nil
}

func (instances Instances) Name() []string {
	name := []string{}
	for _, i := range instances {
		name = append(name, i.Name)
	}
	return name
}

func (instances Instances) InstanceId() []string {
	id := []string{}
	for _, i := range instances {
		id = append(id, i.InstanceId)
	}
	return id
}

func (instances Instances) InstanceType() []string {
	ty := []string{}
	for _, i := range instances {
		ty = append(ty, i.InstanceType)
	}
	return ty
}

func (instances Instances) PrivateIpAddress() []string {
	pip := []string{}
	for _, i := range instances {
		pip = append(pip, i.PrivateIpAddress)
	}
	return pip
}

func (instances Instances) PublicIpAddress() []string {
	gip := []string{}
	for _, i := range instances {
		gip = append(gip, i.PublicIpAddress)
	}
	return gip
}

func (instances Instances) State() []string {
	st := []string{}
	for _, i := range instances {
		st = append(st, i.State)
	}
	return st
}

func (instances Instances) KeyName() []string {
	key := []string{}
	for _, i := range instances {
		key = append(key, i.KeyName)
	}
	return key
}
