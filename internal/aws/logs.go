package aws

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
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

// DescribeLogGroups return []string (logs.DescribeLogGroupsOutput.LogGroupName)
// input logs.DescribeLogGroupsInput
func (c *CloudWatchLogs) DescribeLogGroups(input *logs.DescribeLogGroupsInput) ([]string, error) {
	groups := []string{}
	output := func(page *logs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, i := range page.LogGroups {
			groups = append(groups, *i.LogGroupName)
		}

		return true
	}

	if err := c.Client.DescribeLogGroupsPages(input, output); err != nil {
		return nil, fmt.Errorf("Describe log groups: %v", err)
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	return groups, nil
}

// DescribeLogStreams return []string (logs.DescribeLogStreamsOutput.LogStreamName)
// input logs.DescribeLogStreamsInput
func (c *CloudWatchLogs) DescribeLogStreams(input *logs.DescribeLogStreamsInput) ([]string, error) {
	streams := []string{}
	output := func(page *logs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, i := range page.LogStreams {
			// StoredBytes が 0 のstream は可視性が下がるので返さない
			if *i.StoredBytes != 0 {
				streams = append(streams, *i.LogStreamName)
			}
		}
		sort.Slice(streams, func(i, j int) bool {
			return streams[i] > streams[j]
		})

		return true
	}

	if err := c.Client.DescribeLogStreamsPages(input, output); err != nil {
		return nil, fmt.Errorf("Describe log streams: %v", err)
	}

	if len(streams) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	return streams, nil
}

// GetLogEvents return LogEvents
// input logs.GetLogEventsInput
func (c *CloudWatchLogs) GetLogEvents(input *logs.GetLogEventsInput) (LogEvents, error) {
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
		return nil, fmt.Errorf("Get log events: %v", err)
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	return list, nil
}

func PrintLogEvents(wrt io.Writer, resources LogEvents) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Timestamp",
		"Message",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.LogsTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *LogEvent) LogsTabString() string {
	fields := []string{
		i.Timestamp,
		i.Message,
	}

	return strings.Join(fields, "\t")
}
