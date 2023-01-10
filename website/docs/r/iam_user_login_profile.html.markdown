---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_user_login_profile"
description: |-
  Manages an IAM User Login Profile
---

# Resource: aws_iam_user_login_profile

Manages an IAM User Login Profile with limited support for password creation during this provider resource creation. Uses PGP to encrypt the password for safe transport to the user. PGP keys can be obtained from Keybase.

-> To reset an IAM User login password via this provider, you can use delete and recreate this resource or change any of the arguments.

## Example Usage

```terraform
resource "aws_iam_user" "example" {
  name          = "example"
  path          = "/"
  force_destroy = true
}

resource "aws_iam_user_login_profile" "example" {
  user    = aws_iam_user.example.name
  pgp_key = "keybase:some_person_that_exists"
}

output "password" {
  value = aws_iam_user_login_profile.example.encrypted_password
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) The IAM user's name.
* `pgp_key` - (Optional) Either a base-64 encoded PGP public key, or a keybase username in the form `keybase:username`. Only applies on resource creation. Drift detection is not possible with this argument.
* `password_length` - (Optional) The length of the generated password on resource creation. Only applies on resource creation. Drift detection is not possible with this argument. Default value is `20`.
* `password_reset_required` - (Optional) Whether the user should be forced to reset the generated password on resource creation. Only applies on resource creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `password` - The plain text password, only available when `pgp_key` is not provided.
* `key_fingerprint` - The fingerprint of the PGP key used to encrypt the password. Only available if password was handled on this provider resource creation, not import.
* `encrypted_password` - The encrypted password, base64 encoded. Only available if password was handled on resource creation, not import.

~> **NOTE:** The encrypted password may be decrypted using the command line,
   for example: `pulumi stack output password | base64 --decode | keybase pgp decrypt`.

## Import

IAM User Login Profiles can be imported without password information support via the IAM User name, e.g.,

```sh
$ terraform import aws_iam_user_login_profile.example myusername
```

Since this provider has no method to read the PGP or password information during import, use [`ignore_changes` argument](https://www.pulumi.com/docs/intro/concepts/programming-model/#ignorechanges) to ignore them unless password recreation is desired. e.g.

```terraform
resource "aws_iam_user_login_profile" "example" {
  # ... other configuration ...

  lifecycle {
    ignore_changes = [
      password_length,
      password_reset_required,
      pgp_key,
    ]
  }
}
```
