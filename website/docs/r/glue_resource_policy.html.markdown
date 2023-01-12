---
subcategory: "Glue"
layout: "aws"
page_title: "AWS: aws_glue_resource_policy"
description: |-
  Provides a resource to configure the aws glue resource policy.
---

# Resource: aws_glue_resource_policy

Provides a Glue resource policy. Only one can exist per region.

## Example Usage

```terraform
data "aws_caller_identity" "current" {}

data "aws_partition" "current" {}

data "aws_region" "current" {}

data "aws_iam_policy_document" "glue-example-policy" {
  statement {
    actions = [
      "glue:CreateTable",
    ]
    resources = ["arn:${data.aws_partition.current.partition}:glue:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"]
    principals {
      identifiers = ["*"]
      type        = "AWS"
    }
  }
}

resource "aws_glue_resource_policy" "example" {
  policy = data.aws_iam_policy_document.glue-example-policy.json
}
```

## Argument Reference

The following arguments are supported:

* `policy` – (Required) The policy to be applied to the aws glue data catalog.
* `enable_hybrid` - (Optional) Indicates that you are using both methods to grant cross-account. Valid values are `TRUE` and `FALSE`. Note the provider will not perform drift detetction on this field as its not return on read.

## Attributes Reference

No additional attributes are exported.

## Import

Glue Resource Policy can be imported using the account ID:

```
$ terraform import aws_glue_resource_policy.Test 12356789012
```
