package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getSession(profile string, region string) *session.Session {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		Config:            aws.Config{Region: aws.String(region)},
	}

	return session.Must(session.NewSessionWithOptions(opts))
}
