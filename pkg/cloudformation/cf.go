package cloudformation

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// NewClient returns a cloudformation Client, if url is not empty, a custom endpoint can be specified
func NewClient(url string) *cloudformation.Client {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Couldn't load aws config from file: %v", err)
	}
	if url != "" {
		log.Printf("Using custom resolver with url: %v", url)
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
func GetStacks(client *cloudformation.Client) ([]types.Stack, error) {

	output, err := client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{})
	if err != nil {
		return nil, err
	}
	return output.Stacks, nil
}

// GetStackInfo retrieves Information about a stack by its name
func GetStackInfo(stackname string, client *cloudformation.Client) ([]types.Stack, error) {
	output, err := client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{StackName: &stackname})
	if err != nil {
		return nil, err
	}
	return output.Stacks, nil
}

//GetStackEvents fetches events of a stack by its name
func GetStackEvents(stack string, client *cloudformation.Client) ([]types.StackEvent, error) {
	output, err := client.DescribeStackEvents(context.TODO(), &cloudformation.DescribeStackEventsInput{StackName: &stack})
	if err != nil {
		return nil, err
	}
	return output.StackEvents, nil
}

//GetStackResources returns all resources associated by a stack
func GetStackResources(stack string, client *cloudformation.Client) ([]types.StackResource, error) {
	output, err := client.DescribeStackResources(context.TODO(), &cloudformation.DescribeStackResourcesInput{StackName: &stack})
	if err != nil {
		return nil, err
	}
	return output.StackResources, nil
}
