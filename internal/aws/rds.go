package aws

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sfuruya0612/snatch/internal/util"
)

type DbInstance struct {
	Name             string
	DBInstanceClass  string
	Engine           string
	EngineVersion    string
	DBInstanceStatus string
	Endpoint         string
	EndpointPort     string
}

type DbInstances []DbInstance

func newRdsSess(profile string, region string) *rds.RDS {
	sess := getSession(profile, region)
	return rds.New(sess)
}

func DescribeDBInstances(profile string, region string) error {
	rds := newRdsSess(profile, region)

	res, err := rds.DescribeDBInstances(nil)
	if err != nil {
		return fmt.Errorf("Describe running instances: %v", err)
	}

	list := DbInstances{}
	for _, i := range res.DBInstances {
		port := strconv.FormatInt(*i.Endpoint.Port, 10)

		list = append(list, DbInstance{
			Name:             *i.DBInstanceIdentifier,
			DBInstanceClass:  *i.DBInstanceClass,
			Engine:           *i.Engine,
			EngineVersion:    *i.EngineVersion,
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
		list.DBInstanceStatus(),
		list.Endpoint(),
		list.EndpointPort(),
	)

	sort.Sort(list)
	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.DBInstanceClass,
			i.Engine,
			i.EngineVersion,
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

func (dins DbInstances) Len() int {
	return len(dins)
}

func (dins DbInstances) Swap(i, j int) {
	dins[i], dins[j] = dins[j], dins[i]
}

func (dins DbInstances) Less(i, j int) bool {
	return dins[i].Name < dins[j].Name
}
