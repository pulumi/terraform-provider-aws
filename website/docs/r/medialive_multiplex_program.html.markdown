---
subcategory: "Elemental MediaLive"
layout: "aws"
page_title: "AWS: aws_medialive_multiplex_program"
description: |-
  Resource for managing an AWS MediaLive MultiplexProgram.
---

# Resource: aws_medialive_multiplex_program

Resource for managing an AWS MediaLive MultiplexProgram.

~> **Note** Attribute `statemux_settings` has been deprecated and will be removed in a future major release. Please use `statmux_settings` instead.

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

resource "aws_medialive_multiplex_program" "example" {
  program_name = "example_program"
  multiplex_id = aws_medialive_multiplex.example.id

  multiplex_program_settings {
    program_number             = 1
    preferred_channel_pipeline = "CURRENTLY_ACTIVE"

    video_settings {
      constant_bitrate = 100000
    }
  }
}
```

## Argument Reference

The following arguments are required:

* `multiplex_id` - (Required) Multiplex ID.
* `program_name` - (Required) Unique program name.
* `multiplex_program_settings` - (Required) MultiplexProgram settings. See [Multiplex Program Settings](#multiple-program-settings) for more details.

The following arguments are optional:

### Multiple Program Settings

* `program_number` - (Required) Unique program number.
* `preferred_channel_pipeline` - (Required) Enum for preferred channel pipeline. Options are `CURRENTLY_ACTIVE`, `PIPELINE_0`, or `PIPELINE_1`.
* `service_descriptor` - (Optional) Service Descriptor. See [Service Descriptor](#service-descriptor) for more details.
* `video_settings` - (Optional) Video settings. See [Video Settings](#video-settings) for more details.

### Service Descriptor

* `provider_name` - (Required) Unique provider name.
* `service_name` - (Required) Unique service name.

### Video Settings

* `constant_bitrate` - (Optional) Constant bitrate value.
* `statemux_settings` - (Optional, **Deprecated**) Statemux settings. See [Statmux Settings](#statemux-settings) for more details. Settings from this attribute will apply to `statmux_settings`. Conflicts with `statmux_settings`.
* `statmux_settings` - (Optional) Statmux settings. See [Statmux Settings](#statmux-settings) for more details Conflicts with `statemux_settings`.

### Statmux Settings

* `minimum_bitrate` - (Optional) Minimum bitrate.
* `maximum_bitrate` - (Optional) Maximum bitrate.
* `priority` - (Optional) Priority value.

### Statemux Settings

* `minimum_bitrate` - (Optional) Minimum bitrate.
* `maximum_bitrate` - (Optional) Maximum bitrate.
* `priority` - (Optional) Priority value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the MultiplexProgram.
* `example_attribute` - Concise description.

## Import

MediaLive MultiplexProgram can be imported using the `id`, or a combination of "`program_name`/`multiplex_id`" e.g.,

```
$ terraform import aws_medialive_multiplex_program.example example_program/1234567
```
