---
subcategory: "S3 Control"
layout: "aws"
page_title: "AWS: aws_s3control_access_point_policy"
description: |-
  Provides a resource to manage an S3 Access Point resource policy.
---

# Resource: aws_s3control_access_point_policy

Provides a resource to manage an S3 Access Point resource policy.

~> **NOTE on Access Points and Access Point Policies:** The provider provides both a standalone Access Point Policy resource and an Access Point resource with a resource policy defined in-line. You cannot use an Access Point with in-line resource policy in conjunction with an Access Point Policy resource. Doing so will cause a conflict of policies and will overwrite the access point's resource policy.

## Example Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "example"
}

resource "aws_s3_access_point" "example" {
  bucket = aws_s3_bucket.example.id
  name   = "example"

  public_access_block_configuration {
    block_public_acls       = true
    block_public_policy     = false
    ignore_public_acls      = true
    restrict_public_buckets = false
  }

  lifecycle {
    ignore_changes = [policy]
  }
}

resource "aws_s3control_access_point_policy" "example" {
  access_point_arn = aws_s3_access_point.example.arn

  policy = jsonencode({
    Version = "2008-10-17"
    Statement = [{
      Effect = "Allow"
      Action = "s3:GetObjectTagging"
      Principal = {
        AWS = "*"
      }
      Resource = "${aws_s3_access_point.example.arn}/object/*"
    }]
  })
}
```

## Argument Reference

The following arguments are supported:

* `access_point_arn` - (Required) The ARN of the access point that you want to associate with the specified policy.
* `policy` - (Required) The policy that you want to apply to the specified access point.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `has_public_access_policy` - Indicates whether this access point currently has a policy that allows public access.
* `id` - The AWS account ID and access point name separated by a colon (`:`).

## Import

Access Point policies can be imported using the `access_point_arn`, e.g.

```
$ terraform import aws_s3control_access_point_policy.example arn:aws:s3:us-west-2:123456789012:accesspoint/example
```
