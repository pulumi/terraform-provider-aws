---
subcategory: "ELB Classic"
layout: "aws"
page_title: "AWS: aws_elb_attachment"
description: |-
  Provides an Elastic Load Balancer Attachment resource.
---

# Resource: aws_elb_attachment

Attaches an EC2 instance to an Elastic Load Balancer (ELB). For attaching resources with Application Load Balancer (ALB) or Network Load Balancer (NLB), see the `aws_lb_target_group_attachment` resource.

~> **NOTE on ELB Instances and ELB Attachments:** This provider currently provides
both a standalone ELB Attachment resource (describing an instance attached to
an ELB), and an Elastic Load Balancer resource with
`instances` defined in-line. At this time you cannot use an ELB with in-line
instances in conjunction with an ELB Attachment resource. Doing so will cause a
conflict and will overwrite attachments.

## Example Usage

```terraform
# Create a new load balancer attachment
resource "aws_elb_attachment" "baz" {
  elb      = aws_elb.bar.id
  instance = aws_instance.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `elb` - (Required) The name of the ELB.
* `instance` - (Required) Instance ID to place in the ELB pool.

## Attributes Reference

No additional attributes are exported.
