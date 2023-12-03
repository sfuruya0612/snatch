package aws

import (
	"context"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
)

// GetSessionV2 returns aws.Config structure.
// The received structure is passed to `NewFromConfig` function of each AWS service.
func GetSessionV2(profile string, region string) awsv2.Config {
	cfg, err := configv2.LoadDefaultConfig(context.TODO(),
		configv2.WithSharedConfigProfile(profile),
		configv2.WithRegion(region),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}
