package cloudformation

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// Stack has a name, which resembles the cloudformation stack name and its status
type Stack struct {
	status string
	name string
}

// NewClient returns a cloudformation Client, if url is not empty, a custom endpoint can be specified
func NewClient(url string) *cloudformation.Client {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if url != "" {	
		cfg.EndpointResolver = getCustomEndpointResolver(url, cfg.Region)
	}
	// Create an Amazon CF service client
	return cloudformation.NewFromConfig(cfg)
}

func getCustomEndpointResolver(url string, region string) aws.EndpointResolverFunc {
	return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == cloudformation.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: region,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
}

//GetStacks fetches all Cloudformation stacks from a client and returns them as a list
func GetStacks(client *cloudformation.Client) ([]Stack, error) {
	var stacks []Stack
	output, err := client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{})
	if err != nil {
		return nil, err
	}
	for _, stack := range output.Stacks {
		stacks = append(stacks, Stack{name: *stack.StackName, status: *stack.StackStatusReason})
	}
	return stacks, nil
}

