package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getSession(profile string, region string) *session.Session {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		Config:            aws.Config{Region: aws.String(region)},
	}

	fmt.Printf("AWS_PROFILE: %s , REGION: %s\n", profile, region)

	return session.Must(session.NewSessionWithOptions(opts))
}
