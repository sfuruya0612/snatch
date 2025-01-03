package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "snatch",
	Short: "A CLI tool to interact with AWS resources",
	Long:  `Snatch is a CLI tool that helps you to manage AWS resources like EC2 instances.`,
}

// var tagFilter string

func init() {
	rootCmd.PersistentFlags().StringP("profile", "p", viper.GetString("AWS_PROFILE"), "AWS profile (default $AWS_PROFILE)")
	rootCmd.PersistentFlags().StringP("region", "r", "ap-northeast-1", "AWS region")
	rootCmd.PersistentFlags().StringP("output", "o", "tab", "Output format (tab/csv)")

	rootCmd.AddCommand(ssoCmd, ec2Cmd, rdsCmd)

	// SSO
	ssoCmd.AddCommand(ssoLoginCmd)

	// EC2
	ec2Cmd.AddCommand(ec2ListCmd, ec2StartSessionCmd)

	ec2ListCmd.Flags().BoolP("running", "", false, "Show only running instances")

	ec2StartSessionCmd.Flags().StringP("instance-id", "i", "", "Instance ID")

	// RDS
	rdsCmd.AddCommand(rdsInstanceCmd, rdsClusterCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ドキュメント生成向けにrootCmdを返す
func GetRootCmd() *cobra.Command {
	return rootCmd
}
