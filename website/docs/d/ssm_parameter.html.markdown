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

~> **Note:** The data source is currently following the behavior of the [SSM API](https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#Parameter) to return a string value, regardless of parameter type.


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the parameter.
* `with_decryption` - (Optional) Whether to return decrypted `SecureString` value. Defaults to `true`.


In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the parameter.
* `name` - The name of the parameter.
* `type` - The type of the parameter. Valid types are `String`, `StringList` and `SecureString`.
* `value` - The value of the parameter. This value is always marked as sensitive in the plan output, regardless of `type`.
* `version` - The version of the parameter.
