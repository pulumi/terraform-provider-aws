---
subcategory: "SES (Simple Email)"
layout: "aws"
page_title: "AWS: aws_ses_domain_identity"
description: |-
  Provides an SES domain identity resource
---

# Resource: aws_ses_domain_identity

Provides an SES domain identity resource

## Example Usage

### Basic Usage

```terraform
resource "aws_ses_domain_identity" "example" {
  domain = "example.com"
}
```

### With Route53 Record

```terraform
resource "aws_ses_domain_identity" "example" {
  domain = "example.com"
}

resource "aws_route53_record" "example_amazonses_verification_record" {
  zone_id = "ABCDEFGHIJ123"
  name    = "_amazonses.example.com"
  type    = "TXT"
  ttl     = "600"
  records = [aws_ses_domain_identity.example.verification_token]
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain name to assign to SES

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the domain identity.
* `verification_token` - A code which when added to the domain as a TXT record
  will signal to SES that the owner of the domain has authorised SES to act on
  their behalf. The domain identity will be in state "verification pending"
  until this is done. See the [With Route53 Record](#with-route53-record) example
  for how this might be achieved when the domain is hosted in Route 53 and
  managed by this provider.  Find out more about verifying domains in Amazon
  SES in the [AWS SES
  docs](http://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-domains.html).

## Import

SES domain identities can be imported using the domain name.

```
$ terraform import aws_ses_domain_identity.example example.com
```
