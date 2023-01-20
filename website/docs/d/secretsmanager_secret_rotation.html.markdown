---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "AWS: aws_secretsmanager_secret_rotation"
description: |-
  Retrieve information about a Secrets Manager secret rotation configuration
---

# Data Source: aws_secretsmanager_secret_rotation

Retrieve information about a Secrets Manager secret rotation. To retrieve secret metadata, see the `aws_secretsmanager_secret` data source. To retrieve a secret value, see the `aws_secretsmanager_secret_version` data source.

## Example Usage

### Retrieve Secret Rotation Configuration

```terraform
data "aws_secretsmanager_secret_rotation" "example" {
  secret_id = data.aws_secretsmanager_secret.example.id
}
```

## Argument Reference

* `secret_id` - (Required) Specifies the secret containing the version that you want to retrieve. You can specify either the ARN or the friendly name of the secret.

## Attributes Reference

* `rotation_enabled` - ARN of the secret.
* `rotation_lambda_arn` - Decrypted part of the protected secret information that was originally provided as a string.
* `rotation_rules` - Decrypted part of the protected secret information that was originally provided as a binary. Base64 encoded.
