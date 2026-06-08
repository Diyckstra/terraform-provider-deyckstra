package paas

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestKafkaTopicParseResourceID(t *testing.T) {
	tests := []struct {
		id        string
		serviceId string
		topicId   string
		expectErr bool
	}{
		{
			id:        "fm-cluster-123/topic-456",
			serviceId: "fm-cluster-123",
			topicId:   "topic-456",
		},
		{
			id:        "svc-abc/my-topic",
			serviceId: "svc-abc",
			topicId:   "my-topic",
		},
		{
			id:        "invalid-no-slash",
			expectErr: true,
		},
		{
			id:        "/topic-only",
			expectErr: true,
		},
		{
			id:        "service-only/",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			serviceId, topicId, err := kafkaTopicParseResourceID(tt.id)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for ID %q, got nil", tt.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for ID %q: %s", tt.id, err)
			}
			if serviceId != tt.serviceId {
				t.Errorf("serviceId = %q, want %q", serviceId, tt.serviceId)
			}
			if topicId != tt.topicId {
				t.Errorf("topicId = %q, want %q", topicId, tt.topicId)
			}
		})
	}
}

func TestKafkaTopicCreateResourceID(t *testing.T) {
	id := kafkaTopicCreateResourceID("fm-cluster-123", "topic-456")
	expected := "fm-cluster-123/topic-456"
	if id != expected {
		t.Errorf("got %q, want %q", id, expected)
	}
}

func TestExpandKafkaTopicParameters(t *testing.T) {
	rawConfig := map[string]interface{}{
		"service_id":                      "fm-cluster-123",
		"name":                            "events",
		"partitions":                      3,
		"replication_factor":              3,
		"cleanup_policy":                  "delete",
		"compression_type":                "producer",
		"delete_retention_ms":             0,
		"file_delete_delay_ms":            1000,
		"flush_messages":                  10,
		"flush_ms":                        1000,
		"index_interval_bytes":            4096,
		"max_compaction_lag_ms":           10000,
		"max_message_bytes":               1048576,
		"message_timestamp_after_max_ms":  1000,
		"message_timestamp_before_max_ms": 1000,
		"message_timestamp_type":          "CreateTime",
		"min_cleanable_dirty_ratio":       0.5,
		"min_compaction_lag_ms":           0,
		"min_insync_replicas":             1,
		"preallocate":                     true,
		"retention_bytes":                 -1,
		"retention_ms":                    -1,
		"segment_bytes":                   1048576,
		"segment_index_bytes":             1024,
		"segment_jitter_ms":               0,
		"segment_ms":                      1000,
		"extra_parameters": map[string]interface{}{
			"customFutureOption": "enabled",
		},
		"topic_id": "",
	}

	resourceSchema := ResourceKafkaTopic().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, rawConfig)

	params := expandKafkaTopicParameters(d)

	if v, ok := params["partitions"].(int64); !ok || v != 3 {
		t.Errorf("partitions = %v, want 3", params["partitions"])
	}

	if v, ok := params["replicationFactor"].(int64); !ok || v != 3 {
		t.Errorf("replicationFactor = %v, want 3", params["replicationFactor"])
	}

	if v, ok := params["cleanupPolicy"].(string); !ok || v != "delete" {
		t.Errorf("cleanupPolicy = %v, want 'delete'", params["cleanupPolicy"])
	}

	if v, ok := params["compressionType"].(string); !ok || v != "producer" {
		t.Errorf("compressionType = %v, want 'producer'", params["compressionType"])
	}

	if v, ok := params["messageTimestampType"].(string); !ok || v != "CreateTime" {
		t.Errorf("messageTimestampType = %v, want 'CreateTime'", params["messageTimestampType"])
	}

	if v, ok := params["minCleanableDirtyRatio"].(float64); !ok || v != 0.5 {
		t.Errorf("minCleanableDirtyRatio = %v, want 0.5", params["minCleanableDirtyRatio"])
	}

	if v, ok := params["preallocate"].(bool); !ok || !v {
		t.Errorf("preallocate = %v, want true", params["preallocate"])
	}

	intKeys := map[string]int64{
		"deleteRetentionMs":           0,
		"fileDeleteDelayMs":           1000,
		"flushMessages":               10,
		"flushMs":                     1000,
		"indexIntervalBytes":          4096,
		"maxCompactionLagMs":          10000,
		"maxMessageBytes":             1048576,
		"messageTimestampAfterMaxMs":  1000,
		"messageTimestampBeforeMaxMs": 1000,
		"minCompactionLagMs":          0,
		"minInsyncReplicas":           1,
		"retentionBytes":              -1,
		"retentionMs":                 -1,
		"segmentBytes":                1048576,
		"segmentIndexBytes":           1024,
		"segmentJitterMs":             0,
		"segmentMs":                   1000,
	}

	for key, want := range intKeys {
		if v, ok := params[key].(int64); !ok || v != want {
			t.Errorf("%s = %v, want %d", key, params[key], want)
		}
	}

	if v, ok := params["customFutureOption"].(string); !ok || v != "enabled" {
		t.Errorf("customFutureOption = %v, want 'enabled'", params["customFutureOption"])
	}
}

func TestExpandKafkaTopicParameters_AllFieldsSent(t *testing.T) {
	rawConfig := map[string]interface{}{
		"service_id":                      "fm-cluster-123",
		"name":                            "events",
		"partitions":                      6,
		"replication_factor":              2,
		"cleanup_policy":                  "compact",
		"compression_type":                "lz4",
		"delete_retention_ms":             10,
		"file_delete_delay_ms":            10,
		"flush_messages":                  10,
		"flush_ms":                        10,
		"index_interval_bytes":            4096,
		"max_compaction_lag_ms":           10,
		"max_message_bytes":               1024,
		"message_timestamp_after_max_ms":  10,
		"message_timestamp_before_max_ms": 10,
		"message_timestamp_type":          "LogAppendTime",
		"min_cleanable_dirty_ratio":       0.7,
		"min_compaction_lag_ms":           10,
		"min_insync_replicas":             1,
		"preallocate":                     false,
		"retention_bytes":                 1024,
		"retention_ms":                    1000,
		"segment_bytes":                   1024,
		"segment_index_bytes":             1024,
		"segment_jitter_ms":               10,
		"segment_ms":                      10,
		"extra_parameters": map[string]interface{}{
			"anotherCustomOption": "value",
		},
		"topic_id": "",
	}

	resourceSchema := ResourceKafkaTopic().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, rawConfig)

	params := expandKafkaTopicParameters(d)

	requiredKeys := []string{"partitions", "replicationFactor"}
	for _, key := range requiredKeys {
		if _, ok := params[key]; !ok {
			t.Errorf("expected key %q in parameters, but it was missing", key)
		}
	}

	optionalKeys := []string{
		"cleanupPolicy",
		"compressionType",
		"deleteRetentionMs",
		"fileDeleteDelayMs",
		"flushMessages",
		"flushMs",
		"indexIntervalBytes",
		"maxCompactionLagMs",
		"maxMessageBytes",
		"messageTimestampAfterMaxMs",
		"messageTimestampBeforeMaxMs",
		"messageTimestampType",
		"minCleanableDirtyRatio",
		"minCompactionLagMs",
		"minInsyncReplicas",
		"preallocate",
		"retentionBytes",
		"retentionMs",
		"segmentBytes",
		"segmentIndexBytes",
		"segmentJitterMs",
		"segmentMs",
		"anotherCustomOption",
	}
	for _, key := range optionalKeys {
		if _, ok := params[key]; !ok {
			t.Errorf("expected key %q in parameters, but it was missing", key)
		}
	}
}

func TestFlattenKafkaTopicParameters(t *testing.T) {
	resourceSchema := ResourceKafkaTopic().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"service_id":                      "fm-cluster-123",
		"name":                            "events",
		"partitions":                      0,
		"replication_factor":              0,
		"cleanup_policy":                  "",
		"compression_type":                "",
		"delete_retention_ms":             0,
		"file_delete_delay_ms":            0,
		"flush_messages":                  0,
		"flush_ms":                        0,
		"index_interval_bytes":            0,
		"max_compaction_lag_ms":           0,
		"max_message_bytes":               0,
		"message_timestamp_after_max_ms":  0,
		"message_timestamp_before_max_ms": 0,
		"message_timestamp_type":          "",
		"min_cleanable_dirty_ratio":       0.0,
		"min_compaction_lag_ms":           0,
		"min_insync_replicas":             0,
		"preallocate":                     false,
		"retention_bytes":                 0,
		"retention_ms":                    0,
		"segment_bytes":                   0,
		"segment_index_bytes":             0,
		"segment_jitter_ms":               0,
		"segment_ms":                      0,
		"extra_parameters":                map[string]interface{}{},
		"topic_id":                        "",
	})

	apiParams := map[string]interface{}{
		"partitions":                  float64(3),
		"replicationFactor":           float64(3),
		"cleanupPolicy":               "delete",
		"compressionType":             "producer",
		"deleteRetentionMs":           float64(0),
		"fileDeleteDelayMs":           float64(1000),
		"flushMessages":               float64(10),
		"flushMs":                     float64(1000),
		"indexIntervalBytes":          float64(4096),
		"maxCompactionLagMs":          float64(10000),
		"maxMessageBytes":             float64(1048576),
		"messageTimestampAfterMaxMs":  float64(1000),
		"messageTimestampBeforeMaxMs": float64(1000),
		"messageTimestampType":        "CreateTime",
		"minCleanableDirtyRatio":      float64(0.5),
		"minCompactionLagMs":          float64(0),
		"minInsyncReplicas":           float64(1),
		"preallocate":                 true,
		"retentionBytes":              float64(-1),
		"retentionMs":                 float64(-1),
		"segmentBytes":                float64(1048576),
		"segmentIndexBytes":           float64(1024),
		"segmentJitterMs":             float64(0),
		"segmentMs":                   float64(1000),
		"customFutureOption":          "enabled",
	}

	flattenKafkaTopicParameters(d, apiParams)

	if v := d.Get("partitions").(int); v != 3 {
		t.Errorf("partitions = %d, want 3", v)
	}

	if v := d.Get("replication_factor").(int); v != 3 {
		t.Errorf("replication_factor = %d, want 3", v)
	}

	if v := d.Get("cleanup_policy").(string); v != "delete" {
		t.Errorf("cleanup_policy = %q, want 'delete'", v)
	}

	if v := d.Get("compression_type").(string); v != "producer" {
		t.Errorf("compression_type = %q, want 'producer'", v)
	}

	if v := d.Get("message_timestamp_type").(string); v != "CreateTime" {
		t.Errorf("message_timestamp_type = %q, want 'CreateTime'", v)
	}

	if v := d.Get("min_cleanable_dirty_ratio").(float64); v != 0.5 {
		t.Errorf("min_cleanable_dirty_ratio = %v, want 0.5", v)
	}

	if v := d.Get("preallocate").(bool); !v {
		t.Errorf("preallocate = %v, want true", v)
	}

	intFields := map[string]int{
		"delete_retention_ms":             0,
		"file_delete_delay_ms":            1000,
		"flush_messages":                  10,
		"flush_ms":                        1000,
		"index_interval_bytes":            4096,
		"max_compaction_lag_ms":           10000,
		"max_message_bytes":               1048576,
		"message_timestamp_after_max_ms":  1000,
		"message_timestamp_before_max_ms": 1000,
		"min_compaction_lag_ms":           0,
		"min_insync_replicas":             1,
		"retention_bytes":                 -1,
		"retention_ms":                    -1,
		"segment_bytes":                   1048576,
		"segment_index_bytes":             1024,
		"segment_jitter_ms":               0,
		"segment_ms":                      1000,
	}

	for field, want := range intFields {
		if v := d.Get(field).(int); v != want {
			t.Errorf("%s = %d, want %d", field, v, want)
		}
	}

	extra := d.Get("extra_parameters").(map[string]interface{})
	if v, ok := extra["customFutureOption"].(string); !ok || v != "enabled" {
		t.Errorf("extra_parameters.customFutureOption = %v, want 'enabled'", extra["customFutureOption"])
	}
}

func TestFindKafkaTopicByID_Found(t *testing.T) {
	topics := []*paas.KafkaTopic{
		{
			Id:   aws.String("topic-1"),
			Name: aws.String("events"),
			Parameters: map[string]interface{}{
				"partitions":        float64(3),
				"replicationFactor": float64(3),
			},
		},
		{
			Id:   aws.String("topic-2"),
			Name: aws.String("logs"),
			Parameters: map[string]interface{}{
				"partitions":        float64(1),
				"replicationFactor": float64(1),
			},
		},
	}

	found := findKafkaTopicInList(topics, "topic-2")
	if found == nil {
		t.Fatal("expected to find topic-2, got nil")
	}
	if aws.StringValue(found.Name) != "logs" {
		t.Errorf("name = %q, want 'logs'", aws.StringValue(found.Name))
	}
}

func TestFindKafkaTopicByID_NotFound(t *testing.T) {
	topics := []*paas.KafkaTopic{
		{
			Id:   aws.String("topic-1"),
			Name: aws.String("events"),
		},
	}

	found := findKafkaTopicInList(topics, "topic-missing")
	if found != nil {
		t.Errorf("expected nil for missing topic, got %+v", found)
	}
}
