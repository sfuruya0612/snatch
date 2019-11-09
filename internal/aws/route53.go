package aws

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sfuruya0612/snatch/internal/util"
)

// Route53 client struct
type Route53 struct {
	Client *route53.Route53
}

// NewRoute53Sess return Route53 struct initialized
func NewRoute53Sess(profile, region string) *Route53 {
	return &Route53{
		Client: route53.New(getSession(profile, region)),
	}
}

// Record route53 set record struct
type Record struct {
	ZoneId      string
	DomainName  string
	Type        string
	TTL         string
	DomainValue string
}

// Records Record struct slice
type Records []Record

func (c *Route53) ListHostedZones() error {
	zones, err := c.Client.ListHostedZones(nil)
	if err != nil {
		return fmt.Errorf("List hostedzones sets: %v", err)
	}

	list := Records{}
	for _, h := range zones.HostedZones {
		id := strings.Split(*h.Id, "/")
		zoneid := id[2]

		input := &route53.ListResourceRecordSetsInput{
			HostedZoneId: h.Id,
		}

		output, err := c.Client.ListResourceRecordSets(input)
		if err != nil {
			return fmt.Errorf("List record sets: %v", err)
		}

		for _, r := range output.ResourceRecordSets {

			if r.TTL == nil {
				r.TTL = aws.Int64(0000)
			}
			ttl := strconv.FormatInt(*r.TTL, 10)

			var value string
			if r.AliasTarget == nil {
				var values []string

				for _, rr := range r.ResourceRecords {
					values = append(values, *rr.Value)
				}

				value = strings.Join(values[:], ",")
			} else if r.ResourceRecords == nil {
				value = *r.AliasTarget.DNSName
			}

			list = append(list, Record{
				ZoneId:      zoneid,
				DomainName:  *r.Name,
				Type:        *r.Type,
				TTL:         ttl,
				DomainValue: value,
			})
		}
	}
	f := util.Formatln(
		list.ZoneId(),
		list.DomainName(),
		list.Type(),
		list.TTL(),
		list.DomainValue(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.ZoneId,
			i.DomainName,
			i.Type,
			i.TTL,
			i.DomainValue,
		)
	}

	return nil
}

func (rec Records) ZoneId() []string {
	zid := []string{}
	for _, i := range rec {
		zid = append(zid, i.ZoneId)
	}
	return zid
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
		dvalue = append(dvalue, i.DomainValue)
	}
	return dvalue
}
