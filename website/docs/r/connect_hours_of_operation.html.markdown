---
subcategory: "Connect"
layout: "aws"
page_title: "AWS: aws_connect_hours_of_operation"
description: |-
  Provides details about a specific Amazon Connect Hours of Operation.
---

# Resource: aws_connect_hours_of_operation

Provides an Amazon Connect Hours of Operation resource. For more information see
[Amazon Connect: Getting Started](https://docs.aws.amazon.com/connect/latest/adminguide/amazon-connect-get-started.html)

## Example Usage

```terraform
resource "aws_connect_hours_of_operation" "test" {
  instance_id = "aaaaaaaa-bbbb-cccc-dddd-111111111111"
  name        = "Office Hours"
  description = "Monday office hours"
  time_zone   = "EST"

  config {
    day = "MONDAY"

    end_time {
      hours   = 23
      minutes = 8
    }

    start_time {
      hours   = 8
      minutes = 0
    }
  }

  config {
    day = "TUESDAY"

    end_time {
      hours   = 21
      minutes = 0
    }

    start_time {
      hours   = 9
      minutes = 0
    }
  }

  tags = {
    "Name" = "Example Hours of Operation"
  }
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required) One or more config blocks which define the configuration information for the hours of operation: day, start time, and end time . Config blocks are documented below.
* `description` - (Optional) Specifies the description of the Hours of Operation.
* `instance_id` - (Required) Specifies the identifier of the hosting Amazon Connect Instance.
* `name` - (Required) Specifies the name of the Hours of Operation.
* `tags` - (Optional) Tags to apply to the Hours of Operation. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `time_zone` - (Required) Specifies the time zone of the Hours of Operation.

A `config` block supports the following arguments:

* `day` - (Required) Specifies the day that the hours of operation applies to.
* `end_time` - (Required) A end time block specifies the time that your contact center closes. The `end_time` is documented below.
* `start_time` - (Required) A start time block specifies the time that your contact center opens. The `start_time` is documented below.

A `end_time` block supports the following arguments:

* `hours` - (Required) Specifies the hour of closing.
* `minutes` - (Required) Specifies the minute of closing.

A `start_time` block supports the following arguments:

* `hours` - (Required) Specifies the hour of opening.
* `minutes` - (Required) Specifies the minute of opening.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the Hours of Operation.
* `hours_of_operation_arn` - (**Deprecated**) The Amazon Resource Name (ARN) of the Hours of Operation.
* `hours_of_operation_id` - The identifier for the hours of operation.
* `id` - The identifier of the hosting Amazon Connect Instance and identifier of the Hours of Operation separated by a colon (`:`).
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Amazon Connect Hours of Operations can be imported using the `instance_id` and `hours_of_operation_id` separated by a colon (`:`), e.g.,

```
$ terraform import aws_connect_hours_of_operation.example f1288a1f-6193-445a-b47e-af739b2:c1d4e5f6-1b3c-1b3c-1b3c-c1d4e5f6c1d4e5
```
