---
subcategory: "QLDB (Quantum Ledger Database)"
layout: "aws"
page_title: "AWS: aws_qldb_stream"
description: |-
  Provides a QLDB Stream resource.
---

# Resource: aws_qldb_stream

Provides an AWS Quantum Ledger Database (QLDB) Stream resource

## Example Usage

```terraform
resource "aws_qldb_stream" "example" {
  ledger_name          = "existing-ledger-name"
  stream_name          = "sample-ledger-stream"
  role_arn             = "sample-role-arn"
  inclusive_start_time = "2021-01-01T00:00:00Z"

  kinesis_configuration {
    aggegation_enabled = false
    stream_arn         = "arn:aws:kinesis:us-east-1:xxxxxxxxxxxx:stream/example-kinesis-stream"
  }

  tags = {
    "example" = "tag"
  }
}
```

## Argument Reference

The following arguments are supported:

* `exclusive_end_time` - (Optional) The exclusive date and time that specifies when the stream ends. If you don't define this parameter, the stream runs indefinitely until you cancel it. It must be in ISO 8601 date and time format and in Universal Coordinated Time (UTC). For example: `"2019-06-13T21:36:34Z"`.
* `inclusive_start_time` - (Required) The inclusive start date and time from which to start streaming journal data. This parameter must be in ISO 8601 date and time format and in Universal Coordinated Time (UTC). For example: `"2019-06-13T21:36:34Z"`.  This cannot be in the future and must be before `exclusive_end_time`.  If you provide a value that is before the ledger's `CreationDateTime`, QLDB effectively defaults it to the ledger's `CreationDateTime`.
* `kinesis_configuration` - (Required) The configuration settings of the Kinesis Data Streams destination for your stream request. Documented below.
* `ledger_name` - (Required) The name of the QLDB ledger.
* `role_arn` - (Required) The Amazon Resource Name (ARN) of the IAM role that grants QLDB permissions for a journal stream to write data records to a Kinesis Data Streams resource.
* `stream_name` - (Required) The name that you want to assign to the QLDB journal stream. User-defined names can help identify and indicate the purpose of a stream.  Your stream name must be unique among other active streams for a given ledger. Stream names have the same naming constraints as ledger names, as defined in the [Amazon QLDB Developer Guide](https://docs.aws.amazon.com/qldb/latest/developerguide/limits.html#limits.naming).
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### kinesis_configuration

The `kinesis_configuration` block supports the following arguments:

* `aggregation_enabled` - (Optional) Enables QLDB to publish multiple data records in a single Kinesis Data Streams record, increasing the number of records sent per API call. Default: `true`.
* `stream_arn` - (Required) The Amazon Resource Name (ARN) of the Kinesis Data Streams resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the QLDB Stream.
* `arn` - The ARN of the QLDB Stream.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
