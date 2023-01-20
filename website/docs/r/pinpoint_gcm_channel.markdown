---
subcategory: "Pinpoint"
layout: "aws"
page_title: "AWS: aws_pinpoint_gcm_channel"
description: |-
  Provides a Pinpoint GCM Channel resource.
---

# Resource: aws_pinpoint_gcm_channel

Provides a Pinpoint GCM Channel resource.

## Example Usage

```terraform
resource "aws_pinpoint_gcm_channel" "gcm" {
  application_id = aws_pinpoint_app.app.application_id
  api_key        = "api_key"
}

resource "aws_pinpoint_app" "app" {}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The application ID.
* `api_key` - (Required) Platform credential API key from Google.
* `enabled` - (Optional) Whether the channel is enabled or disabled. Defaults to `true`.

## Attributes Reference

No additional attributes are exported.

## Import

Pinpoint GCM Channel can be imported using the `application-id`, e.g.,

```
$ terraform import aws_pinpoint_gcm_channel.gcm application-id
```
