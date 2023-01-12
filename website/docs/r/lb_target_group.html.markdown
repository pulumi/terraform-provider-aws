---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "AWS: aws_lb_target_group"
description: |-
  Provides a Target Group resource for use with Load Balancers.
---

# Resource: aws_lb_target_group

Provides a Target Group resource for use with Load Balancer resources.

~> **Note:** `aws_alb_target_group` is known as `aws_lb_target_group`. The functionality is identical.

## Example Usage

### Instance Target Group

```terraform
resource "aws_lb_target_group" "test" {
  name     = "tf-example-lb-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
```

### IP Target Group

```terraform
resource "aws_lb_target_group" "ip-example" {
  name        = "tf-example-lb-tg"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
```

### Lambda Target Group

```terraform
resource "aws_lb_target_group" "lambda-example" {
  name        = "tf-example-lb-tg"
  target_type = "lambda"
}
```

### ALB Target Group

```terraform
resource "aws_lb_target_group" "alb-example" {
  name        = "tf-example-lb-alb-tg"
  target_type = "alb"
  port        = 80
  protocol    = "TCP"
  vpc_id      = aws_vpc.main.id
}
```

## Argument Reference

The following arguments are supported:

* `connection_termination` - (Optional) Whether to terminate connections at the end of the deregistration timeout on Network Load Balancers. See [doc](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-target-groups.html#deregistration-delay) for more information. Default is `false`.
* `deregistration_delay` - (Optional) Amount time for Elastic Load Balancing to wait before changing the state of a deregistering target from draining to unused. The range is 0-3600 seconds. The default value is 300 seconds.
* `health_check` - (Optional, Maximum of 1) Health Check configuration block. Detailed below.
* `lambda_multi_value_headers_enabled` - (Optional) Whether the request and response headers exchanged between the load balancer and the Lambda function include arrays of values or strings. Only applies when `target_type` is `lambda`. Default is `false`.
* `load_balancing_algorithm_type` - (Optional) Determines how the load balancer selects targets when routing requests. Only applicable for Application Load Balancer Target Groups. The value is `round_robin` or `least_outstanding_requests`. The default is `round_robin`.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`. Cannot be longer than 6 characters.
* `name` - (Optional, Forces new resource) Name of the target group. If omitted, this provider will assign a random, unique name.
* `port` - (May be required, Forces new resource) Port on which targets receive traffic, unless overridden when registering a specific target. Required when `target_type` is `instance`, `ip` or `alb`. Does not apply when `target_type` is `lambda`.
* `preserve_client_ip` - (Optional) Whether client IP preservation is enabled. See [doc](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-target-groups.html#client-ip-preservation) for more information.
* `protocol_version` - (Optional, Forces new resource) Only applicable when `protocol` is `HTTP` or `HTTPS`. The protocol version. Specify `GRPC` to send requests to targets using gRPC. Specify `HTTP2` to send requests to targets using HTTP/2. The default is `HTTP1`, which sends requests to targets using HTTP/1.1
* `protocol` - (May be required, Forces new resource) Protocol to use for routing traffic to the targets. Should be one of `GENEVE`, `HTTP`, `HTTPS`, `TCP`, `TCP_UDP`, `TLS`, or `UDP`. Required when `target_type` is `instance`, `ip` or `alb`. Does not apply when `target_type` is `lambda`.
* `proxy_protocol_v2` - (Optional) Whether to enable support for proxy protocol v2 on Network Load Balancers. See [doc](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-target-groups.html#proxy-protocol) for more information. Default is `false`.
* `slow_start` - (Optional) Amount time for targets to warm up before the load balancer sends them a full share of requests. The range is 30-900 seconds or 0 to disable. The default value is 0 seconds.
* `stickiness` - (Optional, Maximum of 1) Stickiness configuration block. Detailed below.
* `tags` - (Optional) Map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `target_failover` - (Optional) Target failover block. Only applicable for Gateway Load Balancer target groups. See [target_failover](#target_failover) for more information.
* `target_type` - (May be required, Forces new resource) Type of target that you must specify when registering targets with this target group. See [doc](https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_CreateTargetGroup.html) for supported values. The default is `instance`.

  Note that you can't specify targets for a target group using both instance IDs and IP addresses.

  If the target type is `ip`, specify IP addresses from the subnets of the virtual private cloud (VPC) for the target group, the RFC 1918 range (10.0.0.0/8, 172.16.0.0/12, and 192.168.0.0/16), and the RFC 6598 range (100.64.0.0/10). You can't specify publicly routable IP addresses.

  Network Load Balancers do not support the `lambda` target type.

  Application Load Balancers do not support the `alb` target type.
* `ip_address_type` (Optional, forces new resource) The type of IP addresses used by the target group, only supported when target type is set to `ip`. Possible values are `ipv4` or `ipv6`.
* `vpc_id` - (Optional, Forces new resource) Identifier of the VPC in which to create the target group. Required when `target_type` is `instance`, `ip` or `alb`. Does not apply when `target_type` is `lambda`.

### health_check

~> **Note:** The Health Check parameters you can set vary by the `protocol` of the Target Group. Many parameters cannot be set to custom values for `network` load balancers at this time. See http://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_CreateTargetGroup.html for a complete reference. Keep in mind, that health checks produce actual requests to the backend. The underlying function is invoked when `target_type` is set to `lambda`.

* `enabled` - (Optional) Whether health checks are enabled. Defaults to `true`.
* `healthy_threshold` - (Optional)  Number of consecutive health check successes required before considering a target healthy. The range is 2-10. Defaults to 3.
* `interval` - (Optional) Approximate amount of time, in seconds, between health checks of an individual target. The range is 5-300. For `lambda` target groups, it needs to be greater than the timeout of the underlying `lambda`. Defaults to 30.
* `matcher` (May be required) Response codes to use when checking for a healthy responses from a target. You can specify multiple values (for example, "200,202" for HTTP(s) or "0,12" for GRPC) or a range of values (for example, "200-299" or "0-99"). Required for HTTP/HTTPS/GRPC ALB. Only applies to Application Load Balancers (i.e., HTTP/HTTPS/GRPC) not Network Load Balancers (i.e., TCP).
* `path` - (May be required) Destination for the health check request. Required for HTTP/HTTPS ALB and HTTP NLB. Only applies to HTTP/HTTPS.
* `port` - (Optional) The port the load balancer uses when performing health checks on targets. Default is traffic-port.
* `protocol` - (Optional) Protocol the load balancer uses when performing health checks on targets. Must be either `TCP`, `HTTP`, or `HTTPS`. The TCP protocol is not supported for health checks if the protocol of the target group is HTTP or HTTPS. Defaults to HTTP.
* `timeout` - (optional) Amount of time, in seconds, during which no response from a target means a failed health check. The range is 2–120 seconds. For target groups with a protocol of HTTP, the default is 6 seconds. For target groups with a protocol of TCP, TLS or HTTPS, the default is 10 seconds. For target groups with a protocol of GENEVE, the default is 5 seconds. If the target type is lambda, the default is 30 seconds.
* `unhealthy_threshold` - (Optional) Number of consecutive health check failures required before considering a target unhealthy. The range is 2-10. Defaults to 3.

### stickiness

~> **NOTE:** Currently, an NLB (i.e., protocol of `HTTP` or `HTTPS`) can have an invalid `stickiness` block with `type` set to `lb_cookie` as long as `enabled` is set to `false`. However, please update your configurations to avoid errors in a future version of the provider: either remove the invalid `stickiness` block or set the `type` to `source_ip`.

* `cookie_duration` - (Optional) Only used when the type is `lb_cookie`. The time period, in seconds, during which requests from a client should be routed to the same target. After this time period expires, the load balancer-generated cookie is considered stale. The range is 1 second to 1 week (604800 seconds). The default value is 1 day (86400 seconds).
* `cookie_name` - (Optional) Name of the application based cookie. AWSALB, AWSALBAPP, and AWSALBTG prefixes are reserved and cannot be used. Only needed when type is `app_cookie`.
* `enabled` - (Optional) Boolean to enable / disable `stickiness`. Default is `true`.
* `type` - (Required) The type of sticky sessions. The only current possible values are `lb_cookie`, `app_cookie` for ALBs, `source_ip` for NLBs, and `source_ip_dest_ip`, `source_ip_dest_ip_proto` for GWLBs.

### target_failover
~> **NOTE:** This block is only applicable for a Gateway Load Balancer (GWLB). The two attributes `on_deregistration` and `on_unhealthy` cannot be set independently. The value you set for both attributes must be the same.

* `on_deregistration` - (Optional) Indicates how the GWLB handles existing flows when a target is deregistered. Possible values are `rebalance` and `no_rebalance`. Must match the attribute value set for `on_unhealthy`. Default: `no_rebalance`.
* `on_unhealthy` - Indicates how the GWLB handles existing flows when a target is unhealthy. Possible values are `rebalance` and `no_rebalance`. Must match the attribute value set for `on_deregistration`. Default: `no_rebalance`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn_suffix` - ARN suffix for use with CloudWatch Metrics.
* `arn` - ARN of the Target Group (matches `id`).
* `id` - ARN of the Target Group (matches `arn`).
* `name` - Name of the Target Group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Target Groups can be imported using their ARN, e.g.,

```
$ terraform import aws_lb_target_group.app_front_end arn:aws:elasticloadbalancing:us-west-2:187416307283:targetgroup/app-front-end/20cfe21448b66314
```
