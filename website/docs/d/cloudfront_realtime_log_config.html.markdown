---
subcategory: "CloudFront"
layout: "aws"
page_title: "AWS: aws_cloudfront_realtime_log_config"
description: |-
  Provides a CloudFront real-time log configuration resource.
---

# Data Source: aws_cloudfront_realtime_log_config

Provides a CloudFront real-time log configuration resource.

## Example Usage

```terraform
data "aws_cloudfront_realtime_log_config" "example" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Unique name to identify this real-time log configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN (Amazon Resource Name) of the CloudFront real-time log configuration.
* `endpoint` - (Required) Amazon Kinesis data streams where real-time log data is sent.
* `fields` - (Required) Fields that are included in each real-time log record. See the [AWS documentation](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/real-time-logs.html#understand-real-time-log-config-fields) for supported values.
* `sampling_rate` - (Required) Sampling rate for this real-time log configuration. The sampling rate determines the percentage of viewer requests that are represented in the real-time log data. An integer between `1` and `100`, inclusive.

The `endpoint` object supports the following:

* `kinesis_stream_config` - (Required) Amazon Kinesis data stream configuration.
* `stream_type` - (Required) Type of data stream where real-time log data is sent. The only valid value is `Kinesis`.

The `kinesis_stream_config` object supports the following:

* `role_arn` - (Required) ARN of an IAM role that CloudFront can use to send real-time log data to the Kinesis data stream.
See the [AWS documentation](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/real-time-logs.html#understand-real-time-log-config-iam-role) for more information.
* `stream_arn` - (Required) ARN of the Kinesis data stream.
