---
subcategory: "Backup"
layout: "aws"
page_title: "aws_backup_selection"
description: |-
  Provides information about a backup selection.
---

# Data source: aws_backup_selection

Provides information about a backup selection.

## Example Usage

```terraform
data "aws_backup_plan" "example" {
  plan_id = "01234567-0123-0123-0123-0123456789ab"
}

data "aws_backup_selection" "example" {
  plan_id      = data.aws_backup_plan.example.id
  selection_id = "54321687-1111-2222-3333-9876543210ab"
}
```

## Argument Reference

The following arguments are supported:

* `plan_id` - (Required) The ID of the backup plan.
* `selection_id` - (Required) The ID of the resource selection.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `name` - The selection name.
* `resources` - The list of Amazon Resource Names (ARNs) of the resources included into the backup plan.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `iam_role_arn`.
