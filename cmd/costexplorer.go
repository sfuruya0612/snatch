package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetCost(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := "us-east-1"

	start := c.String("start")
	end := c.String("end")

	if start == "" || end == "" {
		y := time.Now().Year()
		m := time.Now().Month()
		d := time.Now().Day()

		s := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, -1, 0)
		e := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

		layout := "2006-01-02"
		start = s.Format(layout)
		end = e.Format(layout)
	}

	input := &costexplorer.GetCostAndUsageInput{
		Metrics: []*string{
			aws.String("BlendedCost"),
		},
		Granularity: aws.String("MONTHLY"),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("SERVICE"),
				Type: aws.String("DIMENSION"),
			},
		},
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
	}

	if err := input.Validate(); err != nil {
		return fmt.Errorf("%v", err)
	}

	client := saws.NewCeSess(profile, region)
	output, err := client.GetCostAndUsage(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, o := range output {
		fmt.Printf("Start Date: %v\nEnd Date:   %v\n", o.Start, o.End)
		if err := saws.PrintUsage(os.Stdout, o.Usage); err != nil {
			return fmt.Errorf("Failed to print costs")
		}
		fmt.Println("-----")
	}

	return nil
}
