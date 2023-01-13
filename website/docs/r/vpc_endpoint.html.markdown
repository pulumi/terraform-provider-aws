---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_vpc_endpoint"
description: |-
  Provides a VPC Endpoint resource.
---

# Resource: aws_vpc_endpoint

Provides a VPC Endpoint resource.

~> **NOTE on VPC Endpoints and VPC Endpoint Associations:** The provider provides both standalone VPC Endpoint Associations for
Route Tables - (an association between a VPC endpoint and a single `route_table_id`),
Security Groups - (an association between a VPC endpoint and a single `security_group_id`),
and Subnets - (an association between a VPC endpoint and a single `subnet_id`) and
a VPC Endpoint resource with `route_table_ids` and `subnet_ids` attributes.
Do not use the same resource ID in both a VPC Endpoint resource and a VPC Endpoint Association resource.
Doing so will cause a conflict of associations and will overwrite the association.

## Example Usage

### Basic

```terraform
resource "aws_vpc_endpoint" "s3" {
  vpc_id       = aws_vpc.main.id
  service_name = "com.amazonaws.us-west-2.s3"
}
```

### Basic w/ Tags

```terraform
resource "aws_vpc_endpoint" "s3" {
  vpc_id       = aws_vpc.main.id
  service_name = "com.amazonaws.us-west-2.s3"

  tags = {
    Environment = "test"
  }
}
```

### Interface Endpoint Type

```terraform
resource "aws_vpc_endpoint" "ec2" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.us-west-2.ec2"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.sg1.id,
  ]

  private_dns_enabled = true
}
```

### Gateway Load Balancer Endpoint Type

```terraform
data "aws_caller_identity" "current" {}

resource "aws_vpc_endpoint_service" "example" {
  acceptance_required        = false
  allowed_principals         = [data.aws_caller_identity.current.arn]
  gateway_load_balancer_arns = [aws_lb.example.arn]
}

resource "aws_vpc_endpoint" "example" {
  service_name      = aws_vpc_endpoint_service.example.service_name
  subnet_ids        = [aws_subnet.example.id]
  vpc_endpoint_type = aws_vpc_endpoint_service.example.service_type
  vpc_id            = aws_vpc.example.id
}
```

### Non-AWS Service

```terraform
resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id            = var.vpc_id
  service_name      = var.ptfe_service
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.ptfe_service.id,
  ]

  subnet_ids          = [local.subnet_ids]
  private_dns_enabled = false
}

data "aws_route53_zone" "internal" {
  name         = "vpc.internal."
  private_zone = true
  vpc_id       = var.vpc_id
}

resource "aws_route53_record" "ptfe_service" {
  zone_id = data.aws_route53_zone.internal.zone_id
  name    = "ptfe.${data.aws_route53_zone.internal.name}"
  type    = "CNAME"
  ttl     = "300"
  records = [aws_vpc_endpoint.ptfe_service.dns_entry[0]["dns_name"]]
}
```

~> **NOTE The `dns_entry` output is a list of maps:** This provider interpolation support for lists of maps requires the `lookup` and `[]` until full support of lists of maps is available

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name. For AWS services the service name is usually in the form `com.amazonaws.<region>.<service>` (the SageMaker Notebook service is an exception to this rule, the service name is in the form `aws.sagemaker.<region>.notebook`).
* `vpc_id` - (Required) The ID of the VPC in which the endpoint will be used.
* `auto_accept` - (Optional) Accept the VPC endpoint (the VPC endpoint and service need to be in the same AWS account).
* `policy` - (Optional) A policy to attach to the endpoint that controls access to the service. This is a JSON formatted string. Defaults to full access. All `Gateway` and some `Interface` endpoints support policies - see the [relevant AWS documentation](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-endpoints-access.html) for more details.
* `private_dns_enabled` - (Optional; AWS services and AWS Marketplace partner services only) Whether or not to associate a private hosted zone with the specified VPC. Applicable for endpoints of type `Interface`.
Defaults to `false`.
* `dns_options` - (Optional) The DNS options for the endpoint. See dns_options below.
* `ip_address_type` - (Optional) The IP address type for the endpoint. Valid values are `ipv4`, `dualstack`, and `ipv6`.
* `route_table_ids` - (Optional) One or more route table IDs. Applicable for endpoints of type `Gateway`.
* `subnet_ids` - (Optional) The ID of one or more subnets in which to create a network interface for the endpoint. Applicable for endpoints of type `GatewayLoadBalancer` and `Interface`.
* `security_group_ids` - (Optional) The ID of one or more security groups to associate with the network interface. Applicable for endpoints of type `Interface`.
If no security groups are specified, the VPC's [default security group](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html#DefaultSecurityGroup) is associated with the endpoint.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `vpc_endpoint_type` - (Optional) The VPC endpoint type, `Gateway`, `GatewayLoadBalancer`, or `Interface`. Defaults to `Gateway`.

### dns_options

* `dns_record_ip_type` - (Optional) The DNS records created for the endpoint. Valid values are `ipv4`, `dualstack`, `service-defined`, and `ipv6`.

## Timeouts

Configuration options:

- `create` - (Default `10m`)
- `update` - (Default `10m`)
- `delete` - (Default `10m`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPC endpoint.
* `arn` - The Amazon Resource Name (ARN) of the VPC endpoint.
* `cidr_blocks` - The list of CIDR blocks for the exposed AWS service. Applicable for endpoints of type `Gateway`.
* `dns_entry` - The DNS entries for the VPC Endpoint. Applicable for endpoints of type `Interface`. DNS blocks are documented below.
* `network_interface_ids` - One or more network interfaces for the VPC Endpoint. Applicable for endpoints of type `Interface`.
* `owner_id` - The ID of the AWS account that owns the VPC endpoint.
* `prefix_list_id` - The prefix list ID of the exposed AWS service. Applicable for endpoints of type `Gateway`.
* `requester_managed` -  Whether or not the VPC Endpoint is being managed by its service - `true` or `false`.
* `state` - The state of the VPC endpoint.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

DNS blocks (for `dns_entry`) support the following attributes:

* `dns_name` - The DNS name.
* `hosted_zone_id` - The ID of the private hosted zone.

## Import

VPC Endpoints can be imported using the `vpc endpoint id`, e.g.,

```
$ terraform import aws_vpc_endpoint.endpoint1 vpce-3ecf2a57
```
