package paas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

var errKafkaTopicNotFound = errors.New("paas kafka topic not found")

// Kafka topic names: ASCII alphanumerics, '.', '_', '-'; max length 249 (Apache Kafka convention).
var kafkaTopicNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

var kafkaTopicKnownParameterKeys = map[string]struct{}{
	"cleanupPolicy":               {},
	"compressionType":             {},
	"deleteRetentionMs":           {},
	"fileDeleteDelayMs":           {},
	"flushMessages":               {},
	"flushMs":                     {},
	"indexIntervalBytes":          {},
	"maxCompactionLagMs":          {},
	"maxMessageBytes":             {},
	"messageTimestampAfterMaxMs":  {},
	"messageTimestampBeforeMaxMs": {},
	"messageTimestampType":        {},
	"minCleanableDirtyRatio":      {},
	"minCompactionLagMs":          {},
	"minInsyncReplicas":           {},
	"partitions":                  {},
	"preallocate":                 {},
	"replicationFactor":           {},
	"retentionBytes":              {},
	"retentionMs":                 {},
	"segmentBytes":                {},
	"segmentIndexBytes":           {},
	"segmentJitterMs":             {},
	"segmentMs":                   {},
}

func ResourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKafkaTopicCreate,
		ReadContext:   resourceKafkaTopicRead,
		UpdateContext: resourceKafkaTopicUpdate,
		DeleteContext: resourceKafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKafkaTopicImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cleanup_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"compact", "delete", "delete,compact"},
					false,
				),
			},
			"compression_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"uncompressed", "zstd", "lz4", "snappy", "gzip", "producer"},
					false,
				),
			},
			"delete_retention_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"file_delete_delay_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"flush_messages": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"flush_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"index_interval_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 2147483647),
			},
			"max_compaction_lag_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"max_message_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 2147483647),
			},
			"message_timestamp_after_max_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"message_timestamp_before_max_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"message_timestamp_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"CreateTime", "LogAppendTime"},
					false,
				),
			},
			"min_cleanable_dirty_ratio": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.FloatBetween(0, 1),
			},
			"min_compaction_lag_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"min_insync_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},
			"preallocate": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"retention_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(-1, 2147483647),
			},
			"retention_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},
			"segment_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(14, 2147483647),
			},
			"segment_index_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(4, 2147483647),
			},
			"segment_jitter_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"segment_ms": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			// Kafka topic parameters accept mixed scalar values.
			//lintignore:S006
			"extra_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 249),
					validation.StringMatch(
						kafkaTopicNameRegexp,
						`must contain only letters, digits, '.', '_', and '-'`,
					),
				),
			},
			// partitions / replication_factor: minimum 1; maximum is the number of broker nodes
			// (K2 enforces this on the API). See:
			// https://docs.k2.cloud/en/api/paas/parameters/kafka.html
			"partitions": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceId := d.Get("service_id").(string)
	name := d.Get("name").(string)

	parameters := expandKafkaTopicParameters(d)

	input := &paas.CreateKafkaTopicInput{
		ServiceId:  aws.String(serviceId),
		Name:       aws.String(name),
		Parameters: parameters,
	}

	log.Printf("[DEBUG] Creating PaaS Kafka Topic: %s", input)
	output, err := conn.CreateKafkaTopic(input)

	if err != nil {
		return diag.Errorf("error creating PaaS Kafka Topic (%s) on service (%s): %s", name, serviceId, err)
	}

	_, err = waitServiceUpdated(ctx, conn, serviceId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf(
			"error waiting for PaaS Service (%s) to become ready after creating Kafka Topic (%s): %s",
			serviceId,
			name,
			err,
		)
	}

	topicId := aws.StringValue(output.KafkaTopic.Id)
	d.SetId(kafkaTopicCreateResourceID(serviceId, topicId))
	d.Set("topic_id", topicId)

	return resourceKafkaTopicRead(ctx, d, meta)
}

func resourceKafkaTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceId, topicId, err := kafkaTopicParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	topic, err := FindKafkaTopicByID(conn, serviceId, topicId)

	if err != nil {
		if errors.Is(err, errKafkaTopicNotFound) {
			log.Printf("[WARN] PaaS Kafka Topic (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("error reading PaaS Kafka Topic (%s): %s", d.Id(), err)
	}

	d.Set("service_id", serviceId)
	d.Set("topic_id", topicId)
	d.Set("name", topic.Name)

	flattenKafkaTopicParameters(d, topic.Parameters)

	return nil
}

func resourceKafkaTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceId, topicId, err := kafkaTopicParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Per API contract: all editable parameters must be sent on modify;
	// omitted parameters may be cleared.
	parameters := expandKafkaTopicParameters(d)

	input := &paas.ModifyKafkaTopicInput{
		ServiceId:  aws.String(serviceId),
		TopicId:    aws.String(topicId),
		Parameters: parameters,
	}

	log.Printf("[DEBUG] Modifying PaaS Kafka Topic: %s", input)
	_, err = conn.ModifyKafkaTopic(input)

	if err != nil {
		return diag.Errorf("error modifying PaaS Kafka Topic (%s): %s", d.Id(), err)
	}

	_, err = waitServiceUpdated(ctx, conn, serviceId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf(
			"error waiting for PaaS Service (%s) to become ready after modifying Kafka Topic (%s): %s",
			serviceId,
			topicId,
			err,
		)
	}

	return resourceKafkaTopicRead(ctx, d, meta)
}

func resourceKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceId, topicId, err := kafkaTopicParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	input := &paas.DeleteKafkaTopicInput{
		ServiceId: aws.String(serviceId),
		TopicId:   aws.String(topicId),
	}

	log.Printf("[DEBUG] Deleting PaaS Kafka Topic: %s", input)
	_, err = conn.DeleteKafkaTopic(input)

	if err != nil {
		return diag.Errorf("error deleting PaaS Kafka Topic (%s): %s", d.Id(), err)
	}

	_, err = waitServiceUpdated(ctx, conn, serviceId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf(
			"error waiting for PaaS Service (%s) to become ready after deleting Kafka Topic (%s): %s",
			serviceId,
			topicId,
			err,
		)
	}

	return nil
}

func resourceKafkaTopicImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	serviceId, topicId, err := kafkaTopicParseResourceID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("service_id", serviceId)
	d.Set("topic_id", topicId)

	return []*schema.ResourceData{d}, nil
}

func kafkaTopicCreateResourceID(serviceId, topicId string) string {
	return serviceId + "/" + topicId
}

func kafkaTopicParseResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format for ID (%s), expected service_id/topic_id", id)
	}
	return parts[0], parts[1], nil
}

func FindKafkaTopicByID(conn *paas.PaaS, serviceId, topicId string) (*paas.KafkaTopic, error) {
	input := &paas.ListKafkaTopicsInput{
		ServiceId: aws.String(serviceId),
	}

	output, err := conn.ListKafkaTopics(input)
	if err != nil {
		// If the parent service no longer exists, the topic is gone too. Treat it
		// as not found so the topic is removed from state instead of erroring on
		// every plan/refresh/destroy.
		if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
			return nil, fmt.Errorf("%w: %s on service %s", errKafkaTopicNotFound, topicId, serviceId)
		}
		return nil, fmt.Errorf("error listing PaaS Kafka Topics for service (%s): %w", serviceId, err)
	}

	if output == nil {
		return nil, fmt.Errorf("%w: %s on service %s", errKafkaTopicNotFound, topicId, serviceId)
	}

	topic := findKafkaTopicInList(output.KafkaTopics, topicId)
	if topic != nil {
		return topic, nil
	}

	return nil, fmt.Errorf("%w: %s on service %s", errKafkaTopicNotFound, topicId, serviceId)
}

func findKafkaTopicInList(topics []*paas.KafkaTopic, topicId string) *paas.KafkaTopic {
	for _, topic := range topics {
		if topic != nil && aws.StringValue(topic.Id) == topicId {
			return topic
		}
	}

	return nil
}

func expandKafkaTopicParameters(d *schema.ResourceData) map[string]interface{} {
	parameters := map[string]interface{}{
		"partitions":        int64(d.Get("partitions").(int)),
		"replicationFactor": int64(d.Get("replication_factor").(int)),
	}

	if v, ok := d.GetOk("cleanup_policy"); ok {
		parameters["cleanupPolicy"] = v.(string)
	}

	if v, ok := d.GetOk("compression_type"); ok {
		parameters["compressionType"] = v.(string)
	}

	expandKafkaTopicIntParameter(d, parameters, "delete_retention_ms", "deleteRetentionMs")
	expandKafkaTopicIntParameter(d, parameters, "file_delete_delay_ms", "fileDeleteDelayMs")
	expandKafkaTopicIntParameter(d, parameters, "flush_messages", "flushMessages")
	expandKafkaTopicIntParameter(d, parameters, "flush_ms", "flushMs")
	expandKafkaTopicIntParameter(d, parameters, "index_interval_bytes", "indexIntervalBytes")
	expandKafkaTopicIntParameter(d, parameters, "max_compaction_lag_ms", "maxCompactionLagMs")
	expandKafkaTopicIntParameter(d, parameters, "max_message_bytes", "maxMessageBytes")
	expandKafkaTopicIntParameter(d, parameters, "message_timestamp_after_max_ms", "messageTimestampAfterMaxMs")
	expandKafkaTopicIntParameter(d, parameters, "message_timestamp_before_max_ms", "messageTimestampBeforeMaxMs")
	expandKafkaTopicIntParameter(d, parameters, "min_compaction_lag_ms", "minCompactionLagMs")
	expandKafkaTopicIntParameter(d, parameters, "min_insync_replicas", "minInsyncReplicas")
	expandKafkaTopicIntParameter(d, parameters, "retention_bytes", "retentionBytes")
	expandKafkaTopicIntParameter(d, parameters, "retention_ms", "retentionMs")
	expandKafkaTopicIntParameter(d, parameters, "segment_bytes", "segmentBytes")
	expandKafkaTopicIntParameter(d, parameters, "segment_index_bytes", "segmentIndexBytes")
	expandKafkaTopicIntParameter(d, parameters, "segment_jitter_ms", "segmentJitterMs")
	expandKafkaTopicIntParameter(d, parameters, "segment_ms", "segmentMs")

	if v, ok := d.GetOk("message_timestamp_type"); ok {
		parameters["messageTimestampType"] = v.(string)
	}

	// GetOkExists is required (instead of GetOk) so that an explicitly configured
	// zero/false value is still sent. GetOk treats 0 and false as "unset".
	if v, ok := d.GetOkExists("min_cleanable_dirty_ratio"); ok {
		parameters["minCleanableDirtyRatio"] = v.(float64)
	}

	if v, ok := d.GetOkExists("preallocate"); ok {
		parameters["preallocate"] = v.(bool)
	}

	if v, ok := d.GetOk("extra_parameters"); ok {
		for k, raw := range v.(map[string]interface{}) {
			if _, exists := kafkaTopicKnownParameterKeys[k]; exists {
				continue
			}
			parameters[k] = raw
		}
	}

	return parameters
}

func flattenKafkaTopicParameters(d *schema.ResourceData, parameters map[string]interface{}) {
	if parameters == nil {
		return
	}

	if v, ok := parameters["partitions"]; ok {
		switch val := v.(type) {
		case float64:
			d.Set("partitions", int(val))
		case int64:
			d.Set("partitions", int(val))
		case int:
			d.Set("partitions", val)
		}
	}

	if v, ok := parameters["replicationFactor"]; ok {
		switch val := v.(type) {
		case float64:
			d.Set("replication_factor", int(val))
		case int64:
			d.Set("replication_factor", int(val))
		case int:
			d.Set("replication_factor", val)
		}
	}

	if v, ok := parameters["cleanupPolicy"].(string); ok {
		d.Set("cleanup_policy", v)
	}

	if v, ok := parameters["compressionType"].(string); ok {
		d.Set("compression_type", v)
	}

	flattenKafkaTopicIntParameter(d, parameters, "deleteRetentionMs", "delete_retention_ms")
	flattenKafkaTopicIntParameter(d, parameters, "fileDeleteDelayMs", "file_delete_delay_ms")
	flattenKafkaTopicIntParameter(d, parameters, "flushMessages", "flush_messages")
	flattenKafkaTopicIntParameter(d, parameters, "flushMs", "flush_ms")
	flattenKafkaTopicIntParameter(d, parameters, "indexIntervalBytes", "index_interval_bytes")
	flattenKafkaTopicIntParameter(d, parameters, "maxCompactionLagMs", "max_compaction_lag_ms")
	flattenKafkaTopicIntParameter(d, parameters, "maxMessageBytes", "max_message_bytes")
	flattenKafkaTopicIntParameter(d, parameters, "messageTimestampAfterMaxMs", "message_timestamp_after_max_ms")
	flattenKafkaTopicIntParameter(d, parameters, "messageTimestampBeforeMaxMs", "message_timestamp_before_max_ms")
	flattenKafkaTopicIntParameter(d, parameters, "minCompactionLagMs", "min_compaction_lag_ms")
	flattenKafkaTopicIntParameter(d, parameters, "minInsyncReplicas", "min_insync_replicas")
	flattenKafkaTopicIntParameter(d, parameters, "retentionBytes", "retention_bytes")
	flattenKafkaTopicIntParameter(d, parameters, "retentionMs", "retention_ms")
	flattenKafkaTopicIntParameter(d, parameters, "segmentBytes", "segment_bytes")
	flattenKafkaTopicIntParameter(d, parameters, "segmentIndexBytes", "segment_index_bytes")
	flattenKafkaTopicIntParameter(d, parameters, "segmentJitterMs", "segment_jitter_ms")
	flattenKafkaTopicIntParameter(d, parameters, "segmentMs", "segment_ms")

	if v, ok := parameters["messageTimestampType"].(string); ok {
		d.Set("message_timestamp_type", v)
	}

	if v, ok := parameters["minCleanableDirtyRatio"]; ok {
		switch val := v.(type) {
		case float64:
			d.Set("min_cleanable_dirty_ratio", val)
		case int64:
			d.Set("min_cleanable_dirty_ratio", float64(val))
		case int:
			d.Set("min_cleanable_dirty_ratio", float64(val))
		}
	}

	if v, ok := parameters["preallocate"].(bool); ok {
		d.Set("preallocate", v)
	}

	extraParameters := make(map[string]interface{})
	for k, v := range parameters {
		if _, exists := kafkaTopicKnownParameterKeys[k]; exists {
			continue
		}
		extraParameters[k] = v
	}
	d.Set("extra_parameters", extraParameters)
}

func expandKafkaTopicIntParameter(d *schema.ResourceData, parameters map[string]interface{}, tfKey, apiKey string) {
	// GetOkExists is required (instead of GetOk) so that an explicitly configured
	// 0 is still sent. GetOk treats 0 as "unset", which would drop valid values
	// such as delete_retention_ms = 0 or segment_jitter_ms = 0.
	if v, ok := d.GetOkExists(tfKey); ok {
		parameters[apiKey] = int64(v.(int))
	}
}

func flattenKafkaTopicIntParameter(d *schema.ResourceData, parameters map[string]interface{}, apiKey, tfKey string) {
	if v, ok := parameters[apiKey]; ok {
		switch val := v.(type) {
		case float64:
			d.Set(tfKey, int(val))
		case int64:
			d.Set(tfKey, int(val))
		case int:
			d.Set(tfKey, val)
		}
	}
}
