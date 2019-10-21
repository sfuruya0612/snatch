package main

import (
	"fmt"
	"os"

	"github.com/sfuruya0612/snatch/cmd"
	"github.com/urfave/cli"
)

const version = "19.10.1"

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

	app.EnableBashCompletion = true

	app.Name = "snatch"
	app.Usage = "Cli command to get and display Amazon Web Services resources."

	if date != "" || hash != "" || goversion != "" {
		app.Version = fmt.Sprintf("%s %s (Build by: %s)", date, hash, goversion)
	} else {
		app.Version = version
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "profile, p",
			EnvVar: "AWS_PROFILE",
			Value:  "default",
			Usage:  "Specify the AWS profile listed in ~/.aws/config.",
		},
		cli.StringFlag{
			Name:  "region, r",
			Value: "ap-northeast-1",
			Usage: "Specify the AWS region.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "ec2",
			Usage:  "Get a list of EC2 resources. (API: DescribeInstances)",
			Action: cmd.ListEc2,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag, t",
					Usage: "The Key-Value of the tag to filter. (e.g. -t Name:test-ec2)",
				},
			},
		},
		{
			Name:   "rds",
			Usage:  "Get a list of RDS resources. (API: DescribeDbInstances)",
			Action: cmd.ListRds,
		},
		{
			Name:   "ec",
			Usage:  "Get a list of ElastiCache Cluster resources. (API: DescribeCacheClusters)",
			Action: cmd.ListElasticache,
			Subcommands: []cli.Command{
				{
					Name:   "rg",
					Usage:  "Get a list of ElastiCache Node resources. (API: DescribeReplicationGroups)",
					Action: cmd.ListReplicationGroups,
				},
			},
		},
		{
			Name:   "elb",
			Usage:  "Get a list of ELB(Classic) resources. (API: DescribeLoadBalancers)",
			Action: cmd.ListElb,
		},
		{
			Name:   "elbv2",
			Usage:  "Get a list of ELB(Application & Network) resources. (API: DescribeLoadBalancers)",
			Action: cmd.ListElbv2,
		},
		{
			Name:   "route53",
			Usage:  "Get a list of Rotue53 Record resources. (API: ListHostedZones and ListResourceRecordSets)",
			Action: cmd.ListHostedZones,
		},
		{
			Name:   "acm",
			Usage:  "Get a list of ACM resources. (API: ListCertificates and DescribeCertificate)",
			Action: cmd.ListCertificates,
		},
		{
			Name:   "s3",
			Usage:  "Get Objects in selected S3 Bucket at interactive prompt. (API: ListBuckets and ListObjects)",
			Action: cmd.ListBuckets,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "l",
					Usage: "Get Objects list.",
				},
			},
		},
		{
			Name:   "ssm",
			Usage:  "Start a session on your instances by launching bash or shell terminal. (API: StartSession)",
			Action: cmd.StartSession,
			Subcommands: []cli.Command{
				{
					Name:   "run",
					Usage:  "Runs commands on one target instance. (API: SendCommand)",
					Action: cmd.SendCommand,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Usage: "Set execute file.",
						},
						cli.StringFlag{
							Name:  "tag, t",
							Usage: "Set Key-Value of the tag. (e.g. -t Name:test-ec2)",
						},
						cli.StringFlag{
							Name:  "instanceid, i",
							Usage: "Set EC2 instance id.",
						},
					},
				},
			},
		},
		{
			Name:   "logs",
			Usage:  "Display messages for selected log groups and streams at interactive prompt. (API: )",
			Action: cmd.DescribeLogGroups,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f",
					Usage: "Like `tail -f`.",
				},
			},
		},
	}

	return app
}
