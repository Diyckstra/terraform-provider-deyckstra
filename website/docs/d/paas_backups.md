---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_backups"
description: |-
  Provides a list of PaaS service backups IDs.
---

# Data Source: aws_paas_backups

Provides a list of PaaS service backups IDs matching the specified criteria.

## Example Usage

```terraform
data "aws_paas_backups" "selected" {
  service_id = "fm-cluster-12345678"
}

data "aws_paas_backup" "backups" {
  for_each = data.aws_paas_backups.selected.backup_ids
  id       = each.key
}
```

## Argument Reference

The following arguments are supported:

* `service_class` - (Optional) The class of the PaaS service.
    * _Valid values:_ `cacher`, `database`, `message_broker`, `search`
* `service_id` - (Optional) The ID of the PaaS service (e.g. `fm-cluster-12345678`).
* `service_type` - (Optional) The type of the PaaS service.
    * _Valid values:_ `elasticsearch`, `memcached`, `mongodb`, `mysql`, `pgsql`, `rabbitmq`, `redis`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `backup_ids` - List of backup IDs.
* `id` - The region.
