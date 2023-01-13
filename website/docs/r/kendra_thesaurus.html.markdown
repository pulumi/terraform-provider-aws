---
subcategory: "Kendra"
layout: "aws"
page_title: "AWS: aws_kendra_thesaurus"
description: |-
  Resource for managing an AWS Kendra Thesaurus.
---

# Resource: aws_kendra_thesaurus

Resource for managing an AWS Kendra Thesaurus.

## Example Usage

```terraform
resource "aws_kendra_thesaurus" "example" {
  index_id = aws_kendra_index.example.id
  name     = "Example"
  role_arn = aws_iam_role.example.arn

  source_s3_path {
    bucket = aws_s3_bucket.example.id
    key    = aws_s3_object.example.key
  }

  tags = {
    Name = "Example Kendra Thesaurus"
  }
}
```

## Argument Reference

The following arguments are required:

* `index_id`- (Required, Forces new resource) The identifier of the index for a thesaurus.
* `name` - (Required) The name for the thesaurus.
* `role_arn` - (Required) The IAM (Identity and Access Management) role used to access the thesaurus file in S3.
* `source_s3_path` - (Required) The S3 path where your thesaurus file sits in S3. Detailed below.

The `source_s3_path` configuration block supports the following arguments:

* `bucket` - (Required) The name of the S3 bucket that contains the file.
* `key` - (Required) The name of the file.

The following arguments are optional:

* `description` - (Optional) The description for a thesaurus.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the thesaurus.
* `id` - The unique identifiers of the thesaurus and index separated by a slash (`/`).
* `status` - The current status of the thesaurus.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Timeouts

Configuration options:

* `create` - (Default `30m`)
* `update` - (Default `30m`)
* `delete` - (Default `30m`)

## Import

`aws_kendra_thesaurus` can be imported using the unique identifiers of the thesaurus and index separated by a slash (`/`), e.g.,

```
$ terraform import aws_kendra_thesaurus.example thesaurus-123456780/idx-8012925589
```
