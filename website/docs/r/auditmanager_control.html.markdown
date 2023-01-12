---
subcategory: "Audit Manager"
layout: "aws"
page_title: "AWS: aws_auditmanager_control"
description: |-
  Resource for managing an AWS Audit Manager Control.
---

# Resource: aws_auditmanager_control

Resource for managing an AWS Audit Manager Control.

## Example Usage

### Basic Usage

```terraform
resource "aws_auditmanager_control" "example" {
  name = "example"

  control_mapping_sources {
    source_name          = "example"
    source_set_up_option = "Procedural_Controls_Mapping"
    source_type          = "MANUAL"
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Name of the control.
* `control_mapping_sources` - (Required) Data mapping sources. See [`control_mapping_sources`](#control_mapping_sources) below.

The following arguments are optional:

* `action_plan_instructions` - (Optional) Recommended actions to carry out if the control isn't fulfilled.
* `action_plan_title` - (Optional) Title of the action plan for remediating the control.
* `description` - (Optional) Description of the control.
* `tags` - (Optional) A map of tags to assign to the control. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `testing_information` - (Optional) Steps to follow to determine if the control is satisfied.

### control_mapping_sources

The following arguments are required:

* `source_name` - (Required) Name of the source.
* `source_set_up_option` - (Required) The setup option for the data source. This option reflects if the evidence collection is automated or manual. Valid values are `System_Controls_Mapping` (automated) and `Procedural_Controls_Mapping` (manual).
* `source_type` - (Required) Type of data source for evidence collection. If `source_set_up_option` is manual, the only valid value is `MANUAL`. If `source_set_up_option` is automated, valid values are `AWS_Cloudtrail`, `AWS_Config`, `AWS_Security_Hub`, or `AWS_API_Call`.

The following arguments are optional:

* `source_description` - (Optional) Description of the source.
* `source_frequency` - (Optional) Frequency of evidence collection. Valid values are `DAILY`, `WEEKLY`, or `MONTHLY`.
* `source_keyword` - (Optional) The keyword to search for in CloudTrail logs, Config rules, Security Hub checks, and Amazon Web Services API names. See [`source_keyword`](#source_keyword) below.
* `troubleshooting_text` - (Optional) Instructions for troubleshooting the control.

### source_keyword

The following arguments are required:

* `keyword_input_type` - (Required) Input method for the keyword. Valid values are `SELECT_FROM_LIST`.
* `keyword_value` - (Required) The value of the keyword that's used when mapping a control data source. For example, this can be a CloudTrail event name, a rule name for Config, a Security Hub control, or the name of an Amazon Web Services API call. See the [Audit Manager supported control data sources documentation](https://docs.aws.amazon.com/audit-manager/latest/userguide/control-data-sources.html) for more information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the control.
* `control_mapping_sources.*.source_id` - Unique identifier for the source.
* `id` - Unique identifier for the control.
* `type` - Type of control, such as a custom control or a standard control.

## Import

An Audit Manager Control can be imported using the `id`, e.g.,

```
$ terraform import aws_auditmanager_control.example abc123-de45
```
