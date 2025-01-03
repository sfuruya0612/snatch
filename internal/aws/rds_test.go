package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type mockRdsClient struct {
	// rdsApi
}

func (m *mockRdsClient) DescribeDBInstances(ctx context.Context, params *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	return &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test-instance"),
				DBInstanceClass:      aws.String("db.t2.micro"),
				Engine:               aws.String("mysql"),
				EngineVersion:        aws.String("5.7.22"),
				AllocatedStorage:     aws.Int32(20),
				StorageType:          aws.String("gp2"),
				DBInstanceStatus:     aws.String("available"),
			},
		},
	}, nil
}

func (m *mockRdsClient) DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	return &rds.DescribeDBClustersOutput{
		DBClusters: []types.DBCluster{
			{
				DBClusterIdentifier: aws.String("test-cluster"),
				Engine:              aws.String("aurora"),
				EngineVersion:       aws.String("5.6.10a"),
				EngineMode:          aws.String("provisioned"),
				Status:              aws.String("available"),
			},
		},
	}, nil
}

func TestDescribeDBInstances(t *testing.T) {
	mockClient := &mockRdsClient{}
	input := &rds.DescribeDBInstancesInput{}

	result, err := DescribeDBInstances(mockClient, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := [][]string{
		{
			"test-instance",
			"db.t2.micro",
			"mysql",
			"5.7.22",
			"20GB",
			"gp2",
			"available",
		},
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d instances, got %d", len(expected), len(result))
	}

	for i, instance := range result {
		for j, field := range instance {
			if field != expected[i][j] {
				t.Errorf("expected %s, got %s", expected[i][j], field)
			}
		}
	}
}

func TestDescribeDBClusters(t *testing.T) {
	mockClient := &mockRdsClient{}
	input := &rds.DescribeDBClustersInput{}

	result, err := DescribeDBClusters(mockClient, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := [][]string{
		{
			"test-cluster",
			"aurora",
			"5.6.10a",
			"provisioned",
			"available",
		},
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d clusters, got %d", len(expected), len(result))
	}

	for i, cluster := range result {
		for j, field := range cluster {
			if field != expected[i][j] {
				t.Errorf("expected %s, got %s", expected[i][j], field)
			}
		}
	}
}
