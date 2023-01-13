---
subcategory: "Direct Connect"
layout: "aws"
page_title: "AWS: aws_dx_bgp_peer"
description: |-
  Provides a Direct Connect BGP peer resource.
---

# Resource: aws_dx_bgp_peer

Provides a Direct Connect BGP peer resource.

## Example Usage

```terraform
resource "aws_dx_bgp_peer" "peer" {
  virtual_interface_id = aws_dx_private_virtual_interface.foo.id
  address_family       = "ipv6"
  bgp_asn              = 65351
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Required) The address family for the BGP peer. `ipv4 ` or `ipv6`.
* `bgp_asn` - (Required) The autonomous system (AS) number for Border Gateway Protocol (BGP) configuration.
* `virtual_interface_id` - (Required) The ID of the Direct Connect virtual interface on which to create the BGP peer.
* `amazon_address` - (Optional) The IPv4 CIDR address to use to send traffic to Amazon.
Required for IPv4 BGP peers on public virtual interfaces.
* `bgp_auth_key` - (Optional) The authentication key for BGP configuration.
* `customer_address` - (Optional) The IPv4 CIDR destination address to which Amazon should send traffic.
Required for IPv4 BGP peers on public virtual interfaces.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the BGP peer resource.
* `bgp_status` - The Up/Down state of the BGP peer.
* `bgp_peer_id` - The ID of the BGP peer.
* `aws_device` - The Direct Connect endpoint on which the BGP peer terminates.

## Timeouts

Configuration options:

- `create` - (Default `10m`)
- `delete` - (Default `10m`)
