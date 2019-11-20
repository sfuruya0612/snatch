package aws

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sfuruya0612/snatch/internal/util"
)

// RDS client struct
type RDS struct {
	Client *rds.RDS
}

// NewRdsSess return RDS struct initialized
func NewRdsSess(profile, region string) *RDS {
	return &RDS{
		Client: rds.New(getSession(profile, region)),
	}
}

// DbInstance rds db instance struct
type DbInstance struct {
	Name             string
	DBInstanceClass  string
	Engine           string
	EngineVersion    string
	Storage          string
	DBInstanceStatus string
	Endpoint         string
	EndpointPort     string
}

// DbInstances DbInstance struct slice
type DbInstances []DbInstance

func (c *RDS) DescribeDBInstances() error {
	input := &rds.DescribeDBInstancesInput{}

	output, err := c.Client.DescribeDBInstances(input)
	if err != nil {
		return fmt.Errorf("Describe running instances: %v", err)
	}

	list := DbInstances{}
	for _, i := range output.DBInstances {
		port := strconv.FormatInt(*i.Endpoint.Port, 10)

		storage := strconv.FormatInt(*i.AllocatedStorage, 10) + "GB"

		list = append(list, DbInstance{
			Name:             *i.DBInstanceIdentifier,
			DBInstanceClass:  *i.DBInstanceClass,
			Engine:           *i.Engine,
			EngineVersion:    *i.EngineVersion,
			Storage:          storage,
			DBInstanceStatus: *i.DBInstanceStatus,
			Endpoint:         *i.Endpoint.Address,
			EndpointPort:     port,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.DBInstanceClass(),
		list.Engine(),
		list.EngineVersion(),
		list.Storage(),
		list.DBInstanceStatus(),
		list.Endpoint(),
		list.EndpointPort(),
	)

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.DBInstanceClass,
			i.Engine,
			i.EngineVersion,
			i.Storage,
			i.DBInstanceStatus,
			i.Endpoint,
			i.EndpointPort,
		)
	}

	return nil
}

func (dins DbInstances) Name() []string {
	name := []string{}
	for _, i := range dins {
		name = append(name, i.Name)
	}
	return name
}

func (dins DbInstances) DBInstanceClass() []string {
	class := []string{}
	for _, i := range dins {
		class = append(class, i.DBInstanceClass)
	}
	return class
}

func (dins DbInstances) Engine() []string {
	eg := []string{}
	for _, i := range dins {
		eg = append(eg, i.Engine)
	}
	return eg
}

func (dins DbInstances) EngineVersion() []string {
	egv := []string{}
	for _, i := range dins {
		egv = append(egv, i.EngineVersion)
	}
	return egv
}

func (dins DbInstances) Storage() []string {
	s := []string{}
	for _, i := range dins {
		s = append(s, i.Storage)
	}
	return s
}

func (dins DbInstances) DBInstanceStatus() []string {
	st := []string{}
	for _, i := range dins {
		st = append(st, i.DBInstanceStatus)
	}
	return st
}

func (dins DbInstances) Endpoint() []string {
	ep := []string{}
	for _, i := range dins {
		ep = append(ep, i.Endpoint)
	}
	return ep
}

func (dins DbInstances) EndpointPort() []string {
	port := []string{}
	for _, i := range dins {
		port = append(port, i.EndpointPort)
	}
	return port
}
