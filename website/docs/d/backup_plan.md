---
subcategory: "Backup"
layout: "aws"
page_title: "aws_backup_plan"
description: |-
  Provides information about a backup plan.
---

# Data Source: aws_backup_plan

Provides information about a backup plan.

## Example Usage

```terraform
data "aws_backup_plan" "example" {
  plan_id = "01234567-0123-0123-0123-0123456789ab"
}
```

## Argument Reference

The following arguments are supported:

* `plan_id` - (Required) The ID of the backup plan.

## Attribute Reference

### Supported attributes

In addition to the arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the backup plan.
* `name` - The name of the backup plan.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `version`.
