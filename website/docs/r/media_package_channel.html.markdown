---
subcategory: "Elemental MediaPackage"
layout: "aws"
page_title: "AWS: aws_media_package_channel"
description: |-
  Provides an AWS Elemental MediaPackage Channel.
---

# Resource: aws_media_package_channel

Provides an AWS Elemental MediaPackage Channel.

## Example Usage

```terraform
resource "aws_media_package_channel" "kittens" {
  channel_id  = "kitten-channel"
  description = "A channel dedicated to amusing videos of kittens."
}
```

## Argument Reference

The following arguments are supported:

* `channel_id` - (Required) A unique identifier describing the channel
* `description` - (Optional) A description of the channel
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The same as `channel_id`
* `arn` - The ARN of the channel
* `hls_ingest` - A single item list of HLS ingest information
    * `ingest_endpoints` - A list of the ingest endpoints
        * `password` - The password
        * `url` - The URL
        * `username` - The username
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Media Package Channels can be imported via the channel ID, e.g.,

```
$ terraform import aws_media_package_channel.kittens kittens-channel
```
