---
subcategory: "RDS (Relational Database)"
layout: "aws"
page_title: "AWS: aws_db_parameter_group"
description: |-
  Provides an RDS DB parameter group resource.
---

# Resource: aws_db_parameter_group

Provides an RDS DB parameter group resource. Documentation of the available parameters for various RDS engines can be found at:

* [Aurora MySQL Parameters](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AuroraMySQL.Reference.html)
* [Aurora PostgreSQL Parameters](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AuroraPostgreSQL.Reference.html)
* [MariaDB Parameters](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MariaDB.Parameters.html)
* [Oracle Parameters](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ModifyInstance.Oracle.html#USER_ModifyInstance.Oracle.sqlnet)
* [PostgreSQL Parameters](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.PostgreSQL.CommonDBATasks.html#Appendix.PostgreSQL.CommonDBATasks.Parameters)


~> **NOTE:** After applying your changes, you may encounter a perpetual diff in your preview
output for a `parameter` whose `value` remains unchanged but whose `apply_method` is changing
(e.g., from `immediate` to `pending-reboot`, or `pending-reboot` to `immediate`). If only the
apply method of a parameter is changing, the AWS API will not register this change. To change
the `apply_method` of a parameter, its value must also change.

## Example Usage

### Basic Usage

```terraform
resource "aws_db_parameter_group" "default" {
  name   = "rds-pg"
  family = "mysql5.6"

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name  = "character_set_client"
    value = "utf8"
  }
}
```

### `create_before_destroy` Lifecycle Configuration

The `create_before_destroy`
lifecycle configuration is necessary for modifications that force re-creation of an existing,
in-use parameter group. This includes common situations like changing the group `name` or
bumping the `family` version during a major version upgrade. This configuration will prevent destruction
of the deposed parameter group while still in use by the database during upgrade.

```terraform
resource "aws_db_parameter_group" "example" {
  name   = "my-pg"
  family = "postgres13"

  parameter {
    name  = "log_connections"
    value = "1"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_db_instance" "example" {
  # other attributes
  parameter_group_name = aws_db_parameter_group.example.name
  apply_immediately    = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, Forces new resource) The name of the DB parameter group. If omitted, this provider will assign a random, unique name.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `family` - (Required, Forces new resource) The family of the DB parameter group.
* `description` - (Optional, Forces new resource) The description of the DB parameter group. Defaults to "Managed by Pulumi".
* `parameter` - (Optional) A list of DB parameters to apply. Note that parameters may differ from a family to an other. Full list of all parameters can be discovered via [`aws rds describe-db-parameters`](https://docs.aws.amazon.com/cli/latest/reference/rds/describe-db-parameters.html) after initial creation of the group.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

Parameter blocks support the following:

* `name` - (Required) The name of the DB parameter.
* `value` - (Required) The value of the DB parameter.
* `apply_method` - (Optional) "immediate" (default), or "pending-reboot". Some
    engines can't apply some parameters without a reboot, and you will need to
    specify "pending-reboot" here.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The db parameter group name.
* `arn` - The ARN of the db parameter group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

DB Parameter groups can be imported using the `name`, e.g.,

```
$ terraform import aws_db_parameter_group.rds_pg rds-pg
```
