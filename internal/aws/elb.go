package aws

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/sfuruya0612/snatch/internal/util"
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

// Balancer elb struct
type Balancer struct {
	Name      string
	DNSName   string
	Scheme    string
	Instances string
}

// Balancers Balancer struct slice
type Balancers []Balancer

func (c *ELB) DescribeLoadBalancers() error {
	output, err := c.Client.DescribeLoadBalancers(nil)
	if err != nil {
		return fmt.Errorf("No available load balancer: %v", err)
	}

	list := Balancers{}
	for _, i := range output.LoadBalancerDescriptions {

		var instance string
		if i.Instances != nil {
			var instances []string

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
	f := util.Formatln(
		list.Name(),
		list.DNSName(),
		list.Scheme(),
		list.Instances(),
	)

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.DNSName,
			i.Scheme,
			i.Instances,
		)
	}

	return nil
}

func (bal Balancers) Name() []string {
	name := []string{}
	for _, i := range bal {
		name = append(name, i.Name)
	}
	return name
}

func (bal Balancers) DNSName() []string {
	dname := []string{}
	for _, i := range bal {
		dname = append(dname, i.DNSName)
	}
	return dname
}

func (bal Balancers) Scheme() []string {
	scheme := []string{}
	for _, i := range bal {
		scheme = append(scheme, i.Scheme)
	}
	return scheme
}

func (bal Balancers) Instances() []string {
	ins := []string{}
	for _, i := range bal {
		ins = append(ins, i.Instances)
	}
	return ins
}
