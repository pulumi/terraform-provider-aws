---
subcategory: "IVS (Interactive Video)"
layout: "aws"
page_title: "AWS: aws_ivs_recording_configuration"
description: |-
  Resource for managing an AWS IVS (Interactive Video) Recording Configuration.
---

# Resource: aws_ivs_recording_configuration

Resource for managing an AWS IVS (Interactive Video) Recording Configuration.

## Example Usage

### Basic Usage

```terraform
resource "aws_ivs_recording_configuration" "example" {
  name = "recording_configuration-1"
  destination_configuration {
    s3 {
      bucket_name = "ivs-stream-archive"
    }
  }
}
```

## Argument Reference

The following arguments are required:

* `destination_configuration` - Object containing destination configuration for where recorded video will be stored.
    * `s3` - S3 destination configuration where recorded videos will be stored.
        * `bucket_name` - S3 bucket name where recorded videos will be stored.

The following arguments are optional:

* `name` - (Optional) Recording Configuration name.
* `recording_reconnect_window_seconds` - (Optional) If a broadcast disconnects and then reconnects within the specified interval, the multiple streams will be considered a single broadcast and merged together.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `thumbnail_configuration` - (Optional) Object containing information to enable/disable the recording of thumbnails for a live session and modify the interval at which thumbnails are generated for the live session.
    * `recording_mode` - (Optional) Thumbnail recording mode. Valid values: `DISABLED`, `INTERVAL`.
    * `target_interval_seconds` (Configurable [and required] only if `recording_mode` is `INTERVAL`) - The targeted thumbnail-generation interval in seconds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the Recording Configuration.
* `state` -  The current state of the Recording Configuration.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Timeouts

Configuration options:

* `create` - (Default `10m`)
* `delete` - (Default `10m`)

## Import

IVS (Interactive Video) Recording Configuration can be imported using the ARN, e.g.,

```
$ terraform import aws_ivs_recording_configuration.example arn:aws:ivs:us-west-2:326937407773:recording-configuration/KAk1sHBl2L47
```
