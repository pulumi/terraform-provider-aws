---
subcategory: "Route 53"
layout: "aws"
page_title: "AWS: aws_route53_delegation_set"
description: |-
  Provides a Route53 Delegation Set resource.
---

# Resource: aws_route53_delegation_set

Provides a [Route53 Delegation Set](https://docs.aws.amazon.com/Route53/latest/APIReference/API-actions-by-function.html#actions-by-function-reusable-delegation-sets) resource.

## Example Usage

```terraform
resource "aws_route53_delegation_set" "main" {
  reference_name = "DynDNS"
}

resource "aws_route53_zone" "primary" {
  name              = "mydomain.com"
  delegation_set_id = aws_route53_delegation_set.main.id
}

resource "aws_route53_zone" "secondary" {
  name              = "coolcompany.io"
  delegation_set_id = aws_route53_delegation_set.main.id
}
```

## Argument Reference

The following arguments are supported:

* `reference_name` - (Optional) This is a reference name used in Caller Reference
  (helpful for identifying single delegation set amongst others)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the Delegation Set.
* `id` - The delegation set ID
* `name_servers` - A list of authoritative name servers for the hosted zone
  (effectively a list of NS records).

## Import

Route53 Delegation Sets can be imported using the `delegation set id`, e.g.,

```
$ terraform import aws_route53_delegation_set.set1 N1PA6795SAMPLE
```
