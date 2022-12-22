---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_policy"
description: |-
  Provides an IAM policy.
---

# Resource: aws_iam_policy

Provides an IAM policy.

## Example Usage

```terraform
resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"

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
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, Forces new resource) Description of the IAM policy.
* `name` - (Optional, Forces new resource) The name of the policy. If omitted, this provider will assign a random, unique name.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `path` - (Optional, default "/") Path in which to create the policy.
  See [IAM Identifiers](https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_Identifiers.html) for more information.
* `policy` - (Required) The policy document. This is a JSON formatted string.
* `tags` - (Optional) Map of resource tags for the IAM Policy. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ARN assigned by AWS to this policy.
* `arn` - The ARN assigned by AWS to this policy.
* `description` - The description of the policy.
* `name` - The name of the policy.
* `path` - The path of the policy in IAM.
* `policy` - The policy document.
* `policy_id` - The policy's ID.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

IAM Policies can be imported using the `arn`, e.g.,

```
$ terraform import aws_iam_policy.administrator arn:aws:iam::123456789012:policy/UsersManageOwnCredentials
```
