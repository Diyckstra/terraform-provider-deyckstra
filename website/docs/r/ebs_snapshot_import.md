---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot_import"
description: |-
  Imports a disk image from S3 as a snapshot.
---

[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_ebs_snapshot_import

Imports a disk image from S3 as a snapshot.

## Example Usage

```terraform
resource "aws_ebs_snapshot_import" "example" {
  disk_container {
    format = "VHD"
    user_bucket {
      s3_bucket = "disk-images"
      s3_key    = "source.vhd"
    }
  }

  tags = {
    Name = "HelloWorld"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description string for the import snapshot task.
* `disk_container` - (Required) Information about the disk container. Detailed below.
* `tags` - (Optional) Map of tags to assign to the snapshot.

### disk_container Configuration Block

* `description` - (Optional) The description of the disk image being imported.
* `format` - (Required) The format of the disk image being imported. One of `VHD`, `VMDK` or `RAW`.
* `user_bucket` - (Required) The S3 bucket for the disk image.

### user_bucket Configuration Block

* `s3_bucket` - The name of the S3 bucket where the disk image is located.
* `s3_key` - The file name of the disk image.

### Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `60 minutes`) Used for importing the EBS snapshot
- `delete` - (Default `10 minutes`) Used for deleting the EBS snapshot

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the EBS snapshot.
* `id` - The snapshot ID (e.g., snap-12345678).
* `owner_id` - The project ID.
* `owner_alias` - The alias of the EBS snapshot owner.
* `volume_size` - The size of the drive in GiB.
* `tags_all` - Map of tags assigned to the resource.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`client_data`, `data_encryption_key_id`, `disk_container.url`, `encrypted`, `kms_key_id`, `outpost_arn`, `permanent_restore`, `role_name`, `storage_tier`, `temporary_restore_days`.
