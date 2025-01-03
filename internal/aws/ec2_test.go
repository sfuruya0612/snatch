package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockOutput = &ec2.DescribeInstancesOutput{
	Reservations: []types.Reservation{
		{
			Instances: []types.Instance{
				{
					InstanceId:        aws.String("i-1234567890abcdef0"),
					InstanceType:      types.InstanceTypeT2Micro,
					InstanceLifecycle: "OnDemand",
					PrivateIpAddress:  aws.String("192.168.1.1"),
					PublicIpAddress:   aws.String("203.0.113.1"),
					State:             &types.InstanceState{Name: types.InstanceStateNameRunning},
					KeyName:           aws.String("my-key"),
					Placement:         &types.Placement{AvailabilityZone: aws.String("ap-northeast-1a")},
					LaunchTime:        aws.Time(time.Now()),
					Tags: []types.Tag{
						{Key: aws.String("Name"), Value: aws.String("MyInstance")},
					},
				},
			},
		},
	},
}

type mockEc2Api struct {
	mock.Mock
}

func (m *mockEc2Api) DescribeInstances(ctx context.Context, input *ec2.DescribeInstancesInput, opts ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
}

func TestGenerateDescribeInstancesInput(t *testing.T) {
	tests := []struct {
		name    string
		opts    *Ec2Opts
		wantErr bool
	}{
		{
			name: "running instances",
			opts: &Ec2Opts{Running: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := GenerateDescribeInstancesInput(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, input)
			}
		})
	}
}

func TestDescribeInstances(t *testing.T) {
	mockApi := new(mockEc2Api)

	mockApi.On("DescribeInstances", mock.Anything, mock.Anything, mock.Anything).Return(mockOutput, nil)

	input := &ec2.DescribeInstancesInput{}
	result, err := DescribeInstances(mockApi, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "MyInstance", result[0][0])
	assert.Equal(t, "i-1234567890abcdef0", result[0][1])
	assert.Equal(t, "t2.micro", result[0][2])
	assert.Equal(t, "OnDemand", result[0][3])
	assert.Equal(t, "192.168.1.1", result[0][4])
	assert.Equal(t, "203.0.113.1", result[0][5])
	assert.Equal(t, "running", result[0][6])
	assert.Equal(t, "my-key", result[0][7])
	assert.Equal(t, "ap-northeast-1a", result[0][8])
	assert.NotEmpty(t, result[0][9])

	mockApi.AssertExpectations(t)
}

func TestDescribeInstances_Error(t *testing.T) {
	mockApi := new(mockEc2Api)
	mockApi.On("DescribeInstances", mock.Anything, mock.Anything, mock.Anything).Return(mockOutput, errors.New("error"))

	input := &ec2.DescribeInstancesInput{}
	result, err := DescribeInstances(mockApi, input)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockApi.AssertExpectations(t)
}
