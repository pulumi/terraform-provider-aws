---
subcategory: "API Gateway V2"
layout: "aws"
page_title: "AWS: aws_apigatewayv2_integration_response"
description: |-
  Manages an Amazon API Gateway Version 2 integration response.
---

# Resource: aws_apigatewayv2_integration_response

Manages an Amazon API Gateway Version 2 integration response.
More information can be found in the [Amazon API Gateway Developer Guide](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-api.html).

## Example Usage

### Basic

```terraform
resource "aws_apigatewayv2_integration_response" "example" {
  api_id                   = aws_apigatewayv2_api.example.id
  integration_id           = aws_apigatewayv2_integration.example.id
  integration_response_key = "/200/"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) API identifier.
* `integration_id` - (Required) Identifier of the `aws_apigatewayv2_integration`.
* `integration_response_key` - (Required) Integration response key.
* `content_handling_strategy` - (Optional) How to handle response payload content type conversions. Valid values: `CONVERT_TO_BINARY`, `CONVERT_TO_TEXT`.
* `response_templates` - (Optional) Map of Velocity templates that are applied on the request payload based on the value of the Content-Type header sent by the client.
* `template_selection_expression` - (Optional) The [template selection expression](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-api-selection-expressions.html#apigateway-websocket-api-template-selection-expressions) for the integration response.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Integration response identifier.

## Import

`aws_apigatewayv2_integration_response` can be imported by using the API identifier, integration identifier and integration response identifier, e.g.,

```
$ terraform import aws_apigatewayv2_integration_response.example aabbccddee/1122334/998877
```
