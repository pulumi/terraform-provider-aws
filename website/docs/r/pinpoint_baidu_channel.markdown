---
subcategory: "Pinpoint"
layout: "aws"
page_title: "AWS: aws_pinpoint_baidu_channel"
description: |-
  Provides a Pinpoint Baidu Channel resource.
---

# Resource: aws_pinpoint_baidu_channel

Provides a Pinpoint Baidu Channel resource.

## Example Usage

```terraform
resource "aws_pinpoint_app" "app" {}

resource "aws_pinpoint_baidu_channel" "channel" {
  application_id = aws_pinpoint_app.app.application_id
  api_key        = ""
  secret_key     = ""
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The application ID.
* `enabled` - (Optional) Specifies whether to enable the channel. Defaults to `true`.
* `api_key` - (Required) Platform credential API key from Baidu.
* `secret_key` - (Required) Platform credential Secret key from Baidu.

## Attributes Reference

No additional attributes are exported.

## Import

Pinpoint Baidu Channel can be imported using the `application-id`, e.g.,

```
$ terraform import aws_pinpoint_baidu_channel.channel application-id
```
