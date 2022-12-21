---
subcategory: "GameLift"
layout: "aws"
page_title: "AWS: aws_gamelift_matchmaking_rule_set"
description: |-
  Provides a GameLift Matchmaking Rule Set .
---

# Resource: aws_gamelift_matchmaking_rule_set

Provides a GameLift Matchmaking Rule Set resources.

## Example Usage

```terraform
resource "aws_gamelift_matchmaking_rule_set" "example" {
  name          = %[1]q
  rule_set_body = jsonencodes({
	"name": "test",
	"ruleLanguageVersion": "1.0",
	"teams": [{
		"name": "alpha",
		"minPlayers": 1,
		"maxPlayers": 5
	}]
})
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the matchmaking rule set.
* `rule_set_body` - (Required) JSON encoded string containing rule set data.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Rule Set ID.
* `arn` - Rule Set ARN.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

GameLift Matchmaking Rule Sets  can be imported using the ID, e.g.,

```
$ terraform import aws_gamelift_matchmaking_rule_set.example <ruleset-id>
```
