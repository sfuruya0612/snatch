package aws

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// ELB client struct
type ELB struct {
	Client *elb.ELB
}

// NewElbSess return ELB struct initialized
func NewElbSess(profile, region string) *ELB {
	return &ELB{
		Client: elb.New(getSession(profile, region)),
	}
}

// ELBV2 client struct
type ELBV2 struct {
	Client *elbv2.ELBV2
}

// NewElbV2Sess return ELBV2 struct initialized
func NewElbV2Sess(profile, region string) *ELBV2 {
	return &ELBV2{
		Client: elbv2.New(getSession(profile, region)),
	}
}

// Balancer elb struct
type Balancer struct {
	Name      string
	DNSName   string
	Scheme    string
	Instances string
	Type      string
}

// Balancers Balancer struct slice
type Balancers []Balancer

// DescribeLoadBalancers return Balancers
// input elb.DescribeLoadBalancersInput
// Classic Load Balancer
func (c *ELB) DescribeLoadBalancers(input *elb.DescribeLoadBalancersInput) (Balancers, error) {
	output, err := c.Client.DescribeLoadBalancers(input)
	if err != nil {
		return nil, fmt.Errorf("Describe load balancers: %v", err)
	}

	list := Balancers{}
	for _, i := range output.LoadBalancerDescriptions {
		var (
			instances []string
			instance  string
		)

		if i.Instances != nil {
			for _, ii := range i.Instances {
				instances = append(instances, *ii.InstanceId)
			}

			instance = strings.Join(instances[:], ",")
		}

		list = append(list, Balancer{
			Name:      *i.LoadBalancerName,
			DNSName:   *i.DNSName,
			Scheme:    *i.Scheme,
			Instances: instance,
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

// DescribeLoadBalancersV2 return Balancers
// input elbv2.DescribeLoadBalancersInput
// Application and Network Load Balancer
func (c *ELBV2) DescribeLoadBalancersV2(input *elbv2.DescribeLoadBalancersInput) (Balancers, error) {
	output, err := c.Client.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, fmt.Errorf("Describe load balancers v2: %v", err)
	}

	list := Balancers{}
	for _, i := range output.LoadBalancers {

		list = append(list, Balancer{
			Name:    *i.LoadBalancerName,
			DNSName: *i.DNSName,
			Scheme:  *i.Scheme,
			Type:    *i.Type,
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

func PrintBalancers(wrt io.Writer, resources Balancers) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"DNSName",
		"Schema",
		"Instances",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.ElbTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Balancer) ElbTabString() string {
	fields := []string{
		i.Name,
		i.DNSName,
		i.Scheme,
		i.Instances,
	}

	return strings.Join(fields, "\t")
}
func PrintBalancersV2(wrt io.Writer, resources Balancers) error {
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
