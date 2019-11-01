package aws

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/sfuruya0612/snatch/internal/util"
)

type CacheNode struct {
	Name               string
	CacheNodeType      string
	Engine             string
	EngineVersion      string
	CacheClusterStatus string
	Status             string
	Endpoint           string
	Port               string
	CacheClusterID     string
	CacheNodeID        string
	CurrentRole        string
}

type CacheNodes []CacheNode

type GroupNode struct {
	Name string
}

type GroupNodes []GroupNode

func newElasticacheSess(profile string, region string) *elasticache.ElastiCache {
	sess := getSession(profile, region)
	return elasticache.New(sess)
}

func DescribeCacheClusters(profile string, region string) error {
	svc := newElasticacheSess(profile, region)

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

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

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

func DescribeReplicationGroups(profile string, region string) error {
	svc := newElasticacheSess(profile, region)

	res, err := svc.DescribeReplicationGroups(nil)
	if err != nil {
		return fmt.Errorf("Describe running nodes: %v", err)
	}

	list := CacheNodes{}
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

				role := "None"
				if nm.CurrentRole != nil {
					role = *nm.CurrentRole
				}

				list = append(list, CacheNode{
					Name:           *i.ReplicationGroupId,
					Status:         *i.Status,
					Endpoint:       endpoint,
					Port:           port,
					CacheClusterID: *nm.CacheClusterId,
					CacheNodeID:    *nm.CacheNodeId,
					CurrentRole:    role,
				})
			}
		}

	}
	f := util.Formatln(
		list.Name(),
		list.Status(),
		list.Endpoint(),
		list.Port(),
		list.CacheClusterID(),
		list.CacheNodeID(),
		list.CurrentRole(),
	)

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	for _, i := range list {
		fmt.Printf(
			f,
			i.Name,
			i.Status,
			i.Endpoint,
			i.Port,
			i.CacheClusterID,
			i.CacheNodeID,
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

func (cn CacheNodes) Status() []string {
	st := []string{}
	for _, i := range cn {
		st = append(st, i.Status)
	}
	return st
}

func (cn CacheNodes) Endpoint() []string {
	ep := []string{}
	for _, i := range cn {
		ep = append(ep, i.Endpoint)
	}
	return ep
}

func (cn CacheNodes) Port() []string {
	p := []string{}
	for _, i := range cn {
		p = append(p, i.Port)
	}
	return p
}

func (cn CacheNodes) CacheClusterID() []string {
	cc := []string{}
	for _, i := range cn {
		cc = append(cc, i.CacheClusterID)
	}
	return cc
}

func (cn CacheNodes) CacheNodeID() []string {
	cnid := []string{}
	for _, i := range cn {
		cnid = append(cnid, i.CacheNodeID)
	}
	return cnid
}

func (cn CacheNodes) CurrentRole() []string {
	cr := []string{}
	for _, i := range cn {
		cr = append(cr, i.CurrentRole)
	}
	return cr
}
