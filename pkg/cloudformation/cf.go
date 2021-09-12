package cloudformation

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func NewClient() *cloudformation.Client {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon CF service client

	return cloudformation.NewFromConfig(cfg)

}
func getStacks(client *cloudformation.Client) ([]types.Stack, error) {
	output, err := client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{})
	if err != nil {
		return nil, err
	}
	return output.Stacks, nil
}

func getStackStatus(stack types.Stack) string {
	return *stack.StackStatusReason
}
