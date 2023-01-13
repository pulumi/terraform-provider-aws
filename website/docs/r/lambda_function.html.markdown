---
subcategory: "Lambda"
layout: "aws"
page_title: "AWS: aws_lambda_function"
description: |-
  Provides a Lambda Function resource. Lambda allows you to trigger execution of code in response to events in AWS, enabling serverless backend solutions. The Lambda Function itself includes source code and runtime configuration.
---

# Resource: aws_lambda_function

Provides a Lambda Function resource. Lambda allows you to trigger execution of code in response to events in AWS, enabling serverless backend solutions. The Lambda Function itself includes source code and runtime configuration.

For information about Lambda and how to use it, see [What is AWS Lambda?][1]

~> **NOTE:** Due to [AWS Lambda improved VPC networking changes that began deploying in September 2019](https://aws.amazon.com/blogs/compute/announcing-improved-vpc-networking-for-aws-lambda-functions/), EC2 subnets and security groups associated with Lambda Functions can take up to 45 minutes to successfully delete.

-> To give an external source (like an EventBridge Rule, SNS, or S3) permission to access the Lambda function, use the `aws_lambda_permission` resource. See [Lambda Permission Model][4] for more details. On the other hand, the `role` argument of this resource is the function's execution role for identity and access to AWS services and resources.

## Example Usage

### Basic Example

```terraform
resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "test_lambda" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "lambda_function_payload.zip"
  function_name = "lambda_function_name"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.test"

  runtime = "nodejs16.x"

  environment {
    variables = {
      foo = "bar"
    }
  }
}
```

### Lambda Layers

```terraform
resource "aws_lambda_layer_version" "example" {
  # ... other configuration ...
}

resource "aws_lambda_function" "example" {
  # ... other configuration ...
  layers = [aws_lambda_layer_version.example.arn]
}
```

### Lambda Ephemeral Storage

Lambda Function Ephemeral Storage(`/tmp`) allows you to configure the storage upto `10` GB. The default value set to `512` MB.

```terraform
resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "test_lambda" {
  filename      = "lambda_function_payload.zip"
  function_name = "lambda_function_name"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.test"
  runtime       = "nodejs14.x"

  ephemeral_storage {
    size = 10240 # Min 512 MB and the Max 10240 MB
  }
}
```

### Lambda File Systems

Lambda File Systems allow you to connect an Amazon Elastic File System (EFS) file system to a Lambda function to share data across function invocations, access existing data including large files, and save function state.

```terraform
# A lambda function connected to an EFS file system
resource "aws_lambda_function" "example" {
  # ... other configuration ...

  file_system_config {
    # EFS file system access point ARN
    arn = aws_efs_access_point.access_point_for_lambda.arn

    # Local mount path inside the lambda function. Must start with '/mnt/'.
    local_mount_path = "/mnt/efs"
  }

  vpc_config {
    # Every subnet should be able to reach an EFS mount target in the same Availability Zone. Cross-AZ mounts are not permitted.
    subnet_ids         = [aws_subnet.subnet_for_lambda.id]
    security_group_ids = [aws_security_group.sg_for_lambda.id]
  }

  # Explicitly declare dependency on EFS mount target.
  # When creating or updating Lambda functions, mount target must be in 'available' lifecycle state.
  depends_on = [aws_efs_mount_target.alpha]
}

# EFS file system
resource "aws_efs_file_system" "efs_for_lambda" {
  tags = {
    Name = "efs_for_lambda"
  }
}

# Mount target connects the file system to the subnet
resource "aws_efs_mount_target" "alpha" {
  file_system_id  = aws_efs_file_system.efs_for_lambda.id
  subnet_id       = aws_subnet.subnet_for_lambda.id
  security_groups = [aws_security_group.sg_for_lambda.id]
}

# EFS access point used by lambda file system
resource "aws_efs_access_point" "access_point_for_lambda" {
  file_system_id = aws_efs_file_system.efs_for_lambda.id

  root_directory {
    path = "/lambda"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "777"
    }
  }

  posix_user {
    gid = 1000
    uid = 1000
  }
}
```

### Lambda retries

Lambda Functions allow you to configure error handling for asynchronous invocation. The settings that it supports are `Maximum age of event` and `Retry attempts` as stated in [Lambda documentation for Configuring error handling for asynchronous invocation](https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#invocation-async-errors). To configure these settings, refer to the aws_lambda_function_event_invoke_config resource.

## CloudWatch Logging and Permissions

For more information about CloudWatch Logs for Lambda, see the [Lambda User Guide](https://docs.aws.amazon.com/lambda/latest/dg/monitoring-functions-logs.html).

```terraform
variable "lambda_function_name" {
  default = "lambda_function_name"
}

resource "aws_lambda_function" "test_lambda" {
  function_name = var.lambda_function_name

  # ... other configuration ...
  depends_on = [
    aws_iam_role_policy_attachment.lambda_logs,
    aws_cloudwatch_log_group.example,
  ]
}

# This is to optionally manage the CloudWatch Log Group for the Lambda Function.
# If skipping this resource configuration, also add "logs:CreateLogGroup" to the IAM policy below.
resource "aws_cloudwatch_log_group" "example" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = 14
}

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}
```

## Specifying the Deployment Package

AWS Lambda expects source code to be provided as a deployment package whose structure varies depending on which `runtime` is in use. See [Runtimes][6] for the valid values of `runtime`. The expected structure of the deployment package can be found in [the AWS Lambda documentation for each runtime][8].

Once you have created your deployment package you can specify it either directly as a local file (using the `filename` argument) or indirectly via Amazon S3 (using the `s3_bucket`, `s3_key` and `s3_object_version` arguments). When providing the deployment package via S3 it may be useful to use the `aws_s3_object` resource to upload it.

For larger deployment packages it is recommended by Amazon to upload via S3, since the S3 API has better support for uploading large files efficiently.

## Argument Reference

The following arguments are required:

* `function_name` - (Required) Unique name for your Lambda Function.
* `role` - (Required) Amazon Resource Name (ARN) of the function's execution role. The role provides the function's identity and access to AWS services and resources.

The following arguments are optional:

* `architectures` - (Optional) Instruction set architecture for your Lambda function. Valid values are `["x86_64"]` and `["arm64"]`. Default is `["x86_64"]`. Removing this attribute, function's architecture stay the same.
* `code_signing_config_arn` - (Optional) To enable code signing for this function, specify the ARN of a code-signing configuration. A code-signing configuration includes a set of signing profiles, which define the trusted publishers for this function.
* `dead_letter_config` - (Optional) Configuration block. Detailed below.
* `description` - (Optional) Description of what your Lambda Function does.
* `environment` - (Optional) Configuration block. Detailed below.
* `ephemeral_storage` - (Optional) The amount of Ephemeral storage(`/tmp`) to allocate for the Lambda Function in MB. This parameter is used to expand the total amount of Ephemeral storage available, beyond the default amount of `512`MB. Detailed below.
* `file_system_config` - (Optional) Configuration block. Detailed below.
* `filename` - (Optional) Path to the function's deployment package within the local filesystem. Conflicts with `image_uri`, `s3_bucket`, `s3_key`, and `s3_object_version`.
* `handler` - (Optional) Function [entrypoint][3] in your code.
* `image_config` - (Optional) Configuration block. Detailed below.
* `image_uri` - (Optional) ECR image URI containing the function's deployment package. Conflicts with `filename`, `s3_bucket`, `s3_key`, and `s3_object_version`.
* `kms_key_arn` - (Optional) Amazon Resource Name (ARN) of the AWS Key Management Service (KMS) key that is used to encrypt environment variables. If this configuration is not provided when environment variables are in use, AWS Lambda uses a default service key. If this configuration is provided when environment variables are not in use, the AWS Lambda API does not save this configuration and the provider will show a perpetual difference of adding the key. To fix the perpetual difference, remove this configuration.
* `layers` - (Optional) List of Lambda Layer Version ARNs (maximum of 5) to attach to your Lambda Function. See [Lambda Layers][10]
* `memory_size` - (Optional) Amount of memory in MB your Lambda Function can use at runtime. Defaults to `128`. See [Limits][5]
* `package_type` - (Optional) Lambda deployment package type. Valid values are `Zip` and `Image`. Defaults to `Zip`.
* `publish` - (Optional) Whether to publish creation/change as new Lambda Function Version. Defaults to `false`.
* `reserved_concurrent_executions` - (Optional) Amount of reserved concurrent executions for this lambda function. A value of `0` disables lambda from being triggered and `-1` removes any concurrency limitations. Defaults to Unreserved Concurrency Limits `-1`. See [Managing Concurrency][9]
* `runtime` - (Optional) Identifier of the function's runtime. See [Runtimes][6] for valid values.
* `s3_bucket` - (Optional) S3 bucket location containing the function's deployment package. Conflicts with `filename` and `image_uri`. This bucket must reside in the same AWS region where you are creating the Lambda function.
* `s3_key` - (Optional) S3 key of an object containing the function's deployment package. Conflicts with `filename` and `image_uri`.
* `s3_object_version` - (Optional) Object version containing the function's deployment package. Conflicts with `filename` and `image_uri`.
* `source_code_hash` - (Optional) Used to trigger updates. Must be set to a base64-encoded SHA256 hash of the package file specified with either `filename` or `s3_key`.
* `snap_start` - (Optional) Snap start settings block. Detailed below.
* `tags` - (Optional) Map of tags to assign to the object. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `timeout` - (Optional) Amount of time your Lambda Function has to run in seconds. Defaults to `3`. See [Limits][5].
* `tracing_config` - (Optional) Configuration block. Detailed below.
* `vpc_config` - (Optional) Configuration block. Detailed below.

### dead_letter_config

Dead letter queue configuration that specifies the queue or topic where Lambda sends asynchronous events when they fail processing. For more information, see [Dead Letter Queues](https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#dlq).

* `target_arn` - (Required) ARN of an SNS topic or SQS queue to notify when an invocation fails. If this option is used, the function's IAM role must be granted suitable access to write to the target object, which means allowing either the `sns:Publish` or `sqs:SendMessage` action on this ARN, depending on which service is targeted.

### environment

* `variables` - (Optional) Map of environment variables that are accessible from the function code during execution.

### ephemeral_storage

* `size` - (Required) The size of the Lambda function Ephemeral storage(`/tmp`) represented in MB. The minimum supported `ephemeral_storage` value defaults to `512`MB and the maximum supported value is `10240`MB.

### file_system_config

Connection settings for an EFS file system. Before creating or updating Lambda functions with `file_system_config`, EFS mount targets must be in available lifecycle state. Use `depends_on` to explicitly declare this dependency. See [Using Amazon EFS with Lambda][12].

* `arn` - (Required) Amazon Resource Name (ARN) of the Amazon EFS Access Point that provides access to the file system.
* `local_mount_path` - (Required) Path where the function can access the file system, starting with /mnt/.

### image_config

Container image configuration values that override the values in the container image Dockerfile.

* `command` - (Optional) Parameters that you want to pass in with `entry_point`.
* `entry_point` - (Optional) Entry point to your application, which is typically the location of the runtime executable.
* `working_directory` - (Optional) Working directory.

### snap_start

Snap start settings for low-latency startups. This feature is currently only supported for `java11` runtimes. Remove this block to delete the associated settings (rather than setting `apply_on = "None"`).

* `apply_on` - (Required) Conditions where snap start is enabled. Valid values are `PublishedVersions`.

### tracing_config

* `mode` - (Required) Whether to sample and trace a subset of incoming requests with AWS X-Ray. Valid values are `PassThrough` and `Active`. If `PassThrough`, Lambda will only trace the request from an upstream service if it contains a tracing header with "sampled=1". If `Active`, Lambda will respect any tracing header it receives from an upstream service. If no tracing header is received, Lambda will call X-Ray for a tracing decision.

### vpc_config

For network connectivity to AWS resources in a VPC, specify a list of security groups and subnets in the VPC. When you connect a function to a VPC, it can only access resources and the internet through that VPC. See [VPC Settings][7].

~> **NOTE:** If both `subnet_ids` and `security_group_ids` are empty then `vpc_config` is considered to be empty or unset.

* `security_group_ids` - (Required) List of security group IDs associated with the Lambda function.
* `subnet_ids` - (Required) List of subnet IDs associated with the Lambda function.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) identifying your Lambda Function.
* `invoke_arn` - ARN to be used for invoking Lambda Function from API Gateway - to be used in `aws_api_gateway_integration`'s `uri`.
* `last_modified` - Date this resource was last modified.
* `qualified_arn` - ARN identifying your Lambda Function Version (if versioning is enabled via `publish = true`).
* `qualified_invoke_arn` - Qualified ARN (ARN with lambda version number) to be used for invoking Lambda Function from API Gateway - to be used in `aws_api_gateway_integration`'s `uri`.
* `signing_job_arn` - ARN of the signing job.
* `signing_profile_version_arn` - ARN of the signing profile version.
* `snap_start.optimization_status` - Optimization status of the snap start configuration. Valid values are `On` and `Off`.
* `source_code_size` - Size in bytes of the function .zip file.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
* `version` - Latest published version of your Lambda Function.
* `vpc_config.vpc_id` - ID of the VPC.

[1]: https://docs.aws.amazon.com/lambda/latest/dg/welcome.html
[3]: https://docs.aws.amazon.com/lambda/latest/dg/walkthrough-custom-events-create-test-function.html
[4]: https://docs.aws.amazon.com/lambda/latest/dg/intro-permission-model.html
[5]: https://docs.aws.amazon.com/lambda/latest/dg/limits.html
[6]: https://docs.aws.amazon.com/lambda/latest/dg/API_CreateFunction.html#SSS-CreateFunction-request-Runtime
[7]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-vpc.html
[8]: https://docs.aws.amazon.com/lambda/latest/dg/deployment-package-v2.html
[9]: https://docs.aws.amazon.com/lambda/latest/dg/concurrent-executions.html
[10]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html
[12]: https://docs.aws.amazon.com/lambda/latest/dg/services-efs.html

## Timeouts

Configuration options:

* `create` - (Default `10m`)

## Import

Lambda Functions can be imported using the `function_name`, e.g.,

```
$ terraform import aws_lambda_function.test_lambda my_test_lambda_function
```
