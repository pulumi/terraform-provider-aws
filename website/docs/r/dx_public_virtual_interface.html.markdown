---
subcategory: "Direct Connect"
layout: "aws"
page_title: "AWS: aws_dx_public_virtual_interface"
description: |-
  Provides a Direct Connect public virtual interface resource.
---

# Resource: aws_dx_public_virtual_interface

Provides a Direct Connect public virtual interface resource.

## Example Usage

```terraform
resource "aws_dx_public_virtual_interface" "foo" {
  connection_id = "dxcon-zzzzzzzz"

  name           = "vif-foo"
  vlan           = 4094
  address_family = "ipv4"
  bgp_asn        = 65352

  customer_address = "175.45.176.1/30"
  amazon_address   = "175.45.176.2/30"

  route_filter_prefixes = [
    "210.52.109.0/24",
    "175.45.176.0/22",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Required) The address family for the BGP peer. `ipv4 ` or `ipv6`.
* `bgp_asn` - (Required) The autonomous system (AS) number for Border Gateway Protocol (BGP) configuration.
* `connection_id` - (Required) The ID of the Direct Connect connection (or LAG) on which to create the virtual interface.
* `name` - (Required) The name for the virtual interface.
* `vlan` - (Required) The VLAN ID.
* `amazon_address` - (Optional) The IPv4 CIDR address to use to send traffic to Amazon. Required for IPv4 BGP peers.
* `bgp_auth_key` - (Optional) The authentication key for BGP configuration.
* `customer_address` - (Optional) The IPv4 CIDR destination address to which Amazon should send traffic. Required for IPv4 BGP peers.
* `route_filter_prefixes` - (Required) A list of routes to be advertised to the AWS network in this region.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the virtual interface.
* `arn` - The ARN of the virtual interface.
* `aws_device` - The Direct Connect endpoint on which the virtual interface terminates.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Timeouts

Configuration options:

- `create` - (Default `10m`)
- `delete` - (Default `10m`)

## Import

Direct Connect public virtual interfaces can be imported using the `vif id`, e.g.,

```
$ terraform import aws_dx_public_virtual_interface.test dxvif-33cc44dd
```
