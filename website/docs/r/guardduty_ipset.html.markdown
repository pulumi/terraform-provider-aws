---
subcategory: "GuardDuty"
layout: aws
page_title: 'AWS: aws_guardduty_ipset'
description: Provides a resource to manage a GuardDuty IPSet
---

# Resource: aws_guardduty_ipset

Provides a resource to manage a GuardDuty IPSet.

~> **Note:** Currently in GuardDuty, users from member accounts cannot upload and further manage IPSets. IPSets that are uploaded by the primary account are imposed on GuardDuty functionality in its member accounts. See the [GuardDuty API Documentation](https://docs.aws.amazon.com/guardduty/latest/ug/create-ip-set.html)

## Example Usage

```terraform
resource "aws_guardduty_ipset" "example" {
  activate    = true
  detector_id = aws_guardduty_detector.primary.id
  format      = "TXT"
  location    = "https://s3.amazonaws.com/${aws_s3_object.MyIPSet.bucket}/${aws_s3_object.MyIPSet.key}"
  name        = "MyIPSet"
}

resource "aws_guardduty_detector" "primary" {
  enable = true
}

resource "aws_s3_bucket" "bucket" {
  # ... other configuration
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.bucket.id
  acl    = "private"
}

resource "aws_s3_object" "MyIPSet" {
  content = "10.0.0.0/8\n"
  bucket  = aws_s3_bucket.bucket.id
  key     = "MyIPSet"
}
```

## Argument Reference

The following arguments are supported:

* `activate` - (Required) Specifies whether GuardDuty is to start using the uploaded IPSet.
* `detector_id` - (Required) The detector ID of the GuardDuty.
* `format` - (Required) The format of the file that contains the IPSet. Valid values: `TXT` | `STIX` | `OTX_CSV` | `ALIEN_VAULT` | `PROOF_POINT` | `FIRE_EYE`
* `location` - (Required) The URI of the file that contains the IPSet.
* `name` - (Required) The friendly name to identify the IPSet.
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the GuardDuty IPSet.
* `id` - The ID of the GuardDuty IPSet.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

GuardDuty IPSet can be imported using the primary GuardDuty detector ID and IPSet ID, e.g.,

```
$ terraform import aws_guardduty_ipset.MyIPSet 00b00fd5aecc0ab60a708659477e9617:123456789012
```
