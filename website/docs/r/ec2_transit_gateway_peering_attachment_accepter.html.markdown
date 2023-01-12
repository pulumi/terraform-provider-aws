---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "AWS: aws_ec2_transit_gateway_peering_attachment_accepter"
description: |-
  Manages the accepter's side of an EC2 Transit Gateway peering Attachment
---

# Resource: aws_ec2_transit_gateway_peering_attachment_accepter

Manages the accepter's side of an EC2 Transit Gateway Peering Attachment.

## Example Usage

```terraform
resource "aws_ec2_transit_gateway_peering_attachment_accepter" "example" {
  transit_gateway_attachment_id = aws_ec2_transit_gateway_peering_attachment.example.id

  tags = {
    Name = "Example cross-account attachment"
  }
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_attachment_id` - (Required) The ID of the EC2 Transit Gateway Peering Attachment to manage.
* `tags` - (Optional) Key-value tags for the EC2 Transit Gateway Peering Attachment. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - EC2 Transit Gateway Attachment identifier
* `transit_gateway_id` - Identifier of EC2 Transit Gateway.
* `peer_transit_gateway_id` - Identifier of EC2 Transit Gateway to peer with.
* `peer_account_id` - Identifier of the AWS account that owns the EC2 TGW peering.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

`aws_ec2_transit_gateway_peering_attachment_accepter` can be imported by using the EC2 Transit Gateway Attachment identifier, e.g.,

```
$ terraform import aws_ec2_transit_gateway_peering_attachment_accepter.example tgw-attach-12345678
```
