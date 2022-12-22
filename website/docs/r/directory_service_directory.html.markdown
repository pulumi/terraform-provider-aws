---
subcategory: "DS (Directory Service)"
layout: "aws"
page_title: "AWS: aws_directory_service_directory"
description: |-
  Provides a directory in AWS Directory Service.
---

# Resource: aws_directory_service_directory

Provides a Simple or Managed Microsoft directory in AWS Directory Service.

## Example Usage

### SimpleAD

```terraform
resource "aws_directory_service_directory" "bar" {
  name     = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  size     = "Small"

  vpc_settings {
    vpc_id     = aws_vpc.main.id
    subnet_ids = [aws_subnet.foo.id, aws_subnet.bar.id]
  }

  tags = {
    Project = "foo"
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-west-2a"
  cidr_block        = "10.0.1.0/24"
}

resource "aws_subnet" "bar" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-west-2b"
  cidr_block        = "10.0.2.0/24"
}
```

### Microsoft Active Directory (MicrosoftAD)

```terraform
resource "aws_directory_service_directory" "bar" {
  name     = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  edition  = "Standard"
  type     = "MicrosoftAD"

  vpc_settings {
    vpc_id     = aws_vpc.main.id
    subnet_ids = [aws_subnet.foo.id, aws_subnet.bar.id]
  }

  tags = {
    Project = "foo"
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-west-2a"
  cidr_block        = "10.0.1.0/24"
}

resource "aws_subnet" "bar" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-west-2b"
  cidr_block        = "10.0.2.0/24"
}
```

### Microsoft Active Directory Connector (ADConnector)

```terraform
resource "aws_directory_service_directory" "connector" {
  name     = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  size     = "Small"
  type     = "ADConnector"

  connect_settings {
    customer_dns_ips  = ["A.B.C.D"]
    customer_username = "Admin"
    subnet_ids        = [aws_subnet.foo.id, aws_subnet.bar.id]
    vpc_id            = aws_vpc.main.id
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-west-2a"
  cidr_block        = "10.0.1.0/24"
}

resource "aws_subnet" "bar" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-west-2b"
  cidr_block        = "10.0.2.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The fully qualified name for the directory, such as `corp.example.com`
* `password` - (Required) The password for the directory administrator or connector user.
* `size` - (Optional) (For `SimpleAD` and `ADConnector` types) The size of the directory (`Small` or `Large` are accepted values). `Large` by default.
* `vpc_settings` - (Required for `SimpleAD` and `MicrosoftAD`) VPC related information about the directory. Fields documented below.
* `connect_settings` - (Required for `ADConnector`) Connector related information about the directory. Fields documented below.
* `alias` - (Optional) The alias for the directory (must be unique amongst all aliases in AWS). Required for `enable_sso`.
* `description` - (Optional) A textual description for the directory.
* `desired_number_of_domain_controllers` - (Optional) The number of domain controllers desired in the directory. Minimum value of `2`. Scaling of domain controllers is only supported for `MicrosoftAD` directories.
* `short_name` - (Optional) The short name of the directory, such as `CORP`.
* `enable_sso` - (Optional) Whether to enable single-sign on for the directory. Requires `alias`. Defaults to `false`.
* `type` (Optional) - The directory type (`SimpleAD`, `ADConnector` or `MicrosoftAD` are accepted values). Defaults to `SimpleAD`.
* `edition` - (Optional, for type `MicrosoftAD` only) The MicrosoftAD edition (`Standard` or `Enterprise`). Defaults to `Enterprise`.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

**vpc_settings** supports the following:

* `subnet_ids` - (Required) The identifiers of the subnets for the directory servers (2 subnets in 2 different AZs).
* `vpc_id` - (Required) The identifier of the VPC that the directory is in.

**connect_settings** supports the following:

* `customer_username` - (Required) The username corresponding to the password provided.
* `customer_dns_ips` - (Required) The DNS IP addresses of the domain to connect to.
* `subnet_ids` - (Required) The identifiers of the subnets for the directory servers (2 subnets in 2 different AZs).
* `vpc_id` - (Required) The identifier of the VPC that the directory is in.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The directory identifier.
* `access_url` - The access URL for the directory, such as `http://alias.awsapps.com`.
* `dns_ip_addresses` - A list of IP addresses of the DNS servers for the directory or connector.
* `security_group_id` - The ID of the security group created by the directory.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

`connect_settings` (for `ADConnector`) is also exported with the following attributes:

* `connect_ips` - The IP addresses of the AD Connector servers.

## Timeouts

`aws_directory_service_directory` provides the following Timeouts configuration options:

- `create` - (Default `60 minutes`) Used for directory creation
- `update` - (Default `60 minutes`) Used for directory update
- `delete` - (Default `60 minutes`) Used for directory deletion

## Import

DirectoryService directories can be imported using the directory `id`, e.g.,

```
$ terraform import aws_directory_service_directory.sample d-926724cf57
```
