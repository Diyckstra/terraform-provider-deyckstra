package services

import (
	"reflect"
	"testing"
)

func TestKafkaExpandServiceParameters_OptionsPreserveTypes(t *testing.T) {
	t.Parallel()

	in := map[string]interface{}{
		"version": "3.7.0",
		"options": map[string]interface{}{
			"autoCreateTopicsEnable": true,
			"logRetentionHours":      168,
			"compressionType":        "zstd",
		},
	}

	got := Kafka.expandServiceParameters(in)
	want := ServiceParameters{
		"version": "3.7.0",
		"options": map[string]interface{}{
			"autoCreateTopicsEnable": true,
			"logRetentionHours":      168,
			"compressionType":        "zstd",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected kafka parameters.\n got: %#v\nwant: %#v", got, want)
	}
}
