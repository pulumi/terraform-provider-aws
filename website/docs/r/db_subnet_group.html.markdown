---
subcategory: "RDS (Relational Database)"
layout: "aws"
page_title: "AWS: aws_db_subnet_group"
description: |-
  Provides an RDS DB subnet group resource.
---

# Resource: aws_db_subnet_group

Provides an RDS DB subnet group resource.

## Example Usage

```terraform
resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = [aws_subnet.frontend.id, aws_subnet.backend.id]

  tags = {
    Name = "My DB subnet group"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, Forces new resource) The name of the DB subnet group. If omitted, this provider will assign a random, unique name.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `description` - (Optional) The description of the DB subnet group. Defaults to "Managed by Pulumi".
* `subnet_ids` - (Required) A list of VPC subnet IDs.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The db subnet group name.
* `arn` - The ARN of the db subnet group.
* `supported_network_types` - The network type of the db subnet group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

DB Subnet groups can be imported using the `name`, e.g.,

```
$ terraform import aws_db_subnet_group.default production-subnet-group
```
