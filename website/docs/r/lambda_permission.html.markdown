---
subcategory: "Lambda"
layout: "aws"
page_title: "AWS: aws_lambda_permission"
description: |-
  Creates a Lambda function permission.
---

# Resource: aws_lambda_permission

Gives an external source (like an EventBridge Rule, SNS, or S3) permission to access the Lambda function.

## Example Usage

```terraform
resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = "arn:aws:events:eu-west-1:111122223333:rule/RunDaily"
  qualifier     = aws_lambda_alias.test_alias.name
}

resource "aws_lambda_alias" "test_alias" {
  name             = "testalias"
  description      = "a sample description"
  function_name    = aws_lambda_function.test_lambda.function_name
  function_version = "$LATEST"
}

resource "aws_lambda_function" "test_lambda" {
  filename      = "lambdatest.zip"
  function_name = "lambda_function_name"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "exports.handler"
  runtime       = "nodejs16.x"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}
```

## Usage with SNS

```terraform
resource "aws_lambda_permission" "with_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.func.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.default.arn
}

resource "aws_sns_topic" "default" {
  name = "call-lambda-maybe"
}

resource "aws_sns_topic_subscription" "lambda" {
  topic_arn = aws_sns_topic.default.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.func.arn
}

resource "aws_lambda_function" "func" {
  filename      = "lambdatest.zip"
  function_name = "lambda_called_from_sns"
  role          = aws_iam_role.default.arn
  handler       = "exports.handler"
  runtime       = "python3.7"
}

resource "aws_iam_role" "default" {
  name = "iam_for_lambda_with_sns"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}
```

## Specify Lambda permissions for API Gateway REST API

```terraform
resource "aws_api_gateway_rest_api" "MyDemoAPI" {
  name        = "MyDemoAPI"
  description = "This is my API for demonstration purposes"
}

resource "aws_lambda_permission" "lambda_permission" {
  statement_id  = "AllowMyDemoAPIInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "MyDemoFunction"
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.MyDemoAPI.execution_arn}/*/*/*"
}
```

## Usage with CloudWatch log group

```terraform
resource "aws_lambda_permission" "logging" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.logging.function_name
  principal     = "logs.eu-west-1.amazonaws.com"
  source_arn    = "${aws_cloudwatch_log_group.default.arn}:*"
}

resource "aws_cloudwatch_log_group" "default" {
  name = "/default"
}

resource "aws_cloudwatch_log_subscription_filter" "logging" {
  depends_on      = [aws_lambda_permission.logging]
  destination_arn = aws_lambda_function.logging.arn
  filter_pattern  = ""
  log_group_name  = aws_cloudwatch_log_group.default.name
  name            = "logging_default"
}

resource "aws_lambda_function" "logging" {
  filename      = "lamba_logging.zip"
  function_name = "lambda_called_from_cloudwatch_logs"
  handler       = "exports.handler"
  role          = aws_iam_role.default.arn
  runtime       = "python3.7"
}

resource "aws_iam_role" "default" {
  name = "iam_for_lambda_called_from_cloudwatch_logs"

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
```

## Example function URL cross-account invoke policy

```terraform
resource "aws_lambda_function_url" "url" {
  function_name      = aws_lambda_function.example.function_name
  authorization_type = "AWS_IAM"
}

resource "aws_lambda_permission" "url" {
  action        = "lambda:InvokeFunctionUrl"
  function_name = aws_lambda_function.example.function_name
  principal     = "arn:aws:iam::444455556666:role/example"

  source_account         = "444455556666"
  function_url_auth_type = "AWS_IAM"

  # Adds the following condition keys
  # "Condition": {
  #      "StringEquals": {
  #        "AWS:SourceAccount": "444455556666",
  #        "lambda:FunctionUrlAuthType": "AWS_IAM"
  #      }
  #    }

}
```

## Argument Reference

* `action` - (Required) The AWS Lambda action you want to allow in this statement. (e.g., `lambda:InvokeFunction`)
* `event_source_token` - (Optional) The Event Source Token to validate.  Used with [Alexa Skills][1].
* `function_name` - (Required) Name of the Lambda function whose resource policy you are updating
* `function_url_auth_type` - (Optional) Lambda Function URLs [authentication type][3]. Valid values are: `AWS_IAM` or `NONE`. Only supported for `lambda:InvokeFunctionUrl` action.
* `principal` - (Required) The principal who is getting this permission e.g., `s3.amazonaws.com`, an AWS account ID, or AWS IAM principal, or AWS service principal such as `events.amazonaws.com` or `sns.amazonaws.com`.
* `qualifier` - (Optional) Query parameter to specify function version or alias name. The permission will then apply to the specific qualified ARN e.g., `arn:aws:lambda:aws-region:acct-id:function:function-name:2`
* `source_account` - (Optional) This parameter is used when allowing cross-account access, or for S3 and SES. The AWS account ID (without a hyphen) of the source owner.
* `source_arn` - (Optional) When the principal is an AWS service, the ARN of the specific resource within that service to grant permission to.
  Without this, any resource from `principal` will be granted permission – even if that resource is from another account.
  For S3, this should be the ARN of the S3 Bucket.
  For EventBridge events, this should be the ARN of the EventBridge Rule.
  For API Gateway, this should be the ARN of the API, as described [here][2].
* `statement_id` - (Optional) A unique statement identifier. By default generated by the provider.
* `statement_id_prefix` - (Optional) A statement identifier prefix. The provider will generate a unique suffix. Conflicts with `statement_id`.
* `principal_org_id` - (Optional) The identifier for your organization in AWS Organizations. Use this to grant permissions to all the AWS accounts under this organization.

[1]: https://developer.amazon.com/docs/custom-skills/host-a-custom-skill-as-an-aws-lambda-function.html#use-aws-cli
[2]: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
[3]: https://docs.aws.amazon.com/lambda/latest/dg/urls-auth.html

## Attributes Reference

No additional attributes are exported.

## Import

Lambda permission statements can be imported using function_name/statement_id, with an optional qualifier, e.g.,

```
$ terraform import aws_lambda_permission.test_lambda_permission my_test_lambda_function/AllowExecutionFromCloudWatch

$ terraform import aws_lambda_permission.test_lambda_permission my_test_lambda_function:qualifier_name/AllowExecutionFromCloudWatch
```
