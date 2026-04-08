---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_prometheus_scrape_job"
description: |-
  Manages a Prometheus scrape job for a PaaS service.
---

# Resource: aws_paas_prometheus_scrape_job

Manages a Prometheus scrape job for a PaaS service. Scrape jobs define which endpoints Prometheus collects metrics from.

## Example Usage

```terraform
resource "aws_paas_prometheus_scrape_job" "app" {
  service_id = aws_paas_service.prometheus.id

  name = "appmetrics"
  targets = [
    "app-01.internal:9090",
    "app-02.internal:9090",
  ]
  labels = {
    env  = "production"
    team = "platform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, Forces new resource) The ID of the PaaS Prometheus service.
* `name` - (Required, Forces new resource) The name of the scrape job. Must start with a Latin letter,
  contain only Latin letters and digits, and be between 1 and 32 characters.
* `targets` - (Required) A non-empty set of `host:port` addresses to scrape metrics from. Each address must be between 1 and 2048 characters.
* `labels` - (Optional) A map of labels to attach to all metrics collected by this job. A label name must be non-empty,
  can contain only Latin letters, digits, and underscores, and must not start with two underscores.
  Each label value can contain up to 256 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the scrape job.
* `job_id` - The ID of the scrape job (same as `id`).

## Import

PaaS Prometheus scrape jobs can be imported using the PaaS service ID and the scrape job ID in the `service_id/job_id` format, e.g.,

```
$ terraform import aws_paas_prometheus_scrape_job.example fm-cluster-12345678/job-87654321
```
