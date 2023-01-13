---
subcategory: "RDS (Relational Database)"
layout: "aws"
page_title: "AWS: aws_db_proxy_endpoint"
description: |-
  Provides an RDS DB proxy endpoint resource.
---

# Resource: aws_db_proxy_endpoint

Provides an RDS DB proxy endpoint resource. For additional information, see the [RDS User Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-endpoints.html).

## Example Usage

```terraform
resource "aws_db_proxy_endpoint" "example" {
  db_proxy_name          = aws_db_proxy.test.name
  db_proxy_endpoint_name = "example"
  vpc_subnet_ids         = aws_subnet.test[*].id
  target_role            = "READ_ONLY"
}
```

## Argument Reference

The following arguments are supported:

* `db_proxy_endpoint_name` - (Required) The identifier for the proxy endpoint. An identifier must begin with a letter and must contain only ASCII letters, digits, and hyphens; it can't end with a hyphen or contain two consecutive hyphens.
* `db_proxy_name` - (Required) The name of the DB proxy associated with the DB proxy endpoint that you create.
* `vpc_subnet_ids` - (Required) One or more VPC subnet IDs to associate with the new proxy.
* `vpc_security_group_ids` - (Optional) One or more VPC security group IDs to associate with the new proxy.
* `target_role` - (Optional) Indicates whether the DB proxy endpoint can be used for read/write or read-only operations. The default is `READ_WRITE`. Valid values are `READ_WRITE` and `READ_ONLY`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the proxy and proxy endpoint separated by `/`, `DB-PROXY-NAME/DB-PROXY-ENDPOINT-NAME`.
* `arn` - The Amazon Resource Name (ARN) for the proxy endpoint.
* `endpoint` - The endpoint that you can use to connect to the proxy. You include the endpoint value in the connection string for a database client application.
* `is_default` - Indicates whether this endpoint is the default endpoint for the associated DB proxy.
* `vpc_id` - The VPC ID of the DB proxy endpoint.

## Timeouts

Configuration options:

- `create` - (Default `30m`)
- `update` - (Default `30m`)
- `delete` - (Default `60m`)

## Import

DB proxy endpoints can be imported using the `DB-PROXY-NAME/DB-PROXY-ENDPOINT-NAME`, e.g.,

```
$ terraform import aws_db_proxy_endpoint.example example/example
```
