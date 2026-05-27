---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_kafka_topic"
description: |-
  Manages a PaaS Kafka topic.
---

# Resource: aws_paas_kafka_topic

Manages a Kafka topic on a PaaS Kafka service.

## Example Usage

```terraform
resource "aws_paas_service" "kafka" {
  name          = "tf-kafka"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 64
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]

  ssh_key_name     = "<name>"
  additional_roles = ["coordinator"]

  kafka {
    version = "3.7.0"
  }
}

resource "aws_paas_kafka_topic" "events" {
  service_id         = aws_paas_service.kafka.id
  name               = "events"
  partitions         = 1
  replication_factor = 1
  cleanup_policy     = "delete"
  compression_type   = "producer"
  retention_ms       = 604800000
  segment_ms         = 3600000
  preallocate        = false
  extra_parameters = {
    # Use API parameter names in camelCase.
    customFutureOption = "enabled"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, ForceNew) The ID of the PaaS Kafka service.
* `name` - (Required, ForceNew) The name of the Kafka topic.
* `partitions` - (Required) The number of partitions in the topic.
    * _Valid values:_ From 1 to the number of broker nodes.
    * _Constraints:_ The number of partitions cannot be decreased when modifying.
* `replication_factor` - (Required) The number of partition replicas.
    * _Valid values:_ From 1 to the number of broker nodes.
* `cleanup_policy` - (Optional) The cleanup policy for the topic.
    * _Valid values:_ `compact`, `delete`, `delete,compact`
* `compression_type` - (Optional) The compression type for the topic.
    * _Valid values:_ `uncompressed`, `zstd`, `lz4`, `snappy`, `gzip`, `producer`
* `delete_retention_ms` - (Optional) The time in milliseconds to retain delete tombstone markers.
* `file_delete_delay_ms` - (Optional) The time in milliseconds to wait before deleting a file.
* `flush_messages` - (Optional) The number of messages accumulated before forcing an fsync.
* `flush_ms` - (Optional) The maximum time in milliseconds before forcing an fsync.
* `index_interval_bytes` - (Optional) The number of bytes between index entries.
* `max_compaction_lag_ms` - (Optional) The maximum time in milliseconds a message remains uncompacted.
* `max_message_bytes` - (Optional) The largest record batch size accepted by this topic.
* `message_timestamp_after_max_ms` - (Optional) The maximum allowed timestamp difference in milliseconds into the future.
* `message_timestamp_before_max_ms` - (Optional) The maximum allowed timestamp difference in milliseconds into the past.
* `message_timestamp_type` - (Optional) Defines the timestamp type used for messages.
    * _Valid values:_ `CreateTime`, `LogAppendTime`
* `min_cleanable_dirty_ratio` - (Optional) The minimum ratio of dirty log to total log for log compaction.
    * _Valid values:_ From `0.0` to `1.0`
* `min_compaction_lag_ms` - (Optional) The minimum time in milliseconds a message remains uncompacted.
* `min_insync_replicas` - (Optional) The minimum number of in-sync replicas required to accept produce requests.
* `preallocate` - (Optional) Whether to preallocate log segment files on disk.
* `retention_bytes` - (Optional) The maximum size of the log before old segments are discarded.
* `retention_ms` - (Optional) The maximum time in milliseconds to retain a log before old segments are discarded.
* `segment_bytes` - (Optional) The segment file size in bytes.
* `segment_index_bytes` - (Optional) The segment index file size in bytes.
* `segment_jitter_ms` - (Optional) The maximum random jitter subtracted from `segment_ms`.
* `segment_ms` - (Optional) The time in milliseconds after which a new log segment is rolled.
* `extra_parameters` - (Optional, Editable) Map with additional Kafka topic parameters that are not yet described by dedicated Terraform arguments.
  Keys must match API parameter names (camelCase). Values can be strings, numbers, or booleans.
  If the same key is also configured via a dedicated Terraform argument, the dedicated argument takes precedence.

~> **Note** When modifying a Kafka topic, all editable parameters are sent.
If an editable parameter is omitted in the modify request, its value may be cleared.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `topic_id` - The ID of the Kafka topic.

## Import

PaaS Kafka topics can be imported using a composite ID of `service_id/topic_id`, e.g.,

```
$ terraform import aws_paas_kafka_topic.events fm-cluster-12345678/my-topic-id
```
