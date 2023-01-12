---
subcategory: "Elemental MediaLive"
layout: "aws"
page_title: "AWS: aws_medialive_multiplex"
description: |-
  Resource for managing an AWS MediaLive Multiplex.
---

# Resource: aws_medialive_multiplex

Resource for managing an AWS MediaLive Multiplex.

## Example Usage

### Basic Usage

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_medialive_multiplex" "example" {
  name               = "example-multiplex-changed"
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]

  multiplex_settings {
    transport_stream_bitrate                = 1000000
    transport_stream_id                     = 1
    transport_stream_reserved_bitrate       = 1
    maximum_video_buffer_delay_milliseconds = 1000
  }

  start_multiplex = true

  tags = {
    tag1 = "value1"
  }
}
```

## Argument Reference

The following arguments are required:

* `availability_zones` - (Required) A list of availability zones. You must specify exactly two.
* `multiplex_settings`- (Required) Multiplex settings. See [Multiplex Settings](#multiplex-settings) for more details.
* `name` - (Required) name of Multiplex.

The following arguments are optional:

* `start_multiplex` - (Optional) Whether to start the Multiplex. Defaults to `false`.
* `tags` - (Optional) A map of tags to assign to the Multiplex. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### Multiplex Settings

* `transport_stream_bitrate` - (Required) Transport stream bit rate.
* `transport_stream_id` - (Required) Unique ID for each multiplex.
* `transport_stream_reserved_bitrate` - (Optional) Transport stream reserved bit rate.
* `maximum_video_buffer_delay_milliseconds` - (Optional) Maximum video buffer delay.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the Multiplex.

## Timeouts

Configuration options:

* `create` - (Default `30m`)
* `update` - (Default `30m`)
* `delete` - (Default `30m`)

## Import

MediaLive Multiplex can be imported using the `id`, e.g.,

```
$ terraform import aws_medialive_multiplex.example 12345678
```
