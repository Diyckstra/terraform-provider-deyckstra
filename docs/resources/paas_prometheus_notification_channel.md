---
subcategory: "PaaS"
page_title: "aws_paas_prometheus_notification_channel"
description: |-
  Manages a Prometheus notification channel for a PaaS service.
---

# Resource: aws_paas_prometheus_notification_channel

Manages a Prometheus notification channel for a PaaS service.

~> **Note** Configure only the arguments that apply to the selected channel `type`. Telegram, webhook, and email arguments are mutually exclusive.

## Example Usage

### Telegram Channel

```terraform
resource "aws_paas_prometheus_notification_channel" "telegram" {
  service_id = aws_paas_service.prometheus.id

  name      = "opstelegram"
  type      = "telegram"
  bot_token = var.telegram_bot_token
  chat_id   = -1001234567890
}
```

### Webhook Channel

```terraform
resource "aws_paas_prometheus_notification_channel" "webhook" {
  service_id = aws_paas_service.prometheus.id

  name       = "opswebhook"
  type       = "webhook"
  url        = "https://alertmanager.example.com/hook"
  max_alerts = 10
}
```

### Email Channel

```terraform
resource "aws_paas_prometheus_notification_channel" "email" {
  service_id = aws_paas_service.prometheus.id

  name          = "opsemail"
  type          = "email"
  to            = "ops@example.com"
  from          = "alerts@example.com"
  smarthost     = "smtp.example.com:587"
  hello         = "example.com"
  require_tls   = true
  auth_username = "alerts@example.com"
  auth_password = var.smtp_password
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, Forces new resource) The ID of the PaaS Prometheus service.
* `name` - (Required, Forces new resource) The name of the notification channel. Must be between 1 and 256 characters.
* `type` - (Required, Forces new resource) The channel type.
    * _Valid values:_ `telegram`, `webhook`, `email`
* `is_default` - (Optional) Whether this channel is the default receiver. Defaults to `false`.
* `send_resolved` - (Optional) Whether to send notifications when an alert resolves. Defaults to `false`.

### Telegram arguments

The following arguments apply when `type = "telegram"`:

* `bot_token` - (Required) The Telegram bot token. This argument is sensitive; pass it through a sensitive variable or another secret source instead of hardcoding it.
* `chat_id` - (Required) The Telegram chat ID. May be negative for group chats.

### Webhook arguments

The following arguments apply when `type = "webhook"`:

* `url` - (Required) The webhook HTTP or HTTPS URL. Must be between 1 and 2048 characters.
* `max_alerts` - (Optional) Maximum number of alerts per webhook request. Must be a positive integer.

### Email arguments

The following arguments apply when `type = "email"`:

* `to` - (Required) The recipient email address. Must be between 1 and 2048 characters.
* `from` - (Optional) The sender email address. Must be between 1 and 2048 characters.
* `smarthost` - (Required) The SMTP relay host and port (for example, `smtp.example.com:587`).
  Must be between 1 and 2048 characters.
* `hello` - (Optional) The hostname or IP address used in the SMTP `EHLO` command. Must be between 1 and 253 characters.
* `require_tls` - (Optional) Whether to require TLS. Defaults to `false` and must be `true` when SMTP authentication is configured.
* `auth_username` - (Optional) The SMTP authentication username. Must be between 1 and 256 characters and is required together with `auth_password`.
* `auth_password` - (Optional) The SMTP authentication password. Must be between 1 and 256 characters and is required together with `auth_username`.
  This argument is sensitive; pass it through a sensitive variable or another secret source instead of hardcoding it.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the notification channel.
* `channel_id` - The ID of the notification channel (same as `id`).

## Import

PaaS Prometheus notification channels can be imported using the PaaS service ID and the notification channel ID in the
`service_id/channel_id` format, e.g.,

~> **Note** `bot_token` and `auth_password` are write-only: the PaaS API does not return their original values.
After import, specify the applicable secret again in the Terraform configuration before modifying the channel.

```
$ terraform import aws_paas_prometheus_notification_channel.example fm-cluster-12345678/channel-87654321
```
