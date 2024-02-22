package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECS interface {
	DescribeClusters(ctx context.Context, input *ecs.DescribeClustersInput, opts ...func(*ecs.Options)) (*ecs.DescribeClustersOutput, error)
	ListClusters(ctx context.Context, input *ecs.ListClustersInput, opts ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
	ListServices(ctx context.Context, input *ecs.ListServicesInput, opts ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
}

func NewECSClient(profile, region string) ECS {
	return ecs.NewFromConfig(GetSession(profile, region))
}

type Cluster struct {
	Name              string
	Status            string
	ContainerInsights string
}

func GetClusters(api ECS) ([]Cluster, error) {
	list, err := listClusters(api)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	input := &ecs.DescribeClustersInput{
		Clusters: list,
	}

	output, err := api.DescribeClusters(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	clusters := []Cluster{}
	for _, c := range output.Clusters {
		// ci := *c.Settings.

		clusters = append(clusters, Cluster{
			Name:   *c.ClusterName,
			Status: *c.Status,
			// ContainerInsights: ci,
		})
	}

	return clusters, nil
}

func listClusters(api ECS) ([]string, error) {
	input := &ecs.ListClustersInput{}
	output, err := api.ListClusters(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var clusters []string
	for _, cluster := range output.ClusterArns {
		split := strings.Split(cluster, "/")
		clusterName := split[len(split)-1]

		clusters = append(clusters, clusterName)
	}
	return clusters, nil
}

func GetServices(api ECS, cluster string) ([]string, error) {
	input := &ecs.ListServicesInput{
		Cluster: &cluster,
	}

	output, err := api.ListServices(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var services []string
	for _, s := range output.ServiceArns {
		split := strings.Split(s, "/")
		serviceName := split[len(split)-1]

		services = append(services, serviceName)
	}

	return services, nil
}
