---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_connection"
description: |-
  Provides information about a Direct Connect connection.
---

# Data Source: aws_dx_connection

Provides information about a Direct Connect connection.

## Example Usage

```terraform
data "aws_dx_connection" "selected" {
  name = "tf-dxconn-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the connection.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the connection.
* `aws_device` - The ID of the device to which the connection is attached.
* `bandwidth` - The bandwidth of the connection.
* `id` - The ID of the connection.
* `location` - The physical site, where the connection terminates.
* `owner_account_id` - The ID of the project that owns the connection.
* `tags` - Map of tags assigned to the connection.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `provider_name`.
