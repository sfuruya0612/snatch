package aws

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Balancerv2 struct {
	Name    string
	DNSName string
	Scheme  string
	Type    string
}

type Balancersv2 []Balancerv2

func newElbv2Sess(profile string, region string) *elbv2.ELBV2 {
	sess := getSession(profile, region)
	return elbv2.New(sess)
}

func DescribeLoadBalancersv2(profile string, region string) error {
	elbv2 := newElbv2Sess(profile, region)

	res, err := elbv2.DescribeLoadBalancers(nil)
	if err != nil {
		return fmt.Errorf("No available load balancer: %v", err)
	}

	list := Balancersv2{}
	for _, i := range res.LoadBalancers {

		list = append(list, Balancerv2{
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

func (bal Balancersv2) Name() []string {
	name := []string{}
	for _, i := range bal {
		name = append(name, i.Name)
	}
	return name
}

func (bal Balancersv2) DNSName() []string {
	dname := []string{}
	for _, i := range bal {
		dname = append(dname, i.DNSName)
	}
	return dname
}

func (bal Balancersv2) Scheme() []string {
	scheme := []string{}
	for _, i := range bal {
		scheme = append(scheme, i.Scheme)
	}
	return scheme
}

func (bal Balancersv2) Type() []string {
	t := []string{}
	for _, i := range bal {
		t = append(t, i.Type)
	}
	return t
}
