package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var ElastiCache = &cli.Command{
	Name:    "elasticache",
	Aliases: []string{"ec"},
	Usage:   "Get a list of ElastiCache",
	Action: func(c *cli.Context) error {
		return getEcNodeList(c.String("profile"), c.String("region"))
	},
}

func getEcNodeList(profile, region string) error {
	c := saws.NewElastiCacheClient(profile, region)

	clusters, err := c.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// TODO: Get the information of the replication group together. Control the loop of ReplicationGroupId.
	// nodes := []saws.CacheNode{}
	// for _, n := range clusters {
	// 	input := &elasticache.DescribeReplicationGroupsInput{
	// 		ReplicationGroupId: &n.ReplicationGroupId,
	// 	}

	// 	group, err := c.DescribeReplicationGroups(input, n)
	// 	if err != nil {
	// 		return fmt.Errorf("%v", err)
	// 	}
	// 	nodes = append(nodes, group)
	// }

	// sort.Slice(nodes, func(i, j int) bool {
	// 	return nodes[i].ReplicationGroupId < nodes[j].ReplicationGroupId
	// })

	if err := saws.PrintNodes(os.Stdout, clusters); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
