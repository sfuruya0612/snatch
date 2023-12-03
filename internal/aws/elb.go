package aws

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

// ELB structure is elb client.
type ELB struct {
	Client *elb.Client
}

// NewElbClient returns ELB struct initialized.
func NewElbClient(profile, region string) *ELB {
	return &ELB{
		Client: elb.NewFromConfig(GetSessionV2(profile, region)),
	}
}

// Balancer structure is elb information.
type Balancer struct {
	Name    string
	DNSName string
	Scheme  string
	Type    string
}

// DescribeLoadBalancers returns slice Balancer structure.
func (c *ELB) DescribeLoadBalancers(input *elb.DescribeLoadBalancersInput) ([]Balancer, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe load balancers: %v", err)
	}

	if len(output.LoadBalancerDescriptions) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []Balancer{}
	for _, i := range output.LoadBalancerDescriptions {
		list = append(list, Balancer{
			Name:    *i.LoadBalancerName,
			DNSName: *i.DNSName,
			Scheme:  *i.Scheme,
			Type:    "classic",
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

// ELBV2 structure is elb client.
type ELBV2 struct {
	Client *elbv2.Client
}

// NewElbV2Client returns ELBV2 struct initialized.
func NewElbV2Client(profile, region string) *ELBV2 {
	return &ELBV2{
		Client: elbv2.NewFromConfig(GetSessionV2(profile, region)),
	}
}

// DescribeLoadBalancersV2 returns slice Balancer structure.
func (c *ELBV2) DescribeLoadBalancersV2(input *elbv2.DescribeLoadBalancersInput) ([]Balancer, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe load balancers v2: %v", err)
	}

	if len(output.LoadBalancers) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []Balancer{}
	for _, i := range output.LoadBalancers {

		list = append(list, Balancer{
			Name:    *i.LoadBalancerName,
			DNSName: *i.DNSName,
			Scheme:  string(i.Scheme),
			Type:    string(i.Type),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func PrintBalancers(wrt io.Writer, resources []Balancer) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"DNSName",
		"Schema",
		"Type",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.ElbV2TabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Balancer) ElbV2TabString() string {
	fields := []string{
		i.Name,
		i.DNSName,
		i.Scheme,
		i.Type,
	}

	return strings.Join(fields, "\t")
}
