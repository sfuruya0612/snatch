package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSsmClient struct {
	mock.Mock
}

func (m *mockSsmClient) DescribeInstanceInformation(ctx context.Context, params *ssm.DescribeInstanceInformationInput, optFns ...func(*ssm.Options)) (*ssm.DescribeInstanceInformationOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.DescribeInstanceInformationOutput), args.Error(1)
}

func (m *mockSsmClient) StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.StartSessionOutput), args.Error(1)
}

func (m *mockSsmClient) TerminateSession(ctx context.Context, params *ssm.TerminateSessionInput, optFns ...func(*ssm.Options)) (*ssm.TerminateSessionOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.TerminateSessionOutput), args.Error(1)
}

func TestGenerateDescribeInstanceInformationInput(t *testing.T) {
	opts := &SsmOpts{
		PingStatus:   "Online",
		ResourceType: "EC2Instance",
	}

	input := GenerateDescribeInstanceInformationInput(opts)

	assert.Equal(t, 2, len(input.Filters))
	assert.Equal(t, "PingStatus", *input.Filters[0].Key)
	assert.Equal(t, []string{"Online"}, input.Filters[0].Values)
	assert.Equal(t, "ResourceType", *input.Filters[1].Key)
	assert.Equal(t, []string{"EC2Instance"}, input.Filters[1].Values)
}

func TestDescribeInstanceInformation(t *testing.T) {
	mockClient := new(mockSsmClient)
	instanceID := "i-1234567890abcdef0"

	mockOutput := &ssm.DescribeInstanceInformationOutput{
		InstanceInformationList: []types.InstanceInformation{
			{
				InstanceId: aws.String(instanceID),
			},
		},
	}

	mockClient.On("DescribeInstanceInformation",
		context.Background(),
		mock.AnythingOfType("*ssm.DescribeInstanceInformationInput"),
		mock.Anything).
		Return(mockOutput, nil)

	input := &ssm.DescribeInstanceInformationInput{}
	ids, err := DescribeInstanceInformation(mockClient, input)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(ids))
	assert.Equal(t, instanceID, ids[0])
	mockClient.AssertExpectations(t)
}

func TestGenerateStartSessionInput(t *testing.T) {
	opts := &SsmOpts{
		InstanceId: "i-1234567890abcdef0",
	}

	input := GenerateStartSessionInput(opts)

	assert.Equal(t, opts.InstanceId, *input.Target)
}

func TestStartSession(t *testing.T) {
	mockClient := new(mockSsmClient)
	sessionID := "session-1234567890"

	mockOutput := &ssm.StartSessionOutput{
		SessionId: aws.String(sessionID),
	}

	mockClient.On("StartSession",
		context.Background(),
		mock.AnythingOfType("*ssm.StartSessionInput"),
		mock.Anything).
		Return(mockOutput, nil)

	input := &ssm.StartSessionInput{
		Target: aws.String("i-1234567890abcdef0"),
	}

	output, err := StartSession(mockClient, input)

	assert.NoError(t, err)
	assert.Equal(t, sessionID, *output.SessionId)
	mockClient.AssertExpectations(t)
}

func TestGenerateTerminateSessionInput(t *testing.T) {
	opts := &SsmOpts{
		SessionId: "session-1234567890",
	}

	input := GenerateTerminateSessionInput(opts)

	assert.Equal(t, opts.SessionId, *input.SessionId)
}

func TestTerminateSession(t *testing.T) {
	mockClient := new(mockSsmClient)
	sessionID := "session-1234567890"

	mockOutput := &ssm.TerminateSessionOutput{
		SessionId: aws.String(sessionID),
	}

	mockClient.On("TerminateSession",
		context.Background(),
		mock.AnythingOfType("*ssm.TerminateSessionInput"),
		mock.Anything).
		Return(mockOutput, nil)

	input := &ssm.TerminateSessionInput{
		SessionId: aws.String(sessionID),
	}

	output, err := TerminateSession(mockClient, input)

	assert.NoError(t, err)
	assert.Equal(t, sessionID, *output.SessionId)
	mockClient.AssertExpectations(t)
}
