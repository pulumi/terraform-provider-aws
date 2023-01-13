---
subcategory: "SQS (Simple Queue)"
layout: "aws"
page_title: "AWS: aws_sqs_queue_policy"
description: |-
  Provides a SQS Queue Policy resource.
---

# Resource: aws_sqs_queue_policy

Allows you to set a policy of an SQS Queue
while referencing ARN of the queue within the policy.

## Example Usage

```terraform
resource "aws_sqs_queue" "q" {
  name = "examplequeue"
}

resource "aws_sqs_queue_policy" "test" {
  queue_url = aws_sqs_queue.q.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.q.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_sns_topic.example.arn}"
        }
      }
    }
  ]
}
POLICY
}
```

## Argument Reference

The following arguments are supported:

* `queue_url` - (Required) The URL of the SQS Queue to which to attach the policy
* `policy` - (Required) The JSON policy for the SQS queue.

## Attributes Reference

No additional attributes are exported.

## Import

SQS Queue Policies can be imported using the queue URL, e.g.,

```
$ terraform import aws_sqs_queue_policy.test https://queue.amazonaws.com/0123456789012/myqueue
```
