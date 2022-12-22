---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "AWS: aws_s3_bucket_logging"
description: |-
  Provides an S3 bucket (server access) logging resource.
---

# Resource: aws_s3_bucket_logging

Provides an S3 bucket (server access) logging resource. For more information, see [Logging requests using server access logging](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerLogs.html)
in the AWS S3 User Guide.

~> **Note:** Amazon S3 supports server access logging, AWS CloudTrail, or a combination of both. Refer to the [Logging options for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/logging-with-S3.html)
to decide which method meets your requirements.

## Example Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "my-tf-example-bucket"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_s3_bucket" "log_bucket" {
  bucket = "my-tf-log-bucket"
}

resource "aws_s3_bucket_acl" "log_bucket_acl" {
  bucket = aws_s3_bucket.log_bucket.id
  acl    = "log-delivery-write"
}

resource "aws_s3_bucket_logging" "example" {
  bucket = aws_s3_bucket.example.id

  target_bucket = aws_s3_bucket.log_bucket.id
  target_prefix = "log/"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, Forces new resource) The name of the bucket.
* `expected_bucket_owner` - (Optional, Forces new resource) The account ID of the expected bucket owner.
* `target_bucket` - (Required) The name of the bucket where you want Amazon S3 to store server access logs.
* `target_prefix` - (Required) A prefix for all log object keys.
* `target_grant` - (Optional) Set of configuration blocks with information for granting permissions [documented below](#target_grant).

### target_grant

The `target_grant` configuration block supports the following arguments:

* `grantee` - (Required) A configuration block for the person being granted permissions [documented below](#grantee).
* `permission` - (Required) Logging permissions assigned to the grantee for the bucket. Valid values: `FULL_CONTROL`, `READ`, `WRITE`.

### grantee

The `grantee` configuration block supports the following arguments:

* `email_address` - (Optional) Email address of the grantee. See [Regions and Endpoints](https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region) for supported AWS regions where this argument can be specified.
* `id` - (Optional) The canonical user ID of the grantee.
* `type` - (Required) Type of grantee. Valid values: `CanonicalUser`, `AmazonCustomerByEmail`, `Group`.
* `uri` - (Optional) URI of the grantee group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket` or `bucket` and `expected_bucket_owner` separated by a comma (`,`) if the latter is provided.

## Import

S3 bucket logging can be imported in one of two ways.

If the owner (account ID) of the source bucket is the same account used to configure the AWS Provider,
the S3 bucket logging resource should be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_logging.example bucket-name
```

If the owner (account ID) of the source bucket differs from the account used to configure the AWS Provider,
the S3 bucket logging resource should be imported using the `bucket` and `expected_bucket_owner` separated by a comma (`,`) e.g.,

```
$ terraform import aws_s3_bucket_logging.example bucket-name,123456789012
```
