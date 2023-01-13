---
subcategory: "Cloud Map"
layout: "aws"
page_title: "AWS: aws_service_discovery_service"
description: |-
  Retrieves information about a Service Discovery Service
---

# Data Source: aws_service_discovery_service

Retrieves information about a Service Discovery Service.

## Example Usage

```hcl
data "aws_service_discovery_service" "test" {
  name         = "example"
  namespace_id = "NAMESPACE_ID_VALUE"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the service.
* `namespace_id` - (Required) ID of the namespace that the service belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the service.
* `arn` - ARN of the service.
* `description` - Description of the service.
* `dns_config` - Complex type that contains information about the resource record sets that you want Amazon Route 53 to create when you register an instance.
* `health_check_config` - Complex type that contains settings for an optional health check. Only for Public DNS namespaces.
* `health_check_custom_config` -  A complex type that contains settings for ECS managed health checks.
* `tags` - Map of tags to assign to the service. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

### dns_config

The following arguments are supported:

* `namespace_id` - ID of the namespace to use for DNS configuration.
* `dns_records` - An array that contains one DnsRecord object for each resource record set.
* `routing_policy` - Routing policy that you want to apply to all records that Route 53 creates when you register an instance and specify the service. Valid Values: MULTIVALUE, WEIGHTED

#### dns_records

The following arguments are supported:

* `ttl` - Amount of time, in seconds, that you want DNS resolvers to cache the settings for this resource record set.
* `type` - Type of the resource, which indicates the value that Amazon Route 53 returns in response to DNS queries. Valid Values: A, AAAA, SRV, CNAME

### health_check_config

The following arguments are supported:

* `failure_threshold` - Number of consecutive health checks. Maximum value of 10.
* `resource_path` - Path that you want Route 53 to request when performing health checks. Route 53 automatically adds the DNS name for the service. If you don't specify a value, the default value is /.
* `type` -  The type of health check that you want to create, which indicates how Route 53 determines whether an endpoint is healthy. Valid Values: HTTP, HTTPS, TCP

### health_check_custom_config

The following arguments are supported:

* `failure_threshold` -  The number of 30-second intervals that you want service discovery to wait before it changes the health status of a service instance.  Maximum value of 10.
