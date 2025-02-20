---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_key_pair"
description: |-
  Manages a key pair.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_key_pair

Manages an EC2 key pair resource.
Currently, this resource requires an existing user-supplied key pair.
This key pair's public key will be registered to allow logging-in to EC2 instances.

When importing an existing key pair the public key material may be in any format supported by AWS.
Supported public key material formats are:

* OpenSSH public key format (the format in ~/.ssh/authorized_keys)
* Base64 encoded DER format
* SSH public key file format as specified in RFC4716

## Example Usage

```terraform
resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 email@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` – (Optional) The name for the key pair. If neither `key_name` nor `key_name_prefix` is provided, Terraform will create a unique key name using the prefix `terraform-`.
* `key_name_prefix` – (Optional) Creates a unique name beginning with the specified prefix. If neither `key_name` nor `key_name_prefix` is provided, Terraform will create a unique key name using the prefix `terraform-`.
    * _Constraints:_  Conflicts with `key_name`
* `public_key` – (Required) The public key material.
* `tags` – (Optional) Map of tags to assign to the resource. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` – The key pair name.
* `arn` – The key pair ARN.
* `key_name` – The key pair name.
* `key_pair_id` – The key pair ID.
* `fingerprint` – The MD5 public key fingerprint as specified in section 4 of RFC 4716.
* `tags_all` – Map of tags to assign to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

## Import

Key pairs can be imported using the `key_name`, e.g.,

```
$ terraform import aws_key_pair.deployer deployer-key
```
