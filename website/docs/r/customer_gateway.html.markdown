---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "AWS: aws_customer_gateway"
description: |-
  Provides a customer gateway inside a VPC. These objects can be
  connected to VPN gateways via VPN connections, and allow you to
  establish tunnels between your network and the VPC.
---

# Resource: aws_customer_gateway

Provides a customer gateway inside a VPC. These objects can be connected to VPN gateways via VPN connections, and allow you to establish tunnels between your network and the VPC.

## Example Usage

```terraform
resource "aws_customer_gateway" "main" {
  bgp_asn    = 65000
  ip_address = "172.83.124.10"
  type       = "ipsec.1"

  tags = {
    Name = "main-customer-gateway"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bgp_asn` - (Required) The gateway's Border Gateway Protocol (BGP) Autonomous System Number (ASN).
* `certificate_arn` - (Optional) The Amazon Resource Name (ARN) for the customer gateway certificate.
* `device_name` - (Optional) A name for the customer gateway device.
* `ip_address` - (Optional) The IPv4 address for the customer gateway device's outside interface.
* `type` - (Required) The type of customer gateway. The only type AWS
  supports at this time is "ipsec.1".
* `tags` - (Optional) Tags to apply to the gateway. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The amazon-assigned ID of the gateway.
* `arn` - The ARN of the customer gateway.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Customer Gateways can be imported using the `id`, e.g.,

```
$ terraform import aws_customer_gateway.main cgw-b4dc3961
```
