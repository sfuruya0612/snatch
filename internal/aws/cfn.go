package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sfuruya0612/snatch/internal/util"
)

// CloudFormation client struct
type CloudFormation struct {
	Client *cloudformation.CloudFormation
}

// newCfnSess return CloudFormation struct initialized
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

func (c *CloudFormation) DescribeStacks() error {
	output, err := c.Client.DescribeStacks(nil)
	if err != nil {
		return fmt.Errorf("Describe stacks: %v", err)
	}

	list := Stacks{}
	for _, l := range output.Stacks {

		create := l.CreationTime.String()

		update := "None"
		if l.LastUpdatedTime != nil {
			update = l.LastUpdatedTime.String()
		}

		list = append(list, Stack{
			Name:       *l.StackName,
			Status:     *l.StackStatus,
			CreateDate: create,
			UpdateDate: update,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.Status(),
		list.CreateDate(),
		list.UpdateDate(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.Status,
			i.CreateDate,
			i.UpdateDate,
		)
	}

	return nil
}

func (s Stacks) Name() []string {
	n := []string{}
	for _, i := range s {
		n = append(n, i.Name)
	}
	return n
}

func (s Stacks) Status() []string {
	sts := []string{}
	for _, i := range s {
		sts = append(sts, i.Status)
	}
	return sts
}

func (s Stacks) CreateDate() []string {
	c := []string{}
	for _, i := range s {
		c = append(c, i.CreateDate)
	}
	return c
}

func (s Stacks) UpdateDate() []string {
	u := []string{}
	for _, i := range s {
		u = append(u, i.UpdateDate)
	}
	return u
}
