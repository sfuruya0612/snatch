package aws

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSM client struct
type SSM struct {
	Client *ssm.SSM
}

// NewSsmSess return SSM struct initialized
func NewSsmSess(profile, region string) *SSM {
	return &SSM{
		Client: ssm.New(getSession(profile, region)),
	}
}

// Session ssm session history struct
type Session struct {
	SessionId string
	Owner     string
	Target    string
	StartDate string
	EndDate   string
}

// Sessions Session struct slice
type Sessions []Session

// Response sendcommand response struct
type Response struct {
	InstanceId string   `json:"instance_id"`
	Status     string   `json:"status"`
	Output     []string `json:"output"`
}

// Responses Response struct slice
type Responses []Response

// CommandLog sendcommand log struct
type CmdLog struct {
	DocumentName      string
	Commands          string
	Targets           string
	Status            string
	RequestedDateTime string
}

// CommandLogs CommandLog struct slice
type CmdLogs []CmdLog

// DescribeInstanceInformation return []string (ssm.DescribeInstanceInformationOutput.InstanceId)
// input ssm.DescribeInstanceInformationInput
func (c *SSM) DescribeInstanceInformation(input *ssm.DescribeInstanceInformationInput) ([]string, error) {
	ids := []string{}
	output := func(page *ssm.DescribeInstanceInformationOutput, lastPage bool) bool {
		for _, i := range page.InstanceInformationList {
			ids = append(ids, *i.InstanceId)
		}
		return true
	}

	if err := c.Client.DescribeInstanceInformationPages(input, output); err != nil {
		return nil, fmt.Errorf("Describe instance information: %v", err)
	}

	return ids, nil
}

// CreateStartSession return ssm.StartSessionOutput, string ()
// input ssm.DescribeInstanceInformationInput
func (c *SSM) CreateStartSession(input *ssm.StartSessionInput) (*ssm.StartSessionOutput, string, error) {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	sess, err := c.Client.StartSessionWithContext(subctx, input)
	if err != nil {
		return nil, "", err
	}

	return sess, c.Client.Endpoint, nil

}

// DeleteStartSession return none (Only error)
// input ssm.TerminateSessionInput
func (c *SSM) DeleteStartSession(input *ssm.TerminateSessionInput) error {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := c.Client.TerminateSessionWithContext(subctx, input); err != nil {
		return err
	}

	return nil
}

// DescribeSessions return Sessions
// input ssm.DescribeSessionsInput
func (c *SSM) DescribeSessions(input *ssm.DescribeSessionsInput) (Sessions, error) {
	output, err := c.Client.DescribeSessions(input)
	if err != nil {
		return nil, fmt.Errorf("Describe sessions: %v", err)
	}

	list := Sessions{}
	for _, l := range output.Sessions {
		s := strings.Split(*l.Owner, "/")
		owner := s[1]

		list = append(list, Session{
			SessionId: *l.SessionId,
			Owner:     owner,
			Target:    *l.Target,
			StartDate: l.StartDate.String(),
			EndDate:   l.EndDate.String(),
		})
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("No historys")
	}

	return list, nil
}

// SendCommand return ssm.SendCommandOutput
// input ssm.SendCommandInput
func (c *SSM) SendCommand(input *ssm.SendCommandInput) (*ssm.SendCommandOutput, error) {
	output, err := c.Client.SendCommand(input)
	if err != nil {
		return nil, fmt.Errorf("Command send: %v", err)
	}

	return output, nil
}

// ListCommandInvocations return Responses
// input ssm.ListCommandInvocationsInput
func (c *SSM) ListCommandInvocations(input *ssm.ListCommandInvocationsInput) (Responses, error) {
	resp := Responses{}
	for {
		output, err := c.Client.ListCommandInvocations(input)
		if err != nil {
			return nil, fmt.Errorf("List command invocation: %v", err)
		}

		if len(output.CommandInvocations) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		inprogress := false
		for _, ci := range output.CommandInvocations {
			if *ci.Status == "InProgress" {
				inprogress = true
				break
			}
		}

		if inprogress {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, ci := range output.CommandInvocations {
			out := *ci.CommandPlugins[0].Output
			spl := strings.Split(out, "\n")
			if len(spl[len(spl)-1]) < 1 {
				spl = spl[:len(spl)-1]
			}

			res := Response{
				InstanceId: *ci.InstanceId,
				Status:     *ci.Status,
				Output:     spl,
			}
			resp = append(resp, res)

			if len(out) < 2500 {
				continue
			}

			res.Output = spl
		}

		break
	}

	return resp, nil
}

// ListCommands return CmdLogs
// input ssm.ListCommandsInput
func (c *SSM) ListCommands(input *ssm.ListCommandsInput) (CmdLogs, error) {
	list := CmdLogs{}
	output := func(page *ssm.ListCommandsOutput, lastpage bool) bool {
		for _, i := range page.Commands {
			var (
				cmds    []string
				cmd     string = "None"
				targets []string
				target  string = "None"
			)

			if i.Parameters["commands"] != nil {
				for _, c := range i.Parameters["commands"] {
					cmds = append(cmds, *c)
				}
				cmd = strings.Join(cmds[:], " ")
			}

			if i.Targets != nil {
				for _, i := range i.Targets {
					for _, v := range i.Values {
						targets = []string{
							*i.Key,
							*v,
						}
					}
				}
				target = strings.Join(targets[:], ", ")
			}

			if i.InstanceIds != nil {
				for _, i := range i.InstanceIds {
					targets = append(targets, *i)
				}
				target = strings.Join(targets[:], ",")
			}

			list = append(list, CmdLog{
				DocumentName:      *i.DocumentName,
				Commands:          cmd,
				Targets:           target,
				Status:            *i.Status,
				RequestedDateTime: i.RequestedDateTime.String(),
			})
		}
		return true
	}

	if err := c.Client.ListCommandsPages(input, output); err != nil {
		return nil, fmt.Errorf("List commands: %v", err)
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("No command logs")
	}

	return list, nil
}

func PrintSessHist(wrt io.Writer, resources Sessions) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"SessionId",
		"Owner",
		"Target",
		"StartDate",
		"EndDate",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.HistTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Session) HistTabString() string {
	fields := []string{
		i.SessionId,
		i.Owner,
		i.Target,
		i.StartDate,
		i.EndDate,
	}

	return strings.Join(fields, "\t")
}

func PrintCmdLogs(wrt io.Writer, resources CmdLogs) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"DocumentName",
		"Commands",
		"Targets",
		"Status",
		"RequestedDateTime",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.CmdLogTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *CmdLog) CmdLogTabString() string {
	fields := []string{
		i.DocumentName,
		i.Commands,
		i.Targets,
		i.Status,
		i.RequestedDateTime,
	}

	return strings.Join(fields, "\t")
}
