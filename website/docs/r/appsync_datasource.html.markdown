---
subcategory: "AppSync"
layout: "aws"
page_title: "AWS: aws_appsync_datasource"
description: |-
  Provides an AppSync Data Source.
---

# Resource: aws_appsync_datasource

Provides an AppSync Data Source.

## Example Usage

```terraform
resource "aws_dynamodb_table" "example" {
  name           = "example"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "UserId"

  attribute {
    name = "UserId"
    type = "S"
  }
}

resource "aws_iam_role" "example" {
  name = "example"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "appsync.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "example" {
  name = "example"
  role = aws_iam_role.example.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "dynamodb:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_dynamodb_table.example.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_appsync_graphql_api" "example" {
  authentication_type = "API_KEY"
  name                = "my_appsync_example"
}

resource "aws_appsync_datasource" "example" {
  api_id           = aws_appsync_graphql_api.example.id
  name             = "my_appsync_example"
  service_role_arn = aws_iam_role.example.arn
  type             = "AMAZON_DYNAMODB"

  dynamodb_config {
    table_name = aws_dynamodb_table.example.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) API ID for the GraphQL API for the data source.
* `name` - (Required) User-supplied name for the data source.
* `type` - (Required) Type of the Data Source. Valid values: `AWS_LAMBDA`, `AMAZON_DYNAMODB`, `AMAZON_ELASTICSEARCH`, `HTTP`, `NONE`, `RELATIONAL_DATABASE`.
* `description` - (Optional) Description of the data source.
* `service_role_arn` - (Optional) IAM service role ARN for the data source.
* `dynamodb_config` - (Optional) DynamoDB settings. See [below](#dynamodb_config)
* `elasticsearch_config` - (Optional) Amazon Elasticsearch settings. See [below](#elasticsearch_config)
* `http_config` - (Optional) HTTP settings. See [below](#http_config)
* `lambda_config` - (Optional) AWS Lambda settings. See [below](#lambda_config)
* `relational_database_config` (Optional) AWS RDS settings. See [Relational Database Config](#relational_database_config)

### dynamodb_config

The following arguments are supported:

* `table_name` - (Required) Name of the DynamoDB table.
* `region` - (Optional) AWS region of the DynamoDB table. Defaults to current region.
* `use_caller_credentials` - (Optional) Set to `true` to use Amazon Cognito credentials with this data source.

### elasticsearch_config

The following arguments are supported:

* `endpoint` - (Required) HTTP endpoint of the Elasticsearch domain.
* `region` - (Optional) AWS region of Elasticsearch domain. Defaults to current region.

### http_config

The following arguments are supported:

* `endpoint` - (Required) HTTP URL.
* `authorization_config` - (Optional) Authorization configuration in case the HTTP endpoint requires authorization. See [Authorization Config](#authorization_config).

#### authorization_config

The following arguments are supported:

* `authorization_type` - (Optional) Authorization type that the HTTP endpoint requires. Default values is `AWS_IAM`.
* `aws_iam_config` - (Optional) Identity and Access Management (IAM) settings. See [AWS IAM Config](#aws_iam_config).

##### aws_iam_config

The following arguments are supported:

* `signing_region` - (Optional) Signing Amazon Web Services Region for IAM authorization.
* `signing_service_name`- (Optional) Signing service name for IAM authorization.

### relational_database_config

The following arguments are supported:

* `http_endpoint_config` - (Required) Amazon RDS HTTP endpoint configuration. See [HTTP Endpoint Config](#http_endpoint_config).
* `source_type` - (Optional) Source type for the relational database. Valid values: `RDS_HTTP_ENDPOINT`.

#### http_endpoint_config

The following arguments are supported:

* `db_cluster_identifier` - (Required) Amazon RDS cluster identifier.
* `aws_secret_store_arn` - (Required) AWS secret store ARN for database credentials.
* `database_name` - (Optional) Logical database name.
* `region` - (Optional) AWS Region for RDS HTTP endpoint. Defaults to current region.
* `schema` - (Optional) Logical schema name.

### lambda_config

The following arguments are supported:

* `function_arn` - (Required) ARN for the Lambda function.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN

## Import

`aws_appsync_datasource` can be imported with their `api_id`, a hyphen, and `name`, e.g.,

```
$ terraform import aws_appsync_datasource.example abcdef123456-example
```
