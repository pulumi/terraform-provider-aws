---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "AWS: aws_secretsmanager_secret"
description: |-
  Provides a resource to manage AWS Secrets Manager secret metadata
---

# Resource: aws_secretsmanager_secret

Provides a resource to manage AWS Secrets Manager secret metadata. To manage secret rotation, see the `aws_secretsmanager_secret_rotation` resource. To manage a secret value, see the `aws_secretsmanager_secret_version` resource.

## Example Usage

### Basic

```terraform
resource "aws_secretsmanager_secret" "example" {
  name = "example"
}
```

### Rotation Configuration

To enable automatic secret rotation, the Secrets Manager service requires usage of a Lambda function. The [Rotate Secrets section in the Secrets Manager User Guide](https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets.html) provides additional information about deploying a prebuilt Lambda functions for supported credential rotation (e.g., RDS) or deploying a custom Lambda function.

~> **NOTE:** Configuring rotation causes the secret to rotate once as soon as you store the secret. Before you do this, you must ensure that all of your applications that use the credentials stored in the secret are updated to retrieve the secret from AWS Secrets Manager. The old credentials might no longer be usable after the initial rotation and any applications that you fail to update will break as soon as the old credentials are no longer valid.

~> **NOTE:** If you cancel a rotation that is in progress (by removing the `rotation` configuration), it can leave the VersionStage labels in an unexpected state. Depending on what step of the rotation was in progress, you might need to remove the staging label AWSPENDING from the partially created version, specified by the SecretVersionId response value. You should also evaluate the partially rotated new version to see if it should be deleted, which you can do by removing all staging labels from the new version's VersionStage field.

```terraform
resource "aws_secretsmanager_secret" "rotation-example" {
  name                = "rotation-example"
  rotation_lambda_arn = aws_lambda_function.example.arn

  rotation_rules {
    automatically_after_days = 7
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the secret.
* `kms_key_id` - (Optional) ARN or Id of the AWS KMS key to be used to encrypt the secret values in the versions stored in this secret. If you don't specify this value, then Secrets Manager defaults to using the AWS account's default KMS key (the one named `aws/secretsmanager`). If the default KMS key with that name doesn't yet exist, then AWS Secrets Manager creates it for you automatically the first time.
* `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `name` - (Optional) Friendly name of the new secret. The secret name can consist of uppercase letters, lowercase letters, digits, and any of the following characters: `/_+=.@-` Conflicts with `name_prefix`.
* `policy` - (Optional) Valid JSON document representing a [resource policy](https://docs.aws.amazon.com/secretsmanager/latest/userguide/auth-and-access_resource-based-policies.html). Removing `policy` from your configuration or setting `policy` to null or an empty string (i.e., `policy = ""`) _will not_ delete the policy since it could have been set by `aws_secretsmanager_secret_policy`. To delete the `policy`, set it to `"{}"` (an empty JSON document).
* `recovery_window_in_days` - (Optional) Number of days that AWS Secrets Manager waits before it can delete the secret. This value can be `0` to force deletion without recovery or range from `7` to `30` days. The default value is `30`.
* `replica` - (Optional) Configuration block to support secret replication. See details below.
* `force_overwrite_replica_secret` - (Optional) Accepts boolean value to specify whether to overwrite a secret with the same name in the destination Region.
* `rotation_lambda_arn` - (Optional, **DEPRECATED**) ARN of the Lambda function that can rotate the secret. Use the `aws_secretsmanager_secret_rotation` resource to manage this configuration instead. As of version 2.67.0, removal of this configuration will no longer remove rotation due to supporting the new resource. Either import the new resource and remove the configuration or manually remove rotation.
* `rotation_rules` - (Optional, **DEPRECATED**) Configuration block for the rotation configuration of this secret. Defined below. Use the `aws_secretsmanager_secret_rotation` resource to manage this configuration instead. As of version 2.67.0, removal of this configuration will no longer remove rotation due to supporting the new resource. Either import the new resource and remove the configuration or manually remove rotation.
* `tags` - (Optional) Key-value map of user-defined tags that are attached to the secret. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### replica

* `kms_key_id` - (Optional) ARN, Key ID, or Alias of the AWS KMS key within the region secret is replicated to. If one is not specified, then Secrets Manager defaults to using the AWS account's default KMS key (`aws/secretsmanager`) in the region or creates one for use if non-existent.
* `region` - (Required) Region for replicating the secret.

### rotation_rules

* `automatically_after_days` - (Required) Specifies the number of days between automatic scheduled rotations of the secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ARN of the secret.
* `arn` - ARN of the secret.
* `rotation_enabled` - Whether automatic rotation is enabled for this secret.
* `replica` - Attributes of a replica are described below.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

### replica

* `last_accessed_date` - Date that you last accessed the secret in the Region.
* `status` - Status can be `InProgress`, `Failed`, or `InSync`.
* `status_message` - Message such as `Replication succeeded` or `Secret with this name already exists in this region`.

## Import

`aws_secretsmanager_secret` can be imported by using the secret Amazon Resource Name (ARN), e.g.,

```
$ terraform import aws_secretsmanager_secret.example arn:aws:secretsmanager:us-east-1:123456789012:secret:example-123456
```
