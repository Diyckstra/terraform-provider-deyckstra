package paas

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceServiceInfraUpdateSchemaForceNew(t *testing.T) {
	t.Parallel()

	resource := ResourceService()

	if resource.Schema["instance_type"].ForceNew {
		t.Fatal("instance_type should support in-place updates")
	}

	dataVolume := resource.Schema["data_volume"]
	if dataVolume.ForceNew {
		t.Fatal("data_volume should allow selected nested fields to update in-place")
	}

	dataVolumeSchema := dataVolume.Elem.(*schema.Resource).Schema
	if dataVolumeSchema["size"].ForceNew {
		t.Fatal("data_volume.size should support in-place updates")
	}
	if dataVolumeSchema["iops"].ForceNew {
		t.Fatal("data_volume.iops should support in-place updates")
	}
	if !dataVolumeSchema["type"].ForceNew {
		t.Fatal("data_volume.type should still force replacement")
	}

	rootVolume := resource.Schema["root_volume"]
	if !rootVolume.ForceNew {
		t.Fatal("root_volume should still force replacement")
	}

	rootVolumeSchema := rootVolume.Elem.(*schema.Resource).Schema
	if !rootVolumeSchema["size"].ForceNew {
		t.Fatal("root_volume.size should still force replacement")
	}
	if !rootVolumeSchema["iops"].ForceNew {
		t.Fatal("root_volume.iops should still force replacement")
	}
	if !rootVolumeSchema["type"].ForceNew {
		t.Fatal("root_volume.type should still force replacement")
	}
}

func TestValidateIncreaseOnly(t *testing.T) {
	t.Parallel()

	if err := validateIncreaseOnly("data_volume.size", 32, 64); err != nil {
		t.Fatalf("expected increase to be allowed, got: %s", err)
	}

	if err := validateIncreaseOnly("data_volume.size", 32, 32); err != nil {
		t.Fatalf("expected no-op to be allowed, got: %s", err)
	}

	err := validateIncreaseOnly("data_volume.size", 64, 32)
	if err == nil {
		t.Fatal("expected decrease to be rejected")
	}
	if !strings.Contains(err.Error(), "can only be increased in-place") {
		t.Fatalf("expected clear increase-only error, got: %s", err)
	}
}
