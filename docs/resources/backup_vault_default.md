---
subcategory: "Backup"
layout: "aws"
page_title: "aws_backup_vault_default"
description: |-
  Manages the default backup vault.
---

# Resource: aws_backup_vault_default

Manages the default backup vault. If it does not exist, Terraform will create one.

~> **Note** Terraform does not delete the default backup vault when running `terraform destroy`, it will only be removed from the Terraform state.

## Example Usage

```terraform
resource "aws_backup_vault_default" "example" {}
```

## Argument Reference

Currently no arguments are used for the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the backup vault.
* `id` - The ID of the backup vault.
* `name` - The name of the backup vault.
* `recovery_points` - The number of recovery points in the backup vault.


## Import

Backup vault can be imported using `name`, e.g.,

```
$ terraform import aws_backup_vault_default.example Default
```
