---
subcategory: "Cloud Map"
layout: "aws"
page_title: "AWS: aws_service_discovery_instance"
description: |-
  Provides a Service Discovery Instance resource.
---

# Resource: aws_service_discovery_instance

Provides a Service Discovery Instance resource.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_service_discovery_private_dns_namespace" "example" {
  name        = "example.domain.local"
  description = "example"
  vpc         = aws_vpc.example.id
}

resource "aws_service_discovery_service" "example" {
  name = "example"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.example.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

resource "aws_service_discovery_instance" "example" {
  instance_id = "example-instance-id"
  service_id  = aws_service_discovery_service.example.id

  attributes = {
    AWS_INSTANCE_IPV4 = "172.18.0.1"
    custom_attribute  = "custom"
  }
}
```

```terraform
resource "aws_service_discovery_http_namespace" "example" {
  name        = "example.domain.test"
  description = "example"
}

resource "aws_service_discovery_service" "example" {
  name         = "example"
  namespace_id = aws_service_discovery_http_namespace.example.id
}

resource "aws_service_discovery_instance" "example" {
  instance_id = "example-instance-id"
  service_id  = aws_service_discovery_service.example.id

  attributes = {
    AWS_EC2_INSTANCE_ID = "i-0abdg374kd892cj6dl"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the service instance.
* `service_id` - (Required, ForceNew) The ID of the service that you want to use to create the instance.
* `attributes` - (Required) A map contains the attributes of the instance. Check the [doc](https://docs.aws.amazon.com/cloud-map/latest/api/API_RegisterInstance.html#API_RegisterInstance_RequestSyntax) for the supported attributes and syntax.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the instance.

## Import

Service Discovery Instance can be imported using the service ID and instance ID, e.g.,

```
$ terraform import aws_service_discovery_instance.example 0123456789/i-0123
```
