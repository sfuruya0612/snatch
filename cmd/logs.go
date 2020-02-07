package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

func GetCloudWatchLogs(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	// flag := c.GlobalBool("follow")

	client := saws.NewLogsSess(profile, region)

	groups, err := client.DescribeLogGroups(&logs.DescribeLogGroupsInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	group, err := util.Prompt(groups, "Select Log Group")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	si := &logs.DescribeLogStreamsInput{
		LogGroupName: aws.String(group),
	}

	streams, err := client.DescribeLogStreams(si)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	stream, err := util.Prompt(streams, "Select Log Stream")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ei := &logs.GetLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		Limit:         aws.Int64(30),
	}

	// Ctrl+C で終了
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-sc:
			fmt.Println("Stop")
			return nil
		case <-t.C:
			resources, err := client.GetLogEvents(ei)
			if err != nil {
				return fmt.Errorf("%v", err)
			}

			if err := saws.PrintLogEvents(os.Stdout, resources); err != nil {
				return fmt.Errorf("Failed to print resources")
			}
		}
	}
}
