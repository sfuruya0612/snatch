package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Stack struct {
	Name   string
	Status string
}

type Stacks []Stack

func newCfnSess(profile, region string) *cloudformation.CloudFormation {
	sess := getSession(profile, region)
	return cloudformation.New(sess)
}

func DescribeStacks(profile, region string) error {
	client := newCfnSess(profile, region)

	res, err := client.DescribeStacks(nil)
	if err != nil {
		return fmt.Errorf("Describe stacks: %v", err)
	}

	list := Stacks{}
	for _, l := range res.Stacks {
		list = append(list, Stack{
			Name:   *l.StackName,
			Status: *l.StackStatus,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.Status(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.Status,
		)
	}

	return nil
}

func (s Stacks) Name() []string {
	name := []string{}
	for _, i := range s {
		name = append(name, i.Name)
	}
	return name
}

func (s Stacks) Status() []string {
	status := []string{}
	for _, i := range s {
		status = append(status, i.Status)
	}
	return status
}
