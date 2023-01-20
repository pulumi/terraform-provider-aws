---
subcategory: "GameLift"
layout: "aws"
page_title: "AWS: aws_gamelift_alias"
description: |-
  Provides a GameLift Alias resource.
---

# Resource: aws_gamelift_alias

Provides a GameLift Alias resource.

## Example Usage

```terraform
resource "aws_gamelift_alias" "example" {
  name        = "example-alias"
  description = "Example Description"

  routing_strategy {
    message = "Example Message"
    type    = "TERMINAL"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the alias.
* `description` - (Optional) Description of the alias.
* `routing_strategy` - (Required) Specifies the fleet and/or routing type to use for the alias.
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### Nested Fields

#### `routing_strategy`

* `fleet_id` - (Optional) ID of the GameLift Fleet to point the alias to.
* `message` - (Optional) Message text to be used with the `TERMINAL` routing strategy.
* `type` - (Required) Type of routing strategyE.g., `SIMPLE` or `TERMINAL`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Alias ID.
* `arn` - Alias ARN.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

GameLift Aliases can be imported using the ID, e.g.,

```
$ terraform import aws_gamelift_alias.example <alias-id>
```
