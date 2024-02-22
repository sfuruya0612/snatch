package aws

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAM client struct
type IAM struct {
	Client *iam.Client
}

// NewIamSess return IAM struct initialized
func NewIamClient(profile, region string) *IAM {
	return &IAM{
		Client: iam.NewFromConfig(GetSession(profile, region)),
	}
}

// User iam user struct
type User struct {
	Name          string
	ManagedPolicy string
	InlinePolicy  string
	Group         string
	AccessKey     string
	AccessKeyUsed string
	PWLastUsed    string
	CreateDate    string
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
	output, err := c.Client.ListUsers(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("list users: %v", err)
	}

	list := Users{}
	for _, o := range output.Users {
		used := "None"
		if o.PasswordLastUsed != nil {
			used = o.PasswordLastUsed.String()
		}

		username := *o.UserName

		mi := &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(username),
		}

		managed, err := c.listAttachedUserPolicies(mi)
		if err != nil {
			return nil, fmt.Errorf("list attached user policies: %v", err)
		}

		ii := &iam.ListUserPoliciesInput{
			UserName: aws.String(username),
		}

		inline, err := c.listUserPolicies(ii)
		if err != nil {
			return nil, fmt.Errorf("list user policies: %v", err)
		}

		gi := &iam.ListGroupsForUserInput{
			UserName: aws.String(username),
		}

		group, err := c.listGroupsForUser(gi)
		if err != nil {
			return nil, fmt.Errorf("list groups for user: %v", err)
		}

		ai := &iam.ListAccessKeysInput{
			UserName: aws.String(username),
		}

		key, err := c.listAccessKeys(ai)
		if err != nil {
			return nil, fmt.Errorf("list access keys: %v", err)
		}

		list = append(list, User{
			Name:          username,
			ManagedPolicy: managed,
			InlinePolicy:  inline,
			Group:         group,
			AccessKey:     key,
			PWLastUsed:    used,
			CreateDate:    o.CreateDate.String(),
		})
	}

	return list, nil
}

func (c *IAM) listAttachedUserPolicies(input *iam.ListAttachedUserPoliciesInput) (string, error) {
	output, err := c.Client.ListAttachedUserPolicies(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("list attached user policies: %v", err)
	}

	var (
		policies []string
		policy   string
	)
	for _, p := range output.AttachedPolicies {
		policies = append(policies, *p.PolicyName)
	}
	policy = strings.Join(policies[:], ",")

	if len(policy) == 0 {
		policy = "None"
	}

	return policy, nil
}

func (c *IAM) listUserPolicies(input *iam.ListUserPoliciesInput) (string, error) {
	output, err := c.Client.ListUserPolicies(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("list user policies: %v", err)
	}

	var (
		policies []string
		policy   string
	)

	policies = append(policies, output.PolicyNames...)
	policy = strings.Join(policies[:], ",")

	if len(policy) == 0 {
		policy = "None"
	}

	return policy, nil
}

func (c *IAM) listGroupsForUser(input *iam.ListGroupsForUserInput) (string, error) {
	output, err := c.Client.ListGroupsForUser(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("list groups for user: %v", err)
	}

	var (
		groups []string
		group  string
	)
	for _, g := range output.Groups {
		groups = append(groups, *g.GroupName)
	}
	group = strings.Join(groups[:], ",")

	if len(group) == 0 {
		group = "None"
	}

	return group, nil
}

func (c *IAM) listAccessKeys(input *iam.ListAccessKeysInput) (string, error) {
	output, err := c.Client.ListAccessKeys(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("list access keys: %v", err)
	}

	var (
		keys []string
		key  string
	)
	for _, k := range output.AccessKeyMetadata {
		keys = append(keys, *k.AccessKeyId)
	}
	key = strings.Join(keys[:], ",")

	if len(key) == 0 {
		key = "None"
	}

	return key, nil
}

// ListRoles return []string (iam.ListRolesOutput.RoleName)
// input iam.ListRolesInput
func (c *IAM) ListRoles(input *iam.ListRolesInput) ([]string, error) {
	output, err := c.Client.ListRoles(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("list roles: %v", err)
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

		output, err := c.Client.GetRole(context.TODO(), input)
		if err != nil {
			return nil, fmt.Errorf("get role: %v", err)
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
		"ManagedPolicy",
		"InlinePolicy",
		"Group",
		"AccessKey",
		"PWLastUsed",
		"CreateDate",
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
		i.ManagedPolicy,
		i.InlinePolicy,
		i.Group,
		i.AccessKey,
		i.PWLastUsed,
		i.CreateDate,
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
