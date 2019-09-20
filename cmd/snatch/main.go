package main

import (
	"fmt"
	"os"

	"github.com/sfuruya0612/snatch/cmd/snatch/command"
	"github.com/urfave/cli"
)

var (
	date      string
	hash      string
	goversion string
)

func main() {
	snatch := New(date, hash, goversion)
	if err := snatch.Run(os.Args); err != nil {
		fmt.Printf("\n[ERROR]: %v\n", err)
		os.Exit(1)
	}
}

// New returns cli.App
func New(date, hash, goversion string) *cli.App {
	app := cli.NewApp()

	app.Name = "snatch"
	app.Usage = "Show AWS resources cli command."
	app.Version = fmt.Sprintf("%s %s (Build by: %s)", date, hash, goversion)
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile, p",
			Value: "default",
			Usage: "Choose AWS credential.",
		},
		cli.StringFlag{
			Name:  "region, r",
			Value: "ap-northeast-1",
			Usage: "Select Region.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "ec2",
			Usage:  "Show EC2 resources. (default: Describe EC2 instances)",
			Action: command.ListEc2,
		},
		{
			Name:   "rds",
			Usage:  "Show RDS resources. (default: Describe RDS instances)",
			Action: command.ListRds,
		},
		{
			Name:   "ec",
			Usage:  "Show ElastiCache resources. (default: Describe Cache Clusters)",
			Action: command.ListElasticache,
			Subcommands: []cli.Command{
				{
					Name:   "rg",
					Usage:  "Describe Replication Groups.",
					Action: command.ListReplicationGroups,
				},
			},
		},
		{
			Name:   "elb",
			Usage:  "Show Elastic Load Balancer(Classic) resources. (default: Describe Load Balancers)",
			Action: command.ListElb,
		},
		{
			Name:   "elbv2",
			Usage:  "Show Elastic Load Balancer(Application & Network) resources. (default: Describe Load Balancers)",
			Action: command.ListElbv2,
		},
		{
			Name:   "route53",
			Usage:  "Show Rotue53 resources. (default: List hosted zones)",
			Action: command.ListHostedZones,
		},
	}
	return app
}
