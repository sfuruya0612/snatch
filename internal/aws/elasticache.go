package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

type CacheNode struct {
	Name               string
	CacheNodeType      string
	Engine             string
	EngineVersion      string
	CacheClusterStatus string
}

type CacheNodes []CacheNode

type GroupNode struct {
	Name               string
	CacheNodeType      string
	Engine             string
	EngineVersion      string
	CacheClusterStatus string
}

type GroupNodes []GroupNode

func NewEcSess(profile string, region string) *elasticache.ElastiCache {
	sess := getSession(profile, region)
	return elasticache.New(sess)
}

func DescribeCacheClusters(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	svc := NewEcSess(profile, region)

	res, err := svc.DescribeCacheClusters(nil)
	if err != nil {
		return fmt.Errorf("Describe running nodes: %v", err)
	}

	list := CacheNodes{}
	for _, i := range res.CacheClusters {
		list = append(list, CacheNode{
			Name:               *i.CacheClusterId,
			CacheNodeType:      *i.CacheNodeType,
			Engine:             *i.Engine,
			EngineVersion:      *i.EngineVersion,
			CacheClusterStatus: *i.CacheClusterStatus,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.CacheNodeType(),
		list.Engine(),
		list.EngineVersion(),
		list.CacheClusterStatus(),
		[]string{""},
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.CacheNodeType,
			i.Engine,
			i.EngineVersion,
			i.CacheClusterStatus,
		)
	}

	return nil
}

/*
func DescribeReplicationGroups(c *cli.Context) error {
	profile := c.String("profile")
	region := c.String("region")

	svc := NewEcSess(profile, region)

	res, err := svc.DescribeReplicationGroups(nil)
	if err != nil {
		return fmt.Errorf("Describe running nodes: %v", err)
	}

	list := CacheNodes{}
	for _, i := range res.CacheClusters {
		list = append(list, CacheNode{
			Name:               *i.CacheClusterId,
			CacheNodeType:      *i.CacheNodeType,
			Engine:             *i.Engine,
			EngineVersion:      *i.EngineVersion,
			CacheClusterStatus: *i.CacheClusterStatus,
		})
	}
	f := util.Formatln(
		list.Name(),
		list.CacheNodeType(),
		list.Engine(),
		list.EngineVersion(),
		list.CacheClusterStatus(),
		[]string{""},
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.CacheNodeType,
			i.Engine,
			i.EngineVersion,
			i.CacheClusterStatus,
		)
	}

	return nil
}
*/
func (cn CacheNodes) Name() []string {
	name := []string{}
	for _, i := range cn {
		name = append(name, i.Name)
	}
	return name
}

func (cn CacheNodes) CacheNodeType() []string {
	ty := []string{}
	for _, i := range cn {
		ty = append(ty, i.CacheNodeType)
	}
	return ty
}

func (cn CacheNodes) Engine() []string {
	eg := []string{}
	for _, i := range cn {
		eg = append(eg, i.Engine)
	}
	return eg
}

func (cn CacheNodes) EngineVersion() []string {
	egv := []string{}
	for _, i := range cn {
		egv = append(egv, i.EngineVersion)
	}
	return egv
}

func (cn CacheNodes) CacheClusterStatus() []string {
	st := []string{}
	for _, i := range cn {
		st = append(st, i.CacheClusterStatus)
	}
	return st
}
