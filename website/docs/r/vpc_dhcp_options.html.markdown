---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_vpc_dhcp_options"
description: |-
  Provides a VPC DHCP Options resource.
---

# Resource: aws_vpc_dhcp_options

Provides a VPC DHCP Options resource.

## Example Usage

Basic usage:

```terraform
resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name_servers = ["8.8.8.8", "8.8.4.4"]
}
```

Full usage:

```terraform
resource "aws_vpc_dhcp_options" "foo" {
  domain_name          = "service.consul"
  domain_name_servers  = ["127.0.0.1", "10.0.0.2"]
  ntp_servers          = ["127.0.0.1"]
  netbios_name_servers = ["127.0.0.1"]
  netbios_node_type    = 2

  tags = {
    Name = "foo-name"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Optional) the suffix domain name to use by default when resolving non Fully Qualified Domain Names. In other words, this is what ends up being the `search` value in the `/etc/resolv.conf` file.
* `domain_name_servers` - (Optional) List of name servers to configure in `/etc/resolv.conf`. If you want to use the default AWS nameservers you should set this to `AmazonProvidedDNS`.
* `ntp_servers` - (Optional) List of NTP servers to configure.
* `netbios_name_servers` - (Optional) List of NETBIOS name servers.
* `netbios_node_type` - (Optional) The NetBIOS node type (1, 2, 4, or 8). AWS recommends to specify 2 since broadcast and multicast are not supported in their network. For more information about these node types, see [RFC 2132](http://www.ietf.org/rfc/rfc2132.txt).
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Remarks

* Notice that all arguments are optional but you have to specify at least one argument.
* `domain_name_servers`, `netbios_name_servers`, `ntp_servers` are limited by AWS to maximum four servers only.
* To actually use the DHCP Options Set you need to associate it to a VPC using `aws_vpc_dhcp_options_association`.
* If you delete a DHCP Options Set, all VPCs using it will be associated to AWS's `default` DHCP Option Set.
* In most cases unless you're configuring your own DNS you'll want to set `domain_name_servers` to `AmazonProvidedDNS`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the DHCP Options Set.
* `arn` - The ARN of the DHCP Options Set.
* `owner_id` - The ID of the AWS account that owns the DHCP options set.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

You can find more technical documentation about DHCP Options Set in the
official [AWS User Guide](https://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_DHCP_Options.html).

## Import

VPC DHCP Options can be imported using the `dhcp options id`, e.g.,

```
$ terraform import aws_vpc_dhcp_options.my_options dopt-d9070ebb
```
