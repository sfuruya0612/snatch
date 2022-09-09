package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli/v2"
)

var Iam = &cli.Command{
	Name:  "iam",
	Usage: "Get a list of IAM users",
	Action: func(c *cli.Context) error {
		return getUserList(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:  "role",
			Usage: "Get a list of IAM role",
			Action: func(c *cli.Context) error {
				return getRoleList(c.String("profile"), c.String("region"))
			},
		},
	},
}

func getUserList(profile, region string) error {
	client := saws.NewIamSess(profile, region)

	output, err := client.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintUsers(os.Stdout, output); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}

func getRoleList(profile, region string) error {
	client := saws.NewIamSess(profile, region)

	names, err := client.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	output, err := client.GetRole(names)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintRoles(os.Stdout, output); err != nil {
		return fmt.Errorf("failed to print resources")
	}

	return nil
}
