---
subcategory: "Auto Scaling"
layout: "aws"
page_title: "AWS: aws_autoscaling_attachment"
description: |-
  Provides an AutoScaling Group Attachment resource.
---

# Resource: aws_autoscaling_attachment

Provides an Auto Scaling Attachment resource.

~> **NOTE on Auto Scaling Groups and ASG Attachments:** This provider currently provides
both a standalone `aws_autoscaling_attachment` resource
(describing an ASG attached to an ELB or ALB), and an `aws_autoscaling_group`
with `load_balancers` and `target_group_arns` defined in-line. These two methods are not
mutually-exclusive. If `aws_autoscaling_attachment` resources are used, either alone or with inline
`load_balancers` or `target_group_arns`, the `aws_autoscaling_group` resource must be configured
to ignore changes to the `load_balancers` and `target_group_arns` arguments.

## Example Usage

```terraform
# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "asg_attachment_bar" {
  autoscaling_group_name = aws_autoscaling_group.asg.id
  elb                    = aws_elb.bar.id
}
```

```terraform
# Create a new ALB Target Group attachment
resource "aws_autoscaling_attachment" "asg_attachment_bar" {
  autoscaling_group_name = aws_autoscaling_group.asg.id
  lb_target_group_arn    = aws_lb_target_group.test.arn
}
```

## With An AutoScaling Group Resource

```terraform
resource "aws_autoscaling_group" "asg" {
  # ... other configuration ...

  lifecycle {
    ignore_changes = [load_balancers, target_group_arns]
  }
}

resource "aws_autoscaling_attachment" "asg_attachment_bar" {
  autoscaling_group_name = aws_autoscaling_group.asg.id
  elb                    = aws_elb.test.id
}
```

## Argument Reference

The following arguments are supported:

* `autoscaling_group_name` - (Required) Name of ASG to associate with the ELB.
* `elb` - (Optional) Name of the ELB.
* `alb_target_group_arn` - (Optional, **Deprecated** use `lb_target_group_arn` instead) ARN of an ALB Target Group.
* `lb_target_group_arn` - (Optional) ARN of a load balancer target group.

## Attributes Reference

No additional attributes are exported.
