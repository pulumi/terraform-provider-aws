---
subcategory: "Network Firewall"
layout: "aws"
page_title: "AWS: aws_networkfirewall_firewall_policy"
description: |-
  Provides an AWS Network Firewall Policy resource.
---

# Resource: aws_networkfirewall_firewall_policy

Provides an AWS Network Firewall Firewall Policy Resource

## Example Usage

```terraform
resource "aws_networkfirewall_firewall_policy" "example" {
  name = "example"

  firewall_policy {
    stateless_default_actions          = ["aws:pass"]
    stateless_fragment_default_actions = ["aws:drop"]
    stateless_rule_group_reference {
      priority     = 1
      resource_arn = aws_networkfirewall_rule_group.example.arn
    }
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }
}
```

## Policy with a Custom Action for Stateless Inspection

```terraform
resource "aws_networkfirewall_firewall_policy" "test" {
  name = "example"

  firewall_policy {
    stateless_default_actions          = ["aws:pass", "ExampleCustomAction"]
    stateless_fragment_default_actions = ["aws:drop"]

    stateless_custom_action {
      action_definition {
        publish_metric_action {
          dimension {
            value = "1"
          }
        }
      }
      action_name = "ExampleCustomAction"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A friendly description of the firewall policy.

* `encryption_configuration` - (Optional) KMS encryption configuration settings. See [Encryption Configuration](#encryption-configuration) below for details.

* `firewall_policy` - (Required) A configuration block describing the rule groups and policy actions to use in the firewall policy. See [Firewall Policy](#firewall-policy) below for details.

* `name` - (Required, Forces new resource) A friendly name of the firewall policy.

* `tags` - (Optional) Map of resource tags to associate with the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### Encryption Configuration

`encryption_configuration` settings for customer managed KMS keys. Remove this block to use the default AWS-managed KMS encryption (rather than setting `type` to `AWS_OWNED_KMS_KEY`).

* `key_id` - (Optional) The ID of the customer managed key. You can use any of the [key identifiers](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#key-id) that KMS supports, unless you're using a key that's managed by another account. If you're using a key managed by another account, then specify the key ARN.
* `type` - (Required) The type of AWS KMS key to use for encryption of your Network Firewall resources. Valid values are `CUSTOMER_KMS` and `AWS_OWNED_KMS_KEY`.

### Firewall Policy

The `firewall_policy` block supports the following arguments:

* `stateful_default_actions` - (Optional) Set of actions to take on a packet if it does not match any stateful rules in the policy. This can only be specified if the policy has a `stateful_engine_options` block with a `rule_order` value of `STRICT_ORDER`. You can specify one of either or neither values of `aws:drop_strict` or `aws:drop_established`, as well as any combination of `aws:alert_strict` and `aws:alert_established`.

* `stateful_engine_options` - (Optional) A configuration block that defines options on how the policy handles stateful rules. See [Stateful Engine Options](#stateful-engine-options) below for details.

* `stateful_rule_group_reference` - (Optional) Set of configuration blocks containing references to the stateful rule groups that are used in the policy. See [Stateful Rule Group Reference](#stateful-rule-group-reference) below for details.

* `stateless_custom_action` - (Optional) Set of configuration blocks describing the custom action definitions that are available for use in the firewall policy's `stateless_default_actions`. See [Stateless Custom Action](#stateless-custom-action) below for details.

* `stateless_default_actions` - (Required) Set of actions to take on a packet if it does not match any of the stateless rules in the policy. You must specify one of the standard actions including: `aws:drop`, `aws:pass`, or `aws:forward_to_sfe`.
In addition, you can specify custom actions that are compatible with your standard action choice. If you want non-matching packets to be forwarded for stateful inspection, specify `aws:forward_to_sfe`.

* `stateless_fragment_default_actions` - (Required) Set of actions to take on a fragmented packet if it does not match any of the stateless rules in the policy. You must specify one of the standard actions including: `aws:drop`, `aws:pass`, or `aws:forward_to_sfe`.
In addition, you can specify custom actions that are compatible with your standard action choice. If you want non-matching packets to be forwarded for stateful inspection, specify `aws:forward_to_sfe`.

* `stateless_rule_group_reference` - (Optional) Set of configuration blocks containing references to the stateless rule groups that are used in the policy. See [Stateless Rule Group Reference](#stateless-rule-group-reference) below for details.

### Stateful Engine Options
The `stateful_engine_options` block supports the following argument:

~> **NOTE:** If the `STRICT_ORDER` rule order is specified, this firewall policy can only reference stateful rule groups that utilize `STRICT_ORDER`.

* `rule_order` - (Required) Indicates how to manage the order of stateful rule evaluation for the policy. Default value: `DEFAULT_ACTION_ORDER`. Valid values: `DEFAULT_ACTION_ORDER`, `STRICT_ORDER`.

### Stateful Rule Group Reference

The `stateful_rule_group_reference` block supports the following arguments:

* `priority` - (Optional) An integer setting that indicates the order in which to apply the stateful rule groups in a single policy. This argument must be specified if the policy has a `stateful_engine_options` block with a `rule_order` value of `STRICT_ORDER`. AWS Network Firewall applies each stateful rule group to a packet starting with the group that has the lowest priority setting.

* `resource_arn` - (Required) The Amazon Resource Name (ARN) of the stateful rule group.

* `override` - (Optional) Configuration block for override values

#### Override

* `action` - (Optional) The action that changes the rule group from DROP to ALERT . This only applies to managed rule groups.

### Stateless Custom Action

The `stateless_custom_action` block supports the following arguments:

* `action_definition` - (Required) A configuration block describing the custom action associated with the `action_name`. See [Action Definition](#action-definition) below for details.

* `action_name` - (Required, Forces new resource) A friendly name of the custom action.

### Stateless Rule Group Reference

The `stateless_rule_group_reference` block supports the following arguments:

* `priority` - (Required) An integer setting that indicates the order in which to run the stateless rule groups in a single policy. AWS Network Firewall applies each stateless rule group to a packet starting with the group that has the lowest priority setting.

* `resource_arn` - (Required) The Amazon Resource Name (ARN) of the stateless rule group.

### Action Definition

The `action_definition` block supports the following argument:

* `publish_metric_action` - (Required) A configuration block describing the stateless inspection criteria that publishes the specified metrics to Amazon CloudWatch for the matching packet. You can pair this custom action with any of the standard stateless rule actions. See [Publish Metric Action](#publish-metric-action) below for details.

### Publish Metric Action

The `publish_metric_action` block supports the following argument:

* `dimension` - (Required) Set of configuration blocks describing dimension settings to use for Amazon CloudWatch custom metrics. See [Dimension](#dimension) below for more details.

### Dimension

The `dimension` block supports the following argument:

* `value` - (Required) The string value to use in the custom metric dimension.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Amazon Resource Name (ARN) that identifies the firewall policy.

* `arn` - The Amazon Resource Name (ARN) that identifies the firewall policy.

* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

* `update_token` - A string token used when updating a firewall policy.

## Import

Network Firewall Policies can be imported using their `ARN`.

```
$ terraform import aws_networkfirewall_firewall_policy.example arn:aws:network-firewall:us-west-1:123456789012:firewall-policy/example
```
