package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sfuruya0612/snatch/cmd"
	"github.com/urfave/cli/v2"
)

var (
	commit string
)

var Commands = []*cli.Command{
	cmd.Ec2,
	cmd.Rds,
	cmd.ElastiCache,
	// TODO: elb は一緒にしたい
	cmd.Elb,
	cmd.ElbV2,
	cmd.Route53,
	cmd.S3,
	cmd.Ssm,
	cmd.CloudFormation,
	cmd.AutoScaling,
	cmd.Iam,
}

func main() {
	app := cli.NewApp()

	app.EnableBashCompletion = true
	app.Name = "snatch"
	app.Usage = "CLI tool to get AWS resources"
	app.Version = fmt.Sprintf("%s (Build by: %s)", commit, runtime.Version())

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Aliases: []string{"p"},
			EnvVars: []string{"AWS_PROFILE"},
			Value:   "default",
			Usage:   "AWS credential (~/.aws/config) or read AWS_PROFILE environment variable",
		},
		&cli.StringFlag{
			Name:    "region",
			Aliases: []string{"r"},
			EnvVars: []string{"AWS_REGION"},
			Value:   "ap-northeast-1",
			Usage:   "Specify a valid AWS region",
		},
	}

	app.Before = cmd.Before

	app.Commands = Commands

	if err := app.Run(os.Args); err != nil {
		code := 1
		if c, ok := err.(cli.ExitCoder); ok {
			code = c.ExitCode()
		}
		fmt.Printf("\x1b[31mERROR: %v\x1b[0m", err.Error())
		os.Exit(code)
	}
}
