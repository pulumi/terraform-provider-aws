---
subcategory: "IAM Access Analyzer"
layout: "aws"
page_title: "AWS: aws_accessanalyzer_analyzer"
description: |-
  Manages an Access Analyzer Analyzer
---

# Resource: aws_accessanalyzer_analyzer

Manages an Access Analyzer Analyzer. More information can be found in the [Access Analyzer User Guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/what-is-access-analyzer.html).

## Example Usage

### Account Analyzer

```terraform
resource "aws_accessanalyzer_analyzer" "example" {
  analyzer_name = "example"
}
```

### Organization Analyzer

```terraform
resource "aws_organizations_organization" "example" {
  aws_service_access_principals = ["access-analyzer.amazonaws.com"]
}

resource "aws_accessanalyzer_analyzer" "example" {
  depends_on = [aws_organizations_organization.example]

  analyzer_name = "example"
  type          = "ORGANIZATION"
}
```

## Argument Reference

The following arguments are required:

* `analyzer_name` - (Required) Name of the Analyzer.

The following arguments are optional:

* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `type` - (Optional) Type of Analyzer. Valid values are `ACCOUNT` or `ORGANIZATION`. Defaults to `ACCOUNT`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the Analyzer.
* `id` - Analyzer name.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Access Analyzer Analyzers can be imported using the `analyzer_name`, e.g.,

```
$ terraform import aws_accessanalyzer_analyzer.example example
```
