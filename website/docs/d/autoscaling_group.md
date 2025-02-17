---
subcategory: "Auto Scaling"
layout: "aws"
page_title: "aws_autoscaling_group"
description: |-
  Provides information about an Auto Scaling Group.
---

# Data Source: aws_autoscaling_group

Provides information about an existing Auto Scaling Group.

## Example Usage

```terraform
data "aws_autoscaling_group" "example" {
  name = "example-asg"
}
```

## Argument Reference

* `name` - Specify the exact name of the desired Auto Scaling Group.

## Attributes Reference

### Supported attributes

~> **Note** Some values are not always set and may not be available for interpolation.

* `arn` - The Amazon Resource Name (ARN) of the Auto Scaling Group.
* `availability_zones` - One or more Availability Zones for the group.
* `default_cool_down` - The amount of time, in seconds, after a scaling activity completes before another scaling activity can start.
* `desired_capacity` - The desired size of the group.
* `health_check_grace_period` - The amount of time, in seconds, after which Auto Scaling Group can perform a health check on its instances.
* `id` - Name of the Auto Scaling Group.
* `max_size` - The maximum size of the group.
* `min_size` - The minimum size of the group.
* `name` - Name of the Auto Scaling Group.
* `new_instances_protected_from_scale_in` - Indicates whether new instances are protected from deletion when Auto Scaling Group is scaled in.
* `status` -  The status of the Auto Scaling Group when it is deleted.
* `vpc_zone_identifier` - The IDs of the subnets in which instances are created.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enabled_metrics`, `health_check_type`, `launch_configuration`, `load_balancers`, `placement_group`, `service_linked_role_arn`, `target_group_arns`, `termination_policies`.
