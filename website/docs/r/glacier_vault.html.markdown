---
subcategory: "S3 Glacier"
layout: "aws"
page_title: "AWS: aws_glacier_vault"
description: |-
  Provides a Glacier Vault.
---

# Resource: aws_glacier_vault

Provides a Glacier Vault Resource. You can refer to the [Glacier Developer Guide](https://docs.aws.amazon.com/amazonglacier/latest/dev/working-with-vaults.html) for a full explanation of the Glacier Vault functionality

~> **NOTE:** When removing a Glacier Vault, the Vault must be empty.

## Example Usage

```terraform
resource "aws_sns_topic" "aws_sns_topic" {
  name = "glacier-sns-topic"
}

resource "aws_glacier_vault" "my_archive" {
  name = "MyArchive"

  notification {
    sns_topic = aws_sns_topic.aws_sns_topic.arn
    events    = ["ArchiveRetrievalCompleted", "InventoryRetrievalCompleted"]
  }

  access_policy = <<EOF
{
    "Version":"2012-10-17",
    "Statement":[
       {
          "Sid": "add-read-only-perm",
          "Principal": "*",
          "Effect": "Allow",
          "Action": [
             "glacier:InitiateJob",
             "glacier:GetJobOutput"
          ],
          "Resource": "arn:aws:glacier:eu-west-1:432981146916:vaults/MyArchive"
       }
    ]
}
EOF

  tags = {
    Test = "MyArchive"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Vault. Names can be between 1 and 255 characters long and the valid characters are a-z, A-Z, 0-9, '_' (underscore), '-' (hyphen), and '.' (period).
* `access_policy` - (Optional) The policy document. This is a JSON formatted string.
  The heredoc syntax or `file` function is helpful here. Use the [Glacier Developer Guide](https://docs.aws.amazon.com/amazonglacier/latest/dev/vault-access-policy.html) for more information on Glacier Vault Policy
* `notification` - (Optional) The notifications for the Vault. Fields documented below.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

**notification** supports the following:

* `events` - (Required) You can configure a vault to publish a notification for `ArchiveRetrievalCompleted` and `InventoryRetrievalCompleted` events.
* `sns_topic` - (Required) The SNS Topic ARN.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `location` - The URI of the vault that was created.
* `arn` - The ARN of the vault.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Glacier Vaults can be imported using the `name`, e.g.,

```
$ terraform import aws_glacier_vault.archive my_archive
```
