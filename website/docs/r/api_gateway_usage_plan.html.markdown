---
subcategory: "API Gateway"
layout: "aws"
page_title: "AWS: aws_api_gateway_usage_plan"
description: |-
  Provides an API Gateway Usage Plan.
---

# Resource: aws_api_gateway_usage_plan

Provides an API Gateway Usage Plan.

## Example Usage

```terraform
resource "aws_api_gateway_rest_api" "example" {
  body = jsonencode({
    openapi = "3.0.1"
    info = {
      title   = "example"
      version = "1.0"
    }
    paths = {
      "/path1" = {
        get = {
          x-amazon-apigateway-integration = {
            httpMethod           = "GET"
            payloadFormatVersion = "1.0"
            type                 = "HTTP_PROXY"
            uri                  = "https://ip-ranges.amazonaws.com/ip-ranges.json"
          }
        }
      }
    }
  })

  name = "example"
}

resource "aws_api_gateway_deployment" "example" {
  rest_api_id = aws_api_gateway_rest_api.example.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.example.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "development" {
  deployment_id = aws_api_gateway_deployment.example.id
  rest_api_id   = aws_api_gateway_rest_api.example.id
  stage_name    = "development"
}

resource "aws_api_gateway_stage" "production" {
  deployment_id = aws_api_gateway_deployment.example.id
  rest_api_id   = aws_api_gateway_rest_api.example.id
  stage_name    = "production"
}

resource "aws_api_gateway_usage_plan" "example" {
  name         = "my-usage-plan"
  description  = "my description"
  product_code = "MYCODE"

  api_stages {
    api_id = aws_api_gateway_rest_api.example.id
    stage  = aws_api_gateway_stage.development.stage_name
  }

  api_stages {
    api_id = aws_api_gateway_rest_api.example.id
    stage  = aws_api_gateway_stage.production.stage_name
  }

  quota_settings {
    limit  = 20
    offset = 2
    period = "WEEK"
  }

  throttle_settings {
    burst_limit = 5
    rate_limit  = 10
  }
}
```

## Argument Reference

The API Gateway Usage Plan argument layout is a structure composed of several sub-resources - these resources are laid out below.

### Top-Level Arguments

* `name` - (Required) Name of the usage plan.
* `description` - (Optional) Description of a usage plan.
* `api_stages` - (Optional) Associated [API stages](#api-stages-arguments) of the usage plan.
* `quota_settings` - (Optional) The [quota settings](#quota-settings-arguments) of the usage plan.
* `throttle_settings` - (Optional) The [throttling limits](#throttling-settings-arguments) of the usage plan.
* `product_code` - (Optional) AWS Marketplace product identifier to associate with the usage plan as a SaaS product on AWS Marketplace.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

#### Api Stages arguments

* `api_id` (Required) - API Id of the associated API stage in a usage plan.
* `stage` (Required) - API stage name of the associated API stage in a usage plan.
* `throttle` - (Optional) The [throttling limits](#throttle) of the usage plan.

##### Throttle

* `path` (Required) - Method to apply the throttle settings for. Specfiy the path and method, for example `/test/GET`.
* `burst_limit` (Optional) - The API request burst limit, the maximum rate limit over a time ranging from one to a few seconds, depending upon whether the underlying token bucket is at its full capacity.
* `rate_limit` (Optional) - The API request steady-state rate limit.

#### Quota Settings Arguments

* `limit` (Optional) - Maximum number of requests that can be made in a given time period.
* `offset` (Optional) - Number of requests subtracted from the given limit in the initial time period.
* `period` (Optional) - Time period in which the limit applies. Valid values are "DAY", "WEEK" or "MONTH".

#### Throttling Settings Arguments

* `burst_limit` (Optional) - The API request burst limit, the maximum rate limit over a time ranging from one to a few seconds, depending upon whether the underlying token bucket is at its full capacity.
* `rate_limit` (Optional) - The API request steady-state rate limit.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API resource
* `name` - Name of the usage plan.
* `description` - Description of a usage plan.
* `api_stages` - Associated API stages of the usage plan.
* `quota_settings` - Quota of the usage plan.
* `throttle_settings` - Throttling limits of the usage plan.
* `product_code` - AWS Marketplace product identifier to associate with the usage plan as a SaaS product on AWS Marketplace.
* `arn` - ARN
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

AWS API Gateway Usage Plan can be imported using the `id`, e.g.,

```sh
$ terraform import aws_api_gateway_usage_plan.myusageplan <usage_plan_id>
```
