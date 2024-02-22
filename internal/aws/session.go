package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// GetSession returns aws.Config structure.
// The received structure is passed to `NewFromConfig` function of each AWS service.
func GetSession(profile string, region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}
