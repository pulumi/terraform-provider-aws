---
subcategory: "Auto Scaling"
layout: "aws"
page_title: "AWS: aws_autoscaling_schedule"
description: |-
  Provides an AutoScaling Schedule resource.
---

# Resource: aws_autoscaling_schedule

Provides an AutoScaling Schedule resource.

## Example Usage

```terraform
resource "aws_autoscaling_group" "foobar" {
  availability_zones        = ["us-west-2a"]
  name                      = "test-foobar5"
  max_size                  = 1
  min_size                  = 1
  health_check_grace_period = 300
  health_check_type         = "ELB"
  force_delete              = true
  termination_policies      = ["OldestInstance"]
}

resource "aws_autoscaling_schedule" "foobar" {
  scheduled_action_name  = "foobar"
  min_size               = 0
  max_size               = 1
  desired_capacity       = 0
  start_time             = "2016-12-11T18:00:00Z"
  end_time               = "2016-12-12T06:00:00Z"
  autoscaling_group_name = aws_autoscaling_group.foobar.name
}
```

## Argument Reference

The following arguments are supported:

* `autoscaling_group_name` - (Required) Name or ARN of the Auto Scaling group.
* `scheduled_action_name` - (Required) Name of this scaling action.
* `start_time` - (Optional) Time for this action to start, in "YYYY-MM-DDThh:mm:ssZ" format in UTC/GMT only (for example, 2014-06-01T00:00:00Z ).
                            If you try to schedule your action in the past, Auto Scaling returns an error message.
* `end_time` - (Optional) Time for this action to end, in "YYYY-MM-DDThh:mm:ssZ" format in UTC/GMT only (for example, 2014-06-01T00:00:00Z ).
                          If you try to schedule your action in the past, Auto Scaling returns an error message.
* `recurrence` - (Optional) Time when recurring future actions will start. Start time is specified by the user following the Unix cron syntax format.
* `time_zone` - (Optional)  The timezone for the cron expression. Valid values are the canonical names of the IANA time zones (such as Etc/GMT+9 or Pacific/Tahiti).
* `min_size` - (Optional) Minimum size for the Auto Scaling group. Default 0.
Set to -1 if you don't want to change the minimum size at the scheduled time.
* `max_size` - (Optional) Maximum size for the Auto Scaling group. Default 0.
Set to -1 if you don't want to change the maximum size at the scheduled time.
* `desired_capacity` - (Optional) Number of EC2 instances that should be running in the group. Default 0.  Set to -1 if you don't want to change the desired capacity at the scheduled time.

~> **NOTE:** When `start_time` and `end_time` are specified with `recurrence` , they form the boundaries of when the recurring action will start and stop.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN assigned by AWS to the autoscaling schedule.

## Import

AutoScaling ScheduledAction can be imported using the `auto-scaling-group-name` and `scheduled-action-name`, e.g.,

```
$ terraform import aws_autoscaling_schedule.resource-name auto-scaling-group-name/scheduled-action-name
```
