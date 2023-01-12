---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_group_policy"
description: |-
  Provides an IAM policy attached to a group.
---

# Resource: aws_iam_group_policy

Provides an IAM policy attached to a group.

## Example Usage

```terraform
resource "aws_iam_group_policy" "my_developer_policy" {
  name  = "my_developer_policy"
  group = aws_iam_group.my_developers.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_group" "my_developers" {
  name = "developers"
  path = "/users/"
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required) The policy document. This is a JSON formatted string.
* `name` - (Optional) The name of the policy. If omitted, the provider will
assign a random, unique name.
* `name_prefix` - (Optional) Creates a unique name beginning with the specified
  prefix. Conflicts with `name`.
* `group` - (Required) The IAM group to attach to the policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group policy ID.
* `group` - The group to which this policy applies.
* `name` - The name of the policy.
* `policy` - The policy document attached to the group.

## Import

IAM Group Policies can be imported using the `group_name:group_policy_name`, e.g.,

```
$ terraform import aws_iam_group_policy.mypolicy group_of_mypolicy_name:mypolicy_name
```
