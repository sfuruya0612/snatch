package aws

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

// ElastiCache client struct
type ElastiCache struct {
	Client *elasticache.ElastiCache
}

// NewElastiCacheSess return ElastiCache struct initialized
func NewElastiCacheSess(profile, region string) *ElastiCache {
	return &ElastiCache{
		Client: elasticache.New(getSession(profile, region)),
	}
}

// CacheNode elasticache cachenode struct
type CacheNode struct {
	Name               string
	CacheNodeType      string
	Engine             string
	EngineVersion      string
	CacheClusterStatus string
	Status             string
	Endpoint           string
	Port               string
	CacheClusterId     string
	CacheNodeId        string
	CurrentRole        string
}

// CacheNodes CacheNode struct slice
type CacheNodes []CacheNode

// DescribeCacheClusters return CacheNodes
// input elasticache.DescribeCacheClustersInput
func (c *ElastiCache) DescribeCacheClusters(input *elasticache.DescribeCacheClustersInput) (CacheNodes, error) {
	output, err := c.Client.DescribeCacheClusters(input)
	if err != nil {
		return nil, fmt.Errorf("Describe cache cluster: %v", err)
	}

	list := CacheNodes{}
	for _, i := range output.CacheClusters {
		list = append(list, CacheNode{
			Name:               *i.CacheClusterId,
			CacheNodeType:      *i.CacheNodeType,
			Engine:             *i.Engine,
			EngineVersion:      *i.EngineVersion,
			CacheClusterStatus: *i.CacheClusterStatus,
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

// DescribeReplicationGroups return CacheNodes
// input elasticache.DescribeCacheClustersInput
func (c *ElastiCache) DescribeReplicationGroups(input *elasticache.DescribeReplicationGroupsInput) (CacheNodes, error) {
	output, err := c.Client.DescribeReplicationGroups(input)
	if err != nil {
		return nil, fmt.Errorf("Describe replication groups: %v", err)
	}

	list := CacheNodes{}
	for _, i := range output.ReplicationGroups {
		var (
			endpoint string
			port     string
		)

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
					CacheClusterId: *nm.CacheClusterId,
					CacheNodeId:    *nm.CacheNodeId,
					CurrentRole:    role,
				})
			}
		}

	}
	if len(list) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func PrintCacheClusters(wrt io.Writer, resources CacheNodes) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"CacheNodeType",
		"Engine",
		"EngineVersion",
		"CacheClusterStatus",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.CcTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *CacheNode) CcTabString() string {
	fields := []string{
		i.Name,
		i.CacheNodeType,
		i.Engine,
		i.EngineVersion,
		i.CacheClusterStatus,
	}

	return strings.Join(fields, "\t")
}

func PrintRepricationGroups(wrt io.Writer, resources CacheNodes) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"Status",
		"Endpoint",
		"Port",
		"CacheClusterId",
		"CacheNodeId",
		"CurrentRole",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RgTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *CacheNode) RgTabString() string {
	fields := []string{
		i.Name,
		i.Status,
		i.Endpoint,
		i.Port,
		i.CacheClusterId,
		i.CacheNodeId,
		i.CurrentRole,
	}

	return strings.Join(fields, "\t")
}
