package aws

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/imdario/mergo"
)

// ElastiCache structure is elasticache client.
type ElastiCache struct {
	Client *elasticache.Client
}

// NewElastiCacheClient return ElastiCache struct initialized.
func NewElastiCacheClient(profile, region string) *ElastiCache {
	return &ElastiCache{
		Client: elasticache.NewFromConfig(GetSessionV2(profile, region)),
	}
}

// CacheNode structure is elasticache node information.
type CacheNode struct {
	ReplicationGroupId string
	CacheClusterId     string
	CacheNodeId        string
	CacheNodeType      string
	Engine             string
	EngineVersion      string
	CurrentRole        string
	CacheClusterStatus string
	CacheNodeStatus    string
}

// DescribeCacheClusters returns slice CacheNode structure.
func (c *ElastiCache) DescribeCacheClusters(input *elasticache.DescribeCacheClustersInput) ([]CacheNode, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeCacheClusters(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describe cache cluster: %v", err)
	}

	if len(output.CacheClusters) == 0 {
		return nil, fmt.Errorf("no resources")
	}

	list := []CacheNode{}
	for _, cc := range output.CacheClusters {
		replicationGroupId := "None"
		if cc.ReplicationGroupId != nil {
			replicationGroupId = *cc.ReplicationGroupId
		}

		list = append(list, CacheNode{
			ReplicationGroupId: replicationGroupId,
			CacheClusterId:     *cc.CacheClusterId,
			CacheNodeType:      *cc.CacheNodeType,
			Engine:             *cc.Engine,
			EngineVersion:      *cc.EngineVersion,
			CacheClusterStatus: *cc.CacheClusterStatus,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ReplicationGroupId < list[j].ReplicationGroupId
	})

	return list, nil
}

// DescribeReplicationGroups returns CacheNode structure.
func (c *ElastiCache) DescribeReplicationGroups(input *elasticache.DescribeReplicationGroupsInput, node CacheNode) (CacheNode, error) {
	ctx := context.TODO()
	output, err := c.Client.DescribeReplicationGroups(ctx, input)
	if err != nil {
		return CacheNode{}, fmt.Errorf("describe replication groups: %v", err)
	}

	for _, rg := range output.ReplicationGroups {
		for _, ng := range rg.NodeGroups {
			for _, ngm := range ng.NodeGroupMembers {
				role := "None"
				if ngm.CurrentRole != nil {
					role = *ngm.CurrentRole
				}

				dst := CacheNode{
					CacheNodeId:     *ngm.CacheNodeId,
					CurrentRole:     role,
					CacheNodeStatus: *ng.Status,
				}
				if err := mergo.Merge(&node, dst, mergo.WithOverride); err != nil {
					return CacheNode{}, fmt.Errorf("struct merge: %v", err)
				}
				fmt.Printf("%#v\n", node)
			}
		}
	}

	return node, nil
}

func PrintNodes(wrt io.Writer, resources []CacheNode) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"ReplicationGroupId",
		"CacheClusterId",
		// "CacheNodeId",
		"CacheNodeType",
		"Engine",
		"EngineVersion",
		// "CurrentRole",
		"CacheClusterStatus",
		// "CacheNodeStatus",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.NodeTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *CacheNode) NodeTabString() string {
	fields := []string{
		i.ReplicationGroupId,
		i.CacheClusterId,
		// i.CacheNodeId,
		i.CacheNodeType,
		i.Engine,
		i.EngineVersion,
		// i.CurrentRole,
		i.CacheClusterStatus,
		// i.CacheNodeStatus,
	}

	return strings.Join(fields, "\t")
}
