---
subcategory: "API Gateway V2"
layout: "aws"
page_title: "AWS: aws_apigatewayv2_deployment"
description: |-
  Manages an Amazon API Gateway Version 2 deployment.
---

# Resource: aws_apigatewayv2_deployment

Manages an Amazon API Gateway Version 2 deployment.
More information can be found in the [Amazon API Gateway Developer Guide](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-api.html).

-> **Note:** Creating a deployment for an API requires at least one `aws_apigatewayv2_route` resource associated with that API. To avoid race conditions when all resources are being created together, you need to add implicit resource references via the `triggers` argument or explicit resource references using the [resource `dependsOn` meta-argument](https://www.pulumi.com/docs/intro/concepts/programming-model/#dependson).

## Example Usage

### Basic

```terraform
resource "aws_apigatewayv2_deployment" "example" {
  api_id      = aws_apigatewayv2_api.example.id
  description = "Example deployment"

  lifecycle {
    create_before_destroy = true
  }
}
```

### Redeployment Triggers

```terraform
resource "aws_apigatewayv2_deployment" "example" {
  api_id      = aws_apigatewayv2_api.example.id
  description = "Example deployment"

  triggers = {
    redeployment = sha1(join(",", tolist([
      jsonencode(aws_apigatewayv2_integration.example),
      jsonencode(aws_apigatewayv2_route.example),
    ])))
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) API identifier.
* `description` - (Optional) Description for the deployment resource. Must be less than or equal to 1024 characters in length.
* `triggers` - (Optional) Map of arbitrary keys and values that, when changed, will trigger a redeployment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Deployment identifier.
* `auto_deployed` - Whether the deployment was automatically released.

## Import

`aws_apigatewayv2_deployment` can be imported by using the API identifier and deployment identifier, e.g.,

```
$ terraform import aws_apigatewayv2_deployment.example aabbccddee/1122334
```

The `triggers` argument cannot be imported.
