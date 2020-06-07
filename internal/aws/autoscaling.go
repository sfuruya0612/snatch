package aws

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// AutoScaling client struct
type AutoScaling struct {
	Client *autoscaling.AutoScaling
}

// NewAsgSess return AutoScaling struct initialized
func NewAsgSess(profile, region string) *AutoScaling {
	return &AutoScaling{
		Client: autoscaling.New(GetSession(profile, region)),
	}
}

// Group autoscaling group struct
type Group struct {
	Name                string
	Launch              string
	Desired             string
	Min                 string
	Max                 string
	TerminationPolicies string
}

// Groups Group struct slice
type Groups []Group

// DescribeAutoScalingGroups return Groups
// input autoscaling.DescribeAutoScalingGroupsInput
func (c *AutoScaling) DescribeAutoScalingGroups(input *autoscaling.DescribeAutoScalingGroupsInput) (Groups, error) {
	list := Groups{}
	output, err := c.Client.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, fmt.Errorf("Describe autoscaling groups: %v", err)
	}

	for _, i := range output.AutoScalingGroups {
		var (
			launch   string
			policies []string
			policy   string
		)

		// LaunchTemplate, LaunchConfiguration, AutoScaling SpotInstanceでそれぞれパラメータが異なるためswitch文で分岐
		switch {
		// AutoScaling SpotInstance
		case i.LaunchConfigurationName == nil && i.LaunchTemplate == nil:
			launch = *i.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification.LaunchTemplateName
		// LaunchConfiguration
		case i.LaunchTemplate == nil && i.MixedInstancesPolicy == nil:
			launch = *i.LaunchConfigurationName
		// LaunchTemplate
		case i.MixedInstancesPolicy == nil && i.LaunchConfigurationName == nil:
			launch = *i.LaunchTemplate.LaunchTemplateName
		// No match launch pattern
		default:
			launch = "None"
		}

		for _, p := range i.TerminationPolicies {
			policies = append(policies, *p)
		}
		policy = strings.Join(policies, ",")

		list = append(list, Group{
			Name:                *i.AutoScalingGroupName,
			Launch:              launch,
			Desired:             strconv.FormatInt(*i.DesiredCapacity, 10),
			Min:                 strconv.FormatInt(*i.MinSize, 10),
			Max:                 strconv.FormatInt(*i.MaxSize, 10),
			TerminationPolicies: policy,
		})
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

// UpdateAutoScalingGroup return none
// input autoscaling.UpdateAutoScalingGroupInput
func (c *AutoScaling) UpdateAutoScalingGroup(input *autoscaling.UpdateAutoScalingGroupInput) error {
	if _, err := c.Client.UpdateAutoScalingGroup(input); err != nil {
		return fmt.Errorf("Update autoscaling group: %v", err)
	}

	return nil
}

func PrintGroups(wrt io.Writer, resources Groups) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"Launch",
		"Desired",
		"Min",
		"Max",
		"TerminationPolicies",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.GroupTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Group) GroupTabString() string {
	fields := []string{
		i.Name,
		i.Launch,
		i.Desired,
		i.Min,
		i.Max,
		i.TerminationPolicies,
	}

	return strings.Join(fields, "\t")
}
