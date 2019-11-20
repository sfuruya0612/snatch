package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sfuruya0612/snatch/internal/util"
)

// CloudWatchLogs client struct
type CloudWatchLogs struct {
	Client *logs.CloudWatchLogs
}

// NewLogsSess return CloudWatchLogs struct initialized
func NewLogsSess(profile, region string) *CloudWatchLogs {
	return &CloudWatchLogs{
		Client: logs.New(getSession(profile, region)),
	}
}

// LogEvent log event struct
type LogEvent struct {
	Timestamp string
	Message   string
}

// LogEvents LogEvent struct slice
type LogEvents []LogEvent

const limit = 30

func (c *CloudWatchLogs) DescribeLogGroups(flag bool) error {
	input := &logs.DescribeLogGroupsInput{}

	elements := []string{}
	output := func(page *logs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, i := range page.LogGroups {
			elements = append(elements, *i.LogGroupName)
		}

		return true
	}

	if err := c.Client.DescribeLogGroupsPages(input, output); err != nil {
		return fmt.Errorf("Describe log groups: %v", err)
	}

	group, err := util.Prompt(elements, "Select Log Group")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err = c.describeLogStreams(group, flag); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (c *CloudWatchLogs) describeLogStreams(group string, flag bool) error {
	input := &logs.DescribeLogStreamsInput{
		LogGroupName: aws.String(group),
	}

	elements := []string{}
	output := func(page *logs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, i := range page.LogStreams {
			// Expired 等で storedBytes が 0 のstream は返さない
			if *i.StoredBytes != 0 {
				elements = append(elements, *i.LogStreamName)
			}
		}

		return true
	}

	if err := c.Client.DescribeLogStreamsPages(input, output); err != nil {
		return fmt.Errorf("Describe log streams: %v", err)
	}

	stream, err := util.Prompt(elements, "Select Log Stream")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err = c.getLogEvents(group, stream, flag); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (c *CloudWatchLogs) getLogEvents(group, stream string, flag bool) error {
	var limit int64 = limit

	input := &logs.GetLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		Limit:         aws.Int64(limit),
	}

	list := LogEvents{}
	output := func(page *logs.GetLogEventsOutput, lastPage bool) bool {
		for _, i := range page.Events {
			time := aws.SecondsTimeValue(i.Timestamp)
			t := time.String()

			list = append(list, LogEvent{
				Timestamp: t,
				Message:   *i.Message,
			})
		}

		return true
	}

	if err := c.Client.GetLogEventsPages(input, output); err != nil {
		return fmt.Errorf("Get log events: %v", err)
	}

	f := util.Formatln(
		list.Timestamp(),
		list.Message(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Timestamp,
			i.Message,
		)
	}

	return nil
}

func (event LogEvents) Timestamp() []string {
	time := []string{}
	for _, i := range event {
		time = append(time, i.Timestamp)
	}
	return time
}

func (event LogEvents) Message() []string {
	mess := []string{}
	for _, i := range event {
		mess = append(mess, i.Message)
	}
	return mess
}
