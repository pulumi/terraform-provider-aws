---
subcategory: "SNS (Simple Notification)"
layout: "aws"
page_title: "AWS: aws_sns_topic_policy"
description: |-
  Provides an SNS topic policy resource.
---

# Resource: aws_sns_topic_policy

Provides an SNS topic policy resource

~> **NOTE:** If a Principal is specified as just an AWS account ID rather than an ARN, AWS silently converts it to the ARN for the root user, causing future deployments to differ. To avoid this problem, just specify the full ARN, e.g. `arn:aws:iam::123456789012:root`

## Example Usage

```terraform
resource "aws_sns_topic" "test" {
  name = "my-topic-with-policy"
}

resource "aws_sns_topic_policy" "default" {
  arn = aws_sns_topic.test.arn

  policy = data.aws_iam_policy_document.sns_topic_policy.json
}

data "aws_iam_policy_document" "sns_topic_policy" {
  policy_id = "__default_policy_ID"

  statement {
    actions = [
      "SNS:Subscribe",
      "SNS:SetTopicAttributes",
      "SNS:RemovePermission",
      "SNS:Receive",
      "SNS:Publish",
      "SNS:ListSubscriptionsByTopic",
      "SNS:GetTopicAttributes",
      "SNS:DeleteTopic",
      "SNS:AddPermission",
    ]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceOwner"

      values = [
        var.account-id,
      ]
    }

    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      aws_sns_topic.test.arn,
    ]

    sid = "__default_statement_ID"
  }
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Required) The ARN of the SNS topic
* `policy` - (Required) The fully-formed AWS policy as JSON.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `owner` - The AWS Account ID of the SNS topic owner

## Import

SNS Topic Policy can be imported using the topic ARN, e.g.,

```
$ terraform import aws_sns_topic_policy.user_updates arn:aws:sns:us-west-2:0123456789012:my-topic
```
