---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "AWS: aws_s3_bucket_policy"
description: |-
  Attaches a policy to an S3 bucket resource.
---

# Resource: aws_s3_bucket_policy

Attaches a policy to an S3 bucket resource.

## Example Usage

### Basic Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "my-tf-test-bucket"
}

resource "aws_s3_bucket_policy" "allow_access_from_another_account" {
  bucket = aws_s3_bucket.example.id
  policy = data.aws_iam_policy_document.allow_access_from_another_account.json
}

data "aws_iam_policy_document" "allow_access_from_another_account" {
  statement {
    principals {
      type        = "AWS"
      identifiers = ["123456789012"]
    }

    actions = [
      "s3:GetObject",
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.example.arn,
      "${aws_s3_bucket.example.arn}/*",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which to apply the policy.
* `policy` - (Required) The text of the policy. Although this is a bucket policy rather than an IAM policy, the `aws_iam_policy_document` data source may be used, so long as it specifies a principal. Note: Bucket policies are limited to 20 KB in size.

## Attributes Reference

No additional attributes are exported.

## Import

S3 bucket policies can be imported using the bucket name, e.g.,

```
$ terraform import aws_s3_bucket_policy.allow_access_from_another_account my-tf-test-bucket
```
