---
subcategory: "CloudWatch Logs"
layout: "aws"
page_title: "AWS: aws_cloudwatch_log_group"
description: |-
  Provides a CloudWatch Log Group resource.
---

# Resource: aws_cloudwatch_log_group

Provides a CloudWatch Log Group resource.

## Example Usage

```terraform
resource "aws_cloudwatch_log_group" "yada" {
  name = "Yada"

  tags = {
    Environment = "production"
    Application = "serviceA"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, Forces new resource) The name of the log group. If omitted, this provider will assign a random, unique name.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `skip_destroy` - (Optional) Set to true if you do not wish the log group (and any logs it may contain) to be deleted at destroy time, and instead just remove the log group from the state.
* `retention_in_days` - (Optional) Specifies the number of days
  you want to retain log events in the specified log group.  Possible values are: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 2192, 2557, 2922, 3288, 3653, and 0.
  If you select 0, the events in the log group are always retained and never expire.
* `kms_key_id` - (Optional) The ARN of the KMS Key to use when encrypting log data. Please note, after the AWS KMS CMK is disassociated from the log group,
AWS CloudWatch Logs stops encrypting newly ingested data for the log group. All previously ingested data remains encrypted, and AWS CloudWatch Logs requires
permissions for the CMK whenever the encrypted data is requested.
* `tags` - (Optional) A map of tags to assign to the resource. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) specifying the log group. Any `:*` suffix added by the API, denoting all CloudWatch Log Streams under the CloudWatch Log Group, is removed for greater compatibility with other AWS services that do not accept the suffix.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Cloudwatch Log Groups can be imported using the `name`, e.g.,

```
$ terraform import aws_cloudwatch_log_group.test_group yada
```
