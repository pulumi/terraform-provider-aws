---
subcategory: "API Gateway (REST APIs)"
layout: "aws"
page_title: "AWS: aws_api_gateway_deployment"
description: |-
  Manages an API Gateway REST Deployment.
---

# Resource: aws_api_gateway_deployment

Manages an API Gateway REST Deployment. A deployment is a snapshot of the REST API configuration. The deployment can then be published to callable endpoints via the [`aws_api_gateway_stage` resource](api_gateway_stage.html) and optionally managed further with the [`aws_api_gateway_base_path_mapping` resource](api_gateway_base_path_mapping.html), [`aws_api_gateway_domain_name` resource](api_gateway_domain_name.html), and [`aws_api_method_settings` resource](api_gateway_method_settings.html). For more information, see the [API Gateway Developer Guide](https://docs.aws.amazon.com/apigateway/latest/developerguide/how-to-deploy-api.html).

To properly capture all REST API configuration in a deployment, this resource must have dependencies on all prior resources that manage resources/paths, methods, integrations, etc.

* For REST APIs that are configured via OpenAPI specification ([`aws_api_gateway_rest_api` resource](api_gateway_rest_api.html) `body` argument), no special dependency setup is needed beyond referencing the  `id` attribute of that resource unless additional resources have further customized the REST API.
* When the REST API configuration involves other resources (`aws_api_gateway_integration` resource), the dependency setup can be done with implicit resource references in the `triggers` argument or explicit resource references using the [resource `dependsOn` custom option](https://www.pulumi.com/docs/intro/concepts/resources/#dependson). The `triggers` argument should be preferred over `depends_on`, since `depends_on` can only capture dependency ordering and will not cause the resource to recreate (redeploy the REST API) with upstream configuration changes.

!> **WARNING:** It is recommended to use the [`aws_api_gateway_stage` resource](api_gateway_stage.html) instead of managing an API Gateway Stage via the `stage_name` argument of this resource. When this resource is recreated (REST API redeployment) with the `stage_name` configured, the stage is deleted and recreated. This will cause a temporary service interruption, increase provide plan differences, and can require a second apply to recreate any downstream stage configuration such as associated `aws_api_method_settings` resources.

## Example Usage

### OpenAPI Specification

```hcl
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

resource "aws_api_gateway_stage" "example" {
  deployment_id = aws_api_gateway_deployment.example.id
  rest_api_id   = aws_api_gateway_rest_api.example.id
  stage_name    = "example"
}
```

## Argument Reference

The following arguments are supported:

* `rest_api_id` - (Required) REST API identifier.
* `description` - (Optional) Description of the deployment
* `stage_name` - (Optional) Name of the stage to create with this deployment. If the specified stage already exists, it will be updated to point to the new deployment. It is recommended to use the [`aws_api_gateway_stage` resource](api_gateway_stage.html) instead to manage stages.
* `stage_description` - (Optional) Description to set on the stage managed by the `stage_name` argument.
* `triggers` - (Optional) Map of arbitrary keys and values that, when changed, will trigger a redeployment.
* `variables` - (Optional) Map to set on the stage managed by the `stage_name` argument.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the deployment
* `invoke_url` - The URL to invoke the API pointing to the stage,
  e.g. `https://z4675bid1j.execute-api.eu-west-2.amazonaws.com/prod`
* `execution_arn` - The execution ARN to be used in `lambda_permission` resource's `source_arn`
  when allowing API Gateway to invoke a Lambda function,
  e.g. `arn:aws:execute-api:eu-west-2:123456789012:z4675bid1j/prod`
* `created_date` - The creation date of the deployment
