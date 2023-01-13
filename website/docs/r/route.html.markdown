---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_route"
description: |-
  Provides a resource to create a routing entry in a VPC routing table.
---

# Resource: aws_route

Provides a resource to create a routing table entry (a route) in a VPC routing table.

~> **NOTE on Route Tables and Routes:** This provider currently provides both a standalone Route resource and a Route Table resource with routes defined in-line. At this time you cannot use a Route Table with in-line routes in conjunction with any Route resources. Doing so will cause a conflict of rule settings and will overwrite rules.

~> **NOTE on `gateway_id` attribute:** The AWS API is very forgiving with the resource ID passed in the `gateway_id` attribute. For example an `aws_route` resource can be created with an `aws_nat_gateway` or `aws_egress_only_internet_gateway` ID specified for the `gateway_id` attribute. Specifying anything other than an `aws_internet_gateway` or `aws_vpn_gateway` ID will lead to this provider reporting a permanent diff between your configuration and recorded state, as the AWS API returns the more-specific attribute. If you are experiencing constant diffs with an `aws_route` resource, the first thing to check is that the correct attribute is being specified.

## Example Usage

```terraform
resource "aws_route" "r" {
  route_table_id            = "rtb-4fbb3ac4"
  destination_cidr_block    = "10.0.1.0/22"
  vpc_peering_connection_id = "pcx-45ff3dc1"
  depends_on                = [aws_route_table.testing]
}
```

## Example IPv6 Usage

```terraform
resource "aws_vpc" "vpc" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true
}

resource "aws_egress_only_internet_gateway" "egress" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_route" "r" {
  route_table_id              = "rtb-4fbb3ac4"
  destination_ipv6_cidr_block = "::/0"
  egress_only_gateway_id      = aws_egress_only_internet_gateway.egress.id
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required) The ID of the routing table.

One of the following destination arguments must be supplied:

* `destination_cidr_block` - (Optional) The destination CIDR block.
* `destination_ipv6_cidr_block` - (Optional) The destination IPv6 CIDR block.
* `destination_prefix_list_id` - (Optional) The ID of a managed prefix list destination.

One of the following target arguments must be supplied:

* `carrier_gateway_id` - (Optional) Identifier of a carrier gateway. This attribute can only be used when the VPC contains a subnet which is associated with a Wavelength Zone.
* `core_network_arn` - (Optional) The Amazon Resource Name (ARN) of a core network.
* `egress_only_gateway_id` - (Optional) Identifier of a VPC Egress Only Internet Gateway.
* `gateway_id` - (Optional) Identifier of a VPC internet gateway or a virtual private gateway.
* `instance_id` - (Optional, **Deprecated** use `network_interface_id` instead) Identifier of an EC2 instance.
* `nat_gateway_id` - (Optional) Identifier of a VPC NAT gateway.
* `local_gateway_id` - (Optional) Identifier of a Outpost local gateway.
* `network_interface_id` - (Optional) Identifier of an EC2 network interface.
* `transit_gateway_id` - (Optional) Identifier of an EC2 Transit Gateway.
* `vpc_endpoint_id` - (Optional) Identifier of a VPC Endpoint.
* `vpc_peering_connection_id` - (Optional) Identifier of a VPC peering connection.

Note that the default route, mapping the VPC's CIDR block to "local", is created implicitly and cannot be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

~> **NOTE:** Only the arguments that are configured (one of the above) will be exported as an attribute once the resource is created.

* `id` - Route identifier computed from the routing table identifier and route destination.
* `instance_owner_id` - The AWS account ID of the owner of the EC2 instance.
* `origin` - How the route was created - `CreateRouteTable`, `CreateRoute` or `EnableVgwRoutePropagation`.
* `state` - The state of the route - `active` or `blackhole`.

## Timeouts

Configuration options:

- `create` - (Default `5m`)
- `update` - (Default `2m`)
- `delete` - (Default `5m`)

## Import

Individual routes can be imported using `ROUTETABLEID_DESTINATION`.

For example, import a route in route table `rtb-656C65616E6F72` with an IPv4 destination CIDR of `10.42.0.0/16` like this:

```console
$ terraform import aws_route.my_route rtb-656C65616E6F72_10.42.0.0/16
```

Import a route in route table `rtb-656C65616E6F72` with an IPv6 destination CIDR of `2620:0:2d0:200::8/125` similarly:

```console
$ terraform import aws_route.my_route rtb-656C65616E6F72_2620:0:2d0:200::8/125
```

Import a route in route table `rtb-656C65616E6F72` with a managed prefix list destination of `pl-0570a1d2d725c16be` similarly:

```console
$ terraform import aws_route.my_route rtb-656C65616E6F72_pl-0570a1d2d725c16be
```
