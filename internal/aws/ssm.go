package aws

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SSM client struct
type SSM struct {
	Client *ssm.Client
}

// NewSsmSess return SSM struct initialized
func NewSsmClient(profile, region string) *SSM {
	return &SSM{
		Client: ssm.NewFromConfig(GetSession(profile, region)),
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

// CmdLog sendcommand log struct
type CmdLog struct {
	DocumentName      string
	Commands          string
	Targets           string
	Status            string
	RequestedDateTime string
}

// CmdLogs CommandLog struct slice
type CmdLogs []CmdLog

// Parameter parameter store struct
type Parameter struct {
	Name        string
	Value       string
	Description string
}

// Parameters Parameter struct slice
type Parameters []Parameter

// DescribeInstanceInformation return []string (ssm.DescribeInstanceInformationOutput.InstanceId)
// input ssm.DescribeInstanceInformationInput
func (c *SSM) DescribeInstanceInformation(input *ssm.DescribeInstanceInformationInput) ([]string, error) {
	ids := []string{}
	paginator := ssm.NewDescribeInstanceInformationPaginator(c.Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("describe instance information: %v", err)
		}

		for _, i := range page.InstanceInformationList {
			ids = append(ids, *i.InstanceId)
		}
	}

	return ids, nil
}

// CreateStartSession return ssm.StartSessionOutput, string ()
// input ssm.DescribeInstanceInformationInput
func (c *SSM) StartSession(input *ssm.StartSessionInput) (*ssm.StartSessionOutput, error) {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	sess, err := c.Client.StartSession(subctx, input)
	if err != nil {
		return nil, fmt.Errorf("start session: %v", err)
	}

	return sess, nil

}

// DeleteStartSession return none (Only error)
// input ssm.TerminateSessionInput
func (c *SSM) DeleteSession(input *ssm.TerminateSessionInput) error {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := c.Client.TerminateSession(subctx, input); err != nil {
		return fmt.Errorf("terminate session: %v", err)
	}

	return nil
}

// SendCommand return ssm.SendCommandOutput
// input ssm.SendCommandInput
func (c *SSM) SendCommand(input *ssm.SendCommandInput) (*ssm.SendCommandOutput, error) {
	output, err := c.Client.SendCommand(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("command send: %v", err)
	}

	return output, nil
}

// ListCommandInvocations return Responses
// input ssm.ListCommandInvocationsInput
func (c *SSM) ListCommandInvocations(input *ssm.ListCommandInvocationsInput) (Responses, error) {
	resp := Responses{}
	for {
		output, err := c.Client.ListCommandInvocations(context.TODO(), input)
		if err != nil {
			return nil, fmt.Errorf("list command invocation: %v", err)
		}

		if len(output.CommandInvocations) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		inprogress := false
		for _, ci := range output.CommandInvocations {
			if ci.Status == "InProgress" {
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
				Status:     string(ci.Status),
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

// DescribeParameters return []*ssm.Parameters
// input ssm.DescribeParametersInput
func (c *SSM) DescribeParameters(input *ssm.DescribeParametersInput) ([]types.ParameterMetadata, error) {
	var params []types.ParameterMetadata
	paginator := ssm.NewDescribeParametersPaginator(c.Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("describe paramaters: %v", err)
		}

		params = append(params, page.Parameters...)
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no parameters")
	}

	return params, nil
}

// GetParameter return Parameters
// input []*ssm.ParameterMetadata
func (c *SSM) GetParameter(params []types.ParameterMetadata) (Parameters, error) {
	list := Parameters{}
	for _, p := range params {
		input := &ssm.GetParameterInput{
			Name: aws.String(*p.Name),
		}

		output, err := c.Client.GetParameter(context.TODO(), input)
		if err != nil {
			return nil, fmt.Errorf("get parameter: %v", err)
		}

		description := "None"
		if p.Description != nil {
			description = *p.Description
		}

		list = append(list, Parameter{
			Name:        *p.Name,
			Value:       *output.Parameter.Value,
			Description: description,
		})
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

func PrintParameters(wrt io.Writer, resources Parameters) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"Value",
		"Description",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.ParameterTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Parameter) ParameterTabString() string {
	fields := []string{
		i.Name,
		i.Value,
		i.Description,
	}

	return strings.Join(fields, "\t")
}
