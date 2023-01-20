---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "AWS: aws_placement_group"
description: |-
  Provides an EC2 placement group.
---

# Resource: aws_placement_group

Provides an EC2 placement group. Read more about placement groups
in [AWS Docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html).

## Example Usage

```terraform
resource "aws_placement_group" "web" {
  name     = "hunky-dory-pg"
  strategy = "cluster"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the placement group.
* `partition_count` - (Optional) The number of partitions to create in the
  placement group.  Can only be specified when the `strategy` is set to
  `"partition"`.  Valid values are 1 - 7 (default is `2`).
* `spread_level` - (Optional) Determines how placement groups spread instances. Can only be used
   when the `strategy` is set to `"spread"`. Can be `"host"` or `"rack"`. `"host"` can only be used for Outpost placement groups.
* `strategy` - (Required) The placement strategy. Can be `"cluster"`, `"partition"` or `"spread"`.
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the placement group.
* `id` - The name of the placement group.
* `placement_group_id` - The ID of the placement group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Placement groups can be imported using the `name`, e.g.,

```
$ terraform import aws_placement_group.prod_pg production-placement-group
```
