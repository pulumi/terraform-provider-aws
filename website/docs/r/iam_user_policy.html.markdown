---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_user_policy"
description: |-
  Provides an IAM policy attached to a user.
---

# Resource: aws_iam_user_policy

Provides an IAM policy attached to a user.

## Example Usage

```terraform
resource "aws_iam_user_policy" "lb_ro" {
  name = "test"
  user = aws_iam_user.lb.name

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

resource "aws_iam_user" "lb" {
  name = "loadbalancer"
  path = "/system/"
}

resource "aws_iam_access_key" "lb" {
  user = aws_iam_user.lb.name
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required) The policy document. This is a JSON formatted string.
* `name` - (Optional) The name of the policy. If omitted, the provider will assign a random, unique name.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `user` - (Required) IAM user to which to attach this policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The user policy ID, in the form of `user_name:user_policy_name`.
* `name` - The name of the policy (always set).

## Import

IAM User Policies can be imported using the `user_name:user_policy_name`, e.g.,

```
$ terraform import aws_iam_user_policy.mypolicy user_of_mypolicy_name:mypolicy_name
```
