---
subcategory: "MemoryDB for Redis"
layout: "aws"
page_title: "AWS: aws_memorydb_acl"
description: |-
  Provides a MemoryDB ACL.
---

# Resource: aws_memorydb_acl

Provides a MemoryDB ACL.

More information about users and ACL-s can be found in the [MemoryDB User Guide](https://docs.aws.amazon.com/memorydb/latest/devguide/clusters.acls.html).

## Example Usage

```terraform
resource "aws_memorydb_acl" "example" {
  name       = "my-acl"
  user_names = ["my-user-1", "my-user-2"]
}
```

## Argument Reference

The following arguments are optional:

* `name` - (Optional, Forces new resource) Name of the ACL. If omitted, the provider will assign a random, unique name. Conflicts with `name_prefix`.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `user_names` - (Optional) Set of MemoryDB user names to be included in this ACL.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Same as `name`.
* `arn` - The ARN of the ACL.
* `minimum_engine_version` - The minimum engine version supported by the ACL.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Use the `name` to import an ACL. For example:

```
$ terraform import aws_memorydb_acl.example my-acl
```
