---
subcategory: "DynamoDB"
layout: "aws"
page_title: "AWS: aws_dynamodb_table"
description: |-
  Provides a DynamoDB table data source.
---

# Data Source: aws_dynamodb_table

Provides information about a DynamoDB table.

## Example Usage

```terraform
data "aws_dynamodb_table" "tableName" {
  name = "tableName"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the DynamoDB table.

## Attributes Reference

See the DynamoDB Table Resource for details on the
returned attributes - they are identical.
