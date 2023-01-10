---
subcategory: "WAF Classic Regional"
layout: "aws"
page_title: "AWS: aws_wafregional_rule_group"
description: |-
  Provides a AWS WAF Regional Rule Group resource.
---

# Resource: aws_wafregional_rule_group

Provides a WAF Regional Rule Group Resource

## Example Usage

```terraform
resource "aws_wafregional_rule" "example" {
  name        = "example"
  metric_name = "example"
}

resource "aws_wafregional_rule_group" "example" {
  name        = "example"
  metric_name = "example"

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A friendly name of the rule group
* `metric_name` - (Required) A friendly name for the metrics from the rule group
* `activated_rule` - (Optional) A list of activated rules, see below
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Nested Blocks

### `activated_rule`

#### Arguments

* `action` - (Required) Specifies the action that CloudFront or AWS WAF takes when a web request matches the conditions in the rule.
    * `type` - (Required) e.g., `BLOCK`, `ALLOW`, or `COUNT`
* `priority` - (Required) Specifies the order in which the rules are evaluated. Rules with a lower value are evaluated before rules with a higher value.
* `rule_id` - (Required) The ID of a rule
* `type` - (Optional) The rule type, either `REGULAR`, `RATE_BASED`, or `GROUP`. Defaults to `REGULAR`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the WAF Regional Rule Group.
* `arn` - The ARN of the WAF Regional Rule Group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

WAF Regional Rule Group can be imported using the id, e.g.,

```
$ terraform import aws_wafregional_rule_group.example a1b2c3d4-d5f6-7777-8888-9999aaaabbbbcccc
```
