---
subcategory: "Outposts (EC2)"
layout: "aws"
page_title: "AWS: aws_ec2_local_gateway_route_table_vpc_association"
description: |-
  Manages an EC2 Local Gateway Route Table VPC Association
---

# Resource: aws_ec2_local_gateway_route_table_vpc_association

Manages an EC2 Local Gateway Route Table VPC Association. More information can be found in the [Outposts User Guide](https://docs.aws.amazon.com/outposts/latest/userguide/outposts-local-gateways.html#vpc-associations).

## Example Usage

```terraform
data "aws_ec2_local_gateway_route_table" "example" {
  outpost_arn = "arn:aws:outposts:us-west-2:123456789012:outpost/op-1234567890abcdef"
}

resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_ec2_local_gateway_route_table_vpc_association" "example" {
  local_gateway_route_table_id = data.aws_ec2_local_gateway_route_table.example.id
  vpc_id                       = aws_vpc.example.id
}
```

## Argument Reference

The following arguments are required:

* `local_gateway_route_table_id` - (Required) Identifier of EC2 Local Gateway Route Table.
* `vpc_id` - (Required) Identifier of EC2 VPC.

The following arguments are optional:

* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Identifier of EC2 Local Gateway Route Table VPC Association.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

`aws_ec2_local_gateway_route_table_vpc_association` can be imported by using the Local Gateway Route Table VPC Association identifier, e.g.,

```
$ terraform import aws_ec2_local_gateway_route_table_vpc_association.example lgw-vpc-assoc-1234567890abcdef
```
