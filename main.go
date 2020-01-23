package main

import (
	"fmt"
	"os"

	"github.com/sfuruya0612/snatch/cmd"
	"github.com/urfave/cli"
)

const version = "19.11.1"

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
	app.Usage = "Cli command to Amazon Web Services resources"

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
			Usage:  "AWS credential (~/.aws/config) or read AWS_PROFILE environment variable",
		},
		cli.StringFlag{
			Name:  "region, r",
			Value: "ap-northeast-1",
			Usage: "Specify a valid AWS region",
		},
	}

	app.Before = cmd.Before

	app.Commands = []cli.Command{
		{
			Name:   "ec2",
			Usage:  "Get a list of EC2 resources",
			Action: cmd.GetEc2List,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag, t",
					Usage: "The Key-Value of the tag to filter",
				},
				cli.BoolFlag{
					Name:  "short, s",
					Usage: "Desplay fewer items only running instances (Name, PrivateIP, PublicIP)",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:   "log",
					Usage:  "Get the console output for the specified instance",
					Action: cmd.GetEc2SystemLog,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "instanceid, i",
							Usage: "Set EC2 instance id",
						},
					},
				},
				{
					Name:   "terminate",
					Usage:  "Terminate instance",
					Action: cmd.TerminateEc2,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "instanceid, i",
							Usage: "Set EC2 instance id",
						},
					},
				},
			},
		},
		{
			Name:   "rds",
			Usage:  "Get a list of RDS resources",
			Action: cmd.GetRdsList,
		},
		{
			Name:    "elasticache",
			Aliases: []string{"ec"},
			Usage:   "Get a list of ElastiCache Cluster resources",
			Action:  cmd.GetEcClusterList,
			Subcommands: []cli.Command{
				{
					Name:   "node",
					Usage:  "Get a list of ElastiCache Node resources",
					Action: cmd.GetEcGroupsList,
				},
			},
		},
		{
			Name:   "elb",
			Usage:  "Get a list of ELB(Classic) resources.",
			Action: cmd.GetElbList,
		},
		{
			Name:   "elbv2",
			Usage:  "Get a list of ELB(Application & Network) resources",
			Action: cmd.GetElbV2List,
		},
		{
			Name:   "route53",
			Usage:  "Get a list of Rotue53 Record resources",
			Action: cmd.GetRecordsList,
		},
		{
			Name:   "acm",
			Usage:  "Get a list of ACM resources",
			Action: cmd.GetCertificatesList,
		},
		{
			Name:   "s3",
			Usage:  "Get a list of S3 Buckets",
			Action: cmd.GetBucketList,
			Subcommands: []cli.Command{
				{
					Name:   "object",
					Usage:  "Get S3 object list",
					Action: cmd.GetObjectList,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "bucket, b",
							Usage: "Set bucket name",
						},
					},
				},
				{
					Name:   "cat",
					Usage:  "Desplay S3 object file",
					Action: cmd.CatObject,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "bucket, b",
							Usage: "Set bucket name",
						},
						cli.StringFlag{
							Name:  "key, k",
							Usage: "Set object key",
						},
					},
				},
				{
					Name:   "download",
					Usage:  "Download S3 object file",
					Action: cmd.DownloadObject,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "bucket, b",
							Usage: "Set bucket name",
						},
						cli.StringFlag{
							Name:  "key, k",
							Usage: "Set object key",
						},
					},
				},
			},
		},
		{
			Name:   "ssm",
			Usage:  "Start a session on your instances by launching bash or shell terminal",
			Action: cmd.StartSession,
			Subcommands: []cli.Command{
				{
					Name:    "history",
					Aliases: []string{"hist"},
					Usage:   "Get session history",
					Action:  cmd.GetSsmHist,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "active, a",
							Usage: "Get Active history",
						},
					},
				},
				{
					Name:    "command",
					Aliases: []string{"cmd"},
					Usage:   "Runs commands to target instances",
					Action:  cmd.SendCommand,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Usage: "Set execute file",
						},
						cli.StringFlag{
							Name:  "tag, t",
							Usage: "Set Key-Value of the tag (e.g. -t Name:test-ec2)",
						},
						cli.StringFlag{
							Name:  "instanceid, i",
							Usage: "Set EC2 instance id",
						},
					},
					Subcommands: []cli.Command{
						{
							Name:   "log",
							Usage:  "Get send command log",
							Action: cmd.GetCmdLog,
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "active, a",
									Usage: "Get Active history",
								},
							},
						},
					},
				},
				{
					Name:    "parameter",
					Aliases: []string{"param"},
					Usage:   "Get parameter store",
					Action:  cmd.GetParameter,
				},
			},
		},
		{
			Name:   "logs",
			Usage:  "Display messages for selected cloudwatch log groups and streams at interactive prompt",
			Action: cmd.GetCloudWatchLogs,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f",
					Usage: "Like `tail -f`.",
				},
			},
		},
		{
			Name:    "cloudformation",
			Aliases: []string{"cfn"},
			Usage:   "Get a list of stacks",
			Action:  cmd.GetStacksList,
			Subcommands: []cli.Command{
				{
					Name:   "events",
					Usage:  "Get stack events",
					Action: cmd.GetStackEvents,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "Set stack name",
						},
					},
				},
				{
					Name:   "delete",
					Usage:  "Delete stack",
					Action: cmd.DeleteStack,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "Set stack name",
						},
					},
				},
			},
		},
		{
			Name:    "dynamodb",
			Aliases: []string{"dynamo"},
			Usage:   "Scan item from DynamoDB table name",
			Action:  cmd.GetTablesList,
		},
		{
			Name:    "autoscaling",
			Aliases: []string{"as"},
			Usage:   "Get a list of EC2 AutoScalingGroups",
			Action:  cmd.GetASGList,
			Subcommands: []cli.Command{
				{
					Name:    "capacity",
					Aliases: []string{"cap"},
					Usage:   "Update autoscaling group capacity",
					Action:  cmd.UpdateCapacity,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "Set target autoscaling group name",
						},
						cli.Int64Flag{
							Name:  "desired",
							Usage: "Set desired capacity",
						},
						cli.Int64Flag{
							Name:  "min",
							Usage: "Set minimum capacity",
						},
						cli.Int64Flag{
							Name:  "max",
							Usage: "Set maximum capacity",
						},
					},
				},
			},
		},
		{
			Name:   "iam",
			Usage:  "Get a list of IAM users",
			Action: cmd.GetUserList,
			Subcommands: []cli.Command{
				{
					Name:   "role",
					Usage:  "Get a list of IAM role",
					Action: cmd.GetRoleList,
				},
			},
		},
		{
			Name:   "translate",
			Usage:  "Translate [ JP -> EN ] or [ EN -> JP ]",
			Action: cmd.Translate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "text, t",
					Usage: "Set translate text (e.g. -t \"Hello world\")",
				},
			},
		},
		{
			Name:    "costexplorer",
			Aliases: []string{"ce"},
			Usage:   "Get monthly using cost",
			Action:  cmd.GetCost,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "start, s",
					Usage: "Set start date (time format e.g. -s 2019-11-01)",
				},
				cli.StringFlag{
					Name:  "end, e",
					Usage: "Set end date (time format e.g. -s 2020-01-01)",
				},
			},
		},
	}

	return app
}
