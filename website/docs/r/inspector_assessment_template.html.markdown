---
subcategory: "Inspector"
layout: "aws"
page_title: "AWS: aws_inspector_assessment_template"
description: |-
  Provides a Inspector assessment template.
---

# Resource: aws_inspector_assessment_template

Provides a Inspector assessment template

## Example Usage

```terraform
resource "aws_inspector_assessment_template" "example" {
  name       = "example"
  target_arn = aws_inspector_assessment_target.example.arn
  duration   = 3600

  rules_package_arns = [
    "arn:aws:inspector:us-west-2:758058086616:rulespackage/0-9hgA516p",
    "arn:aws:inspector:us-west-2:758058086616:rulespackage/0-H5hpSawc",
    "arn:aws:inspector:us-west-2:758058086616:rulespackage/0-JJOtZiqQ",
    "arn:aws:inspector:us-west-2:758058086616:rulespackage/0-vg5GGHSD",
  ]

  event_subscription {
    event     = "ASSESSMENT_RUN_COMPLETED"
    topic_arn = aws_sns_topic.example.arn
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the assessment template.
* `target_arn` - (Required) The assessment target ARN to attach the template to.
* `duration` - (Required) The duration of the inspector run.
* `rules_package_arns` - (Required) The rules to be used during the run.
* `event_subscription` - (Optional) A block that enables sending notifications about a specified assessment template event to a designated SNS topic. See [Event Subscriptions](#event-subscriptions) for details.
* `tags` - (Optional) Key-value map of tags for the Inspector assessment template. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### Event Subscriptions

The event subscription configuration block supports the following arguments:

* `event` - (Required) The event for which you want to receive SNS notifications. Valid values are `ASSESSMENT_RUN_STARTED`, `ASSESSMENT_RUN_COMPLETED`, `ASSESSMENT_RUN_STATE_CHANGED`, and `FINDING_REPORTED`.
* `topic_arn` - (Required) The ARN of the SNS topic to which notifications are sent.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The template assessment ARN.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

`aws_inspector_assessment_template` can be imported by using the template assessment ARN, e.g.,

```
$ terraform import aws_inspector_assessment_template.example arn:aws:inspector:us-west-2:123456789012:target/0-9IaAzhGR/template/0-WEcjR8CH
```
