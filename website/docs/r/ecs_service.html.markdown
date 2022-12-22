---
subcategory: "ECS (Elastic Container)"
layout: "aws"
page_title: "AWS: aws_ecs_service"
description: |-
  Provides an ECS service.
---

# Resource: aws_ecs_service

-> **Note:** To prevent a race condition during service deletion, make sure to set `depends_on` to the related `aws_iam_role_policy`; otherwise, the policy may be destroyed too soon and the ECS service will then get stuck in the `DRAINING` state.

Provides an ECS service - effectively a task that is expected to run until an error occurs or a user terminates it (typically a webserver or a database).

See [ECS Services section in AWS developer guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html).

## Example Usage

```terraform
resource "aws_ecs_service" "mongo" {
  name            = "mongodb"
  cluster         = aws_ecs_cluster.foo.id
  task_definition = aws_ecs_task_definition.mongo.arn
  desired_count   = 3
  iam_role        = aws_iam_role.foo.arn
  depends_on      = [aws_iam_role_policy.foo]

  ordered_placement_strategy {
    type  = "binpack"
    field = "cpu"
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.foo.arn
    container_name   = "mongo"
    container_port   = 8080
  }

  placement_constraints {
    type       = "memberOf"
    expression = "attribute:ecs.availability-zone in [us-west-2a, us-west-2b]"
  }
}
```

### Ignoring Changes to Desired Count

You can use [`ignoreChanges`](https://www.pulumi.com/docs/intro/concepts/programming-model/#ignorechanges) to create an ECS service with an initial count of running instances, then ignore any changes to that count caused externally (e.g. Application Autoscaling).

```terraform
resource "aws_ecs_service" "example" {
  # ... other configurations ...

  # Example: Create service with 2 instances to start
  desired_count = 2

  # Optional: Allow external changes without this provider plan difference
  lifecycle {
    ignore_changes = [desired_count]
  }
}
```

### Daemon Scheduling Strategy

```terraform
resource "aws_ecs_service" "bar" {
  name                = "bar"
  cluster             = aws_ecs_cluster.foo.id
  task_definition     = aws_ecs_task_definition.bar.arn
  scheduling_strategy = "DAEMON"
}
```

### External Deployment Controller

```terraform
resource "aws_ecs_service" "example" {
  name    = "example"
  cluster = aws_ecs_cluster.example.id

  deployment_controller {
    type = "EXTERNAL"
  }
}
```

### Redeploy Service On Every Apply

The key used with `triggers` is arbitrary.

```terraform
resource "aws_ecs_service" "example" {
  # ... other configurations ...

  force_new_deployment = true

  triggers = {
    redeployment = timestamp()
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Name of the service (up to 255 letters, numbers, hyphens, and underscores)

The following arguments are optional:

* `capacity_provider_strategy` - (Optional) Capacity provider strategies to use for the service. Can be one or more. These can be updated without destroying and recreating the service only if `force_new_deployment = true` and not changing from 0 `capacity_provider_strategy` blocks to greater than 0, or vice versa. See below.
* `cluster` - (Optional) ARN of an ECS cluster.
* `deployment_circuit_breaker` - (Optional) Configuration block for deployment circuit breaker. See below.
* `deployment_controller` - (Optional) Configuration block for deployment controller configuration. See below.
* `deployment_maximum_percent` - (Optional) Upper limit (as a percentage of the service's desiredCount) of the number of running tasks that can be running in a service during a deployment. Not valid when using the `DAEMON` scheduling strategy.
* `deployment_minimum_healthy_percent` - (Optional) Lower limit (as a percentage of the service's desiredCount) of the number of running tasks that must remain running and healthy in a service during a deployment.
* `desired_count` - (Optional) Number of instances of the task definition to place and keep running. Defaults to 0. Do not specify if using the `DAEMON` scheduling strategy.
* `enable_ecs_managed_tags` - (Optional) Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
* `enable_execute_command` - (Optional) Specifies whether to enable Amazon ECS Exec for the tasks within the service.
* `force_new_deployment` - (Optional) Enable to force a new task deployment of the service. This can be used to update tasks to use a newer Docker image with same image/tag combination (e.g., `myimage:latest`), roll Fargate tasks onto a newer platform version, or immediately deploy `ordered_placement_strategy` and `placement_constraints` updates.
* `health_check_grace_period_seconds` - (Optional) Seconds to ignore failing load balancer health checks on newly instantiated tasks to prevent premature shutdown, up to 2147483647. Only valid for services configured to use load balancers.
* `iam_role` - (Optional) ARN of the IAM role that allows Amazon ECS to make calls to your load balancer on your behalf. This parameter is required if you are using a load balancer with your service, but only if your task definition does not use the `awsvpc` network mode. If using `awsvpc` network mode, do not specify this role. If your account has already created the Amazon ECS service-linked role, that role is used by default for your service unless you specify a role here.
* `launch_type` - (Optional) Launch type on which to run your service. The valid values are `EC2`, `FARGATE`, and `EXTERNAL`. Defaults to `EC2`.
* `load_balancer` - (Optional) Configuration block for load balancers. See below.
* `network_configuration` - (Optional) Network configuration for the service. This parameter is required for task definitions that use the `awsvpc` network mode to receive their own Elastic Network Interface, and it is not supported for other network modes. See below.
* `ordered_placement_strategy` - (Optional) Service level strategy rules that are taken into consideration during task placement. List from top to bottom in order of precedence. Updates to this configuration will take effect next task deployment unless `force_new_deployment` is enabled. The maximum number of `ordered_placement_strategy` blocks is `5`. See below.
* `placement_constraints` - (Optional) Rules that are taken into consideration during task placement. Updates to this configuration will take effect next task deployment unless `force_new_deployment` is enabled. Maximum number of `placement_constraints` is `10`. See below.
* `platform_version` - (Optional) Platform version on which to run your service. Only applicable for `launch_type` set to `FARGATE`. Defaults to `LATEST`. More information about Fargate platform versions can be found in the [AWS ECS User Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html).
* `propagate_tags` - (Optional) Specifies whether to propagate the tags from the task definition or the service to the tasks. The valid values are `SERVICE` and `TASK_DEFINITION`.
* `scheduling_strategy` - (Optional) Scheduling strategy to use for the service. The valid values are `REPLICA` and `DAEMON`. Defaults to `REPLICA`. Note that [*Tasks using the Fargate launch type or the `CODE_DEPLOY` or `EXTERNAL` deployment controller types don't support the `DAEMON` scheduling strategy*](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html).
* `service_connect_configuration` - (Optional) The ECS Service Connect configuration for this service to discover and connect to services, and be discovered by, and connected from, other services within a namespace. See below.
* `service_registries` - (Optional) Service discovery registries for the service. The maximum number of `service_registries` blocks is `1`. See below.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `task_definition` - (Optional) Family and revision (`family:revision`) or full ARN of the task definition that you want to run in your service. Required unless using the `EXTERNAL` deployment controller. If a revision is not specified, the latest `ACTIVE` revision is used.
* `triggers` - (Optional) Map of arbitrary keys and values that, when changed, will trigger an in-place update (redeployment). Useful with `timestamp()`. See example above.
* `wait_for_steady_state` - (Optional) If `true`, this provider will wait for the service to reach a steady state (like [`aws ecs wait services-stable`](https://docs.aws.amazon.com/cli/latest/reference/ecs/wait/services-stable.html)) before continuing. Default `false`.

### capacity_provider_strategy

The `capacity_provider_strategy` configuration block supports the following:

* `base` - (Optional) Number of tasks, at a minimum, to run on the specified capacity provider. Only one capacity provider in a capacity provider strategy can have a base defined.
* `capacity_provider` - (Required) Short name of the capacity provider.
* `weight` - (Required) Relative percentage of the total number of launched tasks that should use the specified capacity provider.

### deployment_circuit_breaker

The `deployment_circuit_breaker` configuration block supports the following:

* `enable` - (Required) Whether to enable the deployment circuit breaker logic for the service.
* `rollback` - (Required) Whether to enable Amazon ECS to roll back the service if a service deployment fails. If rollback is enabled, when a service deployment fails, the service is rolled back to the last deployment that completed successfully.

### deployment_controller

The `deployment_controller` configuration block supports the following:

* `type` - (Optional) Type of deployment controller. Valid values: `CODE_DEPLOY`, `ECS`, `EXTERNAL`. Default: `ECS`.

### load_balancer

`load_balancer` supports the following:

* `elb_name` - (Required for ELB Classic) Name of the ELB (Classic) to associate with the service.
* `target_group_arn` - (Required for ALB/NLB) ARN of the Load Balancer target group to associate with the service.
* `container_name` - (Required) Name of the container to associate with the load balancer (as it appears in a container definition).
* `container_port` - (Required) Port on the container to associate with the load balancer.

-> **Version note:** Multiple `load_balancer` configuration block support was added in version 2.22.0 of the provider. This allows configuration of [ECS service support for multiple target groups](https://aws.amazon.com/about-aws/whats-new/2019/07/amazon-ecs-services-now-support-multiple-load-balancer-target-groups/).

### network_configuration

`network_configuration` support the following:

* `subnets` - (Required) Subnets associated with the task or service.
* `security_groups` - (Optional) Security groups associated with the task or service. If you do not specify a security group, the default security group for the VPC is used.
* `assign_public_ip` - (Optional) Assign a public IP address to the ENI (Fargate launch type only). Valid values are `true` or `false`. Default `false`.

For more information, see [Task Networking](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-networking.html)

### ordered_placement_strategy

`ordered_placement_strategy` supports the following:

* `type` - (Required) Type of placement strategy. Must be one of: `binpack`, `random`, or `spread`
* `field` - (Optional) For the `spread` placement strategy, valid values are `instanceId` (or `host`,
 which has the same effect), or any platform or custom attribute that is applied to a container instance.
 For the `binpack` type, valid values are `memory` and `cpu`. For the `random` type, this attribute is not
 needed. For more information, see [Placement Strategy](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_PlacementStrategy.html).

-> **Note:** for `spread`, `host` and `instanceId` will be normalized, by AWS, to be `instanceId`. This means the statefile will show `instanceId` but your config will differ if you use `host`.

### placement_constraints

`placement_constraints` support the following:

* `type` - (Required) Type of constraint. The only valid values at this time are `memberOf` and `distinctInstance`.
* `expression` -  (Optional) Cluster Query Language expression to apply to the constraint. Does not need to be specified for the `distinctInstance` type. For more information, see [Cluster Query Language in the Amazon EC2 Container Service Developer Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cluster-query-language.html).

### service_registries

`service_registries` support the following:

* `registry_arn` - (Required) ARN of the Service Registry. The currently supported service registry is Amazon Route 53 Auto Naming Service(`aws_service_discovery_service`). For more information, see [Service](https://docs.aws.amazon.com/Route53/latest/APIReference/API_autonaming_Service.html)
* `port` - (Optional) Port value used if your Service Discovery service specified an SRV record.
* `container_port` - (Optional) Port value, already specified in the task definition, to be used for your service discovery service.
* `container_name` - (Optional) Container name value, already specified in the task definition, to be used for your service discovery service.

### service_connect_configuration

`service_connect_configuration` supports the following:

* `enabled` - (Required) Specifies whether to use Service Connect with this service.
* `log_configuration` - (Optional) The log configuration for the container. See below.
* `namespace` - (Optional) The namespace name or ARN of the `aws_service_discovery_http_namespace` for use with Service Connect.
* `service` - (Optional) The list of Service Connect service objects. See below.

### log_configuration

`log_configuration` supports the following:

* `log_driver` - (Optional) The log driver to use for the container.
* `options` - (Optional) The configuration options to send to the log driver.
* `secret_option` - (Optional) The secrets to pass to the log configuration. See below.

### secret_option

`secret_option` supports the following:

* `name` - (Optional) The name of the secret.
* `value_from` - (Optional) The secret to expose to the container. The supported values are either the full ARN of the AWS Secrets Manager secret or the full ARN of the parameter in the SSM Parameter Store.

### service

`service` supports the following:

* `client_alias` - (Optional) The list of client aliases for this Service Connect service. You use these to assign names that can be used by client applications. The maximum number of client aliases that you can have in this list is 1. See below.
* `discovery_name` - (Optional) The name of the new AWS Cloud Map service that Amazon ECS creates for this Amazon ECS service.
* `ingress_port_override` - (Optional) The port number for the Service Connect proxy to listen on.
* `port_name` - (Required) The name of one of the `portMappings` from all the containers in the task definition of this Amazon ECS service.

### client_alias

`client_alias` supports the following:

* `dns_name` - (Optional) The name that you use in the applications of client tasks to connect to this service.
* `port` - (Required) The listening port number for the Service Connect proxy. This port is available inside of all of the tasks within the same namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster` - Amazon Resource Name (ARN) of cluster which the service runs on.
* `desired_count` - Number of instances of the task definition.
* `iam_role` - ARN of IAM role used for ELB.
* `id` - ARN that identifies the service.
* `name` - Name of the service.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Timeouts

Configuration options:

- `create` - (Default `20m`)
- `update` - (Default `20m`)
- `delete` - (Default `20m`)

## Import

ECS services can be imported using the `name` together with ecs cluster `name`, e.g.,

```
$ terraform import aws_ecs_service.imported cluster-name/service-name
```
