package cmd

import (
	"fmt"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Ecs = &cli.Command{
	Name:  "ecs",
	Usage: "Get a list of ECS",
	Subcommands: []*cli.Command{
		{
			Name:  "clusters",
			Usage: "Get a list of ECS clusters",
			Action: func(c *cli.Context) error {
				return getClusters(c.String("profile"), c.String("region"))
			},
		},
		{
			Name:  "services",
			Usage: "Get a list of ECS services",
			Action: func(c *cli.Context) error {
				return getServices(c.String("profile"), c.String("region"))
			},
		},
	},
}

func getClusters(profile, region string) error {
	c := saws.NewECSClient(profile, region)
	ecs, err := saws.GetClusters(c)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Println(ecs)

	return nil
}

func getServices(profile, region string) error {
	c := saws.NewECSClient(profile, region)

	clusters, err := saws.GetClusters(c)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, cluster := range clusters {
		services, err := saws.GetServices(c, cluster.Name)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		fmt.Println(services)
	}

	return nil
}
