---
subcategory: "OpsWorks"
layout: "aws"
page_title: "AWS: aws_opsworks_rds_db_instance"
description: |-
  Provides an OpsWorks RDS DB Instance resource.
---

# Resource: aws_opsworks_rds_db_instance

Provides an OpsWorks RDS DB Instance resource.

## Example Usage

```terraform
resource "aws_opsworks_rds_db_instance" "my_instance" {
  stack_id            = aws_opsworks_stack.my_stack.id
  rds_db_instance_arn = aws_db_instance.my_instance.arn
  db_user             = "someUser"
  db_password         = "somePass"
}
```

## Argument Reference

The following arguments are supported:

* `stack_id` - (Required) The stack to register a db instance for. Changing this will force a new resource.
* `rds_db_instance_arn` - (Required) The db instance to register for this stack. Changing this will force a new resource.
* `db_user` - (Required) A db username
* `db_password` - (Required) A db password

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The computed id. Please note that this is only used internally to identify the stack <-> instance relation. This value is not used in aws.
