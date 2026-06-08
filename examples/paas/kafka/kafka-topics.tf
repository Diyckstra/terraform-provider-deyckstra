# Kafka topics managed by the provider.
#
# Limits:
#   * partitions          must be >= 1 and <= number of broker nodes
#                         (3 for HA, 1 for non-HA).
#   * replication_factor  must be >= 1 and <= number of broker nodes.
#
# Optional knobs:
#   * cleanup_policy      one of: "delete", "compact", "delete,compact".
#   * compression_type    one of: "uncompressed", "zstd", "lz4", "snappy",
#                         "gzip", "producer".

resource "aws_paas_kafka_topic" "events" {
  service_id         = aws_paas_service.example.id
  name               = "events"
  partitions         = 3
  replication_factor = 3
  cleanup_policy     = "delete"
  compression_type   = "producer"
}

resource "aws_paas_kafka_topic" "logs" {
  service_id         = aws_paas_service.example.id
  name               = "logs"
  partitions         = 1
  replication_factor = 3
}

resource "aws_paas_kafka_topic" "state" {
  service_id         = aws_paas_service.example.id
  name               = "state"
  partitions         = 2
  replication_factor = 3
  cleanup_policy     = "compact"
  compression_type   = "lz4"
}

output "topic_events_id" {
  description = "ID of the events topic."
  value       = aws_paas_kafka_topic.events.topic_id
}

output "topic_logs_id" {
  description = "ID of the logs topic."
  value       = aws_paas_kafka_topic.logs.topic_id
}

output "topic_state_id" {
  description = "ID of the state topic."
  value       = aws_paas_kafka_topic.state.topic_id
}
