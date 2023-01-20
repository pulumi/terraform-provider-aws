---
subcategory: "KMS (Key Management)"
layout: "aws"
page_title: "AWS: aws_kms_ciphertext"
description: |-
    Provides ciphertext encrypted using a KMS key
---

# Resource: aws_kms_ciphertext

The KMS ciphertext resource allows you to encrypt plaintext into ciphertext
by using an AWS KMS customer master key. The value returned by this resource
is stable across every apply. For a changing ciphertext value each apply, see
the `aws_kms_ciphertext` data source.

## Example Usage

```terraform
resource "aws_kms_key" "oauth_config" {
  description = "oauth config"
  is_enabled  = true
}

resource "aws_kms_ciphertext" "oauth" {
  key_id = aws_kms_key.oauth_config.key_id

  plaintext = <<EOF
{
  "client_id": "e587dbae22222f55da22",
  "client_secret": "8289575d00000ace55e1815ec13673955721b8a5"
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `plaintext` - (Required) Data to be encrypted. Note that this may show up in logs, and it will be stored in the state file.
* `key_id` - (Required) Globally unique key ID for the customer master key.
* `context` - (Optional) An optional mapping that makes up the encryption context.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ciphertext_blob` - Base64 encoded ciphertext
