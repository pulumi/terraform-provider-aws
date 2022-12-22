---
subcategory: "Cloud Control API"
layout: "aws"
page_title: "AWS: aws_cloudcontrolapi_resource"
description: |-
    Provides details for a Cloud Control API Resource.
---

# Data Source: aws_cloudcontrolapi_resource

Provides details for a Cloud Control API Resource. The reading of these resources is proxied through Cloud Control API handlers to the backend service.

## Example Usage

```terraform
data "aws_cloudcontrolapi_resource" "example" {
  identifier = "example"
  type_name  = "AWS::ECS::Cluster"
}
```

## Argument Reference

The following arguments are required:

* `identifier` - (Required) Identifier of the CloudFormation resource type. For example, `vpc-12345678`.
* `type_name` - (Required) CloudFormation resource type name. For example, `AWS::EC2::VPC`.

The following arguments are optional:

* `role_arn` - (Optional) ARN of the IAM Role to assume for operations.
* `type_version_id` - (Optional) Identifier of the CloudFormation resource type version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `properties` - JSON string matching the CloudFormation resource type schema with current configuration.
