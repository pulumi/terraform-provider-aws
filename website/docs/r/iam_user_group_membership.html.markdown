---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_user_group_membership"
description: |-
  Provides a resource for adding an IAM User to IAM Groups without conflicting
  with itself.
---

# Resource: aws_iam_user_group_membership

Provides a resource for adding an IAM User to IAM Groups. This
resource can be used multiple times with the same user for non-overlapping
groups.

To exclusively manage the users in a group, see the
`aws_iam_group_membership` resource.

## Example Usage

```terraform
resource "aws_iam_user_group_membership" "example1" {
  user = aws_iam_user.user1.name

  groups = [
    aws_iam_group.group1.name,
    aws_iam_group.group2.name,
  ]
}

resource "aws_iam_user_group_membership" "example2" {
  user = aws_iam_user.user1.name

  groups = [
    aws_iam_group.group3.name,
  ]
}

resource "aws_iam_user" "user1" {
  name = "user1"
}

resource "aws_iam_group" "group1" {
  name = "group1"
}

resource "aws_iam_group" "group2" {
  name = "group2"
}

resource "aws_iam_group" "group3" {
  name = "group3"
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) The name of the IAM User to add to groups
* `groups` - (Required) A list of IAM Groups to add the user to

## Attributes Reference

No additional attributes are exported.

## Import

IAM user group membership can be imported using the user name and group names separated by `/`.

```
$ terraform import aws_iam_user_group_membership.example1 user1/group1/group2
```
