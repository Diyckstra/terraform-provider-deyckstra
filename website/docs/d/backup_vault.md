---
subcategory: "Backup"
layout: "aws"
page_title: "aws_backup_vault"
description: |-
  Provides information about a backup vault.
---

# Data Source: aws_backup_vault

Provides information about a backup vault.

## Example Usage

```terraform
data "aws_backup_vault" "example" {
  name = "Default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the backup vault.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the backup vault.
* `recovery_points` - The number of recovery points in the vault.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `kms_key_arn`.
