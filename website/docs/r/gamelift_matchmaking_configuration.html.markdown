---
subcategory: "GameLift"
layout: "aws"
page_title: "AWS: aws_gamelift_matchmaking_configuration"
description: |-
  Provides a GameLift Matchmaking Configuration resource.
---

# Resource: aws_gamelift_matchmaking_configuration

Provides a GameLift Alias resource.

## Example Usage

```terraform

resource "aws_gamelift_game_session_queue" "example" {
	name         = example
	destinations = []
	
	player_latency_policy {
		maximum_individual_player_latency_milliseconds = 3
		policy_duration_seconds                        = 7
	}
	
	player_latency_policy {
		maximum_individual_player_latency_milliseconds = 10
	}
	
	timeout_in_seconds = 25
}

resource "aws_gamelift_matchmaking_rule_set" "example" {
	name          = example
	rule_set_body = jsonencode({
		"name": "test",
		"ruleLanguageVersion": "1.0",
		"teams": [{
			"name": "alpha",
			"minPlayers": 1,
			"maxPlayers": 5
		}]
	})
}

resource "aws_gamelift_matchmaking_configuration" "example" {
	name          = example
	acceptance_required = false
	custom_event_data = "pvp"
	game_session_data = "game_session_data"
	backfill_mode = "MANUAL"
	request_timeout_seconds = 30
	rule_set_name = aws_gamelift_matchmaking_rule_set.test.name
	game_session_queue_arns = [aws_gamelift_game_session_queue.test.arn]
	tags = {
		"key1" = "value1"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the matchmaking configuration
* `acceptance_required` - (Required) Specifies if the match that was created with this configuration must be accepted by matched players.
* `additional_player_count` - (Required) The number of player slots in a match to keep open for future players.
* `game_session_data` - (Required) A set of custom game session properties.
* `rule_set_name` - (Required) A rule set names for the matchmaking rule set to use with this configuration.
* `acceptance_timeout_seconds` - (Optional) The length of time (in seconds) to wait for players to accept a proposed match, if acceptance is required.
* `backfill_mode` - (Optional) The method used to backfill game sessions that are created with this matchmaking configuration.
* `custom_event_data` - (Optional) Information to be added to all events related to this matchmaking configuration.
* `description` - (Optional) A human-readable description of the matchmaking configuration.
* `flex_match_mode` - (Optional) Indicates whether this matchmaking configuration is being used with GameLift hosting or as a standalone matchmaking solution.
* `game_property` - (Optional) One or more custom game properties. See below.
* `game_session_queue_arns` - (Optional) The ARNs of the GameLift game session queue resources.
* `notification_target` - (Optional) An SNS topic ARN that is set up to receive matchmaking notifications.
* `request_timeout_seconds` - (Optional) The maximum duration, in seconds, that a matchmaking ticket can remain in process before timing out.
* `tags` - (Optional) Key-value map of resource tags. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.




### Nested Fields

#### `game_property`

* `key` - (Optional) A game property key
* `value` - (Optional) A game property value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Matchmaking Configuration ID.
* `arn` - Matchmaking Configuration ARN.
* `creation_time` - The time when the Matchmaking Configuration was created.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

GameLift Matchmaking Configurations can be imported using the ID, e.g.,

```
$ terraform import aws_gamelift_matchmaking_configuration.example <matchmakingconfiguration-id>
```
