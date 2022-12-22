---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "AWS: aws_customer_gateway"
description: |-
  Get an existing AWS Customer Gateway.
---

# Data Source: aws_customer_gateway

Get an existing AWS Customer Gateway.

## Example Usage

```terraform
data "aws_customer_gateway" "foo" {
  filter {
    name   = "tag:Name"
    values = ["foo-prod"]
  }
}

resource "aws_vpn_gateway" "main" {
  vpc_id          = aws_vpc.main.id
  amazon_side_asn = 7224
}

resource "aws_vpn_connection" "transit" {
  vpn_gateway_id      = aws_vpn_gateway.main.id
  customer_gateway_id = data.aws_customer_gateway.foo.id
  type                = data.aws_customer_gateway.foo.type
  static_routes_only  = false
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the gateway.
* `filter` - (Optional) One or more [name-value pairs][dcg-filters] to filter by.

[dcg-filters]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeCustomerGateways.html

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `arn` - ARN of the customer gateway.
* `bgp_asn` - Gateway's Border Gateway Protocol (BGP) Autonomous System Number (ASN).
* `certificate_arn` - ARN for the customer gateway certificate.
* `device_name` - Name for the customer gateway device.
* `ip_address` - IP address of the gateway's Internet-routable external interface.
* `tags` - Map of key-value pairs assigned to the gateway.
* `type` - Type of customer gateway. The only type AWS supports at this time is "ipsec.1".

## Timeouts

Configuration options:

- `read` - (Default `20m`)
