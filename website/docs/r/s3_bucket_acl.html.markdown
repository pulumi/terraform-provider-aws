---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "AWS: aws_s3_bucket_acl"
description: |-
  Provides an S3 bucket ACL resource.
---

# Resource: aws_s3_bucket_acl

Provides an S3 bucket ACL resource.

~> **Note:** destroy does not delete the S3 Bucket ACL but does remove the resource from state.

## Example Usage

### With ACL

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "my-tf-example-bucket"
}

resource "aws_s3_bucket_acl" "example_bucket_acl" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}
```

### With Grants

```terraform
data "aws_canonical_user_id" "current" {}

resource "aws_s3_bucket" "example" {
  bucket = "my-tf-example-bucket"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  access_control_policy {
    grant {
      grantee {
        id   = data.aws_canonical_user_id.current.id
        type = "CanonicalUser"
      }
      permission = "READ"
    }

    grant {
      grantee {
        type = "Group"
        uri  = "http://acs.amazonaws.com/groups/s3/LogDelivery"
      }
      permission = "READ_ACP"
    }

    owner {
      id = data.aws_canonical_user_id.current.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional, Conflicts with `access_control_policy`) The canned ACL to apply to the bucket.
* `access_control_policy` - (Optional, Conflicts with `acl`) A configuration block that sets the ACL permissions for an object per grantee [documented below](#access_control_policy).
* `bucket` - (Required, Forces new resource) The name of the bucket.
* `expected_bucket_owner` - (Optional, Forces new resource) The account ID of the expected bucket owner.

### access_control_policy

The `access_control_policy` configuration block supports the following arguments:

* `grant` - (Required) Set of `grant` configuration blocks [documented below](#grant).
* `owner` - (Required) Configuration block of the bucket owner's display name and ID [documented below](#owner).

### grant

The `grant` configuration block supports the following arguments:

* `grantee` - (Required) Configuration block for the person being granted permissions [documented below](#grantee).
* `permission` - (Required) Logging permissions assigned to the grantee for the bucket.

### owner

The `owner` configuration block supports the following arguments:

* `id` - (Required) The ID of the owner.
* `display_name` - (Optional) The display name of the owner.

### grantee

The `grantee` configuration block supports the following arguments:

* `email_address` - (Optional) Email address of the grantee. See [Regions and Endpoints](https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region) for supported AWS regions where this argument can be specified.
* `id` - (Optional) The canonical user ID of the grantee.
* `type` - (Required) Type of grantee. Valid values: `CanonicalUser`, `AmazonCustomerByEmail`, `Group`.
* `uri` - (Optional) URI of the grantee group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The `bucket`, `expected_bucket_owner` (if configured), and `acl` (if configured) separated by commas (`,`).

## Import

S3 bucket ACL can be imported in one of four ways.

If the owner (account ID) of the source bucket is the _same_ account used to configure the AWS Provider, and the source bucket is **not configured** with a
[canned ACL][1] (i.e. predefined grant), the S3 bucket ACL resource should be imported using the `bucket` e.g.,

```
$ terraform import aws_s3_bucket_acl.example bucket-name
```

If the owner (account ID) of the source bucket is the _same_ account used to configure the AWS Provider, and the source bucket is **configured** with a
[canned ACL][1] (i.e. predefined grant), the S3 bucket ACL resource should be imported using the `bucket` and `acl` separated by a comma (`,`), e.g.

```
$ terraform import aws_s3_bucket_acl.example bucket-name,private
```

If the owner (account ID) of the source bucket _differs_ from the account used to configure the AWS Provider, and the source bucket is **not configured** with a
[canned ACL][1] (i.e. predefined grant), the S3 bucket ACL resource should be imported using the `bucket` and `expected_bucket_owner` separated by a comma (`,`) e.g.,

```
$ terraform import aws_s3_bucket_acl.example bucket-name,123456789012
```

If the owner (account ID) of the source bucket _differs_ from the account used to configure the AWS Provider, and the source bucket is **configured** with a
[canned ACL][1] (i.e. predefined grant), the S3 bucket ACL resource should be imported using the `bucket`, `expected_bucket_owner`, and `acl` separated by commas (`,`), e.g.,

```
$ terraform import aws_s3_bucket_acl.example bucket-name,123456789012,private
```

[1]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl
