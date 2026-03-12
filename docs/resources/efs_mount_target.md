---
subcategory: "EFS (Elastic File System)"
layout: "aws"
page_title: "aws_efs_mount_target"
description: |-
  Creates an Elastic File System (EFS) mount target.
---

[mount-target]: https://docs.k2.cloud/ru/services/efs/efs.html#id19

# Resource: aws_efs_mount_target

Creates an Elastic File System (EFS) mount target.

## Example usage

```terraform
resource "aws_efs_file_system" "example" {
  creation_token = "my-product"

  tags = {
    Name = "MyProduct"
  }
}

resource "aws_efs_mount_target" "example" {
  file_system_id = aws_efs_file_system.example.id
  subnet_id      = aws_subnet.example.id
}

resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  availability_zone = "ru-msk-vol51"
  cidr_block        = "10.0.1.0/24"
}
```

## Argument reference

The following arguments are supported:

* `file_system_id` - (Required) The ID of the file system for which the mount target is intended.
* `subnet_id` - (Required) The ID of the subnet to add the mount target in.
* `security_groups` - (Optional) A list of up to 5 security group IDs in effect for the mount target.
The security groups must be for the same VPC as the subnet specified.

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone_id` - The unique and consistent identifier of the availability zone that the mount target resides in.
* `dns_name` - The DNS name for the EFS file system.
* `id` - The ID of the mount target.
* `network_interface_id` - The ID of the network interface that EFS created along with the mount target.
* `owner_id` - The ID of the account that owns the resource.

## Import

The EFS mount targets can be imported using `id`, for example:

```
$ terraform import aws_efs_mount_target.alpha fsmt-52a643fb
```
