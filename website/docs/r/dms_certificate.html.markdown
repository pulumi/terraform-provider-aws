---
subcategory: "DMS (Database Migration)"
layout: "aws"
page_title: "AWS: aws_dms_certificate"
description: |-
  Provides a DMS (Data Migration Service) certificate resource.
---

# Resource: aws_dms_certificate

Provides a DMS (Data Migration Service) certificate resource. DMS certificates can be created, deleted, and imported.

## Example Usage

```terraform
# Create a new certificate
resource "aws_dms_certificate" "test" {
  certificate_id  = "test-dms-certificate-tf"
  certificate_pem = "..."

  tags = {
    Name = "test"
  }

}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required) The certificate identifier.

    - Must contain from 1 to 255 alphanumeric characters and hyphens.

* `certificate_pem` - (Optional) The contents of the .pem X.509 certificate file for the certificate. Either `certificate_pem` or `certificate_wallet` must be set.
* `certificate_wallet` - (Optional) The contents of the Oracle Wallet certificate for use with SSL, provided as a base64-encoded String. Either `certificate_pem` or `certificate_wallet` must be set.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certificate_arn` - The Amazon Resource Name (ARN) for the certificate.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Certificates can be imported using the `certificate_id`, e.g.,

```
$ terraform import aws_dms_certificate.test test-dms-certificate-tf
```
