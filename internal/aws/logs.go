package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sfuruya0612/snatch/internal/util"
)

type LogEvent struct {
	Timestamp string
	Message   string
}

type LogEvents []LogEvent

const limit = 10

func newLogsSess(profile, region string) *logs.CloudWatchLogs {
	sess := getSession(profile, region)
	return logs.New(sess)
}

func DescribeLogGroups(profile, region string, flag bool) error {
	client := newLogsSess(profile, region)

	input := &logs.DescribeLogGroupsInput{}
	groups, err := client.DescribeLogGroups(input)
	if err != nil {
		return fmt.Errorf("Describe log groups: %v", err)
	}

	elements := []string{}
	for _, l := range groups.LogGroups {
		item := *l.LogGroupName

		elements = append(elements, item)
	}

	group, err := util.Prompt(elements, "Select Log Group")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	err = describeLogStreams(client, group, flag)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func describeLogStreams(client *logs.CloudWatchLogs, group string, flag bool) error {
	input := &logs.DescribeLogStreamsInput{
		LogGroupName: aws.String(group),
	}

	streams, err := client.DescribeLogStreams(input)
	if err != nil {
		return fmt.Errorf("Describe log streams: %v", err)
	}

	elements := []string{}
	for _, l := range streams.LogStreams {
		item := *l.LogStreamName

		elements = append(elements, item)
	}

	stream, err := util.Prompt(elements, "Select Log Stream")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// tokenがないと新しいログがとれない様子
	// err = GetLogEvents(client, group, stream, token, flag)
	err = GetLogEvents(client, group, stream, flag)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func GetLogEvents(client *logs.CloudWatchLogs, group, stream string, flag bool) error {
	var limit int64 = limit

	input := &logs.GetLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		Limit:         aws.Int64(limit),
	}

	events, err := client.GetLogEvents(input)
	if err != nil {
		return fmt.Errorf("Get log events: %v", err)
	}

	list := LogEvents{}
	for _, e := range events.Events {
		time := aws.SecondsTimeValue(e.Timestamp)
		t := time.String()

		list = append(list, LogEvent{
			Timestamp: t,
			Message:   *e.Message,
		})
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
