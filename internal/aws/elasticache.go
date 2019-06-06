package aws

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
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
	Name           string
	Status         string
	Endpoint       string
	Port           string
	CacheClusterId string
	CacheNodeId    string
	CurrentRole    string
}

type GroupNodes []GroupNode

func NewEcSess(profile string, region string) *elasticache.ElastiCache {
	sess := getSession(profile, region)
	return elasticache.New(sess)
}

func DescribeCacheClusters(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

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

func DescribeReplicationGroups(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	svc := NewEcSess(profile, region)

	res, err := svc.DescribeReplicationGroups(nil)
	if err != nil {
		return fmt.Errorf("Describe running nodes: %v", err)
	}

	list := GroupNodes{}
	var (
		endpoint string
		port     string
	)
	for _, i := range res.ReplicationGroups {
		if i.ConfigurationEndpoint != nil {
			endpoint = *i.ConfigurationEndpoint.Address
			port = strconv.FormatInt(*i.ConfigurationEndpoint.Port, 10)
		}

		for _, n := range i.NodeGroups {
			if n.PrimaryEndpoint != nil {
				endpoint = *n.PrimaryEndpoint.Address
				port = strconv.FormatInt(*n.PrimaryEndpoint.Port, 10)
			}

			for _, nm := range n.NodeGroupMembers {
				if nm.CurrentRole == nil {
					nm.CurrentRole = aws.String("NULL")
				}

				list = append(list, GroupNode{
					Name:           *i.ReplicationGroupId,
					Status:         *i.Status,
					Endpoint:       endpoint,
					Port:           port,
					CacheClusterId: *nm.CacheClusterId,
					CacheNodeId:    *nm.CacheNodeId,
					CurrentRole:    *nm.CurrentRole,
				})
			}
		}

	}
	f := util.Formatln(
		list.Name(),
		list.Status(),
		list.Endpoint(),
		list.Port(),
		list.CacheClusterId(),
		list.CacheNodeId(),
		list.CurrentRole(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.Status,
			i.Endpoint,
			i.Port,
			i.CacheClusterId,
			i.CacheNodeId,
			i.CurrentRole,
		)
	}

	return nil
}

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

func (gn GroupNodes) Name() []string {
	name := []string{}
	for _, i := range gn {
		name = append(name, i.Name)
	}
	return name
}

func (gn GroupNodes) Status() []string {
	st := []string{}
	for _, i := range gn {
		st = append(st, i.Status)
	}
	return st
}

func (gn GroupNodes) Endpoint() []string {
	ep := []string{}
	for _, i := range gn {
		ep = append(ep, i.Endpoint)
	}
	return ep
}

func (gn GroupNodes) Port() []string {
	p := []string{}
	for _, i := range gn {
		p = append(p, i.Port)
	}
	return p
}

func (gn GroupNodes) CacheClusterId() []string {
	cc := []string{}
	for _, i := range gn {
		cc = append(cc, i.CacheClusterId)
	}
	return cc
}

func (gn GroupNodes) CacheNodeId() []string {
	cn := []string{}
	for _, i := range gn {
		cn = append(cn, i.CacheNodeId)
	}
	return cn
}

func (gn GroupNodes) CurrentRole() []string {
	cr := []string{}
	for _, i := range gn {
		cr = append(cr, i.CurrentRole)
	}
	return cr
}
