package gamelift_test

import (
	"fmt"
	"regexp"
	"testing"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	//"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	//"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccGameLiftMatchmakingConfiguration_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_matchmaking_configuration.test"

	queueName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	ruleSetName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	additionalParameters := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, gamelift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameServerMatchmakingConfiguration_basic(rName, queueName, ruleSetName, additionalParameters, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameServerGroupExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "gamelift", regexp.MustCompile(`matchmakingconfiguration/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_event_data", "pvp"),
					resource.TestCheckResourceAttr(resourceName, "game_session_data", "game_session_data"),
					resource.TestCheckResourceAttr(resourceName, "acceptance_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "request_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "backfill_mode", gamelift.BackfillModeManual),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAWSGameliftMatchMakingConfigurationRuleSetBody() string {
	maxPlayers := int64(1)
	return fmt.Sprintf(`{
		"name": "test",
		"ruleLanguageVersion": "1.0",
		"teams": [{
			"name": "alpha",
			"minPlayers": 1,
			"maxPlayers": %[1]d
		}]
	}`, maxPlayers)
}

func testAccGameServerMatchmakingConfiguration_basic(rName string, queueName string, ruleSetName string, additionalParameters string, requestTimeoutSeconds int) string {
	backfillMode := gamelift.BackfillModeManual
	return fmt.Sprintf(`
resource "aws_gamelift_game_session_queue" "test" {
	name         = %[2]q
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

resource "aws_gamelift_matchmaking_rule_set" "test" {
	name          = %[3]q
	rule_set_body = <<RULE_SET_BODY
	%[4]s
	RULE_SET_BODY	
}

resource "aws_gamelift_matchmaking_configuration" "test" {
	name          = %[1]q
	acceptance_required = false
	custom_event_data = "pvp"
	game_session_data = "game_session_data"
	backfill_mode = %[7]q
	request_timeout_seconds = %[6]d
	rule_set_name = aws_gamelift_matchmaking_rule_set.test.name
	game_session_queue_arns = [aws_gamelift_game_session_queue.test.arn]
	%[5]s
}
`, rName, queueName, ruleSetName, testAccAWSGameliftMatchMakingConfigurationRuleSetBody(), additionalParameters, requestTimeoutSeconds, backfillMode)

}
