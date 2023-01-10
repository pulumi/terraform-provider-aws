---
subcategory: "Kinesis Analytics V2"
layout: "aws"
page_title: "AWS: aws_kinesisanalyticsv2_application_snapshot"
description: |-
  Manages a Kinesis Analytics v2 Application Snapshot.
---

# Resource: aws_kinesisanalyticsv2_application_snapshot

Manages a Kinesis Analytics v2 Application Snapshot.
Snapshots are the AWS implementation of [Flink Savepoints](https://ci.apache.org/projects/flink/flink-docs-release-1.11/ops/state/savepoints.html).

## Example Usage

```terraform
resource "aws_kinesisanalyticsv2_application_snapshot" "example" {
  application_name = aws_kinesisanalyticsv2_application.example.name
  snapshot_name    = "example-snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) The name of an existing  Kinesis Analytics v2 Application. Note that the application must be running for a snapshot to be created.
* `snapshot_name` - (Required) The name of the application snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The application snapshot identifier.
* `application_version_id` - The current application version ID when the snapshot was created.
* `snapshot_creation_timestamp` - The timestamp of the application snapshot.

## Timeouts

Configuration options:

- `create` - (Default `10m`)
- `delete` - (Default `10m`)

## Import

`aws_kinesisanalyticsv2_application` can be imported by using `application_name` together with `snapshot_name`, e.g.,

```
$ terraform import aws_kinesisanalyticsv2_application_snapshot.example example-application/example-snapshot
```
