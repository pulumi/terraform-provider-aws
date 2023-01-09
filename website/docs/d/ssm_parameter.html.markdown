---
subcategory: "SSM (Systems Manager)"
layout: "aws"
page_title: "AWS: aws_ssm_parameter"
description: |-
  Provides a SSM Parameter Data Source
---

# Data Source: aws_ssm_parameter

Provides an SSM Parameter data source.

## Example Usage

```terraform
data "aws_ssm_parameter" "foo" {
  name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the parameter.
* `with_decryption` - (Optional) Whether to return decrypted `SecureString` value. Defaults to `true`.

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the parameter.
* `name` - Name of the parameter.
* `type` - Type of the parameter. Valid types are `String`, `StringList` and `SecureString`.
* `value` - Value of the parameter. This value is always marked as sensitive in the plan output, regardless of `type`.
* `version` - Version of the parameter.
