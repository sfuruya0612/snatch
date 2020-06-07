package aws

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// CostExplorer client struct
type CostExplorer struct {
	Client *costexplorer.CostExplorer
}

// NewCeSess return CostExplorer struct initialized
func NewCeSess(profile, region string) *CostExplorer {
	return &CostExplorer{
		Client: costexplorer.New(GetSession(profile, region)),
	}
}

// Cost cost and usage struct
type Cost struct {
	Start string
	End   string
	Usage []usage
}

// usage the cost struct per aws service
type usage struct {
	Key    string
	Amount string
}

// Costs Cost struct slice
type Costs []Cost

// GetCostAndUsage return Costs
// input costexplorer.GetCostAndUsageInput
func (c *CostExplorer) GetCostAndUsage(input *costexplorer.GetCostAndUsageInput) (Costs, error) {
	output, err := c.Client.GetCostAndUsage(input)
	if err != nil {
		return nil, fmt.Errorf("Get cost and usage: %v", err)
	}

	list := Costs{}
	services := []usage{}
	for _, o := range output.ResultsByTime {

		for _, g := range o.Groups {
			services = append(services, usage{
				Key:    *g.Keys[0],
				Amount: *g.Metrics["BlendedCost"].Amount,
			})
		}
		list = append(list, Cost{
			Start: *o.TimePeriod.Start,
			End:   *o.TimePeriod.End,
			Usage: services,
		})
	}

	return list, nil
}

func PrintUsage(wrt io.Writer, resources []usage) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.CostTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *usage) CostTabString() string {
	fields := []string{
		i.Key,
		i.Amount,
	}

	return strings.Join(fields, "\t")
}
