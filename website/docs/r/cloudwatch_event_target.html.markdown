---
subcategory: "EventBridge"
layout: "aws"
page_title: "AWS: aws_cloudwatch_event_target"
description: |-
  Provides an EventBridge Target resource.
---

# Resource: aws_cloudwatch_event_target

Provides an EventBridge Target resource.

~> **Note:** EventBridge was formerly known as CloudWatch Events. The functionality is identical.

## Example Usage

### Kinesis Usage

```terraform
resource "aws_cloudwatch_event_target" "yada" {
  target_id = "Yada"
  rule      = aws_cloudwatch_event_rule.console.name
  arn       = aws_kinesis_stream.test_stream.arn

  run_command_targets {
    key    = "tag:Name"
    values = ["FooBar"]
  }

  run_command_targets {
    key    = "InstanceIds"
    values = ["i-162058cd308bffec2"]
  }
}

resource "aws_cloudwatch_event_rule" "console" {
  name        = "capture-ec2-scaling-events"
  description = "Capture all EC2 scaling events"

  event_pattern = <<PATTERN
{
  "source": [
    "aws.autoscaling"
  ],
  "detail-type": [
    "EC2 Instance Launch Successful",
    "EC2 Instance Terminate Successful",
    "EC2 Instance Launch Unsuccessful",
    "EC2 Instance Terminate Unsuccessful"
  ]
}
PATTERN
}

resource "aws_kinesis_stream" "test_stream" {
  name        = "kinesis-test"
  shard_count = 1
}
```

### SSM Document Usage

```terraform
data "aws_iam_policy_document" "ssm_lifecycle_trust" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ssm_lifecycle" {
  statement {
    effect    = "Allow"
    actions   = ["ssm:SendCommand"]
    resources = ["arn:aws:ec2:eu-west-1:1234567890:instance/*"]

    condition {
      test     = "StringEquals"
      variable = "ec2:ResourceTag/Terminate"
      values   = ["*"]
    }
  }

  statement {
    effect    = "Allow"
    actions   = ["ssm:SendCommand"]
    resources = [aws_ssm_document.stop_instance.arn]
  }
}

resource "aws_iam_role" "ssm_lifecycle" {
  name               = "SSMLifecycle"
  assume_role_policy = data.aws_iam_policy_document.ssm_lifecycle_trust.json
}

resource "aws_iam_policy" "ssm_lifecycle" {
  name   = "SSMLifecycle"
  policy = data.aws_iam_policy_document.ssm_lifecycle.json
}

resource "aws_iam_role_policy_attachment" "ssm_lifecycle" {
  policy_arn = aws_iam_policy.ssm_lifecycle.arn
  role       = aws_iam_role.ssm_lifecycle.name
}

resource "aws_ssm_document" "stop_instance" {
  name          = "stop_instance"
  document_type = "Command"

  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Stop an instance",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["halt"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_cloudwatch_event_rule" "stop_instances" {
  name                = "StopInstance"
  description         = "Stop instances nightly"
  schedule_expression = "cron(0 0 * * ? *)"
}

resource "aws_cloudwatch_event_target" "stop_instances" {
  target_id = "StopInstance"
  arn       = aws_ssm_document.stop_instance.arn
  rule      = aws_cloudwatch_event_rule.stop_instances.name
  role_arn  = aws_iam_role.ssm_lifecycle.arn

  run_command_targets {
    key    = "tag:Terminate"
    values = ["midnight"]
  }
}
```

### RunCommand Usage

```terraform
resource "aws_cloudwatch_event_rule" "stop_instances" {
  name                = "StopInstance"
  description         = "Stop instances nightly"
  schedule_expression = "cron(0 0 * * ? *)"
}

resource "aws_cloudwatch_event_target" "stop_instances" {
  target_id = "StopInstance"
  arn       = "arn:aws:ssm:${var.aws_region}::document/AWS-RunShellScript"
  input     = "{\"commands\":[\"halt\"]}"
  rule      = aws_cloudwatch_event_rule.stop_instances.name
  role_arn  = aws_iam_role.ssm_lifecycle.arn

  run_command_targets {
    key    = "tag:Terminate"
    values = ["midnight"]
  }
}
```

### ECS Run Task with Role and Task Override Usage

```terraform
resource "aws_iam_role" "ecs_events" {
  name = "ecs_events"

  assume_role_policy = <<DOC
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "events.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
DOC
}

resource "aws_iam_role_policy" "ecs_events_run_task_with_any_role" {
  name = "ecs_events_run_task_with_any_role"
  role = aws_iam_role.ecs_events.id

  policy = <<DOC
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "iam:PassRole",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "ecs:RunTask",
            "Resource": "${replace(aws_ecs_task_definition.task_name.arn, "/:\\d+$/", ":*")}"
        }
    ]
}
DOC
}

resource "aws_cloudwatch_event_target" "ecs_scheduled_task" {
  target_id = "run-scheduled-task-every-hour"
  arn       = aws_ecs_cluster.cluster_name.arn
  rule      = aws_cloudwatch_event_rule.every_hour.name
  role_arn  = aws_iam_role.ecs_events.arn

  ecs_target {
    task_count          = 1
    task_definition_arn = aws_ecs_task_definition.task_name.arn
  }

  input = <<DOC
{
  "containerOverrides": [
    {
      "name": "name-of-container-to-override",
      "command": ["bin/console", "scheduled-task"]
    }
  ]
}
DOC
}
```

### API Gateway target

```terraform
resource "aws_cloudwatch_event_target" "example" {
  arn  = "${aws_api_gateway_stage.example.execution_arn}/GET"
  rule = aws_cloudwatch_event_rule.example.id

  http_target {
    query_string_parameters = {
      Body = "$.detail.body"
    }
    header_parameters = {
      Env = "Test"
    }
  }
}

resource "aws_cloudwatch_event_rule" "example" {
  # ...
}

resource "aws_api_gateway_deployment" "example" {
  rest_api_id = aws_api_gateway_rest_api.example.id
  # ...
}

resource "aws_api_gateway_stage" "example" {
  rest_api_id   = aws_api_gateway_rest_api.example.id
  deployment_id = aws_api_gateway_deployment.example.id
  # ...
}
```

### Cross-Account Event Bus target

```terraform
resource "aws_iam_role" "event_bus_invoke_remote_event_bus" {
  name               = "event-bus-invoke-remote-event-bus"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "events.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

data "aws_iam_policy_document" "event_bus_invoke_remote_event_bus" {
  statement {
    effect    = "Allow"
    actions   = ["events:PutEvents"]
    resources = ["arn:aws:events:eu-west-1:1234567890:event-bus/My-Event-Bus"]
  }
}

resource "aws_iam_policy" "event_bus_invoke_remote_event_bus" {
  name   = "event_bus_invoke_remote_event_bus"
  policy = data.aws_iam_policy_document.event_bus_invoke_remote_event_bus.json
}

resource "aws_iam_role_policy_attachment" "event_bus_invoke_remote_event_bus" {
  role       = aws_iam_role.event_bus_invoke_remote_event_bus.name
  policy_arn = aws_iam_policy.event_bus_invoke_remote_event_bus.arn
}

resource "aws_cloudwatch_event_rule" "stop_instances" {
  name                = "StopInstance"
  description         = "Stop instances nightly"
  schedule_expression = "cron(0 0 * * ? *)"
}

resource "aws_cloudwatch_event_target" "stop_instances" {
  target_id = "StopInstance"
  arn       = "arn:aws:events:eu-west-1:1234567890:event-bus/My-Event-Bus"
  rule      = aws_cloudwatch_event_rule.stop_instances.name
  role_arn  = aws_iam_role.event_bus_invoke_remote_event_bus.arn
}
```

### Input Transformer Usage - JSON Object

```terraform
resource "aws_cloudwatch_event_target" "example" {
  arn  = aws_lambda_function.example.arn
  rule = aws_cloudwatch_event_rule.example.id

  input_transformer {
    input_paths = {
      instance = "$.detail.instance",
      status   = "$.detail.status",
    }
    input_template = <<EOF
{
  "instance_id": <instance>,
  "instance_status": <status>
}
EOF
  }
}

resource "aws_cloudwatch_event_rule" "example" {
  # ...
}
```

### Input Transformer Usage - Simple String

```terraform
resource "aws_cloudwatch_event_target" "example" {
  arn  = aws_lambda_function.example.arn
  rule = aws_cloudwatch_event_rule.example.id

  input_transformer {
    input_paths = {
      instance = "$.detail.instance",
      status   = "$.detail.status",
    }
    input_template = "\"<instance> is in state <status>\""
  }
}

resource "aws_cloudwatch_event_rule" "example" {
  # ...
}
```

### Cloudwatch Log Group Usage

```terraform
resource "aws_cloudwatch_log_group" "example" {
  name              = "/aws/events/guardduty/logs"
  retention_in_days = 1
}

data "aws_iam_policy_document" "example_log_policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "${aws_cloudwatch_log_group.example.arn}:*"
    ]

    principals {
      identifiers = ["events.amazonaws.com", "delivery.logs.amazonaws.com"]
      type        = "Service"
    }

    condition {
      test     = "ArnEquals"
      values   = [aws_cloudwatch_event_rule.example.arn]
      variable = "aws:SourceArn"
    }
  }
}

resource "aws_cloudwatch_log_resource_policy" "example" {
  policy_document = data.aws_iam_policy_document.example_log_policy.json
  policy_name     = "guardduty-log-publishing-policy"
}

resource "aws_cloudwatch_event_rule" "example" {
  name        = "guard-duty_event_rule"
  description = "GuardDuty Findings"

  event_pattern = jsonencode(
    {
      "source" : [
        "aws.guardduty"
      ]
    }
  )

  tags = {
    Environment = "example"
  }
}

resource "aws_cloudwatch_event_target" "example" {
  rule = aws_cloudwatch_event_rule.example.name
  arn  = aws_cloudwatch_log_group.example.arn
}
```

## Argument Reference

-> **Note:** In order to be able to have your AWS Lambda function or
   SNS topic invoked by an EventBridge rule, you must set up the right permissions
   using `aws_lambda_permission`
   or `aws_sns_topic.policy`.
   More info [here](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-use-resource-based.html).

The following arguments are supported:

* `rule` - (Required) The name of the rule you want to add targets to.
* `event_bus_name` - (Optional) The event bus to associate with the rule. If you omit this, the `default` event bus is used.
* `target_id` - (Optional) The unique target assignment ID.  If missing, will generate a random, unique id.
* `arn` - (Required) The Amazon Resource Name (ARN) of the target.
* `input` - (Optional) Valid JSON text passed to the target. Conflicts with `input_path` and `input_transformer`.
* `input_path` - (Optional) The value of the [JSONPath](http://goessner.net/articles/JsonPath/) that is used for extracting part of the matched event when passing it to the target. Conflicts with `input` and `input_transformer`.
* `role_arn` - (Optional) The Amazon Resource Name (ARN) of the IAM role to be used for this target when the rule is triggered. Required if `ecs_target` is used or target in `arn` is EC2 instance, Kinesis data stream, Step Functions state machine, or Event Bus in different account or region.
* `run_command_targets` - (Optional) Parameters used when you are using the rule to invoke Amazon EC2 Run Command. Documented below. A maximum of 5 are allowed.
* `ecs_target` - (Optional) Parameters used when you are using the rule to invoke Amazon ECS Task. Documented below. A maximum of 1 are allowed.
* `batch_target` - (Optional) Parameters used when you are using the rule to invoke an Amazon Batch Job. Documented below. A maximum of 1 are allowed.
* `kinesis_target` - (Optional) Parameters used when you are using the rule to invoke an Amazon Kinesis Stream. Documented below. A maximum of 1 are allowed.
* `redshift_target` - (Optional) Parameters used when you are using the rule to invoke an Amazon Redshift Statement. Documented below. A maximum of 1 are allowed.
* `sqs_target` - (Optional) Parameters used when you are using the rule to invoke an Amazon SQS Queue. Documented below. A maximum of 1 are allowed.
* `http_target` - (Optional) Parameters used when you are using the rule to invoke an API Gateway REST endpoint. Documented below. A maximum of 1 is allowed.
* `input_transformer` - (Optional) Parameters used when you are providing a custom input to a target based on certain event data. Conflicts with `input` and `input_path`.
* `retry_policy` - (Optional)  Parameters used when you are providing retry policies. Documented below. A maximum of 1 are allowed.
* `dead_letter_config` - (Optional)  Parameters used when you are providing a dead letter config. Documented below. A maximum of 1 are allowed.

### run_command_targets

* `key` - (Required) Can be either `tag:tag-key` or `InstanceIds`.
* `values` - (Required) If Key is `tag:tag-key`, Values is a list of tag values. If Key is `InstanceIds`, Values is a list of Amazon EC2 instance IDs.

### ecs_target

* `capacity_provider_strategy` - (Optional) The capacity provider strategy to use for the task. If a `capacity_provider_strategy` specified, the `launch_type` parameter must be omitted. If no `capacity_provider_strategy` or `launch_type` is specified, the default capacity provider strategy for the cluster is used. Can be one or more. See below.
* `group` - (Optional) Specifies an ECS task group for the task. The maximum length is 255 characters.
* `launch_type` - (Optional) Specifies the launch type on which your task is running. The launch type that you specify here must match one of the launch type (compatibilities) of the target task. Valid values include: `EC2`, `EXTERNAL`, or `FARGATE`.
* `network_configuration` - (Optional) Use this if the ECS task uses the awsvpc network mode. This specifies the VPC subnets and security groups associated with the task, and whether a public IP address is to be used. Required if launch_type is FARGATE because the awsvpc mode is required for Fargate tasks.
* `platform_version` - (Optional) Specifies the platform version for the task. Specify only the numeric portion of the platform version, such as 1.1.0. This is used only if LaunchType is FARGATE. For more information about valid platform versions, see [AWS Fargate Platform Versions](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html).
* `task_count` - (Optional) The number of tasks to create based on the TaskDefinition. The default is 1.
* `task_definition_arn` - (Required) The ARN of the task definition to use if the event target is an Amazon ECS cluster.
* `tags` - (Optional) A map of tags to assign to ecs resources.
* `propagate_tags` - (Optional) Specifies whether to propagate the tags from the task definition to the task. If no value is specified, the tags are not propagated. Tags can only be propagated to the task during task creation.
* `placement_constraint` - (Optional) An array of placement constraint objects to use for the task. You can specify up to 10 constraints per task (including constraints in the task definition and those specified at runtime). See Below.
* `enable_execute_command` - (Optional) Whether or not to enable the execute command functionality for the containers in this task. If true, this enables execute command functionality on all containers in the task.
* `enable_ecs_managed_tags` - (Optional) Specifies whether to enable Amazon ECS managed tags for the task.

### capacity_provider_strategy

The `capacity_provider_strategy` configuration block supports the following:

* `base` - (Optional) The base value designates how many tasks, at a minimum, to run on the specified capacity provider. Only one capacity provider in a capacity provider strategy can have a base defined. If no value is specified, the default value of 0 is used.
* `capacity_provider` - (Required) Short name of the capacity provider.
* `weight` - (Required) The weight value designates the relative percentage of the total number of tasks launched that should use the specified capacity provider. The weight value is taken into consideration after the base value, if defined, is satisfied.

#### network_configuration

* `subnets` - (Required) The subnets associated with the task or service.
* `security_groups` - (Optional) The security groups associated with the task or service. If you do not specify a security group, the default security group for the VPC is used.
* `assign_public_ip` - (Optional) Assign a public IP address to the ENI (Fargate launch type only). Valid values are `true` or `false`. Default `false`.

For more information, see [Task Networking](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-networking.html)

#### placement_constraint

* `type` - (Required) Type of constraint. The only valid values at this time are `memberOf` and `distinctInstance`.
* `expression` -  (Optional) Cluster Query Language expression to apply to the constraint. Does not need to be specified for the `distinctInstance` type. For more information, see [Cluster Query Language in the Amazon EC2 Container Service Developer Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cluster-query-language.html).

### batch_target

* `job_definition` - (Required) The ARN or name of the job definition to use if the event target is an AWS Batch job. This job definition must already exist.
* `job_name` - (Required) The name to use for this execution of the job, if the target is an AWS Batch job.
* `array_size` - (Optional) The size of the array, if this is an array batch job. Valid values are integers between 2 and 10,000.
* `job_attempts` - (Optional) The number of times to attempt to retry, if the job fails. Valid values are 1 to 10.

### kinesis_target

* `partition_key_path` - (Optional) The JSON path to be extracted from the event and used as the partition key.

### redshift_target

* `database` - (Required) The name of the database.
* `db_user` - (Optional) The database user name.
* `secrets_manager_arn` - (Optional) The name or ARN of the secret that enables access to the database.
* `sql` - (Optional) The SQL statement text to run.
* `statement_name` - (Optional) The name of the SQL statement.
* `with_event` - (Optional) Indicates whether to send an event back to EventBridge after the SQL statement runs.

### sqs_target

* `message_group_id` - (Optional) The FIFO message group ID to use as the target.

`http_target`support the following:

* `path_parameter_values` - (Optional) The list of values that correspond sequentially to any path variables in your endpoint ARN (for example `arn:aws:execute-api:us-east-1:123456:myapi/*/POST/pets/*`).
* `query_string_parameters` - (Optional) Represents keys/values of query string parameters that are appended to the invoked endpoint.
* `header_parameters` - (Optional) Enables you to specify HTTP headers to add to the request.

### input_transformer

* `input_paths` - (Optional) Key value pairs specified in the form of JSONPath (for example, time = $.time)
    * You can have as many as 100 key-value pairs.
    * You must use JSON dot notation, not bracket notation.
    * The keys can't start with "AWS".

* `input_template` - (Required) Template to customize data sent to the target. Must be valid JSON. To send a string value, the string value must include double quotes.

### retry_policy

* `maximum_event_age_in_seconds` - (Optional) The age in seconds to continue to make retry attempts.
* `maximum_retry_attempts` - (Optional) maximum number of retry attempts to make before the request fails

### dead_letter_config

* `arn` - (Optional) - ARN of the SQS queue specified as the target for the dead-letter queue.

## Attributes Reference

No additional attributes are exported.

## Import

EventBridge Targets can be imported using `event_bus_name/rule-name/target-id` (if you omit `event_bus_name`, the `default` event bus will be used).

 ```
$ terraform import aws_cloudwatch_event_target.test-event-target rule-name/target-id
```
