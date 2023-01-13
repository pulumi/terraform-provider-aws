---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "AWS: aws_s3_bucket_server_side_encryption_configuration"
description: |-
  Provides a S3 bucket server-side encryption configuration resource.
---

# Resource: aws_s3_bucket_server_side_encryption_configuration

Provides a S3 bucket server-side encryption configuration resource.

## Example Usage

```terraform
resource "aws_kms_key" "mykey" {
  description             = "This key is used to encrypt bucket objects"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket" "mybucket" {
  bucket = "mybucket"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
  bucket = aws_s3_bucket.mybucket.bucket

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.mykey.arn
      sse_algorithm     = "aws:kms"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, Forces new resource) The name of the bucket.
* `expected_bucket_owner` - (Optional, Forces new resource) The account ID of the expected bucket owner.
* `rule` - (Required) Set of server-side encryption configuration rules. [documented below](#rule). Currently, only a single rule is supported.

### rule

The `rule` configuration block supports the following arguments:

* `apply_server_side_encryption_by_default` - (Optional) A single object for setting server-side encryption by default [documented below](#apply_server_side_encryption_by_default)
* `bucket_key_enabled` - (Optional) Whether or not to use [Amazon S3 Bucket Keys](https://docs.aws.amazon.com/AmazonS3/latest/dev/bucket-key.html) for SSE-KMS.

### apply_server_side_encryption_by_default

The `apply_server_side_encryption_by_default` configuration block supports the following arguments:

* `sse_algorithm` - (Required) The server-side encryption algorithm to use. Valid values are `AES256` and `aws:kms`
* `kms_master_key_id` - (Optional) The AWS KMS master key ID used for the SSE-KMS encryption. This can only be used when you set the value of `sse_algorithm` as `aws:kms`. The default `aws/s3` AWS KMS master key is used if this element is absent while the `sse_algorithm` is `aws:kms`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket` or `bucket` and `expected_bucket_owner` separated by a comma (`,`) if the latter is provided.

## Import

S3 bucket server-side encryption configuration can be imported in one of two ways.

If the owner (account ID) of the source bucket is the same account used to configure the AWS Provider,
the S3 server-side encryption configuration resource should be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_server_side_encryption_configuration.example bucket-name
```

If the owner (account ID) of the source bucket differs from the account used to configure the AWS Provider,
the S3 bucket server-side encryption configuration resource should be imported using the `bucket` and `expected_bucket_owner` separated by a comma (`,`) e.g.,

```
$ terraform import aws_s3_bucket_server_side_encryption_configuration.example bucket-name,123456789012
```
