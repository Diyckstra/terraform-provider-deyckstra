---
subcategory: "CloudWatch"
layout: "aws"
page_title: "aws_cloudwatch_metric_alarm"
description: |-
  Manages a CloudWatch metric alarm.
---

[dimensions]: https://docs.k2.cloud/en/services/monitoring/metrics.html#dimensions
[metrics]: https://docs.k2.cloud/en/services/monitoring/metrics.html

# Resource: aws_cloudwatch_metric_alarm

Manages a CloudWatch metric alarm.

## Example Usage

```terraform
variable instance_id {}

resource "aws_cloudwatch_metric_alarm" "example" {
  alarm_name          = "terraform-test-metric-alarm"
  alarm_description   = "This metric monitors EC2 CPU utilization"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  statistic           = "Average"
  evaluation_periods  = 2
  period              = 120
  threshold           = 80
  alarm_actions       = ["example@mail.com:EMAIL"]
  dimensions = {
    InstanceId = var.instance_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `actions_enabled` - (Optional) Indicates whether actions should be executed during any changes to the alarm state.
    * _Default value:_ `true`
* `alarm_actions` - (Optional) Actions, which must be executed when this alarm transitions to the `alarm` state. Each action must be between 1-1024 characters in length. You can specify a maximum of 5 actions.
* `alarm_description` - (Optional) The alarm description. Must be between 1-255 characters in length.
* `alarm_name` - (Required) The name for the alarm. This name must be unique within the project. Must be between 1-255 characters in length.
* `comparison_operator` - (Required) The arithmetic operation to use when comparing the specified `statistic` and `threshold`.
    * _Valid values:_ `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanThreshold`, `LessThanOrEqualToThreshold`
* `datapoints_to_alarm` - (Optional) The number of datapoints that must be breaching to trigger the alarm. Minimum value is 1.
* `dimensions` - (Required) The alarm dimensions. See docs for [dimensions][dimensions].
* `evaluation_periods` - (Required) The number of periods which is compared to the threshold. Minimum value is 1.
* `insufficient_data_actions` - (Optional) Actions, which must be executed when this alarm transitions to the `insufficient_data` state. Each action must be between 1-1024 characters in length. You can specify a maximum of 5 actions.
* `metric_name` - (Required) The name of the metric that associated with the alarm. Must be between 1-255 characters in length. See docs for [supported metrics][metrics].
* `namespace` - (Required) The namespace of the metric with which the alarm is associated. Must be between 1-255 characters in length. See docs for the [list of namespaces and supported metrics][metrics].
* `ok_actions` - (Optional) Actions, which must be executed when this alarm transitions to the `ok` state. Each action must be between 1-1024 characters in length. You can specify a maximum of 5 actions.
* `period` - (Required) The period in seconds over which the specified `statistic` is applied. Value must be divisible by 60, minimum value is 60.
* `statistic` - (Required) The statistic for the metric.
    * _Valid values:_ `SampleCount`, `Average`, `Sum`, `Minimum`, `Maximum`
* `threshold` - (Required) The value, to which metric values will be compared.
* `treat_missing_data` - (Optional) Defines how periods without values would be interpreted.
    * _Valid values:_ `missing`, `ignore`, `breaching`, `not_breaching`
    * _Default value:_ `missing`
* `unit` - (Optional) The unit of the metric associated with the alarm.
    * _Valid values:_ `Percent`, `Bytes`, `Count`

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the alarm.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn`, `evaluate_low_sample_count_percentiles`, `extended_statistic`, `metric_query`, `threshold_metric_id`.

## Import

CloudWatch metric alarm can be imported using the `alarm_name`, e.g.,

```
$ terraform import aws_cloudwatch_metric_alarm.example terraform-test-metric-alarm
```
