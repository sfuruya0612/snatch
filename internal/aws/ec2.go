package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func NewEc2Sess(profile string, region string) *ec2.EC2 {
	sess := getSession(profile, region)
	return ec2.New(sess)
}

func DescribeEc2(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	svc := NewEc2Sess(profile, region)

	res, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	for _, r := range res.Reservations {
		for _, i := range r.Instances {
			var tag_name string
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					tag_name = *t.Value
				}
			}

			// 取得する情報を記載
			resources := []string{
				tag_name,
				*i.InstanceId,
				*i.InstanceType,
				*i.PrivateIpAddress,
				/*				*i.PublicIpAddress,
				 *i.State.Name,*/
			}

			// Parserに飛ばしてから出力させたい
			fmt.Println(resources)
		}
	}
	return nil
}
