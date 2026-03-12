---
subcategory: "EFS (Elastic File System)"
layout: "aws"
page_title: "aws_efs_mount_target"
description: |-
  Provides information about an Elastic File System (EFS) mount target.
---

# Data Source: aws_efs_mount_target

Provides information about an Elastic File System (EFS) mount target.

## Example usage

```terraform
variable "mount_target_id" {
  type    = string
  default = ""
}

data "aws_efs_mount_target" "by_id" {
  mount_target_id = var.mount_target_id
}
```

## Argument reference

The following arguments are supported:

* `file_system_id` - (Optional) The ID of the file system whose mount target you want to find. It must be included if an `access_point_id` and `mount_target_id` are not included.
    * _Constraints:_ Must be included if a `mount_target_id` is not included
* `mount_target_id` - (Optional) The ID of the mount target that you want to find.

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone_id` - The unique and consistent identifier of the availability zone that the mount target resides in.
* `dns_name` - The DNS name for the EFS file system.
* `ip_address` - The address at which the file system may be mounted via the mount target.
* `network_interface_id` - The ID of the network interface that EFS created along with the mount target.
* `owner_id` - The ID of the account that owns the resource.
* `security_groups` - The list of the IDs of VPC's security groups attached to the mount target.
* `subnet_id` - The ID of the mount target's subnet.
