package cloudformation

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

const cloudformationEndpoint = "http://localhost:4566"

func setupLocalstack() *cloudformation.Client {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
    if service == cloudformation.ServiceID {
        return aws.Endpoint{
            PartitionID:   "aws",
            URL:           cloudformationEndpoint,
            SigningRegion: "eu-central-1",
        }, nil
    }
    // returning EndpointNotFoundError will allow the service to fallback to it's default resolution
    return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(customResolver))
	return cloudformation.NewFromConfig(cfg)
}

func Test_getStacks(t *testing.T) {
	type args struct {
		client *cloudformation.Client
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Stack
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "no stacks",
			args: args{
				client: setupLocalstack(),
			},
			want:    []types.Stack{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStacks(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStacks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_getStackStatus(t *testing.T) {
	type args struct {
		stack types.Stack
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStackStatus(tt.args.stack); got != tt.want {
				t.Errorf("getStackStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
