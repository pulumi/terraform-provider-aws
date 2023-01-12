---
subcategory: "Route 53"
layout: "aws"
page_title: "AWS: aws_route53_zone_association"
description: |-
  Manages a Route53 Hosted Zone VPC association
---

# Resource: aws_route53_zone_association

Manages a Route53 Hosted Zone VPC association. VPC associations can only be made on private zones. See the `aws_route53_vpc_association_authorization` resource for setting up cross-account associations.

~> **NOTE:** Unless explicit association ordering is required (e.g., a separate cross-account association authorization), usage of this resource is not recommended. Use the `vpc` configuration blocks available within the `aws_route53_zone` resource instead.

~> **NOTE:** This provider provides both this standalone Zone VPC Association resource and exclusive VPC associations defined in-line in the `aws_route53_zone` resource via `vpc` configuration blocks. At this time, you cannot use those in-line VPC associations in conjunction with this resource and the same zone ID otherwise it will cause a perpetual difference in plan output. You can optionally use [`ignoreChanges`](https://www.pulumi.com/docs/intro/concepts/programming-model/#ignorechanges) in the `aws_route53_zone` resource to manage additional associations via this resource.

## Example Usage

```terraform
resource "aws_vpc" "primary" {
  cidr_block           = "10.6.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_vpc" "secondary" {
  cidr_block           = "10.7.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_route53_zone" "example" {
  name = "example.com"

  # NOTE: The aws_route53_zone vpc argument accepts multiple configuration
  #       blocks. The below usage of the single vpc configuration, the
  #       lifecycle configuration, and the aws_route53_zone_association
  #       resource is for illustrative purposes (e.g., for a separate
  #       cross-account authorization process, which is not shown here).
  vpc {
    vpc_id = aws_vpc.primary.id
  }

  lifecycle {
    ignore_changes = [vpc]
  }
}

resource "aws_route53_zone_association" "secondary" {
  zone_id = aws_route53_zone.example.zone_id
  vpc_id  = aws_vpc.secondary.id
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The private hosted zone to associate.
* `vpc_id` - (Required) The VPC to associate with the private hosted zone.
* `vpc_region` - (Optional) The VPC's region. Defaults to the region of the AWS provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The calculated unique identifier for the association.
* `owning_account` - The account ID of the account that created the hosted zone.

## Import

Route 53 Hosted Zone Associations can be imported via the Hosted Zone ID and VPC ID, separated by a colon (`:`), e.g.,

```
$ terraform import aws_route53_zone_association.example Z123456ABCDEFG:vpc-12345678
```

If the VPC is in a different region than the provider region configuration, the VPC Region can be added to the end. e.g.

```
$ terraform import aws_route53_zone_association.example Z123456ABCDEFG:vpc-12345678:us-east-2
```
