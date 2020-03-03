package main

import (
	"github.com/sfuruya0612/snatch/cmd"
	"github.com/urfave/cli"
)

// Commands cli.command object
var Commands = []cli.Command{
	commandEc2,
	commandRds,
	commandElastiCache,
	commandElb,
	commandElbV2,
	commandRoute53,
	commandAcm,
	commandS3,
	commandSsm,
	commandLogs,
	commandCloudFormation,
	commandAutoScaling,
	commandIam,
	commandTranslate,
	commandCostExplorer,
	commandMetrics,
}

var commandEc2 = cli.Command{
	Name:      "ec2",
	Usage:     "Get a list of EC2 resources",
	ArgsUsage: "[ --tag | -t ] <Key:Value> [ --short | -s ]",
	Action:    cmd.GetEc2List,
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
			Name:      "log",
			Usage:     "Get the console output for the specified instance",
			ArgsUsage: "[ --id | -i ] <InstanceId>",
			Action:    cmd.GetEc2SystemLog,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id, i",
					Usage: "Set EC2 instance id",
				},
			},
		},
		{
			Name:      "terminate",
			Usage:     "Terminate instance",
			ArgsUsage: "[ --id | -i ] <InstanceId>",
			Action:    cmd.TerminateEc2,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id, i",
					Usage: "Set EC2 instance id",
				},
			},
		},
	},
}

var commandRds = cli.Command{
	Name:   "rds",
	Usage:  "Get a list of RDS resources",
	Action: cmd.GetRdsList,
}

var commandElastiCache = cli.Command{
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
}

var commandElb = cli.Command{
	Name:   "elb",
	Usage:  "Get a list of ELB(Classic) resources.",
	Action: cmd.GetElbList,
}

var commandElbV2 = cli.Command{
	Name:   "elbv2",
	Usage:  "Get a list of ELB(Application & Network) resources",
	Action: cmd.GetElbV2List,
}

var commandRoute53 = cli.Command{
	Name:   "route53",
	Usage:  "Get a list of Rotue53 Record resources",
	Action: cmd.GetRecordsList,
}

var commandAcm = cli.Command{
	Name:   "acm",
	Usage:  "Get a list of ACM resources",
	Action: cmd.GetCertificatesList,
}

var commandS3 = cli.Command{
	Name:   "s3",
	Usage:  "Get a list of S3 Buckets",
	Action: cmd.GetBucketList,
	Subcommands: []cli.Command{
		{
			Name:      "object",
			Usage:     "Get S3 object list",
			ArgsUsage: "[ --bucket | -b ] <BucketName>",
			Action:    cmd.GetObjectList,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bucket, b",
					Usage: "Set bucket name",
				},
			},
		},
		{
			Name:      "cat",
			Usage:     "Desplay S3 object file",
			ArgsUsage: "[ --bucket | -b ] <BucketName> [ --key | -k ] <ObjectKey>",
			Action:    cmd.CatObject,
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
			Name:      "download",
			Usage:     "Download S3 object file",
			ArgsUsage: "[ --bucket | -b ] <BucketName> [ --key | -k ] <ObjectKey>",
			Action:    cmd.DownloadObject,
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
}

var commandSsm = cli.Command{
	Name:   "ssm",
	Usage:  "Start a session on your instances by launching bash or shell terminal",
	Action: cmd.StartSession,
	Subcommands: []cli.Command{
		{
			Name:      "history",
			Aliases:   []string{"hist"},
			Usage:     "Get session history",
			ArgsUsage: "[ --active | -a ]",
			Action:    cmd.GetSsmHist,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "active, a",
					Usage: "Get Active history",
				},
			},
		},
		{
			Name:      "command",
			Aliases:   []string{"cmd"},
			Usage:     "Runs commands to target instances",
			ArgsUsage: "[ --tag | -t ] <Key:Value> [ --id | -i ] <InstanceId> [ --file | -f ] <ScriptFile>",
			Action:    cmd.SendCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag, t",
					Usage: "Set Key-Value of the tag (e.g. -t Name:test-ec2)",
				},
				cli.StringFlag{
					Name:  "id, i",
					Usage: "Set EC2 instance id",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Set execute file",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:      "log",
					Usage:     "Get send command log",
					ArgsUsage: "[ --active | -a ]",
					Action:    cmd.GetCmdLog,
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
}

var commandLogs = cli.Command{
	Name:  "logs",
	Usage: "Display messages for selected cloudwatch log groups and streams at interactive prompt",
	// ArgsUsage: "[ --follow | -f ]",
	Action: cmd.GetCloudWatchLogs,
	// Flags: []cli.Flag{
	// 	cli.BoolFlag{
	// 		Name:  "follow, f",
	// 		Usage: "Like `tail -f`.",
	// 	},
	// },
}

var commandCloudFormation = cli.Command{
	Name:    "cloudformation",
	Aliases: []string{"cfn"},
	Usage:   "Get a list of stacks",
	Action:  cmd.GetStacksList,
	Subcommands: []cli.Command{
		{
			Name:      "events",
			Usage:     "Get stack events",
			ArgsUsage: "[ --name | -n ] <CfnStackName>",
			Action:    cmd.GetStackEvents,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Set stack name",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete stack",
			ArgsUsage: "[ --name | -n ] <CfnStackName>",
			Action:    cmd.DeleteStack,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Set stack name",
				},
			},
		},
	},
}

var commandAutoScaling = cli.Command{
	Name:    "autoscaling",
	Aliases: []string{"as"},
	Usage:   "Get a list of EC2 AutoScalingGroups",
	Action:  cmd.GetASGList,
	Subcommands: []cli.Command{
		{
			Name:      "capacity",
			Aliases:   []string{"cap"},
			Usage:     "Update autoscaling group capacity",
			ArgsUsage: "[ --name | -n ] <AutoScalingGroupName> [ --desired ] <CapacityNum> [ --min ]  <CapacityNum> [ --max ] <CapacityNum> ",
			Action:    cmd.UpdateCapacity,
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
}

var commandIam = cli.Command{
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
}

var commandTranslate = cli.Command{
	Name:      "translate",
	Usage:     "Translate [ JP -> EN ] or [ EN -> JP ]",
	ArgsUsage: "[ --text | -t ] <Text>",
	Action:    cmd.Translate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "text, t",
			Usage: "Set translate text (e.g. -t \"Hello world\")",
		},
	},
}

var commandCostExplorer = cli.Command{
	Name:      "costexplorer",
	Aliases:   []string{"ce"},
	Usage:     "Get monthly using cost",
	ArgsUsage: "[ --start | -s ] <StartDate> [ --end | -e ] <EndDate>",
	Action:    cmd.GetCost,
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
}

var commandMetrics = cli.Command{
	Name:      "metrics",
	Usage:     "Get Cloudwatch metrics",
	ArgsUsage: "[ --service | -s ] <AWSService>",
	Action:    cmd.GetMetrics,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service, s",
			Usage: "Set aws services",
		},
	},
}
