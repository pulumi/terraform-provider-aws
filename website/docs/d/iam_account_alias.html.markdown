---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_account_alias"
description: |-
  Provides the account alias for the AWS account associated with the provider
  connection to AWS.
---

# Data Source: aws_iam_account_alias

The IAM Account Alias data source allows access to the account alias
for the effective account in which this provider is working.

## Example Usage

```terraform
data "aws_iam_account_alias" "current" {}

output "account_id" {
  value = data.aws_iam_account_alias.current.account_alias
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_alias` - Alias associated with the AWS account.
* `id` - Alias associated with the AWS account.
