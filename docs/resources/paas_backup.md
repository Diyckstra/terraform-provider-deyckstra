---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_backup"
description: |-
  Manages a PaaS service backup.
---

[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Resource: aws_paas_backup

Manages a PaaS service backup.

~> **Note** Terraform does not create backups of the PaaS service, but tries to gain control over an already created backup.
Use `backup_settings` section in [`aws_paas_service`](paas_service.md) resource to control the creation of backups.

~> **Note** Terraform does not delete the backup when running `terraform destroy`, it will only be removed from the Terraform state.

## Example Usage

```terraform
data "aws_paas_backups" "selected" {
  service_id = "fm-cluster-12345678"
}

resource "aws_paas_backup" "example" {
  for_each                   = data.aws_paas_backups.selected.backup_ids
  backup_id                  = each.key
  enable_deletion_protection = true
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required) The ID of the existing PaaS service backup (e.g. `paas-backup-12345678`).
* `enable_deletion_protection` - (Optional) Indicates whether the backup must be protected from automatic scheduled deletion.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `databases` - List of databases. The structure of this block is [described below](#databases).
* `id` - The ID of the PaaS service backup.
* `protected` - Indicates whether the backup is protected from automatic scheduled deletion.
* `service_class` - The class of the PaaS service.
* `service_deleted` - Indicates whether the service is deleted.
* `service_id` - The ID of the PaaS service.
* `service_name` - The service name.
* `service_type` - The type of the PaaS service.
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

## Import

PaaS service backup can be imported using `id`, e.g.,

```
$ terraform import aws_paas_backup.example paas-backup-12345678
```
