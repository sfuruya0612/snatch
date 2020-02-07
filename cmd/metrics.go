package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/mapping"
	"github.com/urfave/cli"
)

func GetMetrics(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	service := c.String("service")

	if len(service) == 0 {
		return fmt.Errorf("--service or -s option is required")
	}

	dimentsion := []string{}
	switch service {
	case "EC2", "Ec2", "ec2":
		client := saws.NewEc2Sess(profile, region)
		list, err := client.DescribeInstances(&ec2.DescribeInstancesInput{})
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		for _, i := range list {
			dimentsion = append(dimentsion, i.InstanceId)
		}
	case "RDS", "Rds", "rds":
		client := saws.NewRdsSess(profile, region)
		list, err := client.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		for _, i := range list {
			dimentsion = append(dimentsion, i.Name)
		}
	default:
		return fmt.Errorf("service is not match dimentsion")
	}

	metric := mapping.MetricsMap[service]

	var queries []*cloudwatch.MetricDataQuery
	for _, i := range dimentsion {
		for idx, m := range metric.MetricDetails {
			queries = append(queries, &cloudwatch.MetricDataQuery{
				Id: aws.String(fmt.Sprintf("%s_%d", strings.Replace(i, "-", "_", -1), idx)),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String(metric.Namespace),
						MetricName: aws.String(m.MetricName),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String(metric.DimensionName),
								Value: aws.String(i),
							},
						},
					},
					Period: aws.Int64(metric.Period),
					Stat:   aws.String(m.Statistics),
				},
			})
		}
	}
	// 時間は現在時刻から5分前のレンジで指定している
	input := &cloudwatch.GetMetricDataInput{
		StartTime:         aws.Time(time.Now().Add(-5 * time.Minute)),
		EndTime:           aws.Time(time.Now()),
		MetricDataQueries: queries,
	}

	client := saws.NewMetricsSess(profile, region)
	err := client.GetMetricData(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
