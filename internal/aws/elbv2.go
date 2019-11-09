package aws

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sfuruya0612/snatch/internal/util"
)

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

// BalancerV2 elb struct
type BalancerV2 struct {
	Name    string
	DNSName string
	Scheme  string
	Type    string
}

// BalancersV2 BalancerV2 struct slice
type BalancersV2 []BalancerV2

func (c *ELBV2) DescribeLoadBalancersV2() error {
	output, err := c.Client.DescribeLoadBalancers(nil)
	if err != nil {
		return fmt.Errorf("No available load balancer: %v", err)
	}

	list := BalancersV2{}
	for _, i := range output.LoadBalancers {

		list = append(list, BalancerV2{
			Name:    *i.LoadBalancerName,
			DNSName: *i.DNSName,
			Scheme:  *i.Scheme,
			Type:    *i.Type,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.DNSName(),
		list.Scheme(),
		list.Type(),
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
			i.Type,
		)
	}

	return nil
}

func (bal BalancersV2) Name() []string {
	name := []string{}
	for _, i := range bal {
		name = append(name, i.Name)
	}
	return name
}

func (bal BalancersV2) DNSName() []string {
	dname := []string{}
	for _, i := range bal {
		dname = append(dname, i.DNSName)
	}
	return dname
}

func (bal BalancersV2) Scheme() []string {
	scheme := []string{}
	for _, i := range bal {
		scheme = append(scheme, i.Scheme)
	}
	return scheme
}

func (bal BalancersV2) Type() []string {
	t := []string{}
	for _, i := range bal {
		t = append(t, i.Type)
	}
	return t
}
