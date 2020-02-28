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

	app.Commands = Commands

	if err := app.Run(os.Args); err != nil {
		code := 1
		if c, ok := err.(cli.ExitCoder); ok {
			code = c.ExitCode()
		}
		fmt.Printf("\n\x1b[31mERROR: %v\x1b[0m", err.Error())
		os.Exit(code)
	}
}
