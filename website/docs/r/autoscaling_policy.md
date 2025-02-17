---
subcategory: "Auto Scaling"
layout: "aws"
page_title: "aws_autoscaling_policy"
description: |-
  Manages an Auto Scaling policy resource.
---

# Resource: aws_autoscaling_policy

Manages an Auto Scaling policy resource.

~> **Note** You may want to omit `desired_capacity` attribute from attached `aws_autoscaling_group`
when using autoscaling policies. It's good practice to pick either manual or dynamic (policy-based) scaling.

## Example Usage

```terraform
resource "aws_autoscaling_policy" "example" {
  name                   = "terraform-test"
  scaling_adjustment     = 4
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = "example-asg" # asg is created manually
}
```

## Argument Reference

* `name` - (Required) The name of the policy.
* `autoscaling_group_name` - (Required) The name of the autoscaling group.
* `adjustment_type` - (Optional) Specifies whether the adjustment is an absolute number or a percentage of the current capacity. Valid values are `ChangeInCapacity`, `ExactCapacity`, and `PercentChangeInCapacity`.
* `policy_type` - (Optional) The policy type. Valid value is `SimpleScaling`.
* `min_adjustment_magnitude` - (Optional) Minimum value to scale by when `adjustment_type` is set to `PercentChangeInCapacity`.
* `cooldown` - (Optional) The amount of time, in seconds, after a scaling activity completes and before the next scaling activity can start.
* `scaling_adjustment` - (Optional) The amount by which the Auto Scaling Group is scaled when the scaling policy is executed

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the scaling policy.
* `id` - The scaling policy's name.
* `name` - The scaling policy's name.
* `autoscaling_group_name` - The scaling policy's assigned autoscaling group.
* `adjustment_type` - The scaling policy's adjustment type.
* `policy_type` - The scaling policy's type.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`estimated_instance_warmup`, `metric_aggregation_type`, `predictive_scaling_configuration`, `step_adjustment`, `target_tracking_configuration`.

## Import

AutoScaling scaling policy can be imported using the `autoscaling_group_name` and `name` separated by `/`.

```
$ terraform import aws_autoscaling_policy.test-policy asg-name/policy-name
```
