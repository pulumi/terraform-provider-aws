---
subcategory: "EC2 Image Builder"
layout: "aws"
page_title: "AWS: aws_imagebuilder_image"
description: |-
    Manages an Image Builder Image
---

# Resource: aws_imagebuilder_image

Manages an Image Builder Image.

## Example Usage

```terraform
resource "aws_imagebuilder_image" "example" {
  distribution_configuration_arn   = aws_imagebuilder_distribution_configuration.example.arn
  image_recipe_arn                 = aws_imagebuilder_image_recipe.example.arn
  infrastructure_configuration_arn = aws_imagebuilder_infrastructure_configuration.example.arn
}
```

## Argument Reference

The following arguments are required:

* `infrastructure_configuration_arn` - (Required) Amazon Resource Name (ARN) of the Image Builder Infrastructure Configuration.

The following arguments are optional:

* `container_recipe_arn` - (Optional) - Amazon Resource Name (ARN) of the container recipe.
* `distribution_configuration_arn` - (Optional) Amazon Resource Name (ARN) of the Image Builder Distribution Configuration.
* `enhanced_image_metadata_enabled` - (Optional) Whether additional information about the image being created is collected. Defaults to `true`.
* `image_recipe_arn` - (Optional) Amazon Resource Name (ARN) of the image recipe.
* `image_tests_configuration` - (Optional) Configuration block with image tests configuration. Detailed below.
* `tags` - (Optional) Key-value map of resource tags for the Image Builder Image. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### image_tests_configuration

The following arguments are optional:

* `image_tests_enabled` - (Optional) Whether image tests are enabled. Defaults to `true`.
* `timeout_minutes` - (Optional) Number of minutes before image tests time out. Valid values are between `60` and `1440`. Defaults to `720`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the image.
* `date_created` - Date the image was created.
* `platform` - Platform of the image.
* `os_version` - Operating System version of the image.
* `output_resources` - List of objects with resources created by the image.
    * `amis` - Set of objects with each Amazon Machine Image (AMI) created.
        * `account_id` - Account identifier of the AMI.
        * `description` - Description of the AMI.
        * `image` - Identifier of the AMI.
        * `name` - Name of the AMI.
        * `region` - Region of the AMI.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
* `version` - Version of the image.

## Timeouts

Configuration options:

* `create` - (Default `60m`)

## Import

`aws_imagebuilder_image` resources can be imported using the Amazon Resource Name (ARN), e.g.,

```
$ terraform import aws_imagebuilder_image.example arn:aws:imagebuilder:us-east-1:123456789012:image/example/1.0.0/1
```
