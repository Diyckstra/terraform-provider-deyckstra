---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options"
description: |-
  Provides information about an EC2 DHCP Options configuration
---

# Data Source: aws_vpc_dhcp_options

Provides information about an EC2 DHCP options configuration.

## Example Usage

### Lookup by DHCP options ID

```terraform
variable dopts_id {}

data "aws_vpc_dhcp_options" "example" {
  dhcp_options_id = var.dopts_id
}
```

### Lookup by Filter

```terraform
data "aws_vpc_dhcp_options" "example" {
  filter {
    name   = "key"
    values = ["domain-name"]
  }

  filter {
    name   = "value"
    values = ["example.com"]
  }
}
```

## Argument Reference

* `dhcp_options_id` - (Optional) The EC2 DHCP options ID.
* `filter` - (Optional) List of custom filters as described below.

### filter

For more information about filtering, see the [EC2 API documentation][describe-dhcp-options].

* `name` - (Required) The name of the field to filter.
* `values` - (Required) Set of values for filtering.

## Attributes Reference

### Supported attributes

* `arn` - The ARN of the DHCP options Set.
* `dhcp_options_id` - EC2 DHCP options ID.
* `domain_name` - The suffix domain name to used when resolving non Fully Qualified Domain Names e.g., the `search` value in the `/etc/resolv.conf` file.
* `domain_name_servers` - List of name servers.
* `id` - EC2 DHCP options ID.
* `netbios_name_servers` - List of NETBIOS name servers.
* `netbios_node_type` - The NetBIOS node type (1, 2, 4, or 8). For more information about these node types, see [RFC 2132](http://www.ietf.org/rfc/rfc2132.txt).
* `ntp_servers` - List of NTP servers.
* `tags` - A map of tags assigned to the resource.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `owner_id`.

[describe-dhcp-options]: https://docs.k2.cloud/en/api/ec2/dhcp_options/DescribeDhcpOptions.html
