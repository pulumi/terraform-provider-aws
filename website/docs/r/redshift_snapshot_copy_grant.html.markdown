---
subcategory: "Redshift"
layout: "aws"
page_title: "AWS: aws_redshift_snapshot_copy_grant"
description: |-
  Creates a snapshot copy grant that allows AWS Redshift to encrypt copied snapshots with a customer master key from AWS KMS in a destination region.
---

# Resource: aws_redshift_snapshot_copy_grant

Creates a snapshot copy grant that allows AWS Redshift to encrypt copied snapshots with a customer master key from AWS KMS in a destination region.

Note that the grant must exist in the destination region, and not in the region of the cluster.

## Example Usage

```terraform
resource "aws_redshift_snapshot_copy_grant" "test" {
  snapshot_copy_grant_name = "my-grant"
}

resource "aws_redshift_cluster" "test" {
  # ... other configuration ...
  snapshot_copy {
    destination_region = "us-east-2"
    grant_name         = aws_redshift_snapshot_copy_grant.test.snapshot_copy_grant_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_copy_grant_name` - (Required, Forces new resource) A friendly name for identifying the grant.
* `kms_key_id` - (Optional, Forces new resource) The unique identifier for the customer master key (CMK) that the grant applies to. Specify the key ID or the Amazon Resource Name (ARN) of the CMK. To specify a CMK in a different AWS account, you must use the key ARN. If not specified, the default key is used.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of snapshot copy grant
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Redshift Snapshot Copy Grants support import by name, e.g.,

```console
$ terraform import aws_redshift_snapshot_copy_grant.test my-grant
```
