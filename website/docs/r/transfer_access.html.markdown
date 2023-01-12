---
subcategory: "Transfer Family"
layout: "aws"
page_title: "AWS: aws_transfer_access"
description: |-
  Provides a AWS Transfer Access resource.
---

# Resource: aws_transfer_access

Provides a AWS Transfer Access resource.

## Example Usage

### Basic S3

```terraform
resource "aws_transfer_access" "example" {
  external_id    = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id      = aws_transfer_server.example.id
  role           = aws_iam_role.example.arn
  home_directory = "/${aws_s3_bucket.example.id}/"
}
```

### Basic EFS

```terraform
resource "aws_transfer_access" "test" {
  external_id    = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id      = aws_transfer_server.test.id
  role           = aws_iam_role.test.arn
  home_directory = "/${aws_efs_file_system.test.id}/"
  posix_profile {
    gid = 1000
    uid = 1000
  }
}
```

## Argument Reference

The following arguments are supported:

* `external_id` - (Required) The SID of a group in the directory connected to the Transfer Server (e.g., `S-1-1-12-1234567890-123456789-1234567890-1234`)
* `server_id` - (Required) The Server ID of the Transfer Server (e.g., `s-12345678`)
* `home_directory` - (Optional) The landing directory (folder) for a user when they log in to the server using their SFTP client.  It should begin with a `/`.  The first item in the path is the name of the home bucket (accessible as `${Transfer:HomeBucket}` in the policy) and the rest is the home directory (accessible as `${Transfer:HomeDirectory}` in the policy). For example, `/example-bucket-1234/username` would set the home bucket to `example-bucket-1234` and the home directory to `username`.
* `home_directory_mappings` - (Optional) Logical directory mappings that specify what S3 paths and keys should be visible to your user and how you want to make them visible. See [Home Directory Mappings](#home-directory-mappings) below.
* `home_directory_type` - (Optional) The type of landing directory (folder) you mapped for your users' home directory. Valid values are `PATH` and `LOGICAL`.
* `policy` - (Optional) An IAM JSON policy document that scopes down user access to portions of their Amazon S3 bucket. IAM variables you can use inside this policy include `${Transfer:UserName}`, `${Transfer:HomeDirectory}`, and `${Transfer:HomeBucket}`. These are evaluated on-the-fly when navigating the bucket.
* `posix_profile` - (Optional) Specifies the full POSIX identity, including user ID (Uid), group ID (Gid), and any secondary groups IDs (SecondaryGids), that controls your users' access to your Amazon EFS file systems. See [Posix Profile](#posix-profile) below.
* `role` - (Required) Amazon Resource Name (ARN) of an IAM role that allows the service to controls your user’s access to your Amazon S3 bucket.

### Home Directory Mappings

* `entry` - (Required) Represents an entry and a target.
* `target` - (Required) Represents the map target.

### Posix Profile

* `gid` - (Required) The POSIX group ID used for all EFS operations by this user.
* `uid` - (Required) The POSIX user ID used for all EFS operations by this user.
* `secondary_gids` - (Optional) The secondary POSIX group IDs used for all EFS operations by this user.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id`  - The ID of the resource

## Import

Transfer Accesses can be imported using the `server_id` and `external_id`, e.g.,

```
$ terraform import aws_transfer_access.example s-12345678/S-1-1-12-1234567890-123456789-1234567890-1234
```
