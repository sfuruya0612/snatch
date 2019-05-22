package services

import (
	"fmt"

    //"github.com/urfave/cli"

    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DescribeEc2() error {
//func DescribeEc2(c *cli.Context) error {
    profile_name := "punk"

	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile_name}))

	svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

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
			// Parserに飛ばしたい

			fmt.Println(
				tag_name,
/*				*i.InstanceId,
				*i.InstanceType,
				*i.PrivateIpAddress,
				*i.PublicIpAddress,
                *i.State.Name,*/
			)
		}
	}
    return nil
}
