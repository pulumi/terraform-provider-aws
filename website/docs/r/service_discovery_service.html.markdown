---
subcategory: "Cloud Map"
layout: "aws"
page_title: "AWS: aws_service_discovery_service"
description: |-
  Provides a Service Discovery Service resource.
---

# Resource: aws_service_discovery_service

Provides a Service Discovery Service resource.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_service_discovery_private_dns_namespace" "example" {
  name        = "example.mydomain.local"
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
```

```terraform
resource "aws_service_discovery_public_dns_namespace" "example" {
  name        = "example.mydomain.com"
  description = "example"
}

resource "aws_service_discovery_service" "example" {
  name = "example"

  dns_config {
    namespace_id = aws_service_discovery_public_dns_namespace.example.id

    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_config {
    failure_threshold = 10
    resource_path     = "path"
    type              = "HTTP"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the service.
* `description` - (Optional) The description of the service.
* `dns_config` - (Optional) A complex type that contains information about the resource record sets that you want Amazon Route 53 to create when you register an instance.
* `health_check_config` - (Optional) A complex type that contains settings for an optional health check. Only for Public DNS namespaces.
* `force_destroy` - (Optional, Default:false ) A boolean that indicates all instances should be deleted from the service so that the service can be destroyed without error. These instances are not recoverable.
* `health_check_custom_config` - (Optional, ForceNew) A complex type that contains settings for ECS managed health checks.
* `namespace_id` - (Optional) The ID of the namespace that you want to use to create the service.
* `tags` - (Optional) A map of tags to assign to the service. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### dns_config

The following arguments are supported:

* `namespace_id` - (Required, ForceNew) The ID of the namespace to use for DNS configuration.
* `dns_records` - (Required) An array that contains one DnsRecord object for each resource record set.
* `routing_policy` - (Optional) The routing policy that you want to apply to all records that Route 53 creates when you register an instance and specify the service. Valid Values: MULTIVALUE, WEIGHTED

#### dns_records

The following arguments are supported:

* `ttl` - (Required) The amount of time, in seconds, that you want DNS resolvers to cache the settings for this resource record set.
* `type` - (Required, ForceNew) The type of the resource, which indicates the value that Amazon Route 53 returns in response to DNS queries. Valid Values: A, AAAA, SRV, CNAME

### health_check_config

The following arguments are supported:

* `failure_threshold` - (Optional) The number of consecutive health checks. Maximum value of 10.
* `resource_path` - (Optional) The path that you want Route 53 to request when performing health checks. Route 53 automatically adds the DNS name for the service. If you don't specify a value, the default value is /.
* `type` - (Optional, ForceNew) The type of health check that you want to create, which indicates how Route 53 determines whether an endpoint is healthy. Valid Values: HTTP, HTTPS, TCP

### health_check_custom_config

The following arguments are supported:

* `failure_threshold` - (Optional, ForceNew) The number of 30-second intervals that you want service discovery to wait before it changes the health status of a service instance.  Maximum value of 10.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service.
* `arn` - The ARN of the service.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Service Discovery Service can be imported using the service ID, e.g.,

```
$ terraform import aws_service_discovery_service.example 0123456789
```
