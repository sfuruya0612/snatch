package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
)

type SSOOpts struct {
	AccessToken string
	AccountId   string
}

// type stsApi interface {
// 	GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, opts ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
// }

// func NewSTSClient(profile, region string) stsApi {
// 	return sts.NewFromConfig(GetSession(profile, region))
// }

type ssoApi interface {
	ListAccounts(ctx context.Context, input *sso.ListAccountsInput, opts ...func(*sso.Options)) (*sso.ListAccountsOutput, error)
	ListAccountRoles(ctx context.Context, input *sso.ListAccountRolesInput, opts ...func(*sso.Options)) (*sso.ListAccountRolesOutput, error)
}

func NewSSOClient(profile, region string) ssoApi {
	return sso.NewFromConfig(GetSession(profile, region))
}

// func GenerateGetCallerIdentityInput() *sts.GetCallerIdentityInput {
// 	return &sts.GetCallerIdentityInput{}
// }

// func Login(api stsApi, input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
// 	o, err := api.GetCallerIdentity(context.Background(), input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return o, nil
// }

func GenerateListAccountsInput(opts SSOOpts) *sso.ListAccountsInput {
	return &sso.ListAccountsInput{
		AccessToken: aws.String(opts.AccessToken),
	}
}

func ListAccounts(api ssoApi, input *sso.ListAccountsInput) (*sso.ListAccountsOutput, error) {
	o, err := api.ListAccounts(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func GenerateListAccountRolesInput(opts SSOOpts) *sso.ListAccountRolesInput {
	return &sso.ListAccountRolesInput{
		AccessToken: aws.String(opts.AccessToken),
		AccountId:   aws.String(opts.AccountId),
	}
}

func ListAccountRoles(api ssoApi, input *sso.ListAccountRolesInput) (*sso.ListAccountRolesOutput, error) {
	o, err := api.ListAccountRoles(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return o, nil
}
