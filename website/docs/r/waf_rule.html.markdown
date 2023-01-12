---
subcategory: "WAF Classic"
layout: "aws"
page_title: "AWS: aws_waf_rule"
description: |-
  Provides a AWS WAF rule resource.
---

# Resource: aws_waf_rule

Provides a WAF Rule Resource

## Example Usage

```terraform
resource "aws_waf_ipset" "ipset" {
  name = "tfIPSet"

  ip_set_descriptors {
    type  = "IPV4"
    value = "192.0.7.0/24"
  }
}

resource "aws_waf_rule" "wafrule" {
  depends_on  = [aws_waf_ipset.ipset]
  name        = "tfWAFRule"
  metric_name = "tfWAFRule"

  predicates {
    data_id = aws_waf_ipset.ipset.id
    negated = false
    type    = "IPMatch"
  }
}
```

## Argument Reference

The following arguments are supported:

* `metric_name` - (Required) The name or description for the Amazon CloudWatch metric of this rule. The name can contain only alphanumeric characters (A-Z, a-z, 0-9); the name can't contain whitespace.
* `name` - (Required) The name or description of the rule.
* `predicates` - (Optional) The objects to include in a rule (documented below).
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Nested Blocks

### `predicates`

See the [WAF Documentation](https://docs.aws.amazon.com/waf/latest/APIReference/API_Predicate.html) for more information.

#### Arguments

* `negated` - (Required) Set this to `false` if you want to allow, block, or count requests
  based on the settings in the specified waf_byte_match_set, waf_ipset, aws_waf_size_constraint_set, aws_waf_sql_injection_match_set or aws_waf_xss_match_set.
  For example, if an IPSet includes the IP address `192.0.2.44`, AWS WAF will allow or block requests based on that IP address.
  If set to `true`, AWS WAF will allow, block, or count requests based on all IP addresses except `192.0.2.44`.
* `data_id` - (Required) A unique identifier for a predicate in the rule, such as Byte Match Set ID or IPSet ID.
* `type` - (Required) The type of predicate in a rule. Valid values: `ByteMatch`, `GeoMatch`, `IPMatch`, `RegexMatch`, `SizeConstraint`, `SqlInjectionMatch`, or `XssMatch`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the WAF rule.
* `arn` - The ARN of the WAF rule.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

WAF rules can be imported using the id, e.g.,

```
$ terraform import aws_waf_rule.example a1b2c3d4-d5f6-7777-8888-9999aaaabbbbcccc
```
