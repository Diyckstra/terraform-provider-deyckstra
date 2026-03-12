---
subcategory: "EFS (Elastic File System)"
layout: "aws"
page_title: "aws_efs_file_system"
description: |-
  Provides information about an Elastic File System (EFS) file system.
---

# Data Source: aws_efs_file_system

Provides information about an Elastic File System (EFS) file system.

## Example usage

```terraform
variable "file_system_id" {
  type    = string
  default = ""
}

data "aws_efs_file_system" "by_id" {
  file_system_id = var.file_system_id
}

data "aws_efs_file_system" "by_tag" {
  tags = {
    Environment = "dev"
  }
}
```

## Argument reference

The following arguments are supported:

* `creation_token` - (Optional) Restricts the list to the file system with this creation token.
* `file_system_id` - (Optional) The ID of the file system.
        * _Example:_ `fs-ccfc0d65`
* `tags` - (Optional) Restricts the list to the file systems with these tags.

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `performance_mode` - The file system performance mode.
* `size_in_bytes` - The current byte count used by the file system.
* `tags` - A map of tags to assign to the file system.
