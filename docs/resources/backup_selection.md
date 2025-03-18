---
subcategory: "Backup"
layout: "aws"
page_title: "aws_backup_selection"
description: |-
  Manages selection conditions for backup plan resources.
---

# Resource: aws_backup_selection

Manages selection conditions for backup plan resources.

## Example Usage

### Selecting Backups By Resource

```terraform
resource "aws_backup_vault_default" "example" {}

data "aws_availability_zones" "this" {
  state = "available"
}

resource "aws_ebs_volume" "example" {
  availability_zone = data.aws_availability_zones.this.names[0]
  type              = "gp2"
  size              = 8
}

resource "aws_backup_plan" "example" {
  name = "tf-backup-plan"

  rule {
    rule_name         = "tf-backup-rule"
    target_vault_name = aws_backup_vault_default.example.name
    schedule          = "cron(0 12 * * *)"
    start_window      = 60
    completion_window = 180

    lifecycle {
      delete_after = 30
    }
  }

  rule {
    rule_name         = "tf-backup-rule2"
    target_vault_name = aws_backup_vault_default.example.name
    schedule          = "cron(0 23 * * *)"
    start_window      = 60
    completion_window = 180

    lifecycle {
      delete_after = 60
    }
  }
}

resource "aws_backup_selection" "example" {
  name    = "tf-backup-selection"
  plan_id = aws_backup_plan.example.id

  resources = [
    aws_ebs_volume.example.arn,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Editable) The selection name. The value must be `1` to `50` characters long and must contain only alphanumeric, hyphen, underscore, or periods.
* `plan_id` - (Required) The ID of the backup plan.
* `resources` - (Required, Editable) The list of Amazon Resource Names (ARNs) of the resources included into the backup plan.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource selection.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`condition`, `iam_role_arn`, `not_resources`, `selection_tag`.

## Import

Backup selection can be imported using `plan_id` and `id` separated by `|`.

```
$ terraform import aws_backup_selection.example "01234567-0123-0123-0123-0123456789ab|54321687-1111-2222-3333-9876543210ab"
```
