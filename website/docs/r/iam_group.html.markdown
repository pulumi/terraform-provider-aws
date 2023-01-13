---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_group"
description: |-
  Provides an IAM group.
---

# Resource: aws_iam_group

Provides an IAM group.

~> **NOTE on user management:** Using `aws_iam_group_membership` or `aws_iam_user_group_membership` resources in addition to manually managing user/group membership using the console may lead to configuration drift or conflicts. For this reason, it's recommended to either manage membership entirely with the provider or entirely within the AWS console.

## Example Usage

```terraform
resource "aws_iam_group" "developers" {
  name = "developers"
  path = "/users/"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The group's name. The name must consist of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: `=,.@-_.`. Group names are not distinguished by case. For example, you cannot create groups named both "ADMINS" and "admins".
* `path` - (Optional, default "/") Path in which to create the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group's ID.
* `arn` - The ARN assigned by AWS for this group.
* `name` - The group's name.
* `path` - The path of the group in IAM.
* `unique_id` - The [unique ID][1] assigned by AWS.

  [1]: https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_Identifiers.html#GUIDs

## Import

IAM Groups can be imported using the `name`, e.g.,

```
$ terraform import aws_iam_group.developers developers
```
