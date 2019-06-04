package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

type DbInstance struct {
	Name             string
	DBInstanceClass  string
	Engine           string
	EngineVersion    string
	DBInstanceStatus string
	Endpoint         string
	EndpointPort     int64
}

type DbInstances []DbInstance

func NewRdsSess(profile string, region string) *rds.RDS {
	sess := getSession(profile, region)
	return rds.New(sess)
}

func DescribeDBInstances(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	svc := NewRdsSess(profile, region)

	res, err := svc.DescribeDBInstances(nil)
	if err != nil {
		return fmt.Errorf("Describe running instances: %v", err)
	}

	list := DbInstances{}
	for _, i := range res.DBInstances {
		list = append(list, DbInstance{
			Name:             *i.DBInstanceIdentifier,
			DBInstanceClass:  *i.DBInstanceClass,
			Engine:           *i.Engine,
			EngineVersion:    *i.EngineVersion,
			DBInstanceStatus: *i.DBInstanceStatus,
			Endpoint:         *i.Endpoint.Address,
			EndpointPort:     *i.Endpoint.Port,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.DBInstanceClass(),
		list.Engine(),
		list.EngineVersion(),
		list.DBInstanceStatus(),
		list.Endpoint(),
		//		list.EndpointPort(),
		[]string{""},
	)

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

/*
func (dins DbInstances) EndpointPort() []string {
	port := []string{}
	for _, i := range dins {
		port = append(port, i.EndpointPort)
	}
	return port
}*/
