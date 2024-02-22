package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type EcsMockAPI struct {
	Output *ecs.DescribeClustersOutput
	Error  error
}

func (m *EcsMockAPI) DescribeClusters(ctx context.Context, input *ecs.DescribeClustersInput, opts ...func(*ecs.Options)) (*ecs.DescribeClustersOutput, error) {
	return m.Output, m.Error
}

func TestGetClusters(t *testing.T) {
	var clusterName string = "test_cluster"

	mock := &EcsMockAPI{
		Output: &ecs.DescribeClustersOutput{
			Clusters: []types.Cluster{
				{
					ClusterName: aws.String(clusterName),
				},
			},
		},
		Error: nil,
	}

	clusters, err := GetClusters(mock)
	if err != nil {
		t.Errorf("Error should be nil, but got %v", err)
	}

	if len(clusters) != 1 {
		t.Errorf("Clusters length should be 1, but got %v", len(clusters))
	}

	if clusters[0] != "test_cluster" {
		t.Errorf("ClusterName should be test, but got %v", clusters[0])
	}
}
