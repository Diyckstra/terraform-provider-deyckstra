---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami_ids"
description: |-
  Provides a list of image IDs.
---

# Data Source: aws_ami_ids

Use this data source to get a list of image IDs matching the specified criteria.

## Example Usage

```terraform
data "aws_ami_ids" "example" {
  owners = ["self"]
}
```

## Argument Reference

* `owners` - (Required) List of image owners to limit search. At least one value must be specified.
  Valid items are the project ID (`project@customer`) or `self`.
* `executable_users` - (Optional) Limit search to project with *explicit* launch permission on the image.
  Valid items are the project ID (`project@customer`), `all` or `self`.

* `filter` - (Optional) One or more name/value pairs to filter.

For more information about filtering, see the [EC2 API documentation][describe-images].

* `name_regex` - (Optional) A regex string to apply to the image list returned by the EC2 API.
  This allows more advanced filtering. It is done locally on what the EC2 API returns,
  and could have a performance impact if the result is large.
  It is recommended to combine this with other options to narrow down the list the EC2 API returns.
* `sort_ascending`  - (Defaults to `false`) Used to sort images by creation time.

## Attributes Reference

`ids` is set to the list of image IDs, sorted by creation time according to `sort_ascending`.

[describe-images]: https://docs.k2.cloud/en/api/ec2/images/DescribeImages.html
