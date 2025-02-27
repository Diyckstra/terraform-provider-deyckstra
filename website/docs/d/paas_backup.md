---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_backup"
description: |-
  Provides information about a PaaS service backup.
---

[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Data Source: aws_paas_backup

Provides information about a PaaS service backup.

~> **Note** If more than one backup meets the specified criteria, the most recently created backup is returned.

## Example Usage

```terraform
data "aws_paas_backup" "selected" {
  service_class = "database"
}

output "most-recent-database-backup" {
  value = data.aws_paas_backup.selected.id
}
```

## Argument Reference

The following arguments are supported:

* `age_days` - (Optional) The age of the backup in days.
* `database_name` - (Optional) The database name.
* `id` - (Optional) The ID of the PaaS service backup (e.g. `paas-backup-12345678`).
* `ready_only` - (Optional) Indicates whether to filter only completed backups.
    * _Default value_: `true`
* `service_class` - (Optional) The class of the PaaS service.
    * _Valid values_: `cacher`, `database`, `message_broker`, `search`
* `service_id` - (Optional) The ID of the PaaS service (e.g. `fm-cluster-12345678`).
* `service_type` - (Optional) The type of the PaaS service.
    * _Valid values_: `elasticsearch`, `memcached`, `mongodb`, `mysql`, `pgsql`, `rabbitmq`, `redis`

~> **Note** `id` cannot be specified together with the other parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `databases` - List of databases. The structure of this block is [described below](#databases).
* `id` - The region if `id` is not specified as an argument.
* `protected` - Indicates whether the backup is protected from automatic scheduled deletion.
* `service_deleted` - Indicates whether the service is deleted.
* `service_name` - The service name.
* `status` - The current status of the backup creation process.
* `time` - The backup creation time in [RFC3339 format].

### databases

The `databases` block has the following structure:

* `backup_enabled` - Indicates whether backup is enabled for the database.
* `id` - The ID of the database.
* `location` - The link to the database backup in the bucket in object storage.
* `logfile` - The link to the database backup logfile in the bucket in object storage.
* `name` - The database name.
* `size` - The size of the database backup in bytes.
* `status` - The current status of the database backup creation process.
