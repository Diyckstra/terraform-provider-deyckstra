---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_lag"
description: |-
  Provides information about a Direct Connect link aggregation group (LAG).
---

# Data Source: aws_dx_lag

Provides information about a Direct Connect link aggregation group (LAG).

## Example Usage

```terraform
data "aws_dx_lag" "selected" {
  name = "tf-dxlag-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the LAG.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the LAG.
* `aws_device` - The ID of the device to which the LAG is attached.
* `bandwidth` - The bandwidth of each physical connection in the LAG.
* `id` - The ID of the LAG.
* `location` - The physical site, where the connection terminates.
* `owner_account_id` - The ID of the project that owns the LAG.
* `tags` - Map of tags assigned to the LAG.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `provider_name`.
