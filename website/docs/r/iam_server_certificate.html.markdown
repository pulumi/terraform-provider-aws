---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_server_certificate"
description: |-
  Provides an IAM Server Certificate
---

# Resource: aws_iam_server_certificate

Provides an IAM Server Certificate resource to upload Server Certificates.
Certs uploaded to IAM can easily work with other AWS services such as:

- AWS Elastic Beanstalk
- Elastic Load Balancing
- CloudFront
- AWS OpsWorks

For information about server certificates in IAM, see [Managing Server
Certificates][2] in AWS Documentation.

## Example Usage

**Using certs on file:**

```terraform
resource "aws_iam_server_certificate" "test_cert" {
  name             = "some_test_cert"
  certificate_body = file("self-ca-cert.pem")
  private_key      = file("test-key.pem")
}
```

**Example with cert in-line:**

```terraform
resource "aws_iam_server_certificate" "test_cert_alt" {
  name = "alt_test_cert"

  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
[......] # cert contents
-----END CERTIFICATE-----
EOF

  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
[......] # cert contents
-----END RSA PRIVATE KEY-----
EOF
}
```

**Use in combination with an AWS ELB resource:**

Some properties of an IAM Server Certificates cannot be updated while they are
in use. In order for the provider to effectively manage a Certificate in this situation, it is
recommended you utilize the `name_prefix` attribute and enable the
`create_before_destroy`. This will allow this provider
to create a new, updated `aws_iam_server_certificate` resource and replace it in
dependant resources before attempting to destroy the old version.

```terraform
resource "aws_iam_server_certificate" "test_cert" {
  name_prefix      = "example-cert"
  certificate_body = file("self-ca-cert.pem")
  private_key      = file("test-key.pem")

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_elb" "ourapp" {
  name                      = "asg-deployment-example"
  availability_zones        = ["us-west-2a"]
  cross_zone_load_balancing = true

  listener {
    instance_port      = 8000
    instance_protocol  = "http"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = aws_iam_server_certificate.test_cert.arn
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the Server Certificate. Do not include the
  path in this value. If omitted, this provider will assign a random, unique name.
* `name_prefix` - (Optional) Creates a unique name beginning with the specified
  prefix. Conflicts with `name`.
* `certificate_body` – (Required) The contents of the public key certificate in
  PEM-encoded format.
* `certificate_chain` – (Optional) The contents of the certificate chain.
  This is typically a concatenation of the PEM-encoded public key certificates
  of the chain.
* `private_key` – (Required) The contents of the private key in PEM-encoded format.
* `path` - (Optional) The IAM path for the server certificate.  If it is not
    included, it defaults to a slash (/). If this certificate is for use with
    AWS CloudFront, the path must be in format `/cloudfront/your_path_here`.
    See [IAM Identifiers][1] for more details on IAM Paths.
* `tags` - (Optional) Map of resource tags for the server certificate. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

~> **NOTE:** AWS performs behind-the-scenes modifications to some certificate files if they do not adhere to a specific format. These modifications will result in this provider forever believing that it needs to update the resources since the local and AWS file contents will not match after theses modifications occur. In order to prevent this from happening you must ensure that all your PEM-encoded files use UNIX line-breaks and that `certificate_body` contains only one certificate. All other certificates should go in `certificate_chain`. It is common for some Certificate Authorities to issue certificate files that have DOS line-breaks and that are actually multiple certificates concatenated together in order to form a full certificate chain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) specifying the server certificate.
* `expiration` - Date and time in [RFC3339 format](https://tools.ietf.org/html/rfc3339#section-5.8) on which the certificate is set to expire.
* `id` - The unique Server Certificate name
* `name` - The name of the Server Certificate
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
* `upload_date` - Date and time in [RFC3339 format](https://tools.ietf.org/html/rfc3339#section-5.8) when the server certificate was uploaded.

## Import

IAM Server Certificates can be imported using the `name`, e.g.,

```
$ terraform import aws_iam_server_certificate.certificate example.com-certificate-until-2018
```

[1]: https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_Identifiers.html
[2]: https://docs.aws.amazon.com/IAM/latest/UserGuide/ManagingServerCerts.html
