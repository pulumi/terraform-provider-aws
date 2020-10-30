---
subcategory: "API Gateway (REST APIs)"
layout: "aws"
page_title: "AWS: aws_api_gateway_deployment"
description: |-
  Provides an API Gateway REST Deployment.
---

# Resource: aws_api_gateway_deployment

Provides an API Gateway REST Deployment.

~> **Note:** This resource depends on having at least one `aws_api_gateway_integration` created in the REST API, which 
itself has other dependencies. To avoid race conditions when all resources are being created together, you need to add 
implicit resource references via the `triggers` argument or explicit resource references using the 
[resource `dependsOn` meta-argument](https://www.pulumi.com/docs/intro/concepts/programming-model/#dependson).

## Example Usage

```hcl
resource "aws_api_gateway_rest_api" "MyDemoAPI" {
  name        = "MyDemoAPI"
  description = "This is my API for demonstration purposes"
}

resource "aws_api_gateway_resource" "MyDemoResource" {
  rest_api_id = aws_api_gateway_rest_api.MyDemoAPI.id
  parent_id   = aws_api_gateway_rest_api.MyDemoAPI.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "MyDemoMethod" {
  rest_api_id   = aws_api_gateway_rest_api.MyDemoAPI.id
  resource_id   = aws_api_gateway_resource.MyDemoResource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "MyDemoIntegration" {
  rest_api_id = aws_api_gateway_rest_api.MyDemoAPI.id
  resource_id = aws_api_gateway_resource.MyDemoResource.id
  http_method = aws_api_gateway_method.MyDemoMethod.http_method
  type        = "MOCK"
}

resource "aws_api_gateway_deployment" "MyDemoDeployment" {
  depends_on = [aws_api_gateway_integration.MyDemoIntegration]

  rest_api_id = aws_api_gateway_rest_api.MyDemoAPI.id
  stage_name  = "test"

  variables = {
    "answer" = "42"
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

### Redeployment Triggers

```hcl
resource "aws_api_gateway_deployment" "MyDemoDeployment" {
  rest_api_id = aws_api_gateway_rest_api.MyDemoAPI.id
  stage_name  = "test"

  triggers = {
    redeployment = sha1(join(",", list(
      jsonencode(aws_api_gateway_integration.example),
    )))
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `rest_api_id` - (Required) The ID of the associated REST API
* `stage_name` - (Optional) The name of the stage. If the specified stage already exists, it will be updated to point to the new deployment. If the stage does not exist, a new one will be created and point to this deployment.
* `description` - (Optional) The description of the deployment
* `stage_description` - (Optional) The description of the stage
* `triggers` - (Optional) A map of arbitrary keys and values that, when changed, will trigger a redeployment.
* `variables` - (Optional) A map that defines variables for the stage

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the deployment
* `invoke_url` - The URL to invoke the API pointing to the stage,
  e.g. `https://z4675bid1j.execute-api.eu-west-2.amazonaws.com/prod`
* `execution_arn` - The execution ARN to be used in `lambda_permission` resource's `source_arn`
  when allowing API Gateway to invoke a Lambda function,
  e.g. `arn:aws:execute-api:eu-west-2:123456789012:z4675bid1j/prod`
* `created_date` - The creation date of the deployment
