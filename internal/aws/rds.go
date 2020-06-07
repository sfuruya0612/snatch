package aws

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/rds"
)

// RDS client struct
type RDS struct {
	Client *rds.RDS
}

// NewRdsSess return RDS struct initialized
func NewRdsSess(profile, region string) *RDS {
	return &RDS{
		Client: rds.New(GetSession(profile, region)),
	}
}

// DBInstance rds db instance struct
type DBInstance struct {
	Name             string
	DBInstanceClass  string
	Engine           string
	EngineVersion    string
	Storage          string
	DBInstanceStatus string
	Endpoint         string
	EndpointPort     string
}

// DBInstances DBInstance struct slice
type DBInstances []DBInstance

// DBCluster rds db cluster struct
type DBCluster struct {
	Name          string
	EngineMode    string
	EngineVersion string
	Capacity      string
	Status        string
	Endpoint      string
	EndpointPort  string
}

// DBClusters DBCluster struct slice
type DBClusters []DBCluster

// DescribeDBInstances return DBInstances
// input rds.DescribeDBInstancesInput
func (c *RDS) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (DBInstances, error) {
	output, err := c.Client.DescribeDBInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Describe running instances: %v", err)
	}

	list := DBInstances{}
	for _, i := range output.DBInstances {
		var (
			addr string = "None"
			port string = "None"
		)
		if i.Endpoint != nil {
			addr = *i.Endpoint.Address
			port = strconv.FormatInt(*i.Endpoint.Port, 10)
		}

		storage := strconv.FormatInt(*i.AllocatedStorage, 10) + "GB"

		list = append(list, DBInstance{
			Name:             *i.DBInstanceIdentifier,
			DBInstanceClass:  *i.DBInstanceClass,
			Engine:           *i.Engine,
			EngineVersion:    *i.EngineVersion,
			Storage:          storage,
			DBInstanceStatus: *i.DBInstanceStatus,
			Endpoint:         addr,
			EndpointPort:     port,
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

// DescribeDBClusters return DBClusters
// input rds.DescribeDBClustersInput
func (c *RDS) DescribeDBClusters(input *rds.DescribeDBClustersInput) (DBClusters, error) {
	output, err := c.Client.DescribeDBClusters(input)
	if err != nil {
		return nil, fmt.Errorf("Describe db clusters: %v", err)
	}

	list := DBClusters{}
	for _, i := range output.DBClusters {
		list = append(list, DBCluster{
			Name:          *i.DBClusterIdentifier,
			EngineMode:    *i.EngineMode,
			EngineVersion: *i.EngineVersion,
			Capacity:      strconv.FormatInt(*i.Capacity, 10),
			Status:        *i.ActivityStreamStatus,
			Endpoint:      *i.Endpoint,
			EndpointPort:  strconv.FormatInt(*i.Port, 10),
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

func PrintDBInstances(wrt io.Writer, resources DBInstances) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Endpoint",
		"EndpointPort",
		"DBInstanceClass",
		"Engine",
		"EngineVersion",
		"Storage",
		"DBInstanceStatus",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RdsTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *DBInstance) RdsTabString() string {
	fields := []string{
		i.Endpoint,
		i.EndpointPort,
		i.DBInstanceClass,
		i.Engine,
		i.EngineVersion,
		i.Storage,
		i.DBInstanceStatus,
	}

	return strings.Join(fields, "\t")
}

func PrintDBClusters(wrt io.Writer, resources DBClusters) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Endpoint",
		"EndpointPort",
		"EngineMode",
		"EngineVersion",
		"Capacity",
		"Status",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RdsClusterTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *DBCluster) RdsClusterTabString() string {
	fields := []string{
		i.Endpoint,
		i.EndpointPort,
		i.EngineMode,
		i.EngineVersion,
		i.Capacity,
		i.Status,
	}

	return strings.Join(fields, "\t")
}
