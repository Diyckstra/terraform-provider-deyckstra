---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami_ids"
description: |-
  Provides a list of image IDs.
---

[describe-images]: https://docs.cloud.croc.ru/en/api/ec2/images/DescribeImages.html

# Data Source: aws_ami_ids

Provides a list of image IDs matching the specified criteria.

## Example Usage

```terraform
data "aws_ami_ids" "example" {
  owners = ["self"]
}
```

## Argument Reference

### Required Arguments

* `owners` – (Required) List of image owners to limit search. At least one value must be specified.
    * _Valid values_: Project ID (`project@customer`) or `self`
* `executable_users` – (Optional) Limit search to project with *explicit* launch permission on the image.
    * _Valid values_: Project ID (`project@customer`), `all` or `self`
* `filter` – (Optional) One or more name/value pairs to use as filters.
    * _Valid values_: See valid names and values in [EC2 API documentation][describe-images]
* `name_regex` – (Optional) A regex string to apply to the image list returned by the EC2 API.
  It is recommended to combine this with other options to narrow down the list the EC2 API returns.
* `sort_ascending` – (Optional) Used to sort images by creation time.
    * _Default value_: `false`

## Attributes Reference

`ids` is set to the list of image IDs, sorted by creation time according to `sort_ascending`.
