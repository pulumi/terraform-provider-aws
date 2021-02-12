---
subcategory: "IAM"
layout: "aws"
page_title: "AWS: aws_iam_role"
description: |-
  Provides an IAM role.
---

# Resource: aws_iam_role

Provides an IAM role.

~> *NOTE:* If policies are attached to the role via the `aws_iam_policy_attachment` resource and you are modifying the role `name` or `path`, the `force_detach_policies` argument must be set to `true` and applied before attempting the operation otherwise you will encounter a `DeleteConflict` error. The `aws_iam_role_policy_attachment` resource (recommended) does not have this requirement.

## Example Usage

```hcl
resource "aws_iam_role" "test_role" {
  name = "test_role"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
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

  tags = {
    tag-key = "tag-value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, Forces new resource) The name of the role. If omitted, this provider will assign a random, unique name.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `assume_role_policy` - (Required) The policy that grants an entity permission to assume the role.

* `force_detach_policies` - (Optional) Specifies to force detaching any policies the role has before destroying it. Defaults to `false`.
* `path` - (Optional) The path to the role.
  See [IAM Identifiers](https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_Identifiers.html) for more information.
* `description` - (Optional) The description of the role.

* `max_session_duration` - (Optional) The maximum session duration (in seconds) that you want to set for the specified role. If you do not specify a value for this setting, the default maximum of one hour is applied. This setting can have a value from 1 hour to 12 hours.
* `permissions_boundary` - (Optional) The ARN of the policy that is used to set the permissions boundary for the role.
* `tags` - Key-value map of tags for the IAM role

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) specifying the role.
* `create_date` - The creation date of the IAM role.
* `description` - The description of the role.
* `id` - The name of the role.
* `name` - The name of the role.
* `unique_id` - The stable and unique string identifying the role.

## Example of Using Data Source for Assume Role Policy

```hcl
data "aws_iam_policy_document" "instance-assume-role-policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "instance" {
  name               = "instance_role"
  path               = "/system/"
  assume_role_policy = data.aws_iam_policy_document.instance-assume-role-policy.json
}
```

## Import

IAM Roles can be imported using the `name`, e.g.

```
$ terraform import aws_iam_role.developer developer_name
```
