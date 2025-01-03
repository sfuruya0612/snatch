package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/spf13/cobra"
)

var ssoCmd = &cobra.Command{
	Use:   "sso",
	Short: "Manage SSO",
}

var ssoLoginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"l"},
	Short:   "Login to SSO",
	Run:     login,
}

// func login(cmd *cobra.Command, args []string) {
// 	input := aws.GenerateGetCallerIdentityInput()

// 	client := aws.NewSTSClient(cmd.Flag("profile").Value.String(), cmd.Flag("region").Value.String())
// 	output, err := aws.Login(client, input)
// 	if err != nil {
// 		fmt.Printf("failed to login: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("Account: %s\n", *output.Account)
// 	fmt.Printf("Arn: %s\n", *output.Arn)
// 	fmt.Printf("UserId: %s\n", *output.UserId)

// 	return
// }

func login(cmd *cobra.Command, args []string) {
	profile := cmd.Flag("profile").Value.String()
	region := cmd.Flag("region").Value.String()

	oidcClient := aws.NewSSOOidcClient(profile, region)

	ssoOidcOpts := aws.SSOOidcOpts{
		ClientName: "snatch",
		ClientType: "public",
	}

	registerOutput, err := aws.RegisterClient(oidcClient, aws.GenerateRegisterClientInput(ssoOidcOpts))
	if err != nil {
		fmt.Printf("failed to register client: %v\n", err)
		return
	}

	startUrl := fmt.Sprintf("https://device.sso.%s.amazonaws.com/", region)

	ssoOidcOpts.ClientId = *registerOutput.ClientId
	ssoOidcOpts.ClientSecret = *registerOutput.ClientSecret
	ssoOidcOpts.StartUrl = startUrl

	deviceAuth, err := aws.StartDeviceAuthorization(oidcClient, aws.GenerateStartDeviceAuthorizationInput(ssoOidcOpts))
	if err != nil {
		fmt.Printf("failed to start device authorization: %v\n", err)
		return
	}

	fmt.Printf("Please visit: %s\n", *deviceAuth.VerificationUriComplete)
	fmt.Println("Waiting for authentication...")

	ssoOidcOpts.DeviceCode = *deviceAuth.DeviceCode
	ssoOidcOpts.GrantType = "urn:ietf:params:oauth:grant-type:device_code"

	var accessToken string
	tokenOutput, err := aws.CreateToken(oidcClient, aws.GenerateCreateTokenInput(ssoOidcOpts))
	for {
		if err != nil {
			time.Sleep(time.Minute)
			var slowDown *aws.SlowDown
			if ok := errors.As(err, &slowDown); ok {
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Printf("failed to create token: %v\n", err)
			return
		}
		accessToken = *tokenOutput.AccessToken
		break
	}

	ssoClient := aws.NewSSOClient(profile, region)

	ssoOpts := aws.SSOOpts{
		AccessToken: accessToken,
	}

	accountList, err := aws.ListAccounts(ssoClient, aws.GenerateListAccountsInput(ssoOpts))
	if err != nil {
		fmt.Printf("failed to list accounts: %v\n", err)
		return
	}

	fmt.Println("\nAvailable accounts:")
	for _, account := range accountList.AccountList {
		fmt.Printf("Account ID: %s, Name: %s\n", *account.AccountId, *account.AccountName)

		ssoOpts.AccountId = *account.AccountId

		roles, err := aws.ListAccountRoles(ssoClient, aws.GenerateListAccountRolesInput(ssoOpts))
		if err != nil {
			fmt.Printf("failed to list roles for account %s: %v\n", *account.AccountId, err)
			return
		}

		fmt.Println("Available roles:")
		for _, role := range roles.RoleList {
			fmt.Printf("  - %s\n", *role.RoleName)
		}
		fmt.Println()
	}
}
