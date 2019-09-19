package aws

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Balancer struct {
	Name    string
	DNSName string
	Scheme  string
}

type Balancers []Balancer

func newElbSess(profile string, region string) *elb.ELB {
	sess := getSession(profile, region)
	return elb.New(sess)
}

func DescribeLoadBalancers(profile string, region string) error {
	elb := newElbSess(profile, region)

	res, err := elb.DescribeLoadBalancers(nil)
	if err != nil {
		return fmt.Errorf("No available load balancer: %v", err)
	}

	list := Balancers{}
	for _, i := range res.LoadBalancerDescriptions {

		list = append(list, Balancer{
			Name:    *i.LoadBalancerName,
			DNSName: *i.DNSName,
			Scheme:  *i.Scheme,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.DNSName(),
		list.Scheme(),
	)

	sort.Sort(list)
	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.DNSName,
			i.Scheme,
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

func (bal Balancers) Len() int {
	return len(bal)
}

func (bal Balancers) Swap(i, j int) {
	bal[i], bal[j] = bal[j], bal[i]
}

func (bal Balancers) Less(i, j int) bool {
	return bal[i].Name < bal[j].Name
}
