package aws

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

type Record struct {
	Id          string
	Name        string
	DomainName  string
	Type        string
	TTL         string
	DomainValue []string
}

type Records []Record

type ListResourceRecordSetsInput struct {
	HostedZoneId *string `location:"uri" locationName:"Id" type:"string" required:"true"`
}

func NewRoute53Sess(profile string, region string) *route53.Route53 {
	sess := getSession(profile, region)
	return route53.New(sess)
}

func ListHostedZones(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	svc := NewRoute53Sess(profile, region)

	res, err := svc.ListHostedZones(nil)
	if err != nil {
		return fmt.Errorf("List hostedzones sets: %v", err)
	}

	list := Records{}
	for _, h := range res.HostedZones {
		Id := h.Id

		input := &route53.ListResourceRecordSetsInput{
			HostedZoneId: Id,
		}

		rec, err := svc.ListResourceRecordSets(input)
		if err != nil {
			return fmt.Errorf("List record sets: %v", err)
		}

		for _, r := range rec.ResourceRecordSets {

			if r.TTL == nil {
				r.TTL = aws.Int64(0000)
			}

			ttl := strconv.FormatInt(*r.TTL, 10)

			var value []string
			if r.AliasTarget == nil {
				for _, rr := range r.ResourceRecords {
					value = append(value, *rr.Value)
				}
			} else if r.ResourceRecords == nil {
				value = append(value, *r.AliasTarget.DNSName)
			}

			list = append(list, Record{
				Id:          *h.Id,
				Name:        *h.Name,
				DomainName:  *r.Name,
				Type:        *r.Type,
				TTL:         ttl,
				DomainValue: value,
			})
		}
	}
	f := util.Formatln(
		list.Id(),
		list.Name(),
		list.DomainName(),
		list.Type(),
		list.TTL(),
		list.DomainValue(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Id,
			i.Name,
			i.DomainName,
			i.Type,
			i.TTL,
			i.DomainValue,
		)
	}
	return nil
}

func (rec Records) Id() []string {
	id := []string{}
	for _, i := range rec {
		id = append(id, i.Id)
	}
	return id
}

func (rec Records) Name() []string {
	name := []string{}
	for _, i := range rec {
		name = append(name, i.Name)
	}
	return name
}

func (rec Records) DomainName() []string {
	dname := []string{}
	for _, i := range rec {
		dname = append(dname, i.DomainName)
	}
	return dname

}
func (rec Records) Type() []string {
	ty := []string{}
	for _, i := range rec {
		ty = append(ty, i.Type)
	}
	return ty
}

func (rec Records) TTL() []string {
	ttl := []string{}
	for _, i := range rec {
		ttl = append(ttl, i.TTL)
	}
	return ttl
}

func (rec Records) DomainValue() []string {
	dvalue := []string{}
	for _, i := range rec {
		dvalue = append(dvalue, i.DomainValue...)
	}
	return dvalue
}
