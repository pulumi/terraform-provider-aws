---
subcategory: "DMS (Database Migration)"
layout: "aws"
page_title: "AWS: aws_dms_replication_task"
description: |-
  Provides a DMS (Data Migration Service) replication task resource.
---

# Resource: aws_dms_replication_task

Provides a DMS (Data Migration Service) replication task resource. DMS replication tasks can be created, updated, deleted, and imported.

~> **NOTE:** Changing most arguments will stop the task if it is running. You can set `start_replication_task` to resume the task afterwards.

## Example Usage

```terraform
# Create a new replication task
resource "aws_dms_replication_task" "test" {
  cdc_start_time            = 1484346880
  migration_type            = "full-load"
  replication_instance_arn  = aws_dms_replication_instance.test-dms-replication-instance-tf.replication_instance_arn
  replication_task_id       = "test-dms-replication-task-tf"
  replication_task_settings = "..."
  source_endpoint_arn       = aws_dms_endpoint.test-dms-source-endpoint-tf.endpoint_arn
  table_mappings            = "{\"rules\":[{\"rule-type\":\"selection\",\"rule-id\":\"1\",\"rule-name\":\"1\",\"object-locator\":{\"schema-name\":\"%\",\"table-name\":\"%\"},\"rule-action\":\"include\"}]}"

  tags = {
    Name = "test"
  }

  target_endpoint_arn = aws_dms_endpoint.test-dms-target-endpoint-tf.endpoint_arn
}
```

## Argument Reference

The following arguments are supported:

* `cdc_start_position` - (Optional, Conflicts with `cdc_start_time`) Indicates when you want a change data capture (CDC) operation to start. The value can be in date, checkpoint, or LSN/SCN format depending on the source engine. For more information, see [Determining a CDC native start point](https://docs.aws.amazon.com/dms/latest/userguide/CHAP_Task.CDC.html#CHAP_Task.CDC.StartPoint.Native).
* `cdc_start_time` - (Optional, Conflicts with `cdc_start_position`) The Unix timestamp integer for the start of the Change Data Capture (CDC) operation.
* `migration_type` - (Required) The migration type. Can be one of `full-load | cdc | full-load-and-cdc`.
* `replication_instance_arn` - (Required) The Amazon Resource Name (ARN) of the replication instance.
* `replication_task_id` - (Required) The replication task identifier.

    - Must contain from 1 to 255 alphanumeric characters or hyphens.
    - First character must be a letter.
    - Cannot end with a hyphen.
    - Cannot contain two consecutive hyphens.

* `replication_task_settings` - (Optional) An escaped JSON string that contains the task settings. For a complete list of task settings, see [Task Settings for AWS Database Migration Service Tasks](http://docs.aws.amazon.com/dms/latest/userguide/CHAP_Tasks.CustomizingTasks.TaskSettings.html).
* `source_endpoint_arn` - (Required) The Amazon Resource Name (ARN) string that uniquely identifies the source endpoint.
* `start_replication_task` - (Optional) Whether to run or stop the replication task.
* `status` - Replication Task status.
* `table_mappings` - (Required) An escaped JSON string that contains the table mappings. For information on table mapping see [Using Table Mapping with an AWS Database Migration Service Task to Select and Filter Data](http://docs.aws.amazon.com/dms/latest/userguide/CHAP_Tasks.CustomizingTasks.TableMapping.html)
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `target_endpoint_arn` - (Required) The Amazon Resource Name (ARN) string that uniquely identifies the target endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `replication_task_arn` - The Amazon Resource Name (ARN) for the replication task.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Replication tasks can be imported using the `replication_task_id`, e.g.,

```
$ terraform import aws_dms_replication_task.test test-dms-replication-task-tf
```
