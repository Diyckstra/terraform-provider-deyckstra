---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface_attachment"
description: |-
  Attaches an elastic network interface (ENI) with an EC2 instance.
---

# Resource: aws_network_interface_attachment

Attaches an elastic network interface (ENI) with an EC2 instance.

## Example Usage

```terraform
resource "aws_network_interface_attachment" "test" {
  instance_id          = "i-12345678"
  network_interface_id = "eni-12345678"
  device_index         = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Instance ID to attach.
* `network_interface_id` - (Required) ENI ID to attach.
* `device_index` - (Required) Network interface index (int).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_id` - Instance ID.
* `network_interface_id` - Network interface ID.
* `attachment_id` - The ENI attachment ID.
* `status` - The status of the network interface attachment.
