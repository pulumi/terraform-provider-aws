---
subcategory: "Direct Connect"
layout: "aws"
page_title: "AWS: aws_dx_hosted_transit_virtual_interface"
description: |-
  Provides a Direct Connect hosted transit virtual interface resource.
---

# Resource: aws_dx_hosted_transit_virtual_interface

Provides a Direct Connect hosted transit virtual interface resource.
This resource represents the allocator's side of the hosted virtual interface.
A hosted virtual interface is a virtual interface that is owned by another AWS account.

## Example Usage

```terraform
resource "aws_dx_hosted_transit_virtual_interface" "example" {
  connection_id = aws_dx_connection.example.id

  name           = "tf-transit-vif-example"
  vlan           = 4094
  address_family = "ipv4"
  bgp_asn        = 65352
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Required) The address family for the BGP peer. `ipv4 ` or `ipv6`.
* `bgp_asn` - (Required) The autonomous system (AS) number for Border Gateway Protocol (BGP) configuration.
* `connection_id` - (Required) The ID of the Direct Connect connection (or LAG) on which to create the virtual interface.
* `name` - (Required) The name for the virtual interface.
* `owner_account_id` - (Required) The AWS account that will own the new virtual interface.
* `vlan` - (Required) The VLAN ID.
* `amazon_address` - (Optional) The IPv4 CIDR address to use to send traffic to Amazon. Required for IPv4 BGP peers.
* `bgp_auth_key` - (Optional) The authentication key for BGP configuration.
* `customer_address` - (Optional) The IPv4 CIDR destination address to which Amazon should send traffic. Required for IPv4 BGP peers.
* `mtu` - (Optional) The maximum transmission unit (MTU) is the size, in bytes, of the largest permissible packet that can be passed over the connection. The MTU of a virtual transit interface can be either `1500` or `8500` (jumbo frames). Default is `1500`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the virtual interface.
* `arn` - The ARN of the virtual interface.
* `aws_device` - The Direct Connect endpoint on which the virtual interface terminates.
* `jumbo_frame_capable` - Indicates whether jumbo frames (8500 MTU) are supported.

## Timeouts

Configuration options:

- `create` - (Default `10m`)
- `update` - (Default `10m`)
- `delete` - (Default `10m`)

## Import

Direct Connect hosted transit virtual interfaces can be imported using the `vif id`, e.g.,

```
$ terraform import aws_dx_hosted_transit_virtual_interface.test dxvif-33cc44dd
```
