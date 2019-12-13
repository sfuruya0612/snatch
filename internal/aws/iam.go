package aws

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAM client struct
type IAM struct {
	Client *iam.IAM
}

// NewIamSess return IAM struct initialized
func NewIamSess(profile, region string) *IAM {
	return &IAM{
		Client: iam.New(getSession(profile, region)),
	}
}

// User iam user struct
type User struct {
	Name       string
	PWLastUsed string
}

// Users User struct slice
type Users []User

// Role iam role struct
type Role struct {
	Name string
	Arn  string
}

// Roles Role struct slice
type Roles []Role

// ListUsers return Users
// input iam.ListUsersInput
func (c *IAM) ListUsers(input *iam.ListUsersInput) (Users, error) {
	output, err := c.Client.ListUsers(input)
	if err != nil {
		return nil, fmt.Errorf("List users: %v", err)
	}

	list := Users{}
	for _, o := range output.Users {
		list = append(list, User{
			Name:       *o.UserName,
			PWLastUsed: o.PasswordLastUsed.String(),
		})
	}

	return list, nil
}

// ListRoles return []string (iam.ListRolesOutput.RoleName)
// input iam.ListRolesInput
func (c *IAM) ListRoles(input *iam.ListRolesInput) ([]string, error) {
	output, err := c.Client.ListRoles(input)
	if err != nil {
		return nil, fmt.Errorf("List roles: %v", err)
	}

	names := []string{}
	for _, r := range output.Roles {
		names = append(names, *r.RoleName)
	}

	return names, nil
}

// GetRole return Roles
// input []string (iam.ListRolesOutput.RoleName)
func (c *IAM) GetRole(names []string) (Roles, error) {
	list := Roles{}
	for _, n := range names {
		input := &iam.GetRoleInput{
			RoleName: aws.String(n),
		}

		output, err := c.Client.GetRole(input)
		if err != nil {
			return nil, fmt.Errorf("Get role: %v", err)
		}

		list = append(list, Role{
			Name: *output.Role.RoleName,
			Arn:  *output.Role.Arn,
		})
	}

	return list, nil
}

func PrintUsers(wrt io.Writer, resources Users) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"PWLastUsed",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.UserTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *User) UserTabString() string {
	fields := []string{
		i.Name,
		i.PWLastUsed,
	}

	return strings.Join(fields, "\t")
}

func PrintRoles(wrt io.Writer, resources Roles) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"Name",
		"Arn",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.RoleTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Role) RoleTabString() string {
	fields := []string{
		i.Name,
		i.Arn,
	}

	return strings.Join(fields, "\t")
}
