package paas

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
)

func TestFlattenServiceAdditionalRoles(t *testing.T) {
	testCases := []struct {
		name    string
		service *paas.Service
		want    []string
	}{
		{
			name: "combined coordinator role",
			service: &paas.Service{
				Nodes: &paas.Nodes{
					Main: &paas.Node{
						Role: aws.String("node,coordinator"),
					},
				},
			},
			want: []string{"coordinator"},
		},
		{
			name: "dedicated coordinator node",
			service: &paas.Service{
				Nodes: &paas.Nodes{
					Main: &paas.Node{
						Role: aws.String("node"),
					},
					Coordinator: &paas.Node{
						Role: aws.String("coordinator"),
					},
				},
			},
		},
		{
			name: "service without coordinator role",
			service: &paas.Service{
				Nodes: &paas.Nodes{
					Main: &paas.Node{
						Role: aws.String("node"),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := flattenServiceAdditionalRoles(testCase.service)

			if len(got) != len(testCase.want) {
				t.Fatalf("len(got) = %d, want %d", len(got), len(testCase.want))
			}

			for i := range got {
				if got[i] != testCase.want[i] {
					t.Fatalf("got[%d] = %q, want %q", i, got[i], testCase.want[i])
				}
			}
		})
	}
}

func TestFlattenServiceCoordinator(t *testing.T) {
	service := &paas.Service{
		Nodes: &paas.Nodes{
			Coordinator: &paas.Node{
				InstanceType:   aws.String("c4.large"),
				RootVolumeType: aws.String("gp2"),
				RootVolumeSize: aws.Int64(64),
				RootVolumeIops: aws.Int64(0),
				DataVolumeType: aws.String("st2"),
				DataVolumeSize: aws.Int64(128),
				DataVolumeIops: aws.Int64(0),
			},
		},
	}

	got := flattenServiceCoordinator(service)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	coordinator := got[0]
	if v := coordinator["instance_type"]; v != "c4.large" {
		t.Fatalf("instance_type = %v, want c4.large", v)
	}
	if v := coordinator["root_volume_type"]; v != "gp2" {
		t.Fatalf("root_volume_type = %v, want gp2", v)
	}
	if v := coordinator["root_volume_size"]; v != int64(64) {
		t.Fatalf("root_volume_size = %v, want 64", v)
	}
	if v := coordinator["data_volume_type"]; v != "st2" {
		t.Fatalf("data_volume_type = %v, want st2", v)
	}
	if v := coordinator["data_volume_size"]; v != int64(128) {
		t.Fatalf("data_volume_size = %v, want 128", v)
	}
}
