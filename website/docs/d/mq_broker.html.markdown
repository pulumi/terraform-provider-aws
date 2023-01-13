---
subcategory: "MQ"
layout: "aws"
page_title: "AWS: aws_mq_broker"
description: |-
  Provides a MQ Broker data source.
---

# Data Source: aws_mq_broker

Provides information about a MQ Broker.

## Example Usage

```terraform
variable "broker_id" {
  type    = string
  default = ""
}

variable "broker_name" {
  type    = string
  default = ""
}

data "aws_mq_broker" "by_id" {
  broker_id = var.broker_id
}

data "aws_mq_broker" "by_name" {
  broker_name = var.broker_name
}
```

## Argument Reference

The following arguments are supported:

* `broker_id` - (Optional) Unique id of the mq broker.
* `broker_name` - (Optional) Unique name of the mq broker.

## Attributes Reference

See the `aws_mq_broker` resource for details on the returned attributes.
They are identical except for user password, which is not returned when describing broker.
