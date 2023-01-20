---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_role_policy"
description: |-
  Provides an IAM role policy.
---

# Resource: aws_iam_role_policy

Provides an IAM role inline policy.

~> **NOTE:** For a given role, this resource is incompatible with using the `aws_iam_role` resource `inline_policy` argument. When using that argument and this resource, both will attempt to manage the role's inline policies and the provider will show a permanent difference.

## Example Usage

```terraform
resource "aws_iam_role_policy" "test_policy" {
  name = "test_policy"
  role = aws_iam_role.test_role.id

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

resource "aws_iam_role" "test_role" {
  name = "test_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the role policy. If omitted, this provider will
assign a random, unique name.
* `name_prefix` - (Optional) Creates a unique name beginning with the specified
  prefix. Conflicts with `name`.
* `policy` - (Required) The inline policy document. This is a JSON formatted string. For more information about building IAM policy documents with the provider, see the AWS IAM Policy Document Guide
* `role` - (Required) The name of the IAM role to attach to the policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The role policy ID, in the form of `role_name:role_policy_name`.
* `name` - The name of the policy.
* `policy` - The policy document attached to the role.
* `role` - The name of the role associated with the policy.

## Import

IAM Role Policies can be imported using the `role_name:role_policy_name`, e.g.,

```
$ terraform import aws_iam_role_policy.mypolicy role_of_mypolicy_name:mypolicy_name
```
