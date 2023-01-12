---
subcategory: "CodeStar Connections"
layout: "aws"
page_title: "AWS: aws_codestarconnections_connection"
description: |-
  Provides details about CodeStar Connection
---

# Data Source: aws_codestarconnections_connection

Provides details about CodeStar Connection.

## Example Usage

### By ARN

```terraform
data "aws_codestarconnections_connection" "example" {
  arn = aws_codestarconnections_connection.example.arn
}
```

### By Name

```terraform
data "aws_codestarconnections_connection" "example" {
  name = aws_codestarconnections_connection.example.name
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Optional) CodeStar Connection ARN.
* `name` - (Optional) CodeStar Connection name.

~> **NOTE:** When both `arn` and `name` are specified, `arn` takes precedence.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `connection_status` - CodeStar Connection status. Possible values are `PENDING`, `AVAILABLE` and `ERROR`.
* `id` - CodeStar Connection ARN.
* `host_arn` - ARN of the host associated with the connection.
* `name` - Name of the CodeStar Connection. The name is unique in the calling AWS account.
* `provider_type` - Name of the external provider where your third-party code repository is configured. Possible values are `Bitbucket` and `GitHub`. For connections to a GitHub Enterprise Server instance, you must create an aws_codestarconnections_host resource and use `host_arn` instead.
* `tags` - Map of key-value resource tags to associate with the resource.
