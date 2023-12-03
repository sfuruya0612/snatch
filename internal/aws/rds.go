package aws

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// RDS structure is rds client.
type RDS struct {
	Client *rds.Client
}

// NewRdsClient returns RDS struct initialized.
func NewRdsClient(profile, region string) *RDS {
	return &RDS{
		Client: rds.NewFromConfig(GetSessionV2(profile, region)),
	}
}

// DBInstance structure is rds instance information.
type DBInstance struct {
	Name             string
	DBInstanceClass  string
	Engine           string
	EngineVersion    string
	Storage          string
	StorageType      string
	DBInstanceStatus string
}

// DescribeDBInstances returns slice DBInstance structure.
func (c *RDS) DescribeDBInstances(input *rds.DescribeDBInstancesInput) ([]DBInstance, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeDBInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe db instances: %v", err)
	}

	if len(output.DBInstances) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []DBInstance{}
	for _, i := range output.DBInstances {
		list = append(list, DBInstance{
			Name:             *i.DBInstanceIdentifier,
			DBInstanceClass:  *i.DBInstanceClass,
			Engine:           *i.Engine,
			EngineVersion:    *i.EngineVersion,
			Storage:          strconv.Itoa(int(*i.AllocatedStorage)) + "GB",
			StorageType:      *i.StorageType,
			DBInstanceStatus: *i.DBInstanceStatus,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func PrintDBInstances(wrt io.Writer, resources []DBInstance) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"DBInstanceClass",
		"Engine",
		"EngineVersion",
		"Storage",
		"StrageType",
		"DBInstanceStatus",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("header join: %v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RdsTabString()); err != nil {
			return fmt.Errorf("resources join: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush: %v", err)
	}

	return nil
}

func (i *DBInstance) RdsTabString() string {
	fields := []string{
		i.Name,
		i.DBInstanceClass,
		i.Engine,
		i.EngineVersion,
		i.Storage,
		i.StorageType,
		i.DBInstanceStatus,
	}

	return strings.Join(fields, "\t")
}

// DBCluster structure is rds cluster information.
type DBCluster struct {
	Name          string
	EngineMode    string
	EngineVersion string
	Capacity      string
	Status        string
}

// DescribeDBClusters returns slice DBCluster structure.
func (c *RDS) DescribeDBClusters(input *rds.DescribeDBClustersInput) ([]DBCluster, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeDBClusters(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe db clusters: %v", err)
	}

	if len(output.DBClusters) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []DBCluster{}
	for _, i := range output.DBClusters {
		var cap string = "None"
		if i.Capacity != nil {
			cap = strconv.Itoa(int(*i.Capacity))
		}

		list = append(list, DBCluster{
			Name:          *i.DBClusterIdentifier,
			EngineMode:    *i.EngineMode,
			EngineVersion: *i.EngineVersion,
			Capacity:      cap,
			Status:        string(i.ActivityStreamStatus),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func PrintDBClusters(wrt io.Writer, resources []DBCluster) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"EngineMode",
		"EngineVersion",
		"Capacity",
		"Status",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("header join: %v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RdsClusterTabString()); err != nil {
			return fmt.Errorf("resources join: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush: %v", err)
	}

	return nil
}

func (i *DBCluster) RdsClusterTabString() string {
	fields := []string{
		i.Name,
		i.EngineMode,
		i.EngineVersion,
		i.Capacity,
		i.Status,
	}

	return strings.Join(fields, "\t")
}

// DBClusterEndpoint structure is rds cluster endpoint information.
type DBClusterEndpoint struct {
	Endpoint     string
	EndpointType string
	Status       string
}

// DescribeDBClusterEndpoints returns slice DBInstance structure.
func (c *RDS) DescribeDBClusterEndpoints(input *rds.DescribeDBClusterEndpointsInput) ([]DBClusterEndpoint, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeDBClusterEndpoints(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe db cluster endpoints: %v", err)
	}

	if len(output.DBClusterEndpoints) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []DBClusterEndpoint{}
	for _, i := range output.DBClusterEndpoints {
		list = append(list, DBClusterEndpoint{
			Endpoint:     *i.Endpoint,
			EndpointType: *i.EndpointType,
			Status:       *i.Status,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Endpoint < list[j].Endpoint
	})

	return list, nil
}

func PrintDBClusterEndpoints(wrt io.Writer, resources []DBClusterEndpoint) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Endpoint",
		"EndpointType",
		"Status",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("header join: %v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RdsClusterEndpointTabString()); err != nil {
			return fmt.Errorf("resources join: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush: %v", err)
	}

	return nil
}

func (i *DBClusterEndpoint) RdsClusterEndpointTabString() string {
	fields := []string{
		i.Endpoint,
		i.EndpointType,
		i.Status,
	}

	return strings.Join(fields, "\t")
}
