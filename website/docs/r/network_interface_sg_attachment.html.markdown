---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_network_interface_sg_attachment"
description: |-
  Associates a security group with a network interface.
---

# Resource: aws_network_interface_sg_attachment

This resource attaches a security group to an Elastic Network Interface (ENI).
It can be used to attach a security group to any existing ENI, be it a
secondary ENI or one attached as the primary interface on an instance.

~> **NOTE on instances, interfaces, and security groups:** This provider currently
provides the capability to assign security groups via the [`aws_instance`][1]
and the [`aws_network_interface`][2] resources. Using this resource in
conjunction with security groups provided in-line in those resources will cause
conflicts, and will lead to spurious diffs and undefined behavior - please use
one or the other.

## Example Usage

The following provides a very basic example of setting up an instance (provided
by `instance`) in the default security group, creating a security group
(provided by `sg`) and then attaching the security group to the instance's
primary network interface via the `aws_network_interface_sg_attachment` resource,
named `sg_attachment`:

```terraform
data "aws_ami" "ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-*"]
  }

  owners = ["amazon"]
}

resource "aws_instance" "instance" {
  instance_type = "t2.micro"
  ami           = data.aws_ami.ami.id

  tags = {
    type = "test-instance"
  }
}

resource "aws_security_group" "sg" {
  tags = {
    type = "test-security-group"
  }
}

resource "aws_network_interface_sg_attachment" "sg_attachment" {
  security_group_id    = aws_security_group.sg.id
  network_interface_id = aws_instance.instance.primary_network_interface_id
}
```

In this example, `instance` is provided by the `aws_instance` data source,
fetching an external instance, possibly not managed by this provider.
`sg_attachment` then attaches to the output instance's `network_interface_id`:

```terraform
data "aws_instance" "instance" {
  instance_id = "i-1234567890abcdef0"
}

resource "aws_security_group" "sg" {
  tags = {
    type = "test-security-group"
  }
}

resource "aws_network_interface_sg_attachment" "sg_attachment" {
  security_group_id    = aws_security_group.sg.id
  network_interface_id = data.aws_instance.instance.network_interface_id
}
```

## Argument Reference

* `security_group_id` - (Required) The ID of the security group.
* `network_interface_id` - (Required) The ID of the network interface to attach to.

## Attributes Reference

No additional attributes are exported.

## Import

Network Interface Security Group attachments can be imported using the associated network interface ID and security group ID, separated by an underscore (`_`).

For example:

```
$ terraform import aws_network_interface_sg_attachment.sg_attachment eni-1234567890abcdef0_sg-1234567890abcdef0
```
