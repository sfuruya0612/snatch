package aws

import (
	"context"
	"fmt"
	"strings"
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

// Response sendcommand response struct
type Response struct {
	InstanceId string   `json:"instance_id"`
	Status     string   `json:"status"`
	Output     []string `json:"output"`
}

// Responses Response struct slice
type Responses []Response

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
