---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options"
description: |-
  Provides a VPC DHCP Options resource.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[dhcp-options]: https://docs.k2.cloud/en/services/networking/dhcpattrs.html
[tf-main-route-table-association]: main_route_table_association.html

# Resource: aws_vpc_dhcp_options

Provides a VPC DHCP options resource.

For more information, see the documentation on [DHCP options][dhcp-options].

## Example Usage

Basic usage:

```terraform
resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name_servers = ["8.8.8.8", "8.8.4.4"]
}
```

Full usage:

```terraform
resource "aws_vpc_dhcp_options" "foo" {
  domain_name          = "service.consul"
  domain_name_servers  = ["127.0.0.1", "10.0.0.2"]
  ntp_servers          = ["127.0.0.1"]
  netbios_name_servers = ["127.0.0.1"]
  netbios_node_type    = 2

  tags = {
    Name = "foo-name"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Optional) the suffix domain name to use by default when resolving non Fully Qualified Domain Names. In other words, this is what ends up being the `search` value in the `/etc/resolv.conf` file.
* `domain_name_servers` - (Optional) List of IP addresses of domain name servers or `AmazonProvidedDNS`. We recommend using only one of the two parameters.
* `ntp_servers` - (Optional) List of NTP servers to configure. You can specify up to four IP addresses.
* `netbios_name_servers` - (Optional) List of NETBIOS name servers. You can specify up to four IP addresses.
* `netbios_node_type` - (Optional) The NetBIOS node type (1, 2, 4, or 8). For more information about these node types, see [RFC 2132](http://www.ietf.org/rfc/rfc2132.txt).
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.

## Remarks

* Notice that all arguments are optional, but you have to specify at least one argument.
* `domain_name_servers`, `netbios_name_servers`, `ntp_servers` are limited to maximum four servers only.
* To actually use the DHCP options set you need to associate it to a VPC using [`aws_vpc_dhcp_options_association`][tf-main-route-table-association].
* If you delete a DHCP options set, all VPCs using it will be associated to `default` DHCP Option Set.
* In most cases unless you're configuring your own DNS you'll want to set `domain_name_servers` to `AmazonProvidedDNS`.

## Attributes Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the DHCP options Set.
* `arn` - ARN of the DHCP options Set.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported: `owner_id` .

## Import

VPC DHCP Options can be imported using the `dhcp options id`, e.g.,

```
$ terraform import aws_vpc_dhcp_options.my_options dopt-12345678
```

