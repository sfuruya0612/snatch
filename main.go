package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sfuruya0612/snatch/cmd"
	"github.com/urfave/cli"
)

var (
	date      string
	hash      string
	goversion string
)

func main() {
	app := cli.NewApp()

	app.EnableBashCompletion = true
	app.Name = "snatch"
	app.Usage = "CLI tool to get AWS resources"

	if date != "" || hash != "" || goversion != "" {
		app.Version = fmt.Sprintf("%s %s (Build by: %s)", date, hash, goversion)
	} else {
		v, err := ioutil.ReadFile("VERSION")
		if err != nil {
			fmt.Println("doesn’t read the VERSION file")
		}
		app.Version = string(v)
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
