package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

func GetUserList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

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

func GetRoleList(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

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
