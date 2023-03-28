package gamelift_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfgamelift "github.com/hashicorp/terraform-provider-aws/internal/service/gamelift"
)

func TestAccMatchmakingRuleSet_basic(t *testing.T) {
	ctx := acctest.Context(t)

	var conf gamelift.MatchmakingRuleSet

	ruleSetName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_matchmaking_rule_set.test"
	maxPlayers := 5

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(t, gamelift.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, gamelift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGameServerGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMatchmakingRuleSetBasicConfig(ruleSetName, maxPlayers),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatchmakingRuleSetExists(ctx, resourceName, &conf),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "gamelift", regexp.MustCompile(`matchmakingruleset/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", ruleSetName),
					resource.TestCheckResourceAttr(resourceName, "rule_set_body", testAccMatchmakingRuleSetBody(maxPlayers)+"\n"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccMatchmakingRuleSet_disappears(t *testing.T) {
	ctx := acctest.Context(t)

	var conf gamelift.MatchmakingRuleSet

	ruleSetName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_matchmaking_rule_set.test"
	maxPlayers := 5

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(t, gamelift.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, gamelift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGameServerGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMatchmakingRuleSetBasicConfig(ruleSetName, maxPlayers),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMatchmakingRuleSetExists(ctx, resourceName, &conf),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfgamelift.ResourceMatchmakingRuleSet(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckMatchmakingRuleSetExists(ctx context.Context, n string, res *gamelift.MatchmakingRuleSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Gamelift Matchmaking Rule Set Name is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GameLiftConn()

		name := rs.Primary.Attributes["name"]
		out, err := conn.DescribeMatchmakingRuleSets(&gamelift.DescribeMatchmakingRuleSetsInput{
			Names: aws.StringSlice([]string{name}),
		})
		if err != nil {
			return err
		}
		ruleSets := out.RuleSets
		if len(ruleSets) == 0 {
			return fmt.Errorf("GameLift Matchmaking Rule Set %q not found", name)
		}

		*res = *ruleSets[0]

		return nil
	}
}

func testAccMatchmakingRuleSetBody(maxPlayers int) string {
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

func testAccMatchmakingRuleSetBasicConfig(rName string, maxPlayers int) string {
	return fmt.Sprintf(`
resource "aws_gamelift_matchmaking_rule_set" "test" {
  name          = %[1]q
  rule_set_body = <<RULE_SET_BODY
%[2]s
RULE_SET_BODY
}
`, rName, testAccMatchmakingRuleSetBody(maxPlayers))
}
