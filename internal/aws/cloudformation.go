package aws

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// CloudFormation client struct
type CloudFormation struct {
	Client *cloudformation.CloudFormation
}

// NewCfnSess return CloudFormation struct initialized
func NewCfnSess(profile, region string) *CloudFormation {
	return &CloudFormation{
		Client: cloudformation.New(getSession(profile, region)),
	}
}

// Stack cloudformation stack struct
type Stack struct {
	Name       string
	Status     string
	CreateDate string
	UpdateDate string
}

// Stacks Stack struct slice
type Stacks []Stack

// DescribeStacks return Stacks
// input cloudformation.DescribeStacksInput
func (c *CloudFormation) DescribeStacks(input *cloudformation.DescribeStacksInput) (Stacks, error) {
	output, err := c.Client.DescribeStacks(input)
	if err != nil {
		return nil, fmt.Errorf("Describe stacks: %v", err)
	}

	list := Stacks{}
	for _, l := range output.Stacks {
		update := "None"
		if l.LastUpdatedTime != nil {
			update = l.LastUpdatedTime.String()
		}

		list = append(list, Stack{
			Name:       *l.StackName,
			Status:     *l.StackStatus,
			CreateDate: l.CreationTime.String(),
			UpdateDate: update,
		})
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	return list, nil
}

func PrintStacks(wrt io.Writer, resources Stacks) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"Status",
		"CreateDate",
		"UpdateDate",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.StackTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Stack) StackTabString() string {
	fields := []string{
		i.Name,
		i.Status,
		i.CreateDate,
		i.UpdateDate,
	}

	return strings.Join(fields, "\t")
}
