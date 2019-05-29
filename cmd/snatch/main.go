package main

import (
	"fmt"
	"os"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

var (
	date      string
	hash      string
	goversion string
)

func main() {
	snatch := Exec(date, hash, goversion)
	if err := snatch.Run(os.Args); err != nil {
		fmt.Printf("\n[ERROR] %v\n", err)
		os.Exit(1)
	}
}

func Exec(date, hash, goversion string) *cli.App {
	app := cli.NewApp()

	app.Name = "snatch"
	app.Usage = "Show AWS resources cli command."
	app.Version = fmt.Sprintf("%s %s (%s)", date, hash, goversion)
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:   "ec2",
			Usage:  "Get EC2 list",
			Action: aws.DescribeEc2,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "profile, p",
					Value: "default",
					Usage: "Choose AWS credential.",
				},
				cli.StringFlag{
					Name:  "region",
					Value: "ap-northeast-1",
					Usage: "Select Region.",
				},
				// extra, -e フラグ追加予定
			},
		},
		{
			Name:   "rds",
			Usage:  "Get RDS list",
			Action: aws.DescribeRds,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "profile, p",
					Value: "default",
					Usage: "Choose AWS credential.",
				},
				cli.StringFlag{
					Name:  "region",
					Value: "ap-northeast-1",
					Usage: "Select Region.",
				},
			},
		},
	}
	return app
}
