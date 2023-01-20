---
subcategory: "Organizations"
layout: "aws"
page_title: "AWS: aws_organizations_policy"
description: |-
  Provides a resource to manage an AWS Organizations policy.
---

# Resource: aws_organizations_policy

Provides a resource to manage an [AWS Organizations policy](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies.html).

## Example Usage

```terraform
resource "aws_organizations_policy" "example" {
  name = "example"

  content = <<CONTENT
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
  }
}
CONTENT
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required) The policy content to add to the new policy. For example, if you create a [service control policy (SCP)](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_scp.html), this string must be JSON text that specifies the permissions that admins in attached accounts can delegate to their users, groups, and roles. For more information about the SCP syntax, see the [Service Control Policy Syntax documentation](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_reference_scp-syntax.html) and for more information on the Tag Policy syntax, see the [Tag Policy Syntax documentation](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_example-tag-policies.html).
* `name` - (Required) The friendly name to assign to the policy.
* `description` - (Optional) A description to assign to the policy.
* `type` - (Optional) The type of policy to create. Valid values are `AISERVICES_OPT_OUT_POLICY`, `BACKUP_POLICY`, `SERVICE_CONTROL_POLICY` (SCP), and `TAG_POLICY`. Defaults to `SERVICE_CONTROL_POLICY`.
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (ID) of the policy.
* `arn` - Amazon Resource Name (ARN) of the policy.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

`aws_organizations_policy` can be imported by using the policy ID, e.g.,

```
$ terraform import aws_organizations_policy.example p-12345678
```
