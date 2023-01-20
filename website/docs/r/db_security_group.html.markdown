---
subcategory: "RDS (Relational Database)"
layout: "aws"
page_title: "AWS: aws_db_security_group"
description: |-
  Provides an RDS security group resource.
---

# Resource: aws_db_security_group

Provides an RDS security group resource. This is only for DB instances in the
EC2-Classic Platform. For instances inside a VPC, use the
`aws_db_instance.vpc_security_group_ids`
attribute instead.

!> **WARNING:** With the retirement of EC2-Classic the `aws_db_security_group` resource has been deprecated and will be removed in a future version.

## Example Usage

```terraform
resource "aws_db_security_group" "default" {
  name = "rds_sg"

  ingress {
    cidr = "10.0.0.0/24"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DB security group.
* `description` - (Optional) The description of the DB security group. Defaults to "Managed by Pulumi".
* `ingress` - (Required) A list of ingress rules.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

Ingress blocks support the following:

* `cidr` - The CIDR block to accept
* `security_group_name` - The name of the security group to authorize
* `security_group_id` - The ID of the security group to authorize
* `security_group_owner_id` - The owner Id of the security group provided
  by `security_group_name`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The db security group ID.
* `arn` - The arn of the DB security group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

DB Security groups can be imported using the `name`, e.g.,

```
$ terraform import aws_db_security_group.default aws_rds_sg-1
```
