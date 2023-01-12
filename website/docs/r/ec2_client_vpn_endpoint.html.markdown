---
subcategory: "VPN (Client)"
layout: "aws"
page_title: "AWS: aws_ec2_client_vpn_endpoint"
description: |-
  Provides an AWS Client VPN endpoint for OpenVPN clients.
---

# Resource: aws_ec2_client_vpn_endpoint

Provides an AWS Client VPN endpoint for OpenVPN clients. For more information on usage, please see the
[AWS Client VPN Administrator's Guide](https://docs.aws.amazon.com/vpn/latest/clientvpn-admin/what-is.html).

~> **NOTE on Client VPN endpoint target network security groups:** this provider provides both a standalone Client VPN endpoint network association resource with a (deprecated) `security_groups` argument and a Client VPN endpoint resource with a `security_group_ids` argument. Do not specify security groups in both resources. Doing so will cause a conflict and will overwrite the target network security group association.

## Example Usage

```terraform
resource "aws_ec2_client_vpn_endpoint" "example" {
  description            = "clientvpn-example"
  server_certificate_arn = aws_acm_certificate.cert.arn
  client_cidr_block      = "10.0.0.0/16"

  authentication_options {
    type                       = "certificate-authentication"
    root_certificate_chain_arn = aws_acm_certificate.root_cert.arn
  }

  connection_log_options {
    enabled               = true
    cloudwatch_log_group  = aws_cloudwatch_log_group.lg.name
    cloudwatch_log_stream = aws_cloudwatch_log_stream.ls.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `authentication_options` - (Required) Information about the authentication method to be used to authenticate clients.
* `client_cidr_block` - (Required) The IPv4 address range, in CIDR notation, from which to assign client IP addresses. The address range cannot overlap with the local CIDR of the VPC in which the associated subnet is located, or the routes that you add manually. The address range cannot be changed after the Client VPN endpoint has been created. The CIDR block should be /22 or greater.
* `client_connect_options` - (Optional) The options for managing connection authorization for new client connections.
* `client_login_banner_options` - (Optional) Options for enabling a customizable text banner that will be displayed on AWS provided clients when a VPN session is established.
* `connection_log_options` - (Required) Information about the client connection logging options.
* `description` - (Optional) A brief description of the Client VPN endpoint.
* `dns_servers` - (Optional) Information about the DNS servers to be used for DNS resolution. A Client VPN endpoint can have up to two DNS servers. If no DNS server is specified, the DNS address of the connecting device is used.
* `security_group_ids` - (Optional) The IDs of one or more security groups to apply to the target network. You must also specify the ID of the VPC that contains the security groups.
* `self_service_portal` - (Optional) Specify whether to enable the self-service portal for the Client VPN endpoint. Values can be `enabled` or `disabled`. Default value is `disabled`.
* `server_certificate_arn` - (Required) The ARN of the ACM server certificate.
* `session_timeout_hours` - (Optional) The maximum session duration is a trigger by which end-users are required to re-authenticate prior to establishing a VPN session. Default value is `24` - Valid values: `8 | 10 | 12 | 24`
* `split_tunnel` - (Optional) Indicates whether split-tunnel is enabled on VPN endpoint. Default value is `false`.
* `tags` - (Optional) A mapping of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `transport_protocol` - (Optional) The transport protocol to be used by the VPN session. Default value is `udp`.
* `vpc_id` - (Optional) The ID of the VPC to associate with the Client VPN endpoint. If no security group IDs are specified in the request, the default security group for the VPC is applied.
* `vpn_port` - (Optional) The port number for the Client VPN endpoint. Valid values are `443` and `1194`. Default value is `443`.

### `authentication_options` Argument Reference

One of the following arguments must be supplied:

* `active_directory_id` - (Optional) The ID of the Active Directory to be used for authentication if type is `directory-service-authentication`.
* `root_certificate_chain_arn` - (Optional) The ARN of the client certificate. The certificate must be signed by a certificate authority (CA) and it must be provisioned in AWS Certificate Manager (ACM). Only necessary when type is set to `certificate-authentication`.
* `saml_provider_arn` - (Optional) The ARN of the IAM SAML identity provider if type is `federated-authentication`.
* `self_service_saml_provider_arn` - (Optional) The ARN of the IAM SAML identity provider for the self service portal if type is `federated-authentication`.
* `type` - (Required) The type of client authentication to be used. Specify `certificate-authentication` to use certificate-based authentication, `directory-service-authentication` to use Active Directory authentication, or `federated-authentication` to use Federated Authentication via SAML 2.0.

### `client_connect_options` Argument reference

* `enabled` - (Optional) Indicates whether client connect options are enabled. The default is `false` (not enabled).
* `lambda_function_arn` - (Optional) The Amazon Resource Name (ARN) of the Lambda function used for connection authorization.

### `client_login_banner_options` Argument reference

* `banner_text` - (Optional) Customizable text that will be displayed in a banner on AWS provided clients when a VPN session is established. UTF-8 encoded characters only. Maximum of 1400 characters.
* `enabled` - (Optional) Enable or disable a customizable text banner that will be displayed on AWS provided clients when a VPN session is established. The default is `false` (not enabled).

### `connection_log_options` Argument Reference

One of the following arguments must be supplied:

* `cloudwatch_log_group` - (Optional) The name of the CloudWatch Logs log group.
* `cloudwatch_log_stream` - (Optional) The name of the CloudWatch Logs log stream to which the connection data is published.
* `enabled` - (Required) Indicates whether connection logging is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the Client VPN endpoint.
* `dns_name` - The DNS name to be used by clients when establishing their VPN session.
* `id` - The ID of the Client VPN endpoint.
* `status` - **Deprecated** The current state of the Client VPN endpoint.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

AWS Client VPN endpoints can be imported using the `id` value found via `aws ec2 describe-client-vpn-endpoints`, e.g.,

```
$ terraform import aws_ec2_client_vpn_endpoint.example cvpn-endpoint-0ac3a1abbccddd666
```
