package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/urfave/cli"
)

func NewRdsSess(profile string, region string) *rds.RDS {
	sess := getSession(profile, region)
	return rds.New(sess)
}

func DescribeRds(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	svc := NewRdsSess(profile, region)

	res, err := svc.DescribeDBInstances(nil)
	if err != nil {
		panic(err)
	}

	for _, r := range res.DBInstances {
		// 取得する情報を記載
		resources := []string{
			*r.DBInstanceIdentifier,
		}

		// Parserに飛ばしてから出力させたい
		fmt.Println(resources)
	}
	return nil
}
