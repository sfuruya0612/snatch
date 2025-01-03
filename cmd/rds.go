package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
)

var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Manage RDS",
}

var rdsInstanceCmd = &cobra.Command{
	Use:   "ls",
	Short: "List RDS instances",
	Run:   listRDSInstances,
}

var rdsClusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"c"},
	Short:   "List RDS clusters",
	Run:     listRDSClusters,
}

var rdsInstanceColumns = []util.Column{
	{Header: "Name", Width: 40},
	{Header: "DBInstanceClass", Width: 15},
	{Header: "Engine", Width: 16},
	{Header: "EngineVersion", Width: 22},
	{Header: "Storage", Width: 8},
	{Header: "StorageType", Width: 10},
	{Header: "DBInstanceStatus", Width: 15},
}

var rdsClusterColumns = []util.Column{
	{Header: "Name", Width: 65},
	{Header: "Engine", Width: 16},
	{Header: "EngineVersion", Width: 22},
	{Header: "EngineMode", Width: 12},
	{Header: "Status", Width: 15},
}

func listRDSInstances(cmd *cobra.Command, args []string) {
	client := aws.NewRdsClient(cmd.Flag("profile").Value.String(), cmd.Flag("region").Value.String())

	opts := &aws.RdsOpts{
		// Tag: tagFilter,
	}

	input, err := aws.GenerateDescribeDBInstancesInput(opts)
	if err != nil {
		fmt.Printf("Unable to generate describe DB instances input: %v\n", err)
		return
	}

	instances, err := aws.DescribeDBInstances(client, input)
	if err != nil {
		fmt.Printf("Unable to describe DB instances: %v\n", err)
		return
	}

	if len(instances) == 0 {
		fmt.Println("No DB instances found")
		return
	}

	formatter := util.NewTableFormatter(rdsInstanceColumns, cmd.Flag("output").Value.String())
	formatter.PrintHeader()
	formatter.PrintRows(instances)
}

func listRDSClusters(cmd *cobra.Command, args []string) {
	client := aws.NewRdsClient(cmd.Flag("profile").Value.String(), cmd.Flag("region").Value.String())

	input := aws.GenerateDescribeDBClustersInput(&aws.RdsOpts{})

	clusters, err := aws.DescribeDBClusters(client, input)
	if err != nil {
		fmt.Printf("Unable to describe DB clusters: %v\n", err)
		return
	}

	if len(clusters) == 0 {
		fmt.Println("No DB clusters found")
		return
	}

	formatter := util.NewTableFormatter(rdsClusterColumns, cmd.Flag("output").Value.String())
	formatter.PrintHeader()
	formatter.PrintRows(clusters)
}
