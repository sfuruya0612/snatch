package main

import (
	"fmt"
	"os"

	"github.com/ShoichiFuruya/snatch/cmd/services"

    "github.com/urfave/cli"
)

var (
	date      string
	goversion string
)

func main() {
    //snatch := Exec(date, goversion)
    snatch := Exec()
	if err := snatch.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//func Exec(date, goversion) *cli.App {
func Exec() *cli.App {
	app := cli.NewApp()

    app.Name = "snatch"
	app.Usage = "Show AWS resources."
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile, p",
			Value: "PROFILE",
		},
	}

    app.Commands =[]cli.Command{
        {
            Name:   "ec2",
	        Usage:  "Get EC2 list",
		    Action: func(c *cli.Context) error {
                services.DescribeEc2()
                return nil
            },
        },
	    {
            Name:   "rds",
	        Usage:  "Get RDS list",
		    Action: func(command *cli.Context) error {
                //command.DescribeRds()
                return nil
            },
        },
	}
    return app
}
