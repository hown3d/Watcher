package cloudformation

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

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
				client: NewClient("http://localhost:4566"),
			},
			want:    []types.Stack{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStacks(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStacks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
