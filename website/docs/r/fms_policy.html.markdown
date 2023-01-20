---
subcategory: "FMS (Firewall Manager)"
layout: "aws"
page_title: "AWS: aws_fms_policy"
description: |-
  Provides a resource to create an AWS Firewall Manager policy
---

# Resource: aws_fms_policy

Provides a resource to create an AWS Firewall Manager policy. You need to be using AWS organizations and have enabled the Firewall Manager administrator account.

## Example Usage

```terraform
resource "aws_fms_policy" "example" {
  name                  = "FMS-Policy-Example"
  exclude_resource_tags = false
  remediation_enabled   = false
  resource_type         = "AWS::ElasticLoadBalancingV2::LoadBalancer"

  security_service_policy_data {
    type = "WAF"

    managed_service_data = jsonencode({
      type = "WAF",
      ruleGroups = [{
        id = aws_wafregional_rule_group.example.id
        overrideAction = {
          type = "COUNT"
        }
      }]
      defaultAction = {
        type = "BLOCK"
      }
      overrideCustomerWebACLAssociation = false
    })
  }

  tags = {
    Name = "example-fms-policy"
  }
}

resource "aws_wafregional_rule_group" "example" {
  metric_name = "WAFRuleGroupExample"
  name        = "WAF-Rule-Group-Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The friendly name of the AWS Firewall Manager Policy.
* `delete_all_policy_resources` - (Optional) If true, the request will also perform a clean-up process. Defaults to `true`. More information can be found here [AWS Firewall Manager delete policy](https://docs.aws.amazon.com/fms/2018-01-01/APIReference/API_DeletePolicy.html)
* `delete_unused_fm_managed_resources` - (Optional) If true, Firewall Manager will automatically remove protections from resources that leave the policy scope. Defaults to `false`. More information can be found here [AWS Firewall Manager policy contents](https://docs.aws.amazon.com/fms/2018-01-01/APIReference/API_Policy.html)
* `exclude_map` - (Optional) A map of lists of accounts and OU's to exclude from the policy.
* `exclude_resource_tags` - (Required, Forces new resource) A boolean value, if true the tags that are specified in the `resource_tags` are not protected by this policy. If set to false and resource_tags are populated, resources that contain tags will be protected by this policy.
* `include_map` - (Optional) A map of lists of accounts and OU's to include in the policy.
* `remediation_enabled` - (Required) A boolean value, indicates if the policy should automatically applied to resources that already exist in the account.
* `resource_tags` - (Optional) A map of resource tags, that if present will filter protections on resources based on the exclude_resource_tags.
* `resource_type` - (Optional) A resource type to protect. Conflicts with `resource_type_list`. See the [FMS API Reference](https://docs.aws.amazon.com/fms/2018-01-01/APIReference/API_Policy.html#fms-Type-Policy-ResourceType) for more information about supported values.
* `resource_type_list` - (Optional) A list of resource types to protect. Conflicts with `resource_type`. See the [FMS API Reference](https://docs.aws.amazon.com/fms/2018-01-01/APIReference/API_Policy.html#fms-Type-Policy-ResourceType) for more information about supported values. Lists with only one element are not supported, instead use `resource_type`.
* `security_service_policy_data` - (Required) The objects to include in Security Service Policy Data. Documented below.
* `tags` - (Optional) Key-value mapping of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level

## `exclude_map` Configuration Block

* `account` - (Optional) A list of AWS Organization member Accounts that you want to exclude from this AWS FMS Policy.
* `orgunit` - (Optional) A list of AWS Organizational Units that you want to exclude from this AWS FMS Policy. Specifying an OU is the equivalent of specifying all accounts in the OU and in any of its child OUs, including any child OUs and accounts that are added at a later time.

You can specify inclusions or exclusions, but not both. If you specify an `include_map`, AWS Firewall Manager applies the policy to all accounts specified by the `include_map`, and does not evaluate any `exclude_map` specifications. If you do not specify an `include_map`, then Firewall Manager applies the policy to all accounts except for those specified by the `exclude_map`.

## `include_map` Configuration Block

* `account` - (Optional) A list of AWS Organization member Accounts that you want to include for this AWS FMS Policy.
* `orgunit` - (Optional) A list of AWS Organizational Units that you want to include for this AWS FMS Policy. Specifying an OU is the equivalent of specifying all accounts in the OU and in any of its child OUs, including any child OUs and accounts that are added at a later time.

You can specify inclusions or exclusions, but not both. If you specify an `include_map`, AWS Firewall Manager applies the policy to all accounts specified by the `include_map`, and does not evaluate any `exclude_map` specifications. If you do not specify an `include_map`, then Firewall Manager applies the policy to all accounts except for those specified by the `exclude_map`.

## `security_service_policy_data` Configuration Block

* `managed_service_data` (Optional) Details about the service that are specific to the service type, in JSON format. For service type `SHIELD_ADVANCED`, this is an empty string. Examples depending on `type` can be found in the [AWS Firewall Manager SecurityServicePolicyData API Reference](https://docs.aws.amazon.com/fms/2018-01-01/APIReference/API_SecurityServicePolicyData.html).
* `type` - (Required, Forces new resource) The service that the policy is using to protect the resources. For the current list of supported types, please refer to the [AWS Firewall Manager SecurityServicePolicyData API Type Reference](https://docs.aws.amazon.com/fms/2018-01-01/APIReference/API_SecurityServicePolicyData.html#fms-Type-SecurityServicePolicyData-Type).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The AWS account ID of the AWS Firewall Manager administrator account.
* `policy_update_token` - A unique identifier for each update to the policy.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Firewall Manager policies can be imported using the policy ID, e.g.,

```
$ terraform import aws_fms_policy.example 5be49585-a7e3-4c49-dde1-a179fe4a619a
```
