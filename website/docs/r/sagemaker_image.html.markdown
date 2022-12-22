---
subcategory: "SageMaker"
layout: "aws"
page_title: "AWS: aws_sagemaker_image"
description: |-
  Provides a SageMaker Image resource.
---

# Resource: aws_sagemaker_image

Provides a SageMaker Image resource.

## Example Usage

### Basic usage

```terraform
resource "aws_sagemaker_image" "example" {
  image_name = "example"
  role_arn   = aws_iam_role.test.arn
}
```

## Argument Reference

The following arguments are supported:

* `image_name` - (Required) The name of the image. Must be unique to your account.
* `role_arn` - (Required) The Amazon Resource Name (ARN) of an IAM role that enables Amazon SageMaker to perform tasks on your behalf.
* `display_name` - (Optional) The display name of the image. When the image is added to a domain (must be unique to the domain).
* `description` - (Optional) The description of the image.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the Image.
* `arn` - The Amazon Resource Name (ARN) assigned by AWS to this Image.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

SageMaker Code Images can be imported using the `name`, e.g.,

```
$ terraform import aws_sagemaker_image.test_image my-code-repo
```
