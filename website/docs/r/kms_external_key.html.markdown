---
subcategory: "KMS (Key Management)"
layout: "aws"
page_title: "AWS: aws_kms_external_key"
description: |-
  Manages a single-Region or multi-Region primary KMS key that uses external key material.
---

# Resource: aws_kms_external_key

Manages a single-Region or multi-Region primary KMS key that uses external key material.
To instead manage a single-Region or multi-Region primary KMS key where AWS automatically generates and potentially rotates key material, see the `aws_kms_key` resource.


## Example Usage

```terraform
resource "aws_kms_external_key" "example" {
  description = "KMS EXTERNAL for AMI encryption"
}
```

## Argument Reference

The following arguments are supported:

* `bypass_policy_lockout_safety_check` - (Optional) Specifies whether to disable the policy lockout check performed when creating or updating the key's policy. Setting this value to `true` increases the risk that the key becomes unmanageable. For more information, refer to the scenario in the [Default Key Policy](https://docs.aws.amazon.com/kms/latest/developerguide/key-policies.html#key-policy-default-allow-root-enable-iam) section in the AWS Key Management Service Developer Guide. Defaults to `false`.
* `deletion_window_in_days` - (Optional) Duration in days after which the key is deleted after destruction of the resource. Must be between `7` and `30` days. Defaults to `30`.
* `description` - (Optional) Description of the key.
* `enabled` - (Optional) Specifies whether the key is enabled. Keys pending import can only be `false`. Imported keys default to `true` unless expired.
* `key_material_base64` - (Optional) Base64 encoded 256-bit symmetric encryption key material to import. The CMK is permanently associated with this key material. The same key material can be reimported, but you cannot import different key material.
* `multi_region` - (Optional) Indicates whether the KMS key is a multi-Region (`true`) or regional (`false`) key. Defaults to `false`.
* `policy` - (Optional) A key policy JSON document. If you do not provide a key policy, AWS KMS attaches a default key policy to the CMK.
* `tags` - (Optional) A key-value map of tags to assign to the key. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `valid_to` - (Optional) Time at which the imported key material expires. When the key material expires, AWS KMS deletes the key material and the CMK becomes unusable. If not specified, key material does not expire. Valid values: [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) (`YYYY-MM-DDTHH:MM:SSZ`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the key.
* `expiration_model` - Whether the key material expires. Empty when pending key material import, otherwise `KEY_MATERIAL_EXPIRES` or `KEY_MATERIAL_DOES_NOT_EXPIRE`.
* `id` - The unique identifier for the key.
* `key_state` - The state of the CMK.
* `key_usage` - The cryptographic operations for which you can use the CMK.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

KMS External Keys can be imported using the `id`, e.g.,

```
$ terraform import aws_kms_external_key.a arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab
```
