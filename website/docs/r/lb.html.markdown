---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "AWS: aws_lb"
description: |-
  Provides a Load Balancer resource.
---

# Resource: aws_lb

Provides a Load Balancer resource.

~> **Note:** `aws_alb` is known as `aws_lb`. The functionality is identical.

## Example Usage

### Application Load Balancer

```terraform
resource "aws_lb" "test" {
  name               = "test-lb-tf"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = [for subnet in aws_subnet.public : subnet.id]

  enable_deletion_protection = true

  access_logs {
    bucket  = aws_s3_bucket.lb_logs.bucket
    prefix  = "test-lb"
    enabled = true
  }

  tags = {
    Environment = "production"
  }
}
```

### Network Load Balancer

```terraform
resource "aws_lb" "test" {
  name               = "test-lb-tf"
  internal           = false
  load_balancer_type = "network"
  subnets            = [for subnet in aws_subnet.public : subnet.id]

  enable_deletion_protection = true

  tags = {
    Environment = "production"
  }
}
```

### Specifying Elastic IPs

```terraform
resource "aws_lb" "example" {
  name               = "example"
  load_balancer_type = "network"

  subnet_mapping {
    subnet_id     = aws_subnet.example1.id
    allocation_id = aws_eip.example1.id
  }

  subnet_mapping {
    subnet_id     = aws_subnet.example2.id
    allocation_id = aws_eip.example2.id
  }
}
```

### Specifying private IP addresses for an internal-facing load balancer

```terraform
resource "aws_lb" "example" {
  name               = "example"
  load_balancer_type = "network"

  subnet_mapping {
    subnet_id            = aws_subnet.example1.id
    private_ipv4_address = "10.0.1.15"
  }

  subnet_mapping {
    subnet_id            = aws_subnet.example2.id
    private_ipv4_address = "10.0.2.15"
  }
}
```

## Argument Reference

~> **NOTE:** Please note that internal LBs can only use `ipv4` as the `ip_address_type`. You can only change to `dualstack` `ip_address_type` if the selected subnets are IPv6 enabled.

~> **NOTE:** Please note that one of either `subnets` or `subnet_mapping` is required.

The following arguments are supported:

* `access_logs` - (Optional) An Access Logs block. Access Logs documented below.
* `customer_owned_ipv4_pool` - (Optional) The ID of the customer owned ipv4 pool to use for this load balancer.
* `desync_mitigation_mode` - (Optional) Determines how the load balancer handles requests that might pose a security risk to an application due to HTTP desync. Valid values are `monitor`, `defensive` (default), `strictest`.
* `drop_invalid_header_fields` - (Optional) Indicates whether HTTP headers with header fields that are not valid are removed by the load balancer (true) or routed to targets (false). The default is false. Elastic Load Balancing requires that message header names contain only alphanumeric characters and hyphens. Only valid for Load Balancers of type `application`.
* `enable_cross_zone_load_balancing` - (Optional) If true, cross-zone load balancing of the load balancer will be enabled. This is a `network` load balancer feature. Defaults to `false`.
* `enable_deletion_protection` - (Optional) If true, deletion of the load balancer will be disabled via the AWS API. This will prevent this provider from deleting the load balancer. Defaults to `false`.
* `enable_http2` - (Optional) Indicates whether HTTP/2 is enabled in `application` load balancers. Defaults to `true`.
* `enable_waf_fail_open` - (Optional) Indicates whether to allow a WAF-enabled load balancer to route requests to targets if it is unable to forward the request to AWS WAF. Defaults to `false`.
* `idle_timeout` - (Optional) The time in seconds that the connection is allowed to be idle. Only valid for Load Balancers of type `application`. Default: 60.
* `internal` - (Optional) If true, the LB will be internal.
* `ip_address_type` - (Optional) The type of IP addresses used by the subnets for your load balancer. The possible values are `ipv4` and `dualstack`.
* `load_balancer_type` - (Optional) The type of load balancer to create. Possible values are `application`, `gateway`, or `network`. The default value is `application`.
* `name` - (Optional) The name of the LB. This name must be unique within your AWS account, can have a maximum of 32 characters,
must contain only alphanumeric characters or hyphens, and must not begin or end with a hyphen. If not specified,
this provider will autogenerate a name beginning with `tf-lb`.
* `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `name`.* `security_groups` - (Optional) A list of security group IDs to assign to the LB. Only valid for Load Balancers of type `application`.
* `preserve_host_header` - (Optional) Indicates whether the Application Load Balancer should preserve the Host header in the HTTP request and send it to the target without any change. Defaults to `false`.
* `subnet_mapping` - (Optional) A subnet mapping block as documented below.
* `subnets` - (Optional) A list of subnet IDs to attach to the LB. Subnets
cannot be updated for Load Balancers of type `network`. Changing this value
for load balancers of type `network` will force a recreation of the resource.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### access_logs

* `bucket` - (Required) The S3 bucket name to store the logs in.
* `enabled` - (Optional) Boolean to enable / disable `access_logs`. Defaults to `false`, even when `bucket` is specified.
* `prefix` - (Optional) The S3 bucket prefix. Logs are stored in the root if not configured.

### subnet_mapping

* `subnet_id` - (Required) ID of the subnet of which to attach to the load balancer. You can specify only one subnet per Availability Zone.
* `allocation_id` - (Optional) The allocation ID of the Elastic IP address for an internet-facing load balancer.
* `ipv6_address` - (Optional) The IPv6 address. You associate IPv6 CIDR blocks with your VPC and choose the subnets where you launch both internet-facing and internal Application Load Balancers or Network Load Balancers.
* `private_ipv4_address` - (Optional) The private IPv4 address for an internal load balancer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the load balancer (matches `id`).
* `arn_suffix` - The ARN suffix for use with CloudWatch Metrics.
* `dns_name` - The DNS name of the load balancer.
* `id` - The ARN of the load balancer (matches `arn`).
* `subnet_mapping.*.outpost_id` - ID of the Outpost containing the load balancer.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
* `zone_id` - The canonical hosted zone ID of the load balancer (to be used in a Route 53 Alias record).

## Timeouts

Configuration options:

- `create` - (Default `10m`)
- `update` - (Default `10m`)
- `delete` - (Default `10m`)

## Import

LBs can be imported using their ARN, e.g.,

```
$ terraform import aws_lb.bar arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/50dc6c495c0c9188
```
