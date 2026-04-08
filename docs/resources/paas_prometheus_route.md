---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_prometheus_route"
description: |-
  Manages a Prometheus alerting route for a PaaS service.
---

# Resource: aws_paas_prometheus_route

Manages a Prometheus alerting route for a PaaS service. Routes define how incoming alerts are matched and forwarded to notification channels.

## Example Usage

```terraform
resource "aws_paas_prometheus_notification_channel" "telegram" {
  service_id = aws_paas_service.prometheus.id
  name       = "opstelegram"
  type       = "telegram"
  bot_token  = var.telegram_bot_token
  chat_id    = -1001234567890
}

resource "aws_paas_prometheus_route" "critical" {
  service_id = aws_paas_service.prometheus.id

  name     = "criticalalerts"
  receiver = aws_paas_prometheus_notification_channel.telegram.name
  matchers = ["severity = \"critical\""]

  group_by        = ["alertname", "instance"]
  group_wait      = "30s"
  group_interval  = "5m"
  repeat_interval = "3h"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, Forces new resource) The ID of the PaaS Prometheus service.
* `name` - (Required, Forces new resource) The name of the route. Must be between 1 and 256 characters.
* `receiver` - (Required) The **name** of a notification channel in the same PaaS Prometheus service. Must be between 1 and 256 characters.
* `matchers` - (Required) A non-empty set of label matchers that determines which alerts this route handles.
  Each matcher must use one of the `=`, `!=`, `=~`, or `!~` operators (for example, `severity = "critical"`).
  Multiple matchers are combined with logical AND.
* `continue` - (Optional) Whether to continue matching subsequent routes after this one matches. Defaults to `false`.
* `group_by` - (Optional) A set of label names to group alerts by. A label name must be between 1 and 64 characters,
  start with a Latin letter or underscore, contain only Latin letters, digits, and underscores, and must not start
  with two underscores.
* `group_wait` - (Optional) How long to wait before sending the first notification for a new alert group (for example, `30s`).
* `group_interval` - (Optional) How long to wait before sending a notification for new alerts added to an existing group (for example, `5m`).
* `repeat_interval` - (Optional) How long to wait before re-sending a notification for an alert that has already been sent (for example, `3h`).

The duration arguments use the K2 Prometheus duration ranges: `1-30d`, `1-24h`, `1-60m`, `1-60s`, and
`1-99999ms`. Components must be ordered from longest to shortest and can be combined, for example `30d24h60m60s99999ms`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the route.
* `route_id` - The ID of the route (same as `id`).

## Import

PaaS Prometheus routes can be imported using the PaaS service ID and the route ID in the `service_id/route_id` format, e.g.,

```
$ terraform import aws_paas_prometheus_route.example fm-cluster-12345678/route-87654321
```
