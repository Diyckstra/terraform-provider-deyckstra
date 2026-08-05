---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_logstash_pipeline"
description: |-
  Manages a Logstash data pipeline for a PaaS ELK service.
---

# Resource: aws_paas_logstash_pipeline

Manages a Logstash data pipeline that belongs to a K2 Cloud PaaS ELK service.

## Example Usage

```terraform
resource "aws_paas_logstash_pipeline" "events" {
  service_id    = aws_paas_service.elk.id
  name          = "application-events"
  configuration = "input { http { port => 4567 tags => [\"events\"] } }"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, Forces new resource) The ID of the PaaS ELK service.
* `name` - (Required, Forces new resource) The name of the data pipeline.
  The name `beats-to-elasticsearch` is reserved by K2 Cloud and cannot be used.
* `configuration` - (Required, Sensitive) The Logstash pipeline configuration. Changes are applied in place.
  Syntax validation is performed by the PaaS API. Terraform state still contains
  the value, so state storage must be protected.

~> **Note** The live PaaS API currently rejects literal newline and other
control characters even though its published examples are multiline. Use a
one-line configuration until the cloud behavior is corrected. The provider
continues to delegate Logstash syntax validation to PaaS.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Logstash pipeline.
* `pipeline_id` - The ID of the Logstash pipeline (same as `id`).

## Timeouts

The `timeouts` block allows the following actions to be configured:

* `create` - (Default `60 minutes`) How long to wait for the parent ELK service after pipeline creation.
* `update` - (Default `60 minutes`) How long to wait for the parent ELK service after a pipeline update.
* `delete` - (Default `60 minutes`) How long to wait for the parent ELK service after pipeline deletion.

## Import

PaaS Logstash pipelines can be imported using `service_id/pipeline_id`, e.g.,

```
$ terraform import aws_paas_logstash_pipeline.example fm-cluster-12345678/logstash-pipeline-87654321
```
