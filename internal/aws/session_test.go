package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
)

func TestGetSession(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		region  string
		wantErr bool
	}{
		{
			name:    "Valid profile and region",
			profile: "default",
			region:  "us-west-2",
			wantErr: false,
		},
		{
			name:    "Invalid profile",
			profile: "non-existent-profile",
			region:  "us-west-2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedConfigProfile(tt.profile),
				config.WithRegion(tt.region),
			)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}
