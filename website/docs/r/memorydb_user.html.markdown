---
subcategory: "MemoryDB for Redis"
layout: "aws"
page_title: "AWS: aws_memorydb_user"
description: |-
  Provides a MemoryDB User.
---

# Resource: aws_memorydb_user

Provides a MemoryDB User.

More information about users and ACL-s can be found in the [MemoryDB User Guide](https://docs.aws.amazon.com/memorydb/latest/devguide/clusters.acls.html).

## Example Usage

```terraform
resource "random_password" "example" {
  length = 16
}

resource "aws_memorydb_user" "example" {
  user_name     = "my-user"
  access_string = "on ~* &* +@all"

  authentication_mode {
    type      = "password"
    passwords = [random_password.example.result]
  }
}
```

## Argument Reference

The following arguments are required:

* `access_string` - (Required) The access permissions string used for this user.
* `authentication_mode` - (Required) Denotes the user's authentication properties. Detailed below.
* `user_name` - (Required, Forces new resource) Name of the MemoryDB user. Up to 40 characters.

The following arguments are optional:

* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### authentication_mode Configuration Block

* `passwords` - (Required) The set of passwords used for authentication. You can create up to two passwords for each user.
* `type` - (Required) Indicates whether the user requires a password to authenticate. Must be set to `password`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Same as `user_name`.
* `arn` - The ARN of the user.
* `minimum_engine_version` - The minimum engine version supported for the user.
* `authentication_mode` configuration block
    * `password_count` - The number of passwords belonging to the user.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Use the `user_name` to import a user. For example:

```
$ terraform import aws_memorydb_user.example my-user
```

The `passwords` are not available for imported resources, as this information cannot be read back from the MemoryDB API.
