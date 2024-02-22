package aws

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// Route53 client struct
type Route53 struct {
	Client *route53.Client
}

// NewRoute53Sess return Route53 struct initialized
func NewRoute53Client(profile, region string) *Route53 {
	return &Route53{
		Client: route53.NewFromConfig(GetSession(profile, region)),
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

// ListHostedZones return Records
// input route53.ListHostedZonesInput
func (c *Route53) ListHostedZones(input *route53.ListHostedZonesInput) (Records, error) {
	zones, err := c.Client.ListHostedZones(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("list hostedzones: %v", err)
	}

	list := Records{}
	for _, h := range zones.HostedZones {
		s := strings.Split(*h.Id, "/")
		zoneid := s[2]

		rinput := &route53.ListResourceRecordSetsInput{
			HostedZoneId: h.Id,
		}

		paginator := route53.NewListResourceRecordSetsPaginator(c.Client, rinput)

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.Background())
			if err != nil {
				return nil, fmt.Errorf("list resource record sets: %v", err)
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
					Type:        string(r.Type),
					TTL:         ttl,
					DomainValue: value,
				})
			}
		}
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	return list, nil
}

func PrintRecords(wrt io.Writer, resources Records) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"ZoneId",
		"DomainName",
		"Type",
		"TTL",
		"DomainValue",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RecordTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Record) RecordTabString() string {
	fields := []string{
		i.ZoneId,
		i.DomainName,
		i.Type,
		i.TTL,
		i.DomainValue,
	}

	return strings.Join(fields, "\t")
}
