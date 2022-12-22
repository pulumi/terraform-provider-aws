---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "AWS: aws_instances"
description: |-
  Get information on an Amazon EC2 instances.
---

# Data Source: aws_instances

Use this data source to get IDs or IPs of Amazon EC2 instances to be referenced elsewhere,
e.g., to allow easier migration from another management solution
or to make it easier for an operator to connect through bastion host(s).


~> **Note:** It's strongly discouraged to use this data source for querying ephemeral
instances (e.g., managed via autoscaling group), as the output may change at any time
and you'd need to re-run `apply` every time an instance comes up or dies.

## Example Usage

```terraform
data "aws_instances" "test" {
  instance_tags = {
    Role = "HardWorker"
  }

  filter {
    name   = "instance.group-id"
    values = ["sg-12345678"]
  }

  instance_state_names = ["running", "stopped"]
}

resource "aws_eip" "test" {
  count    = length(data.aws_instances.test.ids)
  instance = data.aws_instances.test.ids[count.index]
}
```

## Argument Reference

* `instance_tags` - (Optional) Map of tags, each pair of which must
exactly match a pair on desired instances.

* `instance_state_names` - (Optional) List of instance states that should be applicable to the desired instances. The permitted values are: `pending, running, shutting-down, stopped, stopping, terminated`. The default value is `running`.

* `filter` - (Optional) One or more name/value pairs to use as filters. There are
several valid keys, for a full reference, check out
[describe-instances in the AWS CLI reference][1].

## Attributes Reference

* `id` - AWS Region.
* `ids` - IDs of instances found through the filter
* `private_ips` - Private IP addresses of instances found through the filter
* `public_ips` - Public IP addresses of instances found through the filter

## Timeouts

Configuration options:

- `read` - (Default `20m`)

[1]: http://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html
